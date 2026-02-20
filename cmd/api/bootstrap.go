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
	redis_adapter "github.com/albin6/api/internal/adapters/redis"
	"github.com/albin6/api/internal/adapters/repo"
	"github.com/albin6/api/internal/adapters/websocket"
	"github.com/albin6/api/internal/core/domain"
	"github.com/albin6/api/internal/core/service"
	"github.com/albin6/api/pkg/database"
	"github.com/albin6/api/pkg/logger"
	"github.com/albin6/api/pkg/scheduler"
	"github.com/albin6/api/pkg/toolapi"
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

	// User, Admin, Student Repos removed
	followUpRepo := repo.NewPostgresFollowUpRepo(db)
	contactLogRepo := repo.NewPostgresContactLogRepo(db)
	meetingRepo := repo.NewPostgresMeetingRepo(db)
	meetingOutcomeRepo := repo.NewPostgresMeetingOutcomeRepo(db)
	reminderRepo := repo.NewPostgresReminderRepo(db)

	hub := websocket.NewHub()
	go hub.Run()

	// Auth, Admin, Student Services removed
	followUpService := service.NewFollowUpService(followUpRepo, contactLogRepo, meetingRepo, meetingOutcomeRepo, reminderRepo, hub)
	reminderService := service.NewReminderService(reminderRepo)

	toolClient := toolapi.NewClient(cfg)
	toolService := service.NewToolService(toolClient, rdb)

	// Auth, Admin, Student Handlers removed
	healthHandler := handler.NewHealthHandler(db, rdb)
	followUpHandler := handler.NewFollowUpHandler(followUpService)
	reminderHandler := handler.NewReminderHandler(reminderService)
	toolHandler := handler.NewToolHandler(toolService)

	db.AutoMigrate(
		// User, Admin, Student removed
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

	r.GET("/ws", websocket.ServeWs(hub, cfg))

	// Auth routes removed
	// Admin routes removed

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		// Student routes removed

		protected.POST("/followups", followUpHandler.CreateFollowUp)
		protected.GET("/followups/:id", followUpHandler.GetFollowUp)
		protected.GET("/followups", followUpHandler.ListFollowUps)
		protected.POST("/followups/:id/restart", followUpHandler.RestartFollowUp)

		protected.POST("/followups/:id/contacts", followUpHandler.AddContactLog)
		protected.GET("/followups/:id/contacts", followUpHandler.GetContactLogs)

		protected.POST("/followups/:id/meetings", followUpHandler.ScheduleMeeting)
		protected.GET("/followups/:id/meetings", followUpHandler.ListMeetings)
		protected.PATCH("/meetings/:id/complete", followUpHandler.CompleteMeeting)
		protected.POST("/meetings/:id/outcome", followUpHandler.SubmitOutcome)

		protected.GET("/reminders/upcoming", reminderHandler.GetUpcomingReminders)

		// User search/profile routes removed
	}

	toolGroup := r.Group("/api/tool")
	toolGroup.Use(middleware.AuthMiddleware(cfg))
	{
		toolGroup.GET("/students", toolHandler.GetStudents)
		toolGroup.GET("/common/page-filters", toolHandler.GetPageFilters)
		toolGroup.GET("/batch", toolHandler.GetBatches)
		toolGroup.GET("/course", toolHandler.GetCourses)
		toolGroup.GET("/domain", toolHandler.GetDomains)
		toolGroup.GET("/employee/roles", toolHandler.GetEmployees)
		toolGroup.GET("/common/status", toolHandler.GetStatusOptions)
	}

	sched := scheduler.NewScheduler(reminderService, log)
	sched.Start()

	// Initialize User Sync Subscriber
	userSubscriber := redis_adapter.NewUserSubscriber(rdb, db)
	go userSubscriber.SubscribeToUserUpdates()

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
