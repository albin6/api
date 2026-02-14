package service

import (
	"context"
	"time"

	"github.com/albin6/api/internal/core/domain"
	"github.com/albin6/api/internal/core/port"
)

type ReminderService struct {
	reminderRepo port.ReminderRepository
}

func NewReminderService(reminderRepo port.ReminderRepository) *ReminderService {
	return &ReminderService{
		reminderRepo: reminderRepo,
	}
}

func (s *ReminderService) GetUpcomingReminders(ctx context.Context, userID uint) ([]domain.FollowUpReminder, error) {
	return s.reminderRepo.GetUpcomingByUserID(ctx, userID)
}

func (s *ReminderService) ProcessPendingReminders(ctx context.Context) error {
	now := time.Now()
	reminders, err := s.reminderRepo.GetPendingReminders(ctx, now)
	if err != nil {
		return err
	}

	for _, reminder := range reminders {
		
		
		
		if err := s.reminderRepo.MarkAsSent(ctx, reminder.ID); err != nil {
			
			continue
		}
	}

	return nil
}
