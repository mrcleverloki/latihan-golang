# Praktikum Pemrograman Backend - Pertemuan 9: Transaction Management & Concurrency Control

Proyek RESTful API Manajemen Mahasiswa dan Kartu Rencana Studi (*Students & Courses API*) berbasis Go, Fiber Framework, dan PostgreSQL (`pgxpool`). Modul ini mengimplementasikan manajemen transaksi ACID terisolasi menggunakan abstraksi **`TxManager`**, kontrol konkurensi **Pessimistic Locking** (`SELECT ... FOR UPDATE`), **Optimistic Locking** (berbasis kolom `version`), serta pembuktian integritas data melalui pengujian beban (*load testing*).

---

## Fitur Utama & Pembaruan Modul 9

1. **Abstraksi Transaksi (`TxManager`):**
   - Mengisolasi dependensi driver database dari domain service menggunakan closure `WithTx(ctx, func(tx DBTX) error)`.
   - Menangani otomatisasi siklus hidup transaksi: `Begin`, `Commit`, dan `Rollback` saat runtime panic atau error.
   - Mekanisme retry otomatis berlandaskan jitter backoff acak untuk mengatasi error konkurensi transien PostgreSQL (Deadlock `40P01` dan Serialization Failure `40001`).

2. **Kontrol Konkurensi Tingkat Lanjut (Pessimistic & Optimistic Locking):**
   - **Pessimistic Locking (`LockStudent` & `LockCourses`):** Mengunci record mahasiswa dan baris kelas secara deterministik berurutan (`ORDER BY id FOR UPDATE`) guna mencegah *race condition* pendaftaran serentak dan mengeliminasi risiko siklus *deadlock* saat *swap* kelas berlawanan arah.
   - **Optimistic Locking (`UpdateQuota`):** Pembaruan kuota kelas memvalidasi keselarasan versi data secara atomik (`WHERE id = $1 AND version = $2`), mendeteksi konflik pembaruan serentak dengan error `VERSION_CONFLICT`.

3. **Penyelesaian 9 Puzzle Bugs:**
   - Memperbaiki kebocoran kuota (*overselling*) pada kueri pembaruan kursi (`< quota` vs `<= quota`).
   - Validasi pembatalan pendaftaran fiktif melalui pengecekan `RowsAffected()`.
   - Mengalihkan eksekusi pelepasan kursi di service swap agar terikat penuh di dalam konteks transaksi (`repo.DecrementEnrolled`).
   - Sinkronisasi kenaikan versi (`version = version + 1`) pada operasi kuota.

4. **Pengujian Beban Serentak (*Load Testing Suite*):**
   - Perkakas CLI pengujian beban multi-goroutine di `tools/loadtest` yang mensimulasikan gerbang pelepasan serentak (*barrier synchronization* via channel Go).
   - Validasi ketat pada skenario *Crowd* (perebutan kuota terbatas), *Credits* (penegakan batas 24 SKS), dan *Swap* (perpindahan kelas silang dua arah).

---

## Skema Basis Data & Relasi Tambahan

### 1. Tabel `courses`
Menyimpan katalog mata kuliah:
- `id` (SERIAL PRIMARY KEY)
- `code` (VARCHAR(20), UNIQUE, NOT NULL)
- `name` (VARCHAR(100), NOT NULL)
- `credits` (INT, NOT NULL, CHECK > 0)
- `quota` (INT, NOT NULL, CHECK > 0)
- `enrolled_count` (INT, DEFAULT 0, CHECK >= 0)
- `version` (INT, DEFAULT 1, NOT NULL)
- `created_at` (TIMESTAMPTZ, DEFAULT NOW())

### 2. Tabel `enrollments`
Pencatatan relasi pendaftaran KRS mahasiswa:
- `id` (SERIAL PRIMARY KEY)
- `course_id` (INT, REFERENCES courses(id) ON DELETE CASCADE)
- `user_id` (INT, REFERENCES users(id) ON DELETE CASCADE)
- `created_at` (TIMESTAMPTZ, DEFAULT NOW())
- Constraint: `UNIQUE(course_id, user_id)` (mencegah pendaftaran ganda).

---

## Matriks Endpoint Tambahan (KRS & Transaksi)

| Method | Endpoint | Izin / Proteksi | Hak Akses Role | Keterangan |
| :--- | :--- | :--- | :--- | :--- |
| GET | /api/v1/courses | Bearer Auth | Semua Role | Katalog kelas yang tersedia |
| GET | /api/v1/courses/mine | Bearer Auth | Mahasiswa (User) | Rincian KRS dan total SKS yang diambil |
| GET | /api/v1/courses/:id | Bearer Auth | Semua Role | Detail kelas |
| POST | /api/v1/courses/:id/enroll | Bearer Auth | Mahasiswa (User) | Pendaftaran kelas (Transaksi ACID & Lock) |
| DELETE | /api/v1/courses/:id/enroll | Bearer Auth | Mahasiswa (User) | Pembatalan kelas (Drop) |
| POST | /api/v1/enrollments/swap | Bearer Auth | Mahasiswa (User) | Perpindahan kelas silang (Deadlock-Free Swap) |
| PATCH | /api/v1/courses/:id | course:manage | Admin, Staff | Pembaruan kuota kelas (Optimistic Lock) |

---

## Panduan Menjalankan Migrasi & Uji Beban

1. Eksekusi Skema dan Data Uji Beban:
   psql -U postgres -d praktikum_backend -f migrations/006_courses.sql
   psql -U postgres -d praktikum_backend -f migrations/007_seed_loadtest.sql

2. Jalankan Server:
   go run main.go

3. Eksekusi Pengujian Beban (*Loadtest CLI*):
   # Mode 1: Uji Perebutan Kuota Terbatas (60 Mhs vs 5 Kuota IF399)
   go run ./tools/loadtest -mode crowd

   # Mode 2: Uji Batas SKS (10 Mhs x 12 Kelas Serentak, Batas 24 SKS)
   go run ./tools/loadtest -mode credits

   # Mode 3: Uji Bebas Deadlock (20 Pasang Mhs Swap Berlawanan Arah)
   go run ./tools/loadtest -mode swap