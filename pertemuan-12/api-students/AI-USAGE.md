# Deklarasi Penggunaan AI (AI-USAGE)

Sesuai dengan prinsip integritas akademik dan transparansi pedoman perkuliahan, dokumen ini mencatat pemanfaatan alat bantu kecerdasan buatan (*Artificial Intelligence*) selama pengerjaan proyek Modul 12 (MongoDB Document Store, Review & Rating System, Compound Unique Index, and BSON Entity Decoding).

---

## 1. Alat Bantu yang Digunakan
- **Nama Alat:** Gemini AI (Google)
- **Model:** Gemini
- **Media Akses:** Obrolan Interaktif

---

## 2. Rincian Pemanfaatan Alat Bantu

| Modul / Komponen | Bentuk Kontribusi AI | Tingkat Bantuan |
| :--- | :--- | :--- |
| **Penyusunan Review Repository (`review_repository.go`)** | Panduan pembuatan implementasi repository MongoDB, perancangan indeks majemuk unik (`course_id`, `user_id`), dan penanganan `mongo.IsDuplicateKeyError`. | Sedang |
| **Debugging Decoding BSON ObjectID (`INTERNAL_ERROR`)** | Identifikasi akar masalah kegagalan unmarshal BSON ObjectID ke struct Go string pada driver v2, serta perancangan struct entity internal (`reviewEntity`) untuk konversi heksadesimal yang aman. | Sedang |
| **Smoke Test & Validasi Bisnis** | Panduan langkah pengujian terminal cURL untuk alur autentikasi JWT, verifikasi keikutsertaan kelas (enroll), pembuatan review, serta pembuktian respon error `CONFLICT` saat duplikasi ulasan. | Rendah |

---

## 3. Pekerjaan Mandiri (Human Contribution)
- Penyiapan service lokal MongoDB dan PostgreSQL pada lingkungan macOS.
- Penyesuaian konfigurasi koneksi MongoDB dan pendaftaran rute review ke dalam Fiber server.
- Eksekusi langsung pengujian fungsional dan verifikasi log server di terminal.
- Pengujian unit test lokal (`go test -v -count=1 ./app/service/...`) untuk memastikan seluruh skenario uji lolos (100% PASS).
- Pengelolaan version control Git (review diff, penyusunan micro-commit, dan pemeliharaan dokumentasi).