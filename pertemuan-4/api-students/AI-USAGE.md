# Pernyataan Penggunaan AI (AI-USAGE)

Dokumen ini menjelaskan pemanfaatan kecerdasan buatan (Artificial Intelligence) dalam penyelesaian Tugas Praktikum Pertemuan 4 (Clean Architecture).

## Alat AI yang Digunakan
- **Model / Tools**: Google Gemini

## Rincian Bagian yang Dibantu oleh AI
1. **Pemahaman Konsep Clean Architecture & Dependency Rule**
   - Diskusi konsep empat layer Clean Architecture, aturan ketergantungan kode (*Dependency Rule*), serta mitigasi penyederhanaan arsitektur pada ekosistem Go dan Fiber.
2. **Restrukturisasi Kode ke Struktur Baku**
   - Panduan pemecahan *presenter* dan pembaca kueri ke dalam package `helper/`.
   - Perancangan pemisahan *business rules* murni (`student_rules.go`) agar terisolasi dari framework Fiber dan basis data.
   - Penyusunan pengujian otomatis (*unit test*) pada layer use case tanpa inisialisasi server.
3. **Konfigurasi Logging & Middleware**
   - Penerapan middleware global, validasi *Content-Type* (`RequireJSON`), dan *structured logger* (`slog`) dengan rotasi file log otomatis menggunakan `lumberjack`.
   - Pembersihan titik masuk `main.go` sehingga murni menangani perakitan dependensi dan *graceful shutdown*.

## Bagian yang Dikerjakan Mandiri
- Pengelolaan struktur folder proyek dan konfigurasi variabel lingkungan (`.env`).
- Eksekusi kompilasi kode, perbaikan *type field* model, dan eksekusi pengujian otomatis `go test`.
- Pengujian fungsional seluruh endpoint REST API menggunakan ekstensi Thunder Client di VS Code.
- Audit pemeriksaan kebocoran layer (*layer leakage inspection*).
- Pengelolaan riwayat *commit* Git bertahap dan penyusunan laporan akhir PDF.
