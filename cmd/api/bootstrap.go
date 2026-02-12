package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/albin6/api/config"
	"github.com/albin6/api/internal/adapters/api/middleware"
	"github.com/albin6/api/internal/adapters/handler"
	"github.com/albin6/api/internal/adapters/repo"
	"github.com/albin6/api/internal/adapters/storage"
	"github.com/albin6/api/internal/core/domain"
	"github.com/albin6/api/internal/core/service"
	"github.com/albin6/api/pkg/database"
	"github.com/albin6/api/pkg/logger"
	"github.com/albin6/api/pkg/scheduler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	Router    *gin.Engine
	DB        *gorm.DB
	Config    *config.Config
	Logger    *slog.Logger
	Scheduler *scheduler.Scheduler
}

func NewServer(cfg *config.Config) *Server {
	log := logger.InitLogger(cfg.Environment)
	db := database.NewPostgresDB(cfg)
	rdb := database.NewRedisClient(cfg)

	userRepo := repo.NewPostgresUserRepo(db)
	adminRepo := repo.NewPostgresAdminRepo(db)
	studentRepo := repo.NewPostgresStudentRepo(db)
	tokenRepo := storage.NewRedisTokenRepo(rdb)
	followUpRepo := repo.NewPostgresFollowUpRepo(db)
	contactLogRepo := repo.NewPostgresContactLogRepo(db)
	meetingRepo := repo.NewPostgresMeetingRepo(db)
	meetingOutcomeRepo := repo.NewPostgresMeetingOutcomeRepo(db)
	reminderRepo := repo.NewPostgresReminderRepo(db)

	authService := service.NewAuthService(userRepo, tokenRepo, cfg)
	adminService := service.NewAdminService(adminRepo)
	studentService := service.NewStudentService(studentRepo)
	followUpService := service.NewFollowUpService(followUpRepo, contactLogRepo, meetingRepo, meetingOutcomeRepo, reminderRepo, studentRepo, userRepo)
	reminderService := service.NewReminderService(reminderRepo)

	authHandler := handler.NewAuthHandler(authService)
	adminHandler := handler.NewAdminHandler(adminService)
	studentHandler := handler.NewStudentHandler(studentService)
	healthHandler := handler.NewHealthHandler(db, rdb)
	followUpHandler := handler.NewFollowUpHandler(followUpService)
	reminderHandler := handler.NewReminderHandler(reminderService)

	db.AutoMigrate(
		&domain.User{},
		&domain.Admin{},
		&domain.Student{},
		&domain.StudentFollowUp{},
		&domain.ContactLog{},
		&domain.Meeting{},
		&domain.MeetingParticipant{},
		&domain.MeetingOutcome{},
		&domain.FollowUpReminder{},
	)

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
		// Student endpoints
		protected.POST("/students", studentHandler.CreateStudent)
		protected.GET("/students", studentHandler.GetStudents)
		protected.GET("/students/search", studentHandler.SearchStudents)

		// Follow-up endpoints
		protected.POST("/followups", followUpHandler.CreateFollowUp)
		protected.GET("/followups/:id", followUpHandler.GetFollowUp)
		protected.GET("/followups", followUpHandler.ListFollowUps)
		protected.POST("/followups/:id/restart", followUpHandler.RestartFollowUp)

		// Contact log endpoints
		protected.POST("/followups/:id/contacts", followUpHandler.AddContactLog)
		protected.GET("/followups/:id/contacts", followUpHandler.GetContactLogs)

		// Meeting endpoints
		protected.POST("/followups/:id/meetings", followUpHandler.ScheduleMeeting)
		protected.GET("/followups/:id/meetings", followUpHandler.ListMeetings)
		protected.PATCH("/meetings/:id/complete", followUpHandler.CompleteMeeting)
		protected.POST("/meetings/:id/outcome", followUpHandler.SubmitOutcome)

		// Reminder endpoints
		protected.GET("/reminders/upcoming", reminderHandler.GetUpcomingReminders)

		// Search endpoints
		protected.GET("/users/search", authHandler.SearchUsers)

		protected.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			role, _ := c.Get("role")
			c.JSON(200, gin.H{"message": "Access granted", "userID": userID, "role": role})
		})
	}

	// Initialize and start scheduler
	sched := scheduler.NewScheduler(reminderService, log)
	sched.Start()

	return &Server{
		Router:    r,
		DB:        db,
		Config:    cfg,
		Logger:    log,
		Scheduler: sched,
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

	// Stop scheduler
	s.Scheduler.Stop()

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
