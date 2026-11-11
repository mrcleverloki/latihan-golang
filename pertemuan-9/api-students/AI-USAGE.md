# Deklarasi Penggunaan AI (AI-USAGE)

Sesuai dengan prinsip integritas akademik dan transparansi pedoman perkuliahan, dokumen ini mencatat pemanfaatan alat bantu kecerdasan buatan (*Artificial Intelligence*) selama pengerjaan proyek Modul 9 (Transaction Management, Pessimistic/Optimistic Concurrency Control, and Load Testing).

---

## 1. Alat Bantu yang Digunakan
- **Nama Alat:** Gemini AI (Google)
- **Model:** Gemini
- **Media Akses:** Obrolan Interaktif

---

## 2. Rincian Pemanfaatan Alat Bantu

| Modul / Komponen | Bentuk Kontribusi AI | Tingkat Bantuan |
| :--- | :--- | :--- |
| **Abstraksi `TxManager` (`app/repository/tx.go`)** | Bantuan penyusunan implementasi closure transaksi atomik, integrasi interface `DBTX`, pemetaan error transien PostgreSQL (`40P01`, `40001`), dan logika retry berlandaskan exponential backoff dengan jitter acak. | Sedang |
| **Penyelidikan 9 Puzzle Bugs (`app/repository/`, `app/service/`)** | Bantuan pelacakan dan pembetulan cacat konkurensi: kondisi batas kuota (`< quota`), evaluasi `cmd.RowsAffected()` pada penghapusan pendaftaran, isolasi transaksi `repo.DecrementEnrolled`, dan inkrementasi `version` untuk *optimistic locking*. | Tinggi |
| **Kontrol Konkurensi & Pencegahan Deadlock (`course_service.go`, `course_repository.go`)** | Formulasi strategi penguncian deterministik berurutan `ORDER BY id FOR UPDATE` pada metode `LockCourses` dan penguncian baris akun mahasiswa via `LockStudent`. | Sedang |
| **Perkakas Pengujian Beban (*Loadtest Tool*)** | Panduan pembuatan skrip pengujian berbasis *goroutine* dan barrier synchronization channel untuk skenario `crowd`, `credits`, dan `swap`. | Rendah |
| **Pemecahan Masalah Eksekusi & Konfigurasi** | Bantuan resolusi impor dependensi yang tidak terpakai pada compiler Go serta penyesuaian target migrasi database `praktikum_backend`. | Rendah |

---

## 3. Pekerjaan Mandiri (Human Contribution)
- Penyiapan database PostgreSQL lokal `praktikum_backend` dan eksekusi skrip DDL migrasi tabel `courses`, `enrollments`, serta skrip pengisian seed data beban `007_seed_loadtest.sql`.
- Penulisan aturan murni pembatasan SKS (*pure rule functions*) pada `app/service/course_rules.go` dan perancangan payload DTO di `app/model/course.go`.
- Perakitan dependency injection di `main.go` untuk menghubungkan `TxManager` ke `CourseService` dan konfigurasi router Fiber.
- Eksekusi langsung uji beban konkurensi di lingkungan terminal lokal macOS pada ketiga mode (`crowd`, `credits`, `swap`) serta verifikasi integritas data akhir di PostgreSQL (memastikan tidak terjadi ketidaksesuaian antara `enrolled_count` dan total baris `enrollments`).
- Pengelolaan version control git (pembuatan micro-commit atomik, pembaruan `go.mod`/`go.sum`, serta sinkronisasi branch utama).