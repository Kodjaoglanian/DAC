package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dac/project-tracker/internal/config"
	"dac/project-tracker/internal/domain/service"
	"dac/project-tracker/internal/infrastructure/database"
	"dac/project-tracker/internal/infrastructure/middleware"
	repoimpl "dac/project-tracker/internal/infrastructure/repository"
	httptransport "dac/project-tracker/internal/transport/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Database
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Redis
	redisClient := database.NewRedisClient(cfg)
	defer redisClient.Close()

	// Repositories
	userRepo := repoimpl.NewUserRepository(db)
	projectRepo := repoimpl.NewProjectRepository(db)
	taskRepo := repoimpl.NewTaskRepository(db)
	memberRepo := repoimpl.NewMemberRepository(db)
	reportRepo := repoimpl.NewReportRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpirationHours)
	userService := service.NewUserService(userRepo)
	projectService := service.NewProjectService(projectRepo, memberRepo)
	taskService := service.NewTaskService(taskRepo, memberRepo, redisClient)
	memberService := service.NewMemberService(memberRepo, projectRepo)
	reportService := service.NewReportService(reportRepo, projectRepo, taskRepo, memberRepo, redisClient)

	// Logger
	logger := middleware.NewLogger(cfg)

	// Router
	router := httptransport.NewRouter(
		authService,
		userService,
		projectService,
		taskService,
		memberService,
		reportService,
		memberRepo,
		db,
		logger,
		cfg,
	)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	fmt.Printf("Server running on http://%s:%s\n", cfg.ServerHost, cfg.ServerPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exited")
}
