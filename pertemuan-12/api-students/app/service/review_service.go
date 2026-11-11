package service

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"
)

type ReviewService struct {
	reviewRepo repository.ReviewRepository
	courseRepo repository.CourseRepository
	perms      *helper.PermissionSet
}

func NewReviewService(reviewRepo repository.ReviewRepository, courseRepo repository.CourseRepository, perms *helper.PermissionSet) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		courseRepo: courseRepo,
		perms:      perms,
	}
}

func (s *ReviewService) Create(c *fiber.Ctx) error {
	courseID, err := strconv.Atoi(c.Params("id"))
	if err != nil || courseID <= 0 {
		return helper.BadRequest("ID kelas tidak valid")
	}

	authUser, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("pengguna belum terautentikasi")
	}

	var req model.CreateReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("Body bukan JSON yang valid")
	}

	if req.Rating < 1 || req.Rating > 5 {
		return helper.NewError(http.StatusUnprocessableEntity, "VALIDATION_ERROR", "rating harus 1 sampai 5")
	}

	_, err = s.courseRepo.FindByID(c.Context(), courseID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.NotFound("kelas tidak ditemukan")
		}
		return helper.Internal(err)
	}

	enrolled, err := s.courseRepo.IsEnrolled(c.Context(), courseID, authUser.UserID)
	if err != nil || !enrolled {
		return helper.Forbidden("bukan peserta kelas")
	}

	tags := NormalizeTags(req.Tags)
	created, err := s.reviewRepo.Create(c.Context(), model.Review{
		CourseID: courseID,
		UserID:   authUser.UserID,
		Rating:   req.Rating,
		Comment:  req.Comment,
		Tags:     tags,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return helper.Conflict("hanya boleh memberikan satu ulasan")
		}
		return helper.Internal(err)
	}

	return helper.Success(c, http.StatusCreated, "ulasan berhasil dibuat", created)
}

func (s *ReviewService) ListByCourse(c *fiber.Ctx) error {
	courseID, err := strconv.Atoi(c.Params("id"))
	if err != nil || courseID <= 0 {
		return helper.BadRequest("ID kelas tidak valid")
	}

	reviews, err := s.reviewRepo.FindByCourse(c.Context(), courseID)
	if err != nil {
		return helper.Internal(err)
	}

	return helper.Success(c, http.StatusOK, "daftar ulasan berhasil diambil", reviews)
}