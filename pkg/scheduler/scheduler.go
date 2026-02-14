package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/albin6/api/internal/core/port"
)

type Scheduler struct {
	reminderService port.ReminderService
	logger          *slog.Logger
	stopChan        chan struct{}
}

func NewScheduler(reminderService port.ReminderService, logger *slog.Logger) *Scheduler {
	return &Scheduler{
		reminderService: reminderService,
		logger:          logger,
		stopChan:        make(chan struct{}),
	}
}


func (s *Scheduler) Start() {
	ticker := time.NewTicker(5 * time.Minute)

	
	go s.processReminders()

	go func() {
		for {
			select {
			case <-ticker.C:
				s.processReminders()
			case <-s.stopChan:
				ticker.Stop()
				s.logger.Info("Scheduler stopped")
				return
			}
		}
	}()

	s.logger.Info("Scheduler started", "interval", "5 minutes")
}


func (s *Scheduler) Stop() {
	close(s.stopChan)
}

func (s *Scheduler) processReminders() {
	s.logger.Info("Processing pending reminders")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.reminderService.ProcessPendingReminders(ctx); err != nil {
		s.logger.Error("Failed to process reminders", "error", err)
		return
	}

	s.logger.Info("Reminders processed successfully")
}
