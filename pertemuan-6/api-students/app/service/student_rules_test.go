package service

import (
	"testing"

	"api-students/app/model"
)

func TestValidateCreate(t *testing.T) {
	// Kasus gagal: NIM kosong
	reqInvalid := model.CreateStudentRequest{
		NIM:   "",
		Name:  "Budi Santoso",
		Grade: "A",
	}
	errs := ValidateCreate(reqInvalid)
	if _, exists := errs["nim"]; !exists {
		t.Errorf("diharapkan error pada NIM, tapi lolos validasi")
	}

	// Kasus sukses
	reqValid := model.CreateStudentRequest{
		NIM:   "187221001",
		Name:  "Budi Santoso",
		Grade: "A",
	}
	errsValid := ValidateCreate(reqValid)
	if len(errsValid) != 0 {
		t.Errorf("seharusnya tidak ada error, dapat: %v", errsValid)
	}
}

func TestApplyPatch(t *testing.T) {
	initial := model.Student{
		ID:       1,
		NIM:      "187221001",
		Name:     "Budi",
		Grade:    "B",
		IsActive: true,
	}
	newName := "Budi Wicaksono"
	newGrade := "A"
	patchReq := model.PatchStudentRequest{
		Name:  &newName,
		Grade: &newGrade,
	}

	updated, errs := ApplyPatch(initial, patchReq)
	if len(errs) != 0 {
		t.Fatalf("tidak diharapkan error: %v", errs)
	}
	if updated.Name != "Budi Wicaksono" {
		t.Errorf("nama gagal diperbarui, dapat: %s", updated.Name)
	}
	if updated.Grade != "A" {
		t.Errorf("grade gagal diperbarui, dapat: %s", updated.Grade)
	}
	if updated.NIM != "187221001" {
		t.Errorf("NIM tidak boleh berubah jika tidak dikirim")
	}
}

func TestCountTotalPages(t *testing.T) {
	cases := []struct {
		total, limit, want int
	}{
		{0, 10, 0},
		{1, 10, 1},
		{10, 10, 1},
		{11, 10, 2},
		{25, 10, 3},
	}

	for _, tc := range cases {
		if got := CountTotalPages(tc.total, tc.limit); got != tc.want {
			t.Errorf("total=%d limit=%d: harap %d, dapat %d", tc.total, tc.limit, tc.want, got)
		}
	}
}
