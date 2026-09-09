package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"api-students/app/model"
)

// Sentinel error untuk diterjemahkan oleh handler
var (
	ErrNotFound  = errors.New("data mahasiswa tidak ditemukan")
	ErrDuplicate = errors.New("NIM mahasiswa sudah terdaftar")
)

// StudentRepository adalah kontrak penyimpanan data mahasiswa
type StudentRepository interface {
	FindAll(ctx context.Context, q model.ListQuery) ([]model.Student, int, error)
	FindByID(ctx context.Context, id int) (model.Student, error)
	Create(ctx context.Context, s model.Student) (model.Student, error)
	Update(ctx context.Context, s model.Student) (model.Student, error)
	Delete(ctx context.Context, id int) error
}

// Whitelist kolom untuk mencegah SQL Injection pada klausa ORDER BY
var kolomUrutStudent = map[string]string{
	"id":         "id",
	"nim":        "nim",
	"name":       "name",
	"grade":      "grade",
	"created_at": "created_at",
}

type studentPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewStudentRepository(pool *pgxpool.Pool) StudentRepository {
	return &studentPostgresRepository{pool: pool}
}

// buildFilter menyusun klausa WHERE dinamis dan argumen berparameter
func buildStudentFilter(q model.ListQuery) (string, []any) {
	where := " WHERE 1=1"
	args := []any{}

	if q.Search != "" {
		where += fmt.Sprintf(" AND (name ILIKE $%d OR nim ILIKE $%d)", len(args)+1, len(args)+2)
		args = append(args, "%"+q.Search+"%", "%"+q.Search+"%")
	}

	if q.IsActive != nil {
		where += fmt.Sprintf(" AND is_active = $%d", len(args)+1)
		args = append(args, *q.IsActive)
	}

	return where, args
}

func (r *studentPostgresRepository) FindAll(ctx context.Context, q model.ListQuery) ([]model.Student, int, error) {
	where, args := buildStudentFilter(q)

	// 1. Hitung total data untuk meta pagination
	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM students"+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("menghitung total mahasiswa: %w", err)
	}

	// 2. Tentukan sorting yang aman melalui whitelist
	kolom, ok := kolomUrutStudent[q.Sort]
	if !ok {
		kolom = "id"
	}

	arah := "ASC"
	if q.Order == "desc" {
		arah = "DESC"
	}

	sqlText := fmt.Sprintf(`
		SELECT id, nim, name, grade, is_active, created_at
		FROM students%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`,
		where, kolom, arah, len(args)+1, len(args)+2,
	)
	args = append(args, q.Limit, q.Offset())

	rows, err := r.pool.Query(ctx, sqlText, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("mengambil daftar mahasiswa: %w", err)
	}
	defer rows.Close()

	hasil := []model.Student{}
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("membaca baris mahasiswa: %w", err)
		}
		hasil = append(hasil, s)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterasi baris mahasiswa: %w", err)
	}

	return hasil, total, nil
}

func (r *studentPostgresRepository) FindByID(ctx context.Context, id int) (model.Student, error) {
	var s model.Student
	err := r.pool.QueryRow(ctx, `
		SELECT id, nim, name, grade, is_active, created_at
		FROM students WHERE id = $1`,
		id,
	).Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		return model.Student{}, fmt.Errorf("mengambil mahasiswa by id: %w", err)
	}

	return s, nil
}

func (r *studentPostgresRepository) Create(ctx context.Context, s model.Student) (model.Student, error) {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO students (nim, name, grade, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`,
		s.NIM, s.Name, s.Grade, s.IsActive,
	).Scan(&s.ID, &s.CreatedAt)

	if err != nil {
		if isUniqueViolation(err) {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, fmt.Errorf("menyimpan mahasiswa: %w", err)
	}

	return s, nil
}

func (r *studentPostgresRepository) Update(ctx context.Context, s model.Student) (model.Student, error) {
	err := r.pool.QueryRow(ctx, `
		UPDATE students 
		SET nim = $1, name = $2, grade = $3, is_active = $4
		WHERE id = $5
		RETURNING id, nim, name, grade, is_active, created_at`,
		s.NIM, s.Name, s.Grade, s.IsActive, s.ID,
	).Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		if isUniqueViolation(err) {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, fmt.Errorf("memperbarui mahasiswa: %w", err)
	}

	return s, nil
}

func (r *studentPostgresRepository) Delete(ctx context.Context, id int) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM students WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("menghapus mahasiswa: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// isUniqueViolation mendeteksi error PostgreSQL 23505 (unique violation)
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}