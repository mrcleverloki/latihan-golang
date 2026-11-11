package service

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"
)

type PortfolioService struct {
	portfolioRepo repository.PortfolioRepository
	studentRepo   repository.StudentRepository
}

func NewPortfolioService(
	portfolioRepo repository.PortfolioRepository,
	studentRepo repository.StudentRepository,
) *PortfolioService {
	return &PortfolioService{
		portfolioRepo: portfolioRepo,
		studentRepo:   studentRepo,
	}
}

// NormalizeTechnologies merapikan teks teknologi (lowercase, trim, hilangkan duplikat)
func NormalizeTechnologies(techs []string) []string {
	result := make([]string, 0, len(techs))
	seen := make(map[string]struct{}, len(techs))
	for _, t := range techs {
		cleaned := strings.ToLower(strings.TrimSpace(t))
		if cleaned == "" {
			continue
		}
		if _, dup := seen[cleaned]; dup {
			continue
		}
		seen[cleaned] = struct{}{}
		result = append(result, cleaned)
	}
	return result
}

// GET /students/:id/portfolio
func (s *PortfolioService) Get(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	studentID, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id mahasiswa harus berupa angka positif")
	}

	// Cek keberadaan mahasiswa di PostgreSQL
	if _, err := s.studentRepo.FindByID(ctx, studentID); err != nil {
		return translatePortfolioError(err)
	}

	portfolio, err := s.portfolioRepo.FindByStudentID(ctx, studentID)
	if err != nil {
		return translatePortfolioError(err)
	}

	return helper.Success(c, fiber.StatusOK, "portofolio berhasil diambil", portfolio)
}

// PUT /students/:id/portfolio
func (s *PortfolioService) Upsert(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	studentID, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id mahasiswa harus berupa angka positif")
	}

	var req model.UpsertPortfolioRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}

	req.Headline = strings.TrimSpace(req.Headline)
	req.Bio = strings.TrimSpace(req.Bio)
	if fieldErrors := helper.ValidateStruct(req); fieldErrors != nil {
		return helper.Validation(fieldErrors)
	}

	// Tolak bila mahasiswa tidak ada di PostgreSQL
	if _, err := s.studentRepo.FindByID(ctx, studentID); err != nil {
		return translatePortfolioError(err)
	}

	portfolio, err := s.portfolioRepo.Upsert(ctx, studentID, req)
	if err != nil {
		return translatePortfolioError(err)
	}

	return helper.Success(c, fiber.StatusOK, "portofolio berhasil diperbarui", portfolio)
}

// POST /students/:id/portfolio/projects
func (s *PortfolioService) AddProject(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	studentID, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id mahasiswa harus berupa angka positif")
	}

	var req model.AddProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Link = strings.TrimSpace(req.Link)
	req.Technologies = NormalizeTechnologies(req.Technologies)

	if fieldErrors := helper.ValidateStruct(req); fieldErrors != nil {
		return helper.Validation(fieldErrors)
	}

	// Cek apakah mahasiswa terdaftar di PostgreSQL
	if _, err := s.studentRepo.FindByID(ctx, studentID); err != nil {
		return translatePortfolioError(err)
	}

	project := model.Project{
		Title:        req.Title,
		Technologies: req.Technologies,
		Link:         req.Link,
		Year:         req.Year,
	}

	updated, err := s.portfolioRepo.AddProject(ctx, studentID, project)
	if err != nil {
		return translatePortfolioError(err)
	}

	return helper.Success(c, fiber.StatusCreated, "proyek berhasil ditambahkan", updated)
}

// DELETE /students/:id/portfolio/projects/:projectId
func (s *PortfolioService) DeleteProject(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	studentID, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id mahasiswa harus berupa angka positif")
	}

	projectID := strings.TrimSpace(c.Params("projectId"))
	if projectID == "" {
		return helper.BadRequest("id proyek tidak boleh kosong")
	}

	// Cek mahasiswa di PostgreSQL
	if _, err := s.studentRepo.FindByID(ctx, studentID); err != nil {
		return translatePortfolioError(err)
	}

	updated, err := s.portfolioRepo.DeleteProject(ctx, studentID, projectID)
	if err != nil {
		return translatePortfolioError(err)
	}

	return helper.Success(c, fiber.StatusOK, "proyek berhasil dihapus", updated)
}

// GET /portfolios/stats/technologies
func (s *PortfolioService) TopTechnologies(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	stats, err := s.portfolioRepo.TopTechnologies(ctx)
	if err != nil {
		return helper.Internal(err)
	}

	return helper.Success(c, fiber.StatusOK, "statistik teknologi berhasil diambil", stats)
}

func translatePortfolioError(err error) error {
	switch {
	case errors.Is(err, repository.ErrInvalidID):
		return helper.BadRequest("format id tidak valid")
	case errors.Is(err, repository.ErrNotFound):
		return helper.NotFound("data mahasiswa tidak ditemukan di database")
	case errors.Is(err, repository.ErrPortfolioNotFound):
		return helper.NotFound("portofolio mahasiswa tidak ditemukan")
	case errors.Is(err, repository.ErrProjectNotFound):
		return helper.NotFound("proyek tidak ditemukan di dalam portofolio")
	case errors.Is(err, repository.ErrProjectLimit):
		return helper.NewError(fiber.StatusConflict, "PROJECT_LIMIT", "jumlah proyek sudah mencapai batas maksimal 20")
	default:
		return helper.Internal(err)
	}
}