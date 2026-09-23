# Praktikum Pemrograman Backend - Pertemuan 6: RBAC & Ownership Authorization

Proyek RESTful API Manajemen Mahasiswa (*Students API*) berbasis Go, Fiber Framework, dan PostgreSQL. Modul ini mengimplementasikan sistem kontrol akses berlapis menggunakan **Role-Based Access Control (RBAC)** dan **Ownership-Based Authorization** untuk melindungi sumber daya data dari akses tidak sah.

---

## Fitur Utama & Mekanisme Keamanan

1. **Role-Based Access Control (RBAC):**
   - Pemetaan hak akses berbasis izin atomik (*permissions*): `student:list`, `student:read`, `student:create`, `student:update`, `student:delete`.
   - Pemisahan 3 tingkatan peran:
     - `user`: Tidak memiliki izin akses kolektif bawaan.
     - `staff`: Memiliki hak baca data kolektif (`student:list`, `student:read`).
     - `admin`: Memiliki seluruh hak akses CRUD.
   - Pengecekan izin via middleware Fiber `RequirePermission`.

2. **Ownership-Based Authorization:**
   - Perlindungan resource individual mahasiswa menggunakan kolom relasi kepemilikan `owner_id`.
   - Pengguna dengan role `user` hanya diizinkan membaca atau menyunting data miliknya sendiri (`owner_id == user_id`).
   - Mitigasi kerentanan *Insecure Direct Object References* (IDOR) pada layer service via `CanAccessStudent`.

3. **Autentikasi Stateless JWT:**
   - Header Authorization Bearer Token.
   - Claims memuat `user_id`, `username`, dan `role`.

---

## Skema Basis Data

### 1. Tabel users
Menyimpan data identitas kredensial dan peran pengguna:
- `id` (SERIAL PRIMARY KEY)
- `username` (VARCHAR(50), UNIQUE, NOT NULL)
- `email` (VARCHAR(100), UNIQUE, NOT NULL)
- `password` (VARCHAR(255), NOT NULL, Hashed bcrypt)
- `role` (VARCHAR(20), DEFAULT 'user')
- `is_active` (BOOLEAN, DEFAULT TRUE)
- `created_at` & `updated_at` (TIMESTAMPTZ)

### 2. Tabel students
Menyimpan data mahasiswa dengan referensi foreign key kepemilikan akun:
- `id` (SERIAL PRIMARY KEY)
- `nim` (VARCHAR(20), UNIQUE, NOT NULL)
- `name` (VARCHAR(100), NOT NULL)
- `grade` (VARCHAR(5), NOT NULL)
- `is_active` (BOOLEAN, DEFAULT TRUE)
- `owner_id` (INT, REFERENCES users(id) ON DELETE SET NULL)
- `created_at` & `updated_at` (TIMESTAMPTZ)

---

## Panduan Menjalankan Aplikasi

1. Setup Database PostgreSQL:
   psql -U postgres -c "CREATE DATABASE praktikum_backend;"

2. Konfigurasi Environment Variable (.env):
   APP_PORT=3000
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=praktikum_backend
   DB_SSLMODE=disable
   JWT_SECRET=supersecretkey

3. Jalankan Aplikasi:
   go run main.go

---

## Matriks Hak Akses Endpoint

| Method | Endpoint | Izin yang Dibutuhkan | Hak Akses Role | Keterangan |
| :--- | :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/login` | Publik | Semua | Penerbitan token akses JWT |
| `GET` | `/api/v1/students` | `student:list` | `staff`, `admin` | Menampilkan seluruh daftar mahasiswa |
| `GET` | `/api/v1/students/:id` | `student:read` / Ownership | `admin`, `staff`, `user` (pemilik) | Detail data mahasiswa |
| `POST` | `/api/v1/students` | `student:create` | `admin` | Penambahan data mahasiswa |
| `PATCH`| `/api/v1/students/:id` | `student:update` / Ownership | `admin`, `user` (pemilik) | Pembaruan data mahasiswa |
| `DELETE`| `/api/v1/students/:id` | `student:delete` | `admin` | Penghapusan data mahasiswa |
