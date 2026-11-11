package service

import (
	"api-students/app/model"
	"api-students/helper"
)

// CanAccessStudent memutuskan apakah seseorang boleh mengakses/mengubah data student.
// Fungsi murni: tidak mengimpor fiber maupun repository.
// Dua jalur yang diizinkan:
//  1. Kepemilikan (ownership) — data itu didaftarkan oleh dirinya sendiri.
//  2. Permission — role-nya memang memiliki permission scope :any yang sesuai.
func CanAccessStudent(
	current model.AuthUser,
	ownerID int,
	perms *helper.PermissionSet,
	anyPermission string,
) bool {
	if current.UserID == ownerID {
		return true
	}
	return perms.Can(current.Role, anyPermission)
}
