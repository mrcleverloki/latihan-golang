# Pernyataan Penggunaan AI (AI-USAGE)

Dokumen ini menjelaskan pemanfaatan kecerdasan buatan (Artificial Intelligence) dalam penyelesaian Tugas Praktikum Pertemuan 3 (Database & Repository Pattern).

## Alat AI yang Digunakan
- **Model / Tools**: Google Gemini

## Rincian Bagian yang Dibantu oleh AI
1. **Pemahaman Konsep Teori & Database Persistence**
   - Diskusi konsep *connection pooling* (`pgxpool`), pencegahan *SQL injection* melalui query berparameter, serta pemindahan operasi filter, sorting, dan pagination dari memori ke basis data.
2. **Perancangan Skema & Migrasi SQL**
   - Panduan penyusunan berkas migrasi `001_create_students.sql`, pembuatan *unique index* pada `LOWER(nim)`, dan penambahan indeks pencarian pada `LOWER(name)`.
3. **Scaffolding Pola Repository & Integrasi Layer**
   - Pemisahan kontrak *interface* dan implementasi basis data pada `app/repository/student_repository.go`.
   - Pemetaan *error* basis data (*unique constraint violation* 23505 dan `pgx.ErrNoRows`) menjadi *sentinel error* (`ErrNotFound`, `ErrDuplicate`) serta translasi status HTTP (404, 409, 503).
   - Pengaturan batas waktu operasi database menggunakan `context.WithTimeout`.

## Bagian yang Dikerjakan Mandiri
- Instalasi dan konfigurasi PostgreSQL serta Postgres.app pada macOS.
- Penyesuaian struktur folder proyek dan pengelolaan variabel environment (`.env`).
- Eksekusi migrasi tabel dan perakitan *dependency injection* pada `main.go`.
- Pengujian langsung skenario endpoint API via Thunder Client serta pengambilan bukti tangkapan layar (*screenshot*).
- Pengelolaan riwayat *commit* Git bertahap dan penyusunan laporan akhir PDF.