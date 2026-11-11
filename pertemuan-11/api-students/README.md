# Praktikum Pemrograman Backend - Pertemuan 10: Multi-Database Integration & Student Portfolio

Proyek RESTful API Manajemen Mahasiswa (*Students API*) berbasis Go, Fiber Framework, PostgreSQL (`pgxpool`), dan MongoDB. Modul ini mengimplementasikan arsitektur *polyglot persistence* (multi-database) dengan memadukan basis data relasional PostgreSQL untuk entitas inti mahasiswa dan otentikasi pengguna, serta basis data dokumen MongoDB untuk mengelola data portofolio yang dinamis (*nested sub-documents*), pipeline agregasi data (*aggregation pipeline*), dan penegakan integritas data lintas basis data (*cross-database integrity*).

---

## Fitur Utama & Pembaruan Modul 10

1. **Integrasi Basis Data Dokumen (MongoDB Go Driver):**
   - Inisialisasi koneksi `NewMongoClient` dengan timeout kontekstual, verifikasi ping, pembacaan konfigurasi database via environment variable (`MONGO_DATABASE`), serta penanganan *graceful disconnect* saat aplikasi berhenti.
   - Penjaminan indeks otomatis (`EnsureIndexes`) pada koleksi `portfolios` dengan *compound unique index* pada `student_id` guna memastikan hubungan 1-to-1 antar data mahasiswa dan portofolionya.

2. **Integritas Data Lintas Basis Data (*Cross-Database Integrity*):**
   - Validasi keberadaan entitas mahasiswa di PostgreSQL (`studentRepo.FindByID`) sebelum menjalankan operasi mutasi atau kueri pada koleksi MongoDB.
   - Penanganan terpusat jika ID mahasiswa tidak ditemukan (`404 NOT_FOUND`), mencegah terbentuknya *orphan documents* di dalam koleksi MongoDB.

3. **Manajemen Portofolio & Sub-Dokumen Proyek:**
   - **Upsert Portofolio (`PUT /students/:id/portfolio`):** Operasi idempoten menggunakan opsi `SetUpsert(true)` dengan pembaruan otomatis kolom timestamp `updated_at`.
   - **Penambahan Proyek Tersemat (`POST /students/:id/portfolio/projects`):** Menambahkan data proyek ke dalam array `projects` menggunakan operator atomik `$push`, pembuatan `ObjectID` unik per entitas proyek, dan penegakan batas kuota maksimal 20 proyek per mahasiswa (`409 CONFLICT: PROJECT_LIMIT`).
   - **Normalisasi Tag Teknologi:** Penyeragaman otomatis nama teknologi ke format huruf kecil (*lowercase*) sebelum disimpan ke basis data untuk mempermudah pencarian dan agregasi yang konsisten.
   - **Penghapusan Proyek Tertentu (`DELETE /students/:id/portfolio/projects/:projectId`):** Menghapus sub-dokumen dari array tersemat secara atomik menggunakan operator `$pull`.

4. **Pipeline Agregasi Statistik Teknologi (`GET /portfolios/stats/technologies`):**
   - Pemanfaatan tahapan agregasi MongoDB teroptimasi:
     - `$unwind`: Memecah array sub-dokumen `projects`.
     - `$unwind`: Memecah array string `technologies`.
     - `$group`: Mengelompokkan berdasarkan nama teknologi dan menghitung total frekuensi kemunculan (`$sum: 1`).
     - `$sort`: Mengurutkan popularitas teknologi dari frekuensi terbanyak (descending).
     - `$limit`: Membatasi jumlah keluaran teknologi teratas.

---

## Skema Dokumen MongoDB (`portfolios`)

```json
{
  "_id": "ObjectID",
  "student_id": 3,
  "headline": "Backend Engineer & Cloud Enthusiast",
  "bio": "Mahasiswa Teknik Informatika yang fokus pada high-concurrency backend dan distributed systems.",
  "projects": [
    {
      "_id": "ObjectID",
      "title": "High Performance Student Management API",
      "technologies": ["go", "fiber", "postgresql", "mongodb", "docker"],
      "link": "[https://github.com/mrcleverloki/api-students](https://github.com/mrcleverloki/api-students)",
      "year": 2026
    }
  ],
  "certificates": [],
  "created_at": "ISODate",
  "updated_at": "ISODate"
}

## Matriks Endpoint Tambahan (Portofolio Mahasiswa)

| Method | Endpoint | Izin / Proteksi | Hak Akses Role | Keterangan |
| :--- | :--- | :--- | :--- | :--- |
| PUT | /api/v1/students/:id/portfolio | Bearer Auth | Admin, Staff, Owner | Membuat atau memperbarui data dasar portofolio |
| GET | /api/v1/students/:id/portfolio | Bearer Auth | Semua Role | Mengambil data portofolio lengkap mahasiswa |
| POST | /api/v1/students/:id/portfolio/projects | Bearer Auth | Admin, Staff, Owner | Menambahkan proyek baru (Maks. 20 proyek) |
| DELETE | /api/v1/students/:id/portfolio/projects/:projectId | Bearer Auth | Admin, Staff, Owner | Menghapus proyek spesifik dari array portofolio |
| GET | /api/v1/portfolios/stats/technologies | Bearer Auth | Semua Role | Agregasi statistik teknologi proyek terpopuler |

---

## Panduan Menjalankan Layanan

1. **Jalankan Service Basis Data:**
   - PostgreSQL: `brew services start postgresql@16`
   - MongoDB: `brew services start mongodb-community`

2. **Jalankan Server API:**
   - `go run main.go`

3. **Verifikasi Status Koneksi Basis Data:**
   - `curl -i http://localhost:3000/api/v1/health`
