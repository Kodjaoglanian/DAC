package http

import (
	"dac/project-tracker/internal/config"
	"dac/project-tracker/internal/domain/repository"
	"dac/project-tracker/internal/domain/service"
	"dac/project-tracker/internal/infrastructure/middleware"
	"dac/project-tracker/internal/transport/http/handler"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// NewRouter creates and configures the Gin router with all routes
func NewRouter(
	authService *service.AuthService,
	userService *service.UserService,
	projectService *service.ProjectService,
	taskService *service.TaskService,
	memberService *service.MemberService,
	reportService *service.ReportService,
	memberRepo repository.MemberRepository,
	db *gorm.DB,
	logger *zap.Logger,
	cfg *config.Config,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware(cfg))
	r.Use(middleware.LoggerMiddleware(logger))

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	projectHandler := handler.NewProjectHandler(projectService, memberRepo)
	taskHandler := handler.NewTaskHandler(taskService)
	memberHandler := handler.NewMemberHandler(memberService)
	reportHandler := handler.NewReportHandler(reportService)

	// Middleware
	authMiddleware := middleware.AuthMiddleware(cfg)
	projectMemberMiddleware := middleware.ProjectMemberMiddleware(memberRepo, db)

	api := r.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/me", authMiddleware, authHandler.Me)
	}

	// User routes (protected)
	users := api.Group("/users", authMiddleware)
	{
		users.GET("", middleware.RoleMiddleware("admin", "manager"), userHandler.List)
		users.GET("/:id", userHandler.Get)
		users.PUT("/:id", userHandler.Update)
		users.DELETE("/:id", middleware.RoleMiddleware("admin"), userHandler.Delete)
	}

	// Project routes (protected)
	projects := api.Group("/projects", authMiddleware)
	{
		projects.GET("", projectHandler.List)
		projects.POST("", projectHandler.Create)
		projects.GET("/:id", projectMemberMiddleware, projectHandler.Get)
		projects.PUT("/:id", projectMemberMiddleware, projectHandler.Update)
		projects.DELETE("/:id", projectMemberMiddleware, projectHandler.Delete)
		projects.PATCH("/:id/status", projectMemberMiddleware, projectHandler.UpdateStatus)
	}

	// Task routes (protected)
	tasks := api.Group("/projects/:id/tasks", authMiddleware, projectMemberMiddleware)
	{
		tasks.GET("", taskHandler.List)
		tasks.POST("", taskHandler.Create)
	}

	api.GET("/tasks/:id", authMiddleware, projectMemberMiddleware, taskHandler.Get)
	api.PUT("/tasks/:id", authMiddleware, projectMemberMiddleware, taskHandler.Update)
	api.DELETE("/tasks/:id", authMiddleware, projectMemberMiddleware, taskHandler.Delete)
	api.PATCH("/tasks/:id/status", authMiddleware, projectMemberMiddleware, taskHandler.UpdateStatus)

	// Member routes (protected)
	members := api.Group("/projects/:id/members", authMiddleware, projectMemberMiddleware)
	{
		members.GET("", memberHandler.List)
		members.POST("", memberHandler.Add)
		members.PATCH("/:userId", memberHandler.UpdateRole)
		members.DELETE("/:userId", memberHandler.Remove)
	}

	// Report routes (protected)
	reports := api.Group("/reports", authMiddleware)
	{
		reports.GET("/dashboard", reportHandler.Dashboard)
		reports.GET("/projects/:id", projectMemberMiddleware, reportHandler.ProjectReport)
		reports.GET("/projects/by-status", reportHandler.ProjectsByStatus)
		reports.GET("/tasks/by-status", reportHandler.TasksByStatus)
	}

	return r
}
