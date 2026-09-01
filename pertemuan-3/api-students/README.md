# Praktikum Pemrograman Backend Lanjut - Pertemuan 3: Database & Repository Pattern

Proyek REST API Manajemen Mahasiswa (*Students API*) menggunakan Go, Fiber Framework, PostgreSQL, dan arsitektur Repository Pattern.

## Skema Basis Data
Tabel `students` dibuat menggunakan migrasi SQL pada `migrations/001_create_students.sql`:
- `id` (SERIAL PRIMARY KEY)
- `nim` (VARCHAR(20), NOT NULL, UNIQUE via Lowercase Index)
- `name` (VARCHAR(100), NOT NULL)
- `grade` (VARCHAR(5), NOT NULL)
- `is_active` (BOOLEAN, DEFAULT TRUE)
- `created_at` (TIMESTAMPTZ, DEFAULT NOW())

Indeks:
- `students_nim_lower_key` ON `LOWER(nim)` (UNIQUE)
- `students_name_lower_idx` ON `LOWER(name)` (Index pencarian ILIKE)

## Panduan Menyiapkan Basis Data dari Nol
1. Pastikan PostgreSQL sudah aktif di lokal Anda.
2. Buat database baru:
   ```bash
   psql -U postgres -c "CREATE DATABASE praktikum_backend;"