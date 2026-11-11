# Praktikum Pemrograman Backend - Pertemuan 11: Asynchronous Notifications & Transactional Outbox Pattern

Proyek RESTful API Manajemen Mahasiswa (*Students API*) berbasis Go, Fiber Framework, PostgreSQL (`pgxpool`), MongoDB, dan Redis (`asynq`). 

Modul 11 memperluas arsitektur sistem menjadi terdistribusi dan *event-driven* melalui implementasi pola **Transactional Outbox**, pemrosesan latar belakang (*background workers*) menggunakan Asynq dan Redis, penjaminan idempoten pengiriman notifikasi (*idempotent notification delivery*), penanganan *transient & permanent failure* dengan antrean *Dead Letter Queue* (DLQ), serta pemisahan peran eksekusi biner (*multi-role entrypoint*).

---

## Fitur Utama & Pembaruan Modul 11

1. **Transactional Outbox Pattern (PostgreSQL):**
   - Penjaminan atomisitas data transaksi bisnis (pendaftaran dan pembatalan kelas) dengan pencatatan event outbox (`outbox`) dalam satu transaksi database tunggal (`pgx.Tx`).
   - Mencegah fenomena *dual-write problem*: data bisnis tidak akan tersimpan jika event gagal ditulis, dan event tidak akan dibuat jika transaksi bisnis rollback.
   - Pemanfaatan *row-level locking* `FOR UPDATE SKIP LOCKED` pada kueri relay untuk mencegah *race condition* dan pengiriman duplikat saat dijalankan dengan relai konkurensi tinggi.

2. **Outbox Relay Service:**
   - Polling berkala yang membaca event outbox berstatus pending (`published_at IS NULL`) dan mengirimkannya ke broker antrean Redis menggunakan Asynq.
   - Penandaan timestamp `published_at` secara transaksional setelah task berhasil diterima broker Redis.
   - Toleransi kegagalan antrean: transaksi pendaftaran tetap sukses meskipun broker Redis sedang tidak aktif (*downtime*).

3. **Background Worker & Antrean Asynq:**
   - Penguraian beban kerja notifikasi (*offloading*) ke proses latar belakang dengan konfigurasi antrean dedicated `notifications`.
   - Penerapan *retry policy* bertahap (*exponential backoff*) untuk kegagalan sementara (*transient failure*) seperti server email sibuk.
   - Penerapan `asynq.SkipRetry` untuk kegagalan permanen (*permanent failure*, seperti domain email invalid) yang langsung menghentikan retry dan memindahkan task ke antrean arsip (*Dead Letter Queue / archived*).

4. **Penjaminan Idempotensi Notifikasi (Consumer-Side):**
   - Mekanisme klaim idempoten (`Claim`) pada tabel `notifications` berbasis `outbox_id` (unique constraint).
   - Menjamin bahwa jika sebuah task terkirim lebih dari satu kali (*at-least-once delivery* akibat redelivery broker atau *network partition*), pengiriman email hanya dieksekusi tepat satu kali (*exactly-once side effect*).

5. **Multi-Role Application Entrypoint:**
   - Flag CLI `-role` (`api`, `worker`, `all`) untuk memisahkan siklus hidup HTTP web server dan background worker relay dalam satu *codebase*.
   - *Graceful shutdown* terkoordinasi yang memastikan seluruh worker menyelesaikan task aktif dan koneksi resource (PostgreSQL, MongoDB, Redis Inspector) ditutup dengan rapi.

6. **Endpoint Monitoring Antrean (Admin):**
   - Endpoint `GET /api/v1/admin/queue` terproteksi RBAC (`queue:monitor`) untuk memantau metrik outbox pending, usia event tertua, notifikasi terkirim, serta status real-time antrean Redis Asynq (*pending*, *active*, *retry*, *archived*).

## Arsitektur Alur Notifikasi

```
[Client] 
   │ HTTP Request (Enroll / Drop)
   ▼
[Fiber API] ──────── (Single DB Transaction) ────────► [PostgreSQL]
                                                         ├── courses / enrollments
                                                         └── outbox (published_at = NULL)
                                                                    ▲
[Outbox Relay] ── Polling (FOR UPDATE SKIP LOCKED) ─────────────────┘
   │
   │ asynq.Enqueue (Task: notification:send)
   ▼
[Redis Broker] (Queue: notifications)
   │
   │ Pull Task
   ▼
[Notification Worker] 
   │
   ├── 1. Claim outbox_id (Idempotency Check)
   ├── 2. Render Template & Send Email (LogMailer)
   └── 3. Simpan Riwayat ─────────────────────────────► [PostgreSQL]
                                                         └── notifications
```

---

## Matriks Endpoint Tambahan (Modul 11)

| Method | Endpoint | Izin / Proteksi | Hak Akses Role | Keterangan |
| :--- | :--- | :--- | :--- | :--- |
| GET | `/api/v1/notifications/mine` | Bearer Auth | Mahasiswa / Semua Role | Mengambil daftar riwayat notifikasi milik pengguna login |
| GET | `/api/v1/admin/queue` | Bearer Auth (`queue:monitor`) | Admin | Monitoring metrik antrean outbox dan status Asynq Redis |

---

## Hasil Verifikasi Skenario Pengujian (Bagian C)

| Skenario | Deskripsi Pengujian | Hasil Observasi | Status |
| :--- | :--- | :--- | :--- |
| **C.1: Jalur Normal** | Pendaftaran kelas saat worker mati, lalu worker dinyalakan | Outbox tersimpan dengan `published_at = NULL`. Saat worker hidup, relay mempublish event, email tercatat di `logs/mail.log`, dan notifikasi muncul di DB. | **VALID** |
| **C.2: Crash Redis** | Pendaftaran kelas saat service Redis dihentikan (`brew services stop redis`) | Pendaftaran kelas tetap berhasil (`201 Created`). Outbox tertahan di DB. Begitu Redis dinyalakan kembali, relay otomatis memproses sisa antrean hingga selesai tanpa kehilangan data. | **VALID** |
| **C.3: Temporary Failure** | Simulasi mailer sibuk (`MAILER_FAIL_EVERY=2`) | Worker mendeteksi kegagalan sementara, Asynq menahan task dengan delay retry bertahap (`(n+1)*2` detik), dan task berhasil terkirim pada percobaan berikutnya. | **VALID** |
| **C.4: Permanent Failure** | Pendaftaran kelas akun email tidak valid (`salah@kampus.invalid`) | Error dibungkus dengan `asynq.SkipRetry`. Worker tidak melakukan retry berulang kali, task langsung dialihkan ke status `archived` (DLQ), dan data `notifications` tetap bersih. | **VALID** |
| **C.5: Idempotency** | Pengiriman 5 salinan task duplikat outbox yang sama via `tools/redeliver` | Worker hanya mengirim email 1 kali. Salinan task lainnya ditandai `already_sent=true` dan diabaikan (*no-op*). Jumlah data notifikasi di PostgreSQL tetap tepat 1 baris. | **VALID** |

---

## Panduan Menjalankan Layanan

1. **Jalankan Service Prasyarat:**
   ```zsh
   brew services start postgresql@16
   brew services start mongodb-community
   brew services start redis
   ```

2. **Jalankan Aplikasi Berdasarkan Role:**
   - Menjalankan seluruh sistem (API & Worker bersamaan):
     ```zsh
     go run . -role=all
     ```
   - Menjalankan server HTTP API saja:
     ```zsh
     go run . -role=api
     ```
   - Menjalankan Background Worker & Outbox Relay saja:
     ```zsh
     go run . -role=worker
     ```

3. **Menjalankan Stress Test Idempotensi:**
   ```zsh
   go run tools/redeliver/main.go -outbox <ID_OUTBOX> -copies 5
   ```

