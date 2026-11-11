# Deklarasi Penggunaan AI (AI-USAGE)

Sesuai dengan prinsip integritas akademik dan transparansi pedoman perkuliahan, dokumen ini mencatat pemanfaatan alat bantu kecerdasan buatan (*Artificial Intelligence*) selama pengerjaan proyek Modul 11 (Asynchronous Notifications, Transactional Outbox Pattern, Asynq Task Queue, Idempotency Verification, and Multi-Role Orchestration).

---

## 1. Alat Bantu yang Digunakan
- **Nama Alat:** Gemini AI (Google)
- **Model:** Gemini
- **Media Akses:** Obrolan Interaktif

---

## 2. Rincian Pemanfaatan Alat Bantu

| Modul / Komponen | Bentuk Kontribusi AI | Tingkat Bantuan |
| :--- | :--- | :--- |
| **Analisis & Perbaikan Puzzle Bug Modul (Tugas D.1)** | Bantuan identifikasi 8 bug bawaan modul, mencakup perbaikan transaksi `DBTX` pada outbox repository, locking `FOR UPDATE SKIP LOCKED`, opsi timeout Asynq, integrasi signature handler `asynq.Task`, pembungkusan `asynq.SkipRetry`, serta perbaikan mismatch nama queue pada `config/queue.go`. | Sedang |
| **Penyusunan Arsitektur Transaksional Outbox (`outbox_repository.go` & `outbox_relay.go`)** | Panduan struktur *polling consumer*, locking batch outbox, dan penanganan graceful stop relay via konteks `signal.NotifyContext`. | Sedang |
| **Refactoring Lifecycle & Multi-Role Entrypoint (`main.go` & `worker.go`)** | Bantuan penyusunan abstraksi CLI flag `-role` (`api`, `worker`, `all`), inisialisasi `asynq.Inspector`, serta sinkronisasi graceful shutdown antara Fiber HTTP server dan Asynq worker server. | Sedang |
| **Pengembangan Tool Uji Idempotensi (`tools/redeliver/main.go`)** | Perancangan CLI tool injeksi task Asynq duplikat dengan penomoran `TaskID` dinamis berbasis timestamp nano guna menghindari benturan `ErrTaskIDConflict` di Redis. | Sedang |
| **Perumusan Skenario Uji Ketahanan Sistem (Skenario C.1 - C.5)** | Panduan langkah verifikasi pengujian terminal (simulasi crash Redis, temporary failure, penolakan domain invalid, dan pengujian idempotensi). | Rendah |

---

## 3. Pekerjaan Mandiri (Human Contribution)
- Penyiapan dan pemeliharaan service lokal macOS via Homebrew (PostgreSQL, MongoDB, dan Redis).
- Pemeriksaan dan penyesuaian dependensi kode eksisting di `route/route.go` dan `main.go` agar tetap mempertahankan integrasi `StudentService`, `PortfolioService` (MongoDB), dan `CourseService`.
- Pembuatan dan migrasi skema tabel `outbox` dan `notifications` di PostgreSQL.
- Eksekusi manual seluruh skenario pengujian Bagian C (C.1 sampai C.5) menggunakan terminal macOS (zsh), cURL, dan psql.
- Debugging kendala autentikasi (token JWT expired, penyesuaian header `Content-Type: application/json`, dan validasi password policy).
- Pengelolaan version control git (penyusunan micro-commit terstruktur per modul, pembaruan dokumentasi teknis `README.md` dan `AI-USAGE.md`, serta sinkronisasi branch utama).