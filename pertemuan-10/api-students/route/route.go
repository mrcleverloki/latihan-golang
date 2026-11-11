package route

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"api-students/app/service"
	"api-students/helper"
	"api-students/middleware"
)

type Dependencies struct {
	Pool             *pgxpool.Pool
	MongoClient      *mongo.Client
	JWT              *helper.JWTManager
	Permissions      *helper.PermissionSet
	StudentService   *service.StudentService
	AuthService      *service.AuthService
	CourseService    *service.CourseService
	PortfolioService *service.PortfolioService
}

func Register(app *fiber.App, deps Dependencies) {
	api := app.Group("/api/v1")

	// 1. Endpoint publik
	api.Get("/health", healthCheck(deps.Pool, deps.MongoClient))

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
	students.Get("/:id", middleware.RequirePermission(deps.Permissions, "student:read"), deps.StudentService.Get)
	students.Put("/:id", middleware.RequirePermission(deps.Permissions, "student:update"), deps.StudentService.Replace)
	students.Patch("/:id", middleware.RequirePermission(deps.Permissions, "student:update"), deps.StudentService.Patch)
	students.Delete("/:id", middleware.RequirePermission(deps.Permissions, "student:delete"), deps.StudentService.Delete)

	// Endpoint Portofolio Mahasiswa (MongoDB)
	students.Get("/:id/portfolio", deps.PortfolioService.Get)
	students.Put("/:id/portfolio", deps.PortfolioService.Upsert)
	students.Post("/:id/portfolio/projects", deps.PortfolioService.AddProject)
	students.Delete("/:id/portfolio/projects/:projectId", deps.PortfolioService.DeleteProject)

	// 4. Endpoint courses
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

	// 6. Endpoint statistik agregasi portofolio (publik / terproteksi token)
	portfolios := api.Group("/portfolios",
		middleware.RequireAuth(deps.JWT),
	)
	portfolios.Get("/stats/technologies", deps.PortfolioService.TopTechnologies)
}

func healthCheck(pool *pgxpool.Pool, mongoClient *mongo.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			return helper.ServiceUnavailable("PostgreSQL tidak dapat dihubungi")
		}

		if err := mongoClient.Ping(ctx, nil); err != nil {
			return helper.ServiceUnavailable("MongoDB tidak dapat dihubungi")
		}

		return helper.Success(c, fiber.StatusOK, "server, PostgreSQL, dan MongoDB berjalan normal", fiber.Map{
			"status":   "ok",
			"postgres": "connected",
			"mongodb":  "connected",
		})
	}
}