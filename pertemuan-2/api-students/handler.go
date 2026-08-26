package main

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// In-memory database (data hilang saat server mati)
var students []Student
var nextID = 1

// Helper internal: mencari index mahasiswa berdasarkan ID
func findStudentIndex(id int) int {
	for i := range students {
		if students[i].ID == id {
			return i
		}
	}
	return -1
}

// Helper internal: validasi parameter ID di path URL
func paramID(c *fiber.Ctx) (int, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

// 1. GET /api/v1/students (List dengan Search, Filter, Sort, & Pagination)
func listStudents(c *fiber.Ctx) error {
	q := parseListQuery(c)

	// Saring data (Search & Filter)
	hasil := []Student{}
	for _, s := range students {
		if q.IsActive != nil && s.IsActive != *q.IsActive {
			continue
		}
		if q.Search != "" && !strings.Contains(strings.ToLower(s.Name), strings.ToLower(q.Search)) {
			continue
		}
		hasil = append(hasil, s)
	}

	// Pengurutan (Sort)
	sort.SliceStable(hasil, func(i, j int) bool {
		var lebihKecil bool
		switch q.Sort {
		case "nim":
			lebihKecil = hasil[i].NIM < hasil[j].NIM
		case "name":
			lebihKecil = hasil[i].Name < hasil[j].Name
		case "grade":
			lebihKecil = hasil[i].Grade < hasil[j].Grade
		case "created_at":
			lebihKecil = hasil[i].CreatedAt.Before(hasil[j].CreatedAt)
		default:
			lebihKecil = hasil[i].ID < hasil[j].ID
		}
		if q.Order == "desc" {
			return !lebihKecil
		}
		return lebihKecil
	})

	// Paginasi
	total := len(hasil)
	totalPages := (total + q.Limit - 1) / q.Limit
	mulai := (q.Page - 1) * q.Limit
	if mulai > total {
		mulai = total
	}
	akhir := mulai + q.Limit
	if akhir > total {
		akhir = total
	}

	return okList(c, "daftar mahasiswa berhasil diambil", hasil[mulai:akhir], &Meta{
		Page:       q.Page,
		Limit:      q.Limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

// 2. GET /api/v1/students/:id (Ambil 1 data)
func getStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	idx := findStudentIndex(id)
	if idx == -1 {
		return fail(c, fiber.StatusNotFound, "data mahasiswa tidak ditemukan")
	}

	return ok(c, "data mahasiswa ditemukan", students[idx])
}

// 3. POST /api/v1/students (Tambah data)
func createStudent(c *fiber.Ctx) error {
	var req CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	errs := map[string]string{}
	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	if req.NIM == "" {
		errs["nim"] = "wajib diisi"
	}
	if req.Name == "" {
		errs["name"] = "wajib diisi"
	}
	if req.Grade < 0 || req.Grade > 4.0 {
		errs["grade"] = "harus berada di rentang 0.0 - 4.0"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	// Validasi NIM unik -> Status 409 Conflict
	for _, s := range students {
		if s.NIM == req.NIM {
			return fail(c, fiber.StatusConflict, "NIM sudah terdaftar")
		}
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	baru := Student{
		ID:        nextID,
		NIM:       req.NIM,
		Name:      req.Name,
		Grade:     req.Grade,
		IsActive:  isActive,
		CreatedAt: time.Now(),
	}
	students = append(students, baru)
	nextID++

	return created(c, "mahasiswa berhasil ditambahkan", baru, "/api/v1/students/"+strconv.Itoa(baru.ID))
}

// 4. PUT /api/v1/students/:id (Ganti seluruh field)
func replaceStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	idx := findStudentIndex(id)
	if idx == -1 {
		return fail(c, fiber.StatusNotFound, "data mahasiswa tidak ditemukan")
	}

	var req ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	errs := map[string]string{}
	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	if req.NIM == "" {
		errs["nim"] = "wajib diisi pada PUT"
	}
	if req.Name == "" {
		errs["name"] = "wajib diisi pada PUT"
	}
	if req.Grade < 0 || req.Grade > 4.0 {
		errs["grade"] = "harus berada di rentang 0.0 - 4.0"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	// Cek jika NIM baru bentrok dengan mahasiswa lain
	for _, s := range students {
		if s.NIM == req.NIM && s.ID != id {
			return fail(c, fiber.StatusConflict, "NIM sudah digunakan oleh mahasiswa lain")
		}
	}

	students[idx].NIM = req.NIM
	students[idx].Name = req.Name
	students[idx].Grade = req.Grade
	students[idx].IsActive = req.IsActive

	return ok(c, "data mahasiswa berhasil diganti seluruhnya", students[idx])
}

// 5. PATCH /api/v1/students/:id (Update sebagian field yang dikirim saja)
func patchStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	idx := findStudentIndex(id)
	if idx == -1 {
		return fail(c, fiber.StatusNotFound, "data mahasiswa tidak ditemukan")
	}

	var req PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	// Kalau body kosong sama sekali
	if req.NIM == nil && req.Name == nil && req.Grade == nil && req.IsActive == nil {
		return fail(c, fiber.StatusBadRequest, "tidak ada field yang diubah")
	}

	if req.NIM != nil {
		nim := strings.TrimSpace(*req.NIM)
		if nim == "" {
			return failValidation(c, map[string]string{"nim": "tidak boleh kosong"})
		}
		for _, s := range students {
			if s.NIM == nim && s.ID != id {
				return fail(c, fiber.StatusConflict, "NIM sudah digunakan oleh mahasiswa lain")
			}
		}
		students[idx].NIM = nim
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return failValidation(c, map[string]string{"name": "tidak boleh kosong"})
		}
		students[idx].Name = name
	}

	if req.Grade != nil {
		if *req.Grade < 0 || *req.Grade > 4.0 {
			return failValidation(c, map[string]string{"grade": "harus berada di rentang 0.0 - 4.0"})
		}
		students[idx].Grade = *req.Grade
	}

	if req.IsActive != nil {
		students[idx].IsActive = *req.IsActive
	}

	return ok(c, "data mahasiswa berhasil diperbarui sebagian", students[idx])
}

// 6. DELETE /api/v1/students/:id (Hapus data -> 204 No Content)
func deleteStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	idx := findStudentIndex(id)
	if idx == -1 {
		return fail(c, fiber.StatusNotFound, "data mahasiswa tidak ditemukan")
	}

	students = append(students[:idx], students[idx+1:]...)
	return noContent(c)
}