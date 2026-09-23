# Deklarasi Penggunaan AI (AI-USAGE)

Sesuai dengan integritas akademik dan pedoman tugas praktikum, dokumen ini mencatat rincian transparansi penggunaan alat bantu kecerdasan buatan (Artificial Intelligence) selama penyelesaian proyek Modul 6 (Role-Based Access Control & Ownership Authorization).

---

## 1. Alat Bantu yang Digunakan
* **Nama Alat:** Gemini AI (Google)
* **Model:** Gemini
* **Media Akses:** Obrolan Interaktif Web

---

## 2. Rincian Pemanfaatan Alat Bantu

| Modul / Komponen | Bentuk Kontribusi AI | Tingkat Bantuan |
| :--- | :--- | :--- |
| **Helper RBAC (`helper/rbac.go`)** | Bantuan penyusunan struktur `PermissionSet` dan pemetaan matriks izin berbasis role (`user`, `staff`, `admin`). | Sedang |
| **Middleware RBAC (`middleware/rbac.go`)** | Bantuan formulasi middleware otorisasi dinamis (`RequirePermission`) untuk validasi token claims dan role checking. | Sedang |
| **Service Layer (`service/`)** | Diskusi perancangan logika otorisasi kepemilikan data (`CanAccessStudent`) untuk validasi kepemilikan (`owner_id`) vs role hierarki. | Sedang |
| **Routing & Dependency Injection (`route/`, `main.go`)** | Panduan penyusunan route handler mahasiswa dengan proteksi middleware RBAC terpadu. | Rendah |
| **Skenario Pengujian** | Penyusunan urutan skenario pengujian bertahap di Thunder Client (Role User ditolak, Staff read-only/forbidden delete, Admin full access, dan verifikasi Ownership). | Sedang |

---

## 3. Pekerjaan Mandiri (Human Contribution)
* Migrasi dan eksekusi skema basis data PostgreSQL (penambahan kolom `owner_id` dan foreign key constraint pada tabel `students`).
* Implementasi query data layer di repository (`FindByID`, `Update`, `Delete`, `List`) dengan integrasi `owner_id`.
* Konfigurasi database lokal, manajemen role user via CLI `psql`, dan pengujian data real.
* Eksekusi langsung seluruh skenario HTTP testing via Thunder Client dan dokumentasi tangkapan layar respons status code (`200 OK`, `201 Created`, `204 No Content`, `401 Unauthorized`, `403 Forbidden`).
* Debugging kendala teknis (penanganan token expired, header `Content-Type: application/json`, penyesuaian payload JSON).
* Pengelolaan version control git (commit bertahap dan sinkronisasi branch repository).
