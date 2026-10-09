# Praktikum Pemrograman Backend - Pertemuan 7: Advanced API Design

Proyek RESTful API Manajemen Mahasiswa (*Students API*) berbasis Go, Fiber Framework, dan PostgreSQL. Modul ini merupakan penyempurnaan arsitektural dari modul sebelumnya dengan menerapkan **Centralized Error Handling**, **Declarative Validation**, **Keyset (Cursor-Based) Pagination**, dan **Content Negotiation**.

---

## Fitur Utama & Pembaruan Modul 7

1. **Sentralisasi Penanganan Error (Centralized Error Handling):**
   - Eliminasi fungsi `helper.Fail` pada seluruh layer service dan middleware.
   - Pemanfaatan custom error struct `AppError` yang mengimplementasikan interface standar Go `error` dan method `Unwrap()`.
   - Pola *fail-closed* di mana seluruh error diteruskan ke middleware terpusat `ErrorHandler` di `config/app.go` untuk menghasilkan respons seragam bertipe `ErrorResponse` (`success`, `code`, `message`, `request_id`).

2. **Validasi Deklaratif (Declarative Validation):**
   - Penggantian validasi manual imperatif dengan tag struct deklaratif menggunakan pustaka `go-playground/validator/v10`.
   - Pendaftaran custom validator:
     - `nim_format`: Memvalidasi bahwa NIM harus tepat terdiri dari 9 karakter digit numerik.
     - `strongpassword`: Memvalidasi kompleksitas kata sandi pengguna (minimal 8 karakter, memuat huruf besar, huruf kecil, angka, dan simbol).
   - Penanganan pembaruan parsial (`PATCH`) menggunakan pointer field (`*string`, `*bool`) berpadu dengan tag `omitnil`.

3. **Keyset (Cursor-Based) Pagination:**
   - Menggantikan offset pagination tradisional guna mengatasi penurunan performa O(N) dan anomali *window drift* pada data berskala besar.
   - Menggunakan filtering predikat tuple `(created_at, id) < ($1, $2)` dengan kursor terenkode Base64.
   - Didukung oleh indeks komposit b-tree PostgreSQL `idx_students_created_at_id_desc` untuk pencarian berbasis *index scan* O(log N).

4. **Content Negotiation (RFC 9110 & RFC 4180):**
   - Evaluasi header `Accept` pada endpoint pengambilan koleksi mahasiswa (`GET /api/v1/students`).
   - Format `application/json` merespons payload data terpaginasi beserta metadata cursor.
   - Format `text/csv` mengalirkan data baris tabel secara streaming menggunakan `csv.NewWriter`.
   - Permintaan format yang tidak didukung secara otomatis ditolak dengan status HTTP `406 Not Acceptable`.

5. **Mekanisme Keamanan Lanjutan (RBAC & Rate Limiting):**
   - Autentikasi stateless berbasis JWT Bearer Token dengan penanganan refresh token rotasional.
   - Kontrol akses berbasis izin atomik (`RequirePermission`) dan proteksi endpoint autentikasi menggunakan `LoginRateLimiter`.

---

## Skema Basis Data & Indeks

### 1. Tabel users
Menyimpan kredensial dan peran pengguna:
- `id` (SERIAL PRIMARY KEY)
- `username` (VARCHAR(50), UNIQUE, NOT NULL)
- `email` (VARCHAR(100), UNIQUE, NOT NULL)
- `password` (VARCHAR(255), NOT NULL, Hashed bcrypt)
- `role` (VARCHAR(20), DEFAULT 'user')
- `is_active` (BOOLEAN, DEFAULT TRUE)
- `created_at` & `updated_at` (TIMESTAMPTZ)

### 2. Tabel students
Menyimpan data entitas mahasiswa:
- `id` (SERIAL PRIMARY KEY)
- `nim` (VARCHAR(20), UNIQUE, NOT NULL)
- `name` (VARCHAR(100), NOT NULL)
- `grade` (VARCHAR(5), NOT NULL)
- `is_active` (BOOLEAN, DEFAULT TRUE)
- `owner_id` (INT, REFERENCES users(id) ON DELETE SET NULL)
- `created_at` & `updated_at` (TIMESTAMPTZ)

### 3. Indeks Khusus (Modul 7)
- `idx_students_created_at_id_desc` pada tabel `students(created_at DESC, id DESC)` untuk mengoptimalkan keyset pagination.

---

## Panduan Menjalankan Aplikasi

1. Setup Database & Migrasi Indeks:
   psql -U postgres -d praktikum_backend -f migrations/005_cursor_index.sql

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

| Method | Endpoint | Izin / Proteksi | Hak Akses Role | Keterangan |
| :--- | :--- | :--- | :--- | :--- |
| POST | /api/v1/auth/login | Rate Limiter | Publik | Autentikasi & penerbitan token |
| POST | /api/v1/auth/register | Validasi Deklaratif | Publik | Registrasi akun baru |
| GET | /api/v1/students | student:list + Conneg | staff, admin | Daftar mahasiswa (Keyset Pagination & CSV/JSON) |
| GET | /api/v1/students/:id | student:read / Ownership | admin, staff, user (pemilik) | Detail data mahasiswa |
| POST | /api/v1/students | student:create | admin | Tambah mahasiswa (Validasi deklaratif & NIM kustom) |
| PATCH | /api/v1/students/:id | student:update / Ownership | admin, user (pemilik) | Pembaruan data parsial (omitnil) |
| DELETE | /api/v1/students/:id | student:delete | admin | Penghapusan data mahasiswa |