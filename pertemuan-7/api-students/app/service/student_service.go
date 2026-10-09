package service

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"
)

type StudentService struct {
	repo  repository.StudentRepository
	perms *helper.PermissionSet
}

func NewStudentService(repo repository.StudentRepository, perms *helper.PermissionSet) *StudentService {
	return &StudentService{
		repo:  repo,
		perms: perms,
	}
}

// List mendukung Keyset Pagination & Content Negotiation (JSON dan CSV)
func (s *StudentService) List(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	// Content Negotiation (RFC 9110)
	format, err := helper.Negotiate(c, helper.FormatJSON, helper.FormatCSV)
	if err != nil {
		return err
	}

	q, err := helper.ParseCursorQuery(c)
	if err != nil {
		return err
	}

	rows, err := s.repo.FindAfterCursor(ctx, q)
	if err != nil {
		return helper.Internal(err)
	}

	// Ekspor CSV bila diminta
	if format == helper.FormatCSV {
		return helper.WriteStudentsCSV(c, rows)
	}

	// Logika Limit + 1 pada Cursor
	hasMore := len(rows) > q.Limit
	if hasMore {
		rows = rows[:q.Limit]
	}

	meta := &model.CursorMeta{Limit: q.Limit, HasMore: hasMore}
	if hasMore && len(rows) > 0 {
		last := rows[len(rows)-1]
		meta.NextCursor = helper.EncodeCursor(last.CreatedAt, last.ID)
	}

	return helper.SuccessCursor(c, "daftar mahasiswa berhasil diambil", rows, meta)
}

func (s *StudentService) Get(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	currentUser, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id harus berupa angka positif")
	}

	student, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return translateStudentError(err)
	}

	ownerID := 0
	if student.OwnerID != nil {
		ownerID = *student.OwnerID
	}

	if !CanAccessStudent(currentUser, ownerID, s.perms, "student:read:any") {
		return helper.Forbidden("tidak memiliki hak akses ke data mahasiswa ini")
	}

	return helper.Success(c, fiber.StatusOK, "mahasiswa ditemukan", student)
}

func (s *StudentService) Create(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	currentUser, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}

	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}

	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)
	req.Grade = strings.TrimSpace(req.Grade)

	// Validasi Deklaratif
	if fieldErrors := helper.ValidateStruct(req); fieldErrors != nil {
		return helper.Validation(fieldErrors)
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	newStudent, err := s.repo.Create(ctx, model.Student{
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		IsActive: isActive,
		OwnerID:  &currentUser.UserID,
	})
	if err != nil {
		return translateStudentError(err)
	}

	return helper.Created(c, "mahasiswa berhasil ditambahkan", newStudent,
		"/api/v1/students/"+strconv.Itoa(newStudent.ID))
}

func (s *StudentService) Replace(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	currentUser, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id harus berupa angka positif")
	}

	currentStudent, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return translateStudentError(err)
	}

	ownerID := 0
	if currentStudent.OwnerID != nil {
		ownerID = *currentStudent.OwnerID
	}

	if !CanAccessStudent(currentUser, ownerID, s.perms, "student:update:any") {
		return helper.Forbidden("tidak memiliki hak akses mengubah data mahasiswa ini")
	}

	var req model.ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}

	// Validasi Deklaratif
	if fieldErrors := helper.ValidateStruct(req); fieldErrors != nil {
		return helper.Validation(fieldErrors)
	}

	result, err := s.repo.Update(ctx, model.Student{
		ID:       id,
		NIM:      strings.TrimSpace(req.NIM),
		Name:     strings.TrimSpace(req.Name),
		Grade:    strings.TrimSpace(req.Grade),
		IsActive: req.IsActive,
		OwnerID:  currentStudent.OwnerID,
	})
	if err != nil {
		return translateStudentError(err)
	}

	return helper.Success(c, fiber.StatusOK, "data mahasiswa berhasil diganti seluruhnya", result)
}

func (s *StudentService) Patch(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	currentUser, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id harus berupa angka positif")
	}

	current, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return translateStudentError(err)
	}

	ownerID := 0
	if current.OwnerID != nil {
		ownerID = *current.OwnerID
	}

	if !CanAccessStudent(currentUser, ownerID, s.perms, "student:update:any") {
		return helper.Forbidden("tidak memiliki hak akses mengubah data mahasiswa ini")
	}

	var req model.PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}

	if IsEmptyPatch(req) {
		return helper.BadRequest("tidak ada field yang diubah")
	}

	// Validasi Deklaratif (dengan omitnil pada pointer)
	if fieldErrors := helper.ValidateStruct(req); fieldErrors != nil {
		return helper.Validation(fieldErrors)
	}

	updated := ApplyPatch(current, req)

	result, err := s.repo.Update(ctx, updated)
	if err != nil {
		return translateStudentError(err)
	}

	return helper.Success(c, fiber.StatusOK, "data mahasiswa berhasil diperbarui sebagian", result)
}

func (s *StudentService) Delete(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id harus berupa angka positif")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return translateStudentError(err)
	}

	return helper.NoContent(c)
}

// PERBAIKAN BUG #2: translateStudentError mengembalikan helper.Internal(err) pada default case (fail closed)
func translateStudentError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return helper.NotFound("data mahasiswa tidak ditemukan")
	case errors.Is(err, repository.ErrDuplicate):
		return helper.Conflict("NIM sudah terdaftar")
	default:
		return helper.Internal(err)
	}
}
