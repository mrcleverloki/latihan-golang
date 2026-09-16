# Deklarasi Penggunaan AI (AI-USAGE)

Sesuai dengan integritas akademik dan pedoman tugas praktikum, dokumen ini mencatat rincian transparansi penggunaan alat bantu kecerdasan buatan (Artificial Intelligence) selama penyelesaian proyek.

---

## 1. Alat Bantu yang Digunakan
* **Nama Alat:** Gemini AI (Google)
* **Model:** Gemini
* **Media Akses:** Obrolan Interaktif Web

---

## 2. Rincian Pemanfaatan Alat Bantu

| Modul / Komponen | Bentuk Kontribusi AI | Tingkat Bantuan |
| :--- | :--- | :--- |
| **Helper & Security (`helper/`)** | Pembuatan boilerplate fungsi token JWT, hashing password (bcrypt), komputasi hash SHA-256 untuk refresh token, dan fungsi `VerifyDummyPassword`. | Sedang |
| **Middleware (`middleware/`)** | Pembuatan logika parsing header Bearer token pada `RequireAuth` dan konfigurasi `LoginRateLimiter` Fiber. | Sedang |
| **Routing & Wireup (`route/`, `main.go`)** | Penyusunan Dependency Injection dan pemetaan route auth/protected di Fiber. | Rendah |
| **Skenario Pengujian** | Pemberian panduan langkah uji coba manual (payload JSON dan headers) pada Thunder Client. | Rendah |
| **Penyusunan Dokumentasi** | Pembuatan draf narasi teknis untuk laporan praktikum (Bagian 1–5). | Sedang |

---

## 3. Pekerjaan Mandiri (Human Contribution)
* Penulisan dan eksekusi skema DDL basis data (`users`, `refresh_tokens`) pada PostgreSQL.
* Perakitan, review, dan penyesuaian logika bisnis pada layer `repository` dan `service`.
* Eksekusi unit testing lokal dan verifikasi kelolosan uji (`go test`).
* Pengujian manual seluruh skenario API melalui Thunder Client serta pengambilan tangkapan layar bukti respon HTTP.
* Pengelolaan branching, commit bertahap (*micro-commits*), dan konfigurasi `.gitignore`.