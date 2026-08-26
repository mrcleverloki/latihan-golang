# Praktikum Pemrograman Backend Lanjut - Pertemuan 1
Repositori ini berisi kode latihan dan tugas mandiri untuk praktikum lingkungan kerja Go dan framework Fiber v2.

---

## 📁 Struktur Repositori

```text
├── latihan-fiber/
│   ├── go.mod
│   ├── go.sum
│   └── main.go         # Endpoint HTTP dasar menggunakan Fiber v2
└── tugas-mandiri/
    ├── go.mod
    ├── soal2.go        # Variabel dan operasi Map
    ├── soal3.go        # Pointer dan perbandingan pass by value vs pointer
    └── soal4.go        # Struct Student dan implementasi method receiver
🚀 Cara Menjalankan Program
1. Menjalankan Server Fiber
Bash
cd latihan-fiber
go run main.go
Buka browser di http://localhost:3000.

2. Menjalankan Tugas Mandiri
Masuk ke folder tugas-mandiri:

Bash
cd tugas-mandiri
Jalankan file spesifik sesuai tugas:

Soal 2 (Variabel & Map):

Bash
go run soal2.go
Soal 3 (Pointer & Pass by Reference):

Bash
go run soal3.go
Soal 4 (Struct & Receiver):

Bash
go run soal4.go
📝 Catatan Implementasi Tugas
Variabel & Struktur Data (soal2.go): Implementasi tipe dasar (string, int, float64, bool, slice) serta operasi map (insert, lookup dengan pengecekan exists, delete, dan iterasi).

Pointer (soal3.go): Demonstrasi mutasi nilai menggunakan pointer pada integer dan slice, serta perbandingan mekanisme pass by value vs pass by pointer.

Struct & Method (soal4.go): Penggunaan Value Receiver pada method baca (GetInfo) dan Pointer Receiver pada method mutasi data (UpdateGrade, Activate, Deactivate).

🤖 Pernyataan Bantuan AI
Pengembangan dan penyusunan modul tugas ini dibantu oleh Gemini AI untuk panduan setup environment, pembuatan scaffolding kode, dan troubleshooting konfigurasi Git.