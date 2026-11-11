package route

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"api-students/app/service"
	"api-students/helper"
	"api-students/middleware"
)

type Dependencies struct {
	Pool           *pgxpool.Pool
	JWT            *helper.JWTManager
	Permissions    *helper.PermissionSet
	StudentService *service.StudentService
	AuthService    *service.AuthService
	CourseService  *service.CourseService
}

func Register(app *fiber.App, deps Dependencies) {
	api := app.Group("/api/v1")

	// 1. Endpoint publik
	api.Get("/health", healthCheck(deps.Pool))

	// 2. Endpoint autentikasi
	auth := api.Group("/auth", middleware.RequireJSON)
	auth.Post("/register", deps.AuthService.Register)
	auth.Post("/login", middleware.LoginRateLimiter(), deps.AuthService.Login)
	auth.Post("/refresh", deps.AuthService.Refresh)
	auth.Post("/logout", deps.AuthService.Logout)
	auth.Get("/me", middleware.RequireAuth(deps.JWT), deps.AuthService.Me)

	// 3. Endpoint students (terlindungi token JWT dan RBAC)
	students := api.Group("/students",
		middleware.RequireJSON,
		middleware.RequireAuth(deps.JWT),
	)

	// List dan Create menggunakan RBAC middleware
	students.Get("/", middleware.RequirePermission(deps.Permissions, "student:list"), deps.StudentService.List)
	students.Post("/", middleware.RequirePermission(deps.Permissions, "student:create"), deps.StudentService.Create)

	// Get, Replace, Patch mengandalkan ownership/role check di service layer
	students.Get("/:id", deps.StudentService.Get)
	students.Put("/:id", deps.StudentService.Replace)
	students.Patch("/:id", deps.StudentService.Patch)

	// Delete menggunakan RBAC middleware
	students.Delete("/:id", middleware.RequirePermission(deps.Permissions, "student:delete"), deps.StudentService.Delete)

	// 4. Endpoint kelas dan KRS
	courses := api.Group("/courses",
		middleware.RequireJSON,
		middleware.RequireAuth(deps.JWT),
	)

	courses.Get("/", deps.CourseService.List)
	courses.Get("/mine", deps.CourseService.Mine)
	courses.Get("/:id", deps.CourseService.Get)
	courses.Post("/:id/enroll", deps.CourseService.Enroll)
	courses.Delete("/:id/enroll", deps.CourseService.Drop)
	courses.Patch("/:id", middleware.RequirePermission(deps.Permissions, "course:manage"), deps.CourseService.UpdateQuota)

	// 5. Endpoint perpindahan kelas (swap)
	enrollments := api.Group("/enrollments",
		middleware.RequireJSON,
		middleware.RequireAuth(deps.JWT),
	)
	enrollments.Post("/swap", deps.CourseService.Swap)

	// "/mine" didaftarkan SEBELUM "/:id". Fiber mencocokkan route sesuai
	// urutan pendaftaran; bila dibalik, "mine" akan dianggap sebagai :id.
	courses.Get("/", deps.CourseService.List)
	courses.Get("/mine", deps.CourseService.Mine)
	courses.Get("/:id", deps.CourseService.Get)
}

func healthCheck(pool *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
		defer cancel()

		dbStatus := "connected"
		if err := pool.Ping(ctx); err != nil {
			dbStatus = "disconnected"
		}

		return helper.Success(c, fiber.StatusOK, "layanan berjalan normal", fiber.Map{
			"status":   "ok",
			"database": dbStatus,
		})
	}
}