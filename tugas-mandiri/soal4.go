package main

import "fmt"

// Definisi struct Student
type Student struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Grade    float64 `json:"grade"`
	IsActive bool    `json:"is_active"`
}

// 1. GetInfo menggunakan Value Receiver
// Alasan: Method ini hanya membaca data untuk ditampilkan tanpa perlu mengubah nilai asli struct.
func (s Student) GetInfo() string {
	return fmt.Sprintf("ID: %d | Nama: %s | Nilai: %.2f | Aktif: %t", s.ID, s.Name, s.Grade, s.IsActive)
}

// 2. UpdateGrade menggunakan Pointer Receiver
// Alasan: Method ini perlu mengubah field Grade langsung pada alamat memori objek aslinya.
func (s *Student) UpdateGrade(grade float64) {
	s.Grade = grade
}

// 3. Activate menggunakan Pointer Receiver
// Alasan: Mengubah field IsActive menjadi true secara permanen pada objek aslinya.
func (s *Student) Activate() {
	s.IsActive = true
}

// 4. Deactivate menggunakan Pointer Receiver
// Alasan: Mengubah field IsActive menjadi false secara permanen pada objek aslinya.
func (s *Student) Deactivate() {
	s.IsActive = false
}

func main() {
	fmt.Println("=== DEMONSTRASI STRUCT & METHOD STUDENT ===")

	// Inisialisasi objek mahasiswa
	mhs := Student{
		ID:       101,
		Name:     "Achmad Derajat",
		Grade:    82.5,
		IsActive: false,
	}

	// 1. Tampilkan info awal
	fmt.Println("Info Awal:")
	fmt.Println(mhs.GetInfo())
	fmt.Println()

	// 2. Aktivasi mahasiswa
	fmt.Println("Mengaktifkan status mahasiswa...")
	mhs.Activate()
	fmt.Println(mhs.GetInfo())
	fmt.Println()

	// 3. Update nilai mahasiswa
	fmt.Println("Memperbarui nilai mahasiswa menjadi 93.75...")
	mhs.UpdateGrade(93.75)
	fmt.Println(mhs.GetInfo())
	fmt.Println()

	// 4. Nonaktifkan mahasiswa
	fmt.Println("Menonaktifkan status mahasiswa...")
	mhs.Deactivate()
	fmt.Println(mhs.GetInfo())
}