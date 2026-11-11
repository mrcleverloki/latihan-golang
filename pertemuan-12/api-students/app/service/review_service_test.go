package service_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
	"api-students/app/service"
	"api-students/helper"
)

func setupReviewApp(reviewRepo *fakeReviewRepo, courseRepo *fakeCourseRepo) *fiber.App {
	app := fiber.New()
	svc := service.NewReviewService(reviewRepo, courseRepo, nil)

	api := app.Group("/api/v1")
	api.Use(func(c *fiber.Ctx) error {
		// Mock context auth user
		c.Locals(helper.LocalsAuthUser, model.AuthUser{
			UserID: 100,
			Role:   "mahasiswa",
		})
		return c.Next()
	})

	api.Post("/courses/:id/reviews", svc.Create)
	api.Get("/courses/:id/reviews", svc.ListByCourse)

	return app
}

func TestCreateReview(t *testing.T) {
	// Setiap test membuat fake repository baru (isolated state)
	reviewRepo := newFakeReviewRepo()

	courses := map[int]model.Course{
		10: {ID: 10, Code: "CS101", Name: "Pemrograman Go"},
	}
	enrolled := map[[2]int]bool{
		{10, 100}: true,
	}

	courseRepo := &fakeCourseRepo{courses: courses, enrolled: enrolled}
	app := setupReviewApp(reviewRepo, courseRepo)

	body, _ := json.Marshal(model.CreateReviewRequest{
		Rating:  5,
		Comment: "Materi sangat jelas dan praktis!",
		Tags:    []string{"Golang", "Backend"},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/courses/10/reviews", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}
}

func TestListReviews(t *testing.T) {
	// Buat repo baru yang terisolasi
	reviewRepo := newFakeReviewRepo()
	courseRepo := &fakeCourseRepo{}
	app := setupReviewApp(reviewRepo, courseRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/courses/10/reviews", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}