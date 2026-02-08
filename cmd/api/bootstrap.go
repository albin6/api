package main

import (
	"context"
	"github.com/albin6/api/config"
	"github.com/albin6/api/internal/adapters/api/middleware"
	"github.com/albin6/api/internal/adapters/handler"
	"github.com/albin6/api/internal/adapters/repo"
	"github.com/albin6/api/internal/adapters/storage"
	"github.com/albin6/api/internal/core/domain"
	"github.com/albin6/api/internal/core/service"
	"github.com/albin6/api/pkg/database"
	"github.com/albin6/api/pkg/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Router *gin.Engine
	DB     *gorm.DB
	Config *config.Config
	Logger *slog.Logger
}

func NewServer(cfg *config.Config) *Server {
	log := logger.InitLogger(cfg.Environment)
	db := database.NewPostgresDB(cfg)
	rdb := database.NewRedisClient(cfg)

	userRepo := repo.NewPostgresUserRepo(db)
	adminRepo := repo.NewPostgresAdminRepo(db)
	studentRepo := repo.NewPostgresStudentRepo(db)
	tokenRepo := storage.NewRedisTokenRepo(rdb)

	authService := service.NewAuthService(userRepo, tokenRepo, cfg)
	adminService := service.NewAdminService(adminRepo)
	studentService := service.NewStudentService(studentRepo)

	authHandler := handler.NewAuthHandler(authService)
	adminHandler := handler.NewAdminHandler(adminService)
	studentHandler := handler.NewStudentHandler(studentService)
	healthHandler := handler.NewHealthHandler(db, rdb)

	db.AutoMigrate(&domain.User{}, &domain.Admin{}, &domain.Student{})

	if cfg.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimitMiddleware(rdb))

	r.GET("/health", healthHandler.HealthCheck)


	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", authHandler.Signup)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.Refresh)
		authGroup.POST("/logout", authHandler.Logout)
	}

	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.AdminKeyMiddleware(cfg))
	{
		adminGroup.POST("/create", adminHandler.CreateAdmin)
	}

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.POST("/students", studentHandler.CreateStudent)
		protected.GET("/students", studentHandler.GetStudents)
		protected.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			role, _ := c.Get("role")
			c.JSON(200, gin.H{"message": "Access granted", "userID": userID, "role": role})
		})
	}

	return &Server{
		Router: r,
		DB:     db,
		Config: cfg,
		Logger: log,
	}
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:    ":" + s.Config.Port,
		Handler: s.Router,
	}

	go func() {
		s.Logger.Info("Starting server", "port", s.Config.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Error("listen", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.Logger.Error("Server forced to shutdown", "error", err)
		return err
	}

	sqlDB, _ := s.DB.DB()
	if err := sqlDB.Close(); err != nil {
		s.Logger.Error("Database close error", "error", err)
	}

	s.Logger.Info("Server exiting")
	return nil
}
