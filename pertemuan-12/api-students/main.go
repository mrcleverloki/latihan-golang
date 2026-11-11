package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"

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

	// Multi-role CLI flag: api | worker | all
	role := flag.String("role", config.GetEnv("APP_ROLE", "all"), "api | worker | all")
	flag.Parse()
	runAPI := *role == "api" || *role == "all"
	runWorker := *role == "worker" || *role == "all"
	if !runAPI && !runWorker {
		logger.Error("role tidak dikenal", slog.String("role", *role))
		os.Exit(1)
	}

	// 2. Koneksi ke PostgreSQL
	pool, err := database.NewPool(context.Background())
	if err != nil {
		logger.Error("gagal terhubung ke database postgresql", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var app *fiber.App
	var closeAPIResources func()

	if runAPI {
		app, closeAPIResources = startAPI(logger, pool)
		if app == nil {
			os.Exit(1)
		}
	}

	var stopWorker func()
	if runWorker {
		stopWorker = startWorker(ctx, logger, pool)
		if stopWorker == nil {
			os.Exit(1)
		}
	}

	logger.Info("aplikasi berjalan", slog.String("role", *role))
	<-ctx.Done()

	logger.Info("sinyal berhenti diterima, menutup proses")

	// Urutan shutdown: matikan API dulu, lalu worker
	if app != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			logger.Error("gagal menutup server HTTP dengan rapi", slog.String("error", err.Error()))
		}
		if closeAPIResources != nil {
			closeAPIResources()
		}
	}

	if stopWorker != nil {
		stopWorker()
	}

	logger.Info("aplikasi berhenti dengan rapi")
}

// startAPI menyiapkan seluruh komponen HTTP Fiber
func startAPI(logger *slog.Logger, pool *pgxpool.Pool) (*fiber.App, func()) {
	// Koneksi ke MongoDB
	mongoClient, err := database.NewMongoClient(context.Background())
	if err != nil {
		logger.Error("gagal terhubung ke database mongodb", slog.String("error", err.Error()))
		return nil, nil
	}

	mongoDatabaseName := config.GetEnv("MONGO_DATABASE", "praktikum_backend")
	mongoDB := mongoClient.Database(mongoDatabaseName)

	// Auth & JWT
	jwtSecret := config.GetEnv("JWT_SECRET", "kunci-rahasia-minimal-32-karakter-harus-aman-banget")
	jwtIssuer := config.GetEnv("JWT_ISSUER", "api-students")
	accessTTL := 15 * time.Minute
	refreshTTL := 7 * 24 * time.Hour
	jwtManager := helper.NewJWTManager(jwtSecret, jwtIssuer, accessTTL)

	// Inisialisasi Repository
	studentRepo := repository.NewStudentRepository(pool)
	userRepo := repository.NewUserRepository(pool)
	tokenRepo := repository.NewTokenRepository(pool)
	roleRepo := repository.NewRoleRepository(pool)
	courseRepo := repository.NewCourseRepository(pool)
	outboxRepo := repository.NewOutboxRepository(pool)
	notificationRepo := repository.NewNotificationRepository(pool)
	txManager := repository.NewTxManager(pool)

	portfolioRepo := repository.NewPortfolioRepository(mongoDB)
	if err := portfolioRepo.EnsureIndexes(context.Background()); err != nil {
		logger.Error("gagal membuat index mongodb portfolio", slog.String("error", err.Error()))
		return nil, nil
	}

	// Muat Permissions
	rawPermissions, err := roleRepo.LoadPermissions(context.Background())
	if err != nil {
		logger.Error("gagal memuat permission", slog.String("error", err.Error()))
		return nil, nil
	}
	perms := helper.NewPermissionSet(rawPermissions)

	// Inspector Asynq untuk endpoint monitoring admin
	inspector := asynq.NewInspector(config.RedisOpt())

	// Inisialisasi Service
	studentService := service.NewStudentService(studentRepo, perms)
	authService := service.NewAuthService(userRepo, tokenRepo, jwtManager, perms, refreshTTL)
	courseService := service.NewCourseService(courseRepo, outboxRepo, txManager)
	portfolioService := service.NewPortfolioService(portfolioRepo, studentRepo)
	notificationService := service.NewNotificationService(notificationRepo, outboxRepo, inspector)

	deps := route.Dependencies{
		Pool:                pool,
		MongoClient:         mongoClient,
		JWT:                 jwtManager,
		Permissions:         perms,
		StudentService:      studentService,
		AuthService:         authService,
		CourseService:       courseService,
		PortfolioService:    portfolioService,
		NotificationService: notificationService,
	}

	app := config.NewApp(logger, deps)
	port := config.GetEnv("APP_PORT", "3000")

	go func() {
		if err := app.Listen(":" + port); err != nil {
			logger.Error("server HTTP berhenti", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()
	logger.Info("server HTTP berjalan", slog.String("port", port))

	closeFn := func() {
		_ = inspector.Close()
		disconnectCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(disconnectCtx); err != nil {
			logger.Error("gagal disconnect mongodb", slog.String("error", err.Error()))
		}
	}

	return app, closeFn
}