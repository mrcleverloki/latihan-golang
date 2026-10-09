package service

import (
	"testing"

	"api-students/app/model"
	"api-students/helper"
)

func TestValidateCreateStudent(t *testing.T) {
	// Kasus gagal: NIM tidak valid / kosong
	reqInvalid := model.CreateStudentRequest{
		NIM:   "",
		Name:  "Budi Santoso",
		Grade: "A",
	}
	errs := helper.ValidateStruct(reqInvalid)
	if _, exists := errs["nim"]; !exists {
		t.Errorf("diharapkan error pada NIM, tapi lolos validasi")
	}

	// Kasus gagal: NIM bukan 9 digit angka (custom validation nim_format)
	reqNIMInvalid := model.CreateStudentRequest{
		NIM:   "18722ABC",
		Name:  "Budi Santoso",
		Grade: "A",
	}
	errsNIM := helper.ValidateStruct(reqNIMInvalid)
	if _, exists := errsNIM["nim"]; !exists {
		t.Errorf("diharapkan error custom nim_format pada NIM bukan digit")
	}

	// Kasus sukses
	reqValid := model.CreateStudentRequest{
		NIM:   "187221001",
		Name:  "Budi Santoso",
		Grade: "A",
	}
	errsValid := helper.ValidateStruct(reqValid)
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

	updated := ApplyPatch(initial, patchReq)
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
