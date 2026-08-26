package main

import "fmt"

// Fungsi menukar nilai dua integer menggunakan pointer
func swap(a, b *int) {
	*a, *b = *b, *a
}

// Fungsi menambah item ke slice menggunakan pointer
func updateSlice(s *[]string, newItem string) {
	*s = append(*s, newItem)
}

// Fungsi pembuktian Pass by Value (tidak mengubah data asli)
func passByValue(x int) {
	x = 999
}

// Fungsi pembuktian Pass by Pointer (mengubah data asli)
func passByPointer(x *int) {
	*x = 999
}

func main() {
	// ==========================================
	// 1. FUNGSI SWAP (POINTER INT)
	// ==========================================
	fmt.Println("=== 1. UJI COBA SWAP POINTER ===")
	angka1, angka2 := 10, 50
	fmt.Printf("Sebelum swap: angka1 = %d, angka2 = %d\n", angka1, angka2)

	swap(&angka1, &angka2)
	fmt.Printf("Setelah swap : angka1 = %d, angka2 = %d\n\n", angka1, angka2)

	// ==========================================
	// 2. FUNGSI UPDATESLICE (POINTER SLICE)
	// ==========================================
	fmt.Println("=== 2. UJI COBA UPDATE SLICE POINTER ===")
	daftarTools := []string{"VS Code", "Golang"}
	fmt.Println("Slice awal   :", daftarTools)

	updateSlice(&daftarTools, "Fiber v2")
	updateSlice(&daftarTools, "Git")
	fmt.Println("Slice update :", daftarTools)
	fmt.Println()

	// ==========================================
	// 3. PERBANDINGAN PASS BY VALUE VS POINTER
	// ==========================================
	fmt.Println("=== 3. PERBANDINGAN PASS BY VALUE VS PASS BY POINTER ===")
	nilaiAsli := 100
	fmt.Println("Nilai mula-mula        :", nilaiAsli)

	// Uji Pass by Value
	passByValue(nilaiAsli)
	fmt.Println("Setelah passByValue    :", nilaiAsli, "(Tetap, karena menerima salinan data)")

	// Uji Pass by Pointer
	passByPointer(&nilaiAsli)
	fmt.Println("Setelah passByPointer  :", nilaiAsli, "(Berubah, karena mengakses alamat memori)")
}