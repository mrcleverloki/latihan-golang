# Deklarasi Penggunaan AI (AI-USAGE)

Sesuai dengan integritas akademik dan pedoman tugas praktikum, dokumen ini mencatat rincian transparansi penggunaan alat bantu kecerdasan buatan (Artificial Intelligence) selama penyelesaian proyek Modul 7 (Advanced API Design: Centralized Error Handling, Declarative Validation, Keyset Pagination, & Content Negotiation).

---

## 1. Alat Bantu yang Digunakan
- **Nama Alat:** Gemini AI (Google)
- **Model:** Gemini
- **Media Akses:** Obrolan Interaktif Web

---

## 2. Rincian Pemanfaatan Alat Bantu

| Modul / Komponen | Bentuk Kontribusi AI | Tingkat Bantuan |
| :--- | :--- | :--- |
| **Penyelidikan 9 Bug Modul (`config/`, `service/`, `helper/`)** | Bantuan analisis akar masalah (*root cause*) dan perbaikan logika (status code logger Fiber, logika inversi password validator, query order keyset pagination, dan fail-open error translation). | Tinggi |
| **Sentralisasi Error Handling (`helper/errors.go`, `config/app.go`)** | Bantuan penyusunan arsitektur `AppError` (implementasi interface `error` dan `Unwrap()`), standarisasi payload respons error, serta eliminasi seluruh `helper.Fail` di service dan middleware. | Sedang |
| **Validasi Deklaratif (`helper/validator.go`, `model/`)** | Panduan integrasi pustaka `go-playground/validator/v10`, konfigurasi tag custom `nim_format` dan `strongpassword`, serta penanganan pointer `omitnil` untuk operasi `PATCH`. | Sedang |
| **Keyset Pagination & Indexing (`repository/`, `migrations/`)** | Panduan penyusunan kueri predikat tuple `(created_at, id) < ($1, $2)`, pembuatan DDL indeks komposit b-tree, dan analisis optimasi kueri melalui `EXPLAIN ANALYZE`. | Sedang |
| **Content Negotiation (`helper/negotiate.go`, `service/`)** | Formulasi evaluasi header `Accept` (RFC 9110), streaming format `text/csv` (RFC 4180), dan penolakan format tidak didukung via status `406 Not Acceptable`. | Sedang |
| **Skenario Pengujian API** | Panduan penyusunan urutan pengujian bertahap di Thunder Client (POST/PATCH validasi gagal 422, cursor pagination halaman 1-2, dan negosiasi format respons CSV/JSON). | Rendah |

---

## 3. Pekerjaan Mandiri (Human Contribution)
- Migrasi dan eksekusi skema basis data PostgreSQL lokal (pembuatan indeks komposit `idx_students_created_at_id_desc` via CLI `psql`).
- Refaktorisasi menyeluruh pada unit testing (`student_rules_test.go`) dan pembersihan file rule lama yang terdepresiasi.
- Verifikasi bebas `helper.Fail` di seluruh lapisan arsitektur melalui pengujian `grep -rn "helper.Fail" .`.
- Eksekusi langsung seluruh skenario HTTP testing via Thunder Client pada database `praktikum_backend` serta dokumentasi tangkapan layar respons status code (`200 OK`, `201 Created`, `406 Not Acceptable`, `422 Unprocessable Entity`).
- Pengujian langsung efisiensi indeks PostgreSQL menggunakan kueri `EXPLAIN ANALYZE` (`Index Scan`).
- Pengelolaan version control git (commit berkala bertema atomik dan sinkronisasi branch repository).
