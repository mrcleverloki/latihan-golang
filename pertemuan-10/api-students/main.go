package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-students/app/repository"
	"api-students/app/service"
	"api-students/config"
	"api-students/database"
	"api-students/helper"
	"api-students/route"
)

func main() {
	// 1. Inisialisasi environment dan logger terstruktur
	config.LoadEnv()
	logger := config.NewLogger()

	// 2. Koneksi ke PostgreSQL
	pool, err := database.NewPool(context.Background())
	if err != nil {
		logger.Error("gagal terhubung ke database postgresql", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	// 3. Koneksi ke MongoDB
	mongoClient, err := database.NewMongoClient(context.Background())
	if err != nil {
		logger.Error("gagal terhubung ke database mongodb", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		disconnectCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(disconnectCtx); err != nil {
			logger.Error("gagal disconnect mongodb", slog.String("error", err.Error()))
		}
	}()

	mongoDatabaseName := config.GetEnv("MONGO_DATABASE", "praktikum_backend")
	mongoDB := mongoClient.Database(mongoDatabaseName)

	// 4. Konfigurasi Auth & JWT
	jwtSecret := config.GetEnv("JWT_SECRET", "kunci-rahasia-minimal-32-karakter-harus-aman-banget")
	jwtIssuer := config.GetEnv("JWT_ISSUER", "api-students")
	accessTTL := 15 * time.Minute
	refreshTTL := 7 * 24 * time.Hour

	jwtManager := helper.NewJWTManager(jwtSecret, jwtIssuer, accessTTL)

	// 5. Dependency Injection: Repository & TxManager
	studentRepo := repository.NewStudentRepository(pool)
	userRepo := repository.NewUserRepository(pool)
	tokenRepo := repository.NewTokenRepository(pool)
	roleRepo := repository.NewRoleRepository(pool)
	courseRepo := repository.NewCourseRepository(pool)
	txManager := repository.NewTxManager(pool)

	portfolioRepo := repository.NewPortfolioRepository(mongoDB)
	if err := portfolioRepo.EnsureIndexes(context.Background()); err != nil {
		logger.Error("gagal membuat index mongodb portfolio", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// 6. Muat Role Permissions dari DB (Fail Closed)
	rawPermissions, err := roleRepo.LoadPermissions(context.Background())
	if err != nil {
		logger.Error("gagal memuat permission", slog.String("error", err.Error()))
		os.Exit(1)
	}
	perms := helper.NewPermissionSet(rawPermissions)
	logger.Info("permission dimuat", slog.Any("roles", perms.KnownRoles()))

	// 7. Services
	studentService := service.NewStudentService(studentRepo, perms)
	authService := service.NewAuthService(userRepo, tokenRepo, jwtManager, perms, refreshTTL)
	courseService := service.NewCourseService(courseRepo, txManager)
	portfolioService := service.NewPortfolioService(portfolioRepo, studentRepo)

	// 8. Inisialisasi App dengan Dependencies
	deps := route.Dependencies{
		Pool:             pool,
		MongoClient:      mongoClient,
		JWT:              jwtManager,
		Permissions:      perms,
		StudentService:   studentService,
		AuthService:      authService,
		CourseService:    courseService,
		PortfolioService: portfolioService,
	}
	app := config.NewApp(logger, deps)
	port := config.GetEnv("APP_PORT", "3000")

	go func() {
		if err := app.Listen(":" + port); err != nil {
			logger.Error("server berhenti", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	logger.Info("server berjalan", slog.String("port", port))

	// 9. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("sinyal berhenti diterima, menutup server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("gagal menutup server dengan rapi", slog.String("error", err.Error()))
	}

	logger.Info("server berhenti dengan rapi")
}