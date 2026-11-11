# Praktikum Pemrograman Backend - Pertemuan 12: Document Database Integration & Course Review System

Proyek RESTful API Manajemen Mahasiswa (*Students API*) berbasis Go, Fiber Framework, PostgreSQL (`pgxpool`), MongoDB (`mongo-driver/v2`), dan Redis (`asynq`).

Modul 12 mengintegrasikan MongoDB sebagai *document database* untuk menangani data ulasan dan penilaian mata kuliah (*course reviews*), memanfaatkan skema fleksibel dokumen NoSQL, penegakan integritas data melalui *compound unique index*, serta pemetaan entitas internal untuk decoding BSON `ObjectID`.

---

## Fitur Utama & Pembaruan Modul 12

1. **MongoDB Review Repository:**
   - Integrasi MongoDB menggunakan official driver `go.mongodb.org/mongo-driver/v2`.
   - Pemisahan model entitas repository (`reviewEntity`, `replyEntity`) untuk menangani decoding BSON native `bson.ObjectID` ke representasi string hex pada JSON API (`model.Review`).
   - Implementasi kueri terurut dengan sorting descending berdasarkan waktu pembuatan (`created_at: -1`).

2. **Course Review & Rating Validation:**
   - Aturan bisnis: ulasan hanya dapat dikirimkan oleh mahasiswa yang telah terdaftar aktif (*enrolled*) pada mata kuliah terkait di PostgreSQL.
   - Validasi input rating (1-5), batas panjang komentar, dan sanitasi array tag ulasan (`helper.NormalizeTags`).

3. **Compound Unique Index:**
   - Penegakan integritas data ulasan di MongoDB melalui index unik majemuk:
     ```json
     { "course_id": 1, "user_id": 1 } (Unique: true)
     ```
   - Mencegah mahasiswa memberikan lebih dari satu ulasan pada kelas yang sama dan mengembalikan kode respon `CONFLICT` (`409`).

4. **Koneksi Dual-Database Terkoordinasi:**
   - PostgreSQL menangani entitas relasional transaksional (akun mahasiswa, mata kuliah, pendaftaran, dan outbox).
   - MongoDB menangani data berbasis dokumen ulasan bersarang (*reviews and replies*).

---

## Verifikasi Endpoint Ulasan

### 1. Registrasi & Login Mahasiswa
```bash
RESP=$(curl -s -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "tester01", "password": "KucingTerbang#2026!SangatKuat"}')
TOKEN=$(echo "$RESP" | jq -r '.data.access_token')
```

### 2. Pendaftaran Kelas (PostgreSQL)
```bash
curl -s -X POST http://localhost:3000/api/v1/courses/1/enroll \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" | jq .
```

### 3. Pembuatan Ulasan Kelas (MongoDB)
```bash
curl -s -X POST http://localhost:3000/api/v1/courses/1/reviews \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "rating": 5,
    "comment": "Materi Go Fiber dan MongoDB sangat aplikatif!",
    "tags": ["backend", "recommended"]
  }' | jq .
```

### 4. Ambil Daftar Ulasan Kelas
```bash
curl -s -X GET http://localhost:3000/api/v1/courses/1/reviews \
  -H "Authorization: Bearer $TOKEN" | jq .
```