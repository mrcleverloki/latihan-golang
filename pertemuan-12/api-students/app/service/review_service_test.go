package service_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
	"api-students/app/service"
	"api-students/helper"
)

func setupReviewApp(reviewRepo *fakeReviewRepo, courseRepo *fakeCourseRepo) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			var appErr *helper.AppError
			if errors.As(err, &appErr) {
				return c.Status(appErr.Status).JSON(fiber.Map{
					"success": false,
					"code":    appErr.Code,
					"message": appErr.Message,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})
	svc := service.NewReviewService(reviewRepo, courseRepo, nil)

	api := app.Group("/api/v1")
	api.Use(func(c *fiber.Ctx) error {
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
	tests := []struct {
		name           string
		courseID       string
		courses        map[int]model.Course
		enrolled       map[[2]int]bool
		existingReview []model.Review
		payload        interface{}
		expectedStatus int
	}{
		{
			name:     "sukses membuat ulasan",
			courseID: "10",
			courses:  map[int]model.Course{10: {ID: 10, Code: "CS101", Name: "Pemrograman Go"}},
			enrolled: map[[2]int]bool{{10, 100}: true},
			payload: model.CreateReviewRequest{
				Rating:  5,
				Comment: "Materi sangat jelas dan praktis!",
				Tags:    []string{"Golang", "Backend"},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:     "rating di luar batas",
			courseID: "10",
			courses:  map[int]model.Course{10: {ID: 10, Code: "CS101", Name: "Pemrograman Go"}},
			enrolled: map[[2]int]bool{{10, 100}: true},
			payload: model.CreateReviewRequest{
				Rating:  6,
				Comment: "Rating tidak valid",
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:     "kelas tidak ditemukan",
			courseID: "99",
			courses:  map[int]model.Course{},
			enrolled: map[[2]int]bool{},
			payload: model.CreateReviewRequest{
				Rating:  4,
				Comment: "Kelas ini tidak ada",
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:     "belum terdaftar di kelas",
			courseID: "10",
			courses:  map[int]model.Course{10: {ID: 10, Code: "CS101", Name: "Pemrograman Go"}},
			enrolled: map[[2]int]bool{},
			payload: model.CreateReviewRequest{
				Rating:  4,
				Comment: "Mencoba mengulas tanpa enroll",
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:     "ulasan duplikat",
			courseID: "10",
			courses:  map[int]model.Course{10: {ID: 10, Code: "CS101", Name: "Pemrograman Go"}},
			enrolled: map[[2]int]bool{{10, 100}: true},
			existingReview: []model.Review{
				{CourseID: 10, UserID: 100, Rating: 5, Comment: "Review pertama"},
			},
			payload: model.CreateReviewRequest{
				Rating:  4,
				Comment: "Mencoba review kedua kalinya",
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reviewRepo := newFakeReviewRepo()
			for _, r := range tc.existingReview {
				_, _ = reviewRepo.Create(t.Context(), r)
			}

			courseRepo := &fakeCourseRepo{courses: tc.courses, enrolled: tc.enrolled}
			app := setupReviewApp(reviewRepo, courseRepo)

			body, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/courses/"+tc.courseID+"/reviews", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("request gagal: %v", err)
			}

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("status tidak sesuai, ekspektasi %d dapat %d", tc.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestListReviews(t *testing.T) {
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