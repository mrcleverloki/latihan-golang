package main

import "fmt"

func main() {
	// ==========================================
	// 1. DEKLARASI 5 VARIABEL TIPE BERBEDA
	// ==========================================
	var nama string = "Achmad Derajat"
	var umur int = 21
	var ipk float64 = 3.85
	var isAktif bool = true
	hobi := []string{"Coding", "Membaca", "Desain"} // tipe slice

	fmt.Println("=== 1. DEKLARASI VARIABEL ===")
	fmt.Printf("Nama      : %s (Tipe: string)\n", nama)
	fmt.Printf("Umur      : %d tahun (Tipe: int)\n", umur)
	fmt.Printf("IPK       : %.2f (Tipe: float64)\n", ipk)
	fmt.Printf("Status    : %t (Tipe: bool)\n", isAktif)
	fmt.Printf("Hobi      : %v (Tipe: slice []string)\n\n", hobi)

	// ==========================================
	// 2. MAP DATA MAHASISWA (Nama -> Nilai)
	// ==========================================
	fmt.Println("=== 2. OPERASI MAP DATA MAHASISWA ===")
	
	// Inisialisasi map
	dataNilai := make(map[string]int)

	// a. Menambah data ke map
	dataNilai["Achmad"] = 95
	dataNilai["Budi"] = 80
	dataNilai["Siti"] = 88
	dataNilai["Doni"] = 75
	fmt.Println("Data awal setelah penambahan:")
	for k, v := range dataNilai {
		fmt.Printf("- %s: %d\n", k, v)
	}
	fmt.Println()

	// b. Membaca data dengan pengecekan keberadaan (two-value assignment)
	namaCari := "Achmad"
	if nilai, ada := dataNilai[namaCari]; ada {
		fmt.Printf("[CEK] Data %s ditemukan dengan nilai: %d\n", namaCari, nilai)
	} else {
		fmt.Printf("[CEK] Data %s tidak ditemukan dalam map!\n", namaCari)
	}

	namaTidakAda := "Riko"
	if nilai, ada := dataNilai[namaTidakAda]; ada {
		fmt.Printf("[CEK] Data %s ditemukan dengan nilai: %d\n", namaTidakAda, nilai)
	} else {
		fmt.Printf("[CEK] Data %s tidak ditemukan dalam map!\n", namaTidakAda)
	}
	fmt.Println()

	// c. Menghapus data dari map
	fmt.Println("Menghapus data 'Doni' dari map...")
	delete(dataNilai, "Doni")
	fmt.Println()

	// d. Menelusuri (looping) seluruh isi map yang tersisa
	fmt.Println("Data akhir mahasiswa dalam map:")
	for mhs, skor := range dataNilai {
		fmt.Printf("Mahasiswa: %s | Nilai: %d\n", mhs, skor)
	}
}