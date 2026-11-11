# Deklarasi Penggunaan AI (AI-USAGE)

Sesuai dengan prinsip integritas akademik dan transparansi pedoman perkuliahan, dokumen ini mencatat pemanfaatan alat bantu kecerdasan buatan (*Artificial Intelligence*) selama pengerjaan proyek Modul 10 (Multi-Database Integration with MongoDB, Embedded Document Lifecycle, Aggregation Pipeline, and Cross-Database Validation).

---

## 1. Alat Bantu yang Digunakan
- **Nama Alat:** Gemini AI (Google)
- **Model:** Gemini
- **Media Akses:** Obrolan Interaktif

---

## 2. Rincian Pemanfaatan Alat Bantu

| Modul / Komponen | Bentuk Kontribusi AI | Tingkat Bantuan |
| :--- | :--- | :--- |
| **Integrasi Klien MongoDB (`database/mongo.go` & `main.go`)** | Bantuan penyusunan perakitan koneksi driver MongoDB resmi, penanganan *timeout context*, *graceful shutdown disconnect*, serta pendaftaran *compound index* pada startup aplikasi. | Sedang |
| **Pipeline Agregasi Statistik (`portfolio_repository.go`)** | Bantuan penyusunan tahapan agregasi multi-stage MongoDB (`$unwind` ganda untuk nested sub-documents, `$group` untuk penghitungan total, `$sort` descending, dan `$limit`). | Sedang |
| **Validasi Lintas Basis Data (*Cross-Database Integrity*)** | Arahan penerapan pengecekan keberadaan ID mahasiswa di PostgreSQL sebelum mutasi dokumen MongoDB pada layer service (`PortfolioService`). | Rendah |
| **Penyusunan Perintah Pengujian Terminal (*cURL Test Suite*)** | Penyusunan skrip cURL untuk pengujian berurutan: registrasi akun, login JWT, validasi error 404 (ID fiktif), *upsert* portofolio, penambahan proyek dengan normalisasi huruf kecil, agregasi statistik, dan penghapusan proyek atomik. | Rendah |
| **Troubleshooting Driver & Service MongoDB macOS** | Panduan mengatasi error `connection refused` dan konfigurasi *tap* formula `mongodb/brew` pada package manager Homebrew macOS. | Rendah |

---

## 3. Pekerjaan Mandiri (Human Contribution)
- Penyiapan dan konfigurasi service MongoDB lokal (`mongodb-community`) menggunakan Homebrew di macOS.
- Perancangan model DTO dan entitas BSON/JSON portofolio (`app/model/portfolio.go`) serta pemetaan tag validasi Fiber.
- Implementasi logika bisnis di `PortfolioService` mencakup aturan normalisasi tag teknologi dan pembatasan maksimal 20 proyek per portofolio.
- Konfigurasi registrasi routing di `route/route.go` dan perakitan dependency injection pada `main.go`.
- Pengujian langsung seluruh alur integrasi end-to-end melalui terminal menggunakan cURL dan verifikasi struktur BSON dokumen MongoDB melalui terminal.
- Pengelolaan version control git (pembuatan micro-commit terstruktur, pembaharuan berkas dokumentasi `README.md` dan `AI-USAGE.md`, serta sinkronisasi commit ke remote repository).