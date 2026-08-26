# Pernyataan Penggunaan AI (AI-USAGE)

Dokumen ini menjelaskan pemanfaatan kecerdasan buatan (Artificial Intelligence) dalam penyelesaian Tugas Praktikum Pertemuan 2 (REST API & HTTP Deep Dive).

## Alat AI yang Digunakan
- **Model / Tools**: Google Gemini

## Rincian Bagian yang Dibantu oleh AI
1. **Pemahaman Konsep Teori & HTTP Semantics**
   - Diskusi konsep safe & idempotent methods, perbedaan PUT vs PATCH (penggunaan pointer di Go), serta penentuan status code HTTP (400, 409, 415, 422).
2. **Scaffolding Struktur Kode**
   - Panduan pembagian berkas arsitektur (`model.go`, `helper.go`, `handler.go`, `main.go`).
   - Pembuatan format seragam struct envelope JSON (`WebResponse` dan `Meta`).
3. **Penyusunan Kontrak API & Validasi Query**
   - Penulisan tabel spesifikasi kontrak API di `README.md`.
   - Logika sanitasi query string (whitelist sort parameter dan pembatasan limit pagination).

## Bagian yang Dikerjakan Mandiri
- Inisialisasi modul Go dan dependensi Fiber lokal.
- Penulisan, penyesuaian logika slice in-memory, dan kompilasi proyek.
- Pengujian langsung seluruh skenario HTTP via cURL terminal dan pengambilan bukti tangkapan layar (screenshot).
- Pengelolaan riwayat commit Git dan penyusunan laporan PDF.