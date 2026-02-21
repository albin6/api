package port

import (
	"context"
	"time"

	"github.com/albin6/api/internal/core/domain"
)

type FollowUpRepository interface {
	Create(ctx context.Context, followUp *domain.StudentFollowUp) error
	GetByID(ctx context.Context, id uint) (*domain.StudentFollowUp, error)
	GetAll(ctx context.Context, stage *domain.FollowUpStage, assignedTo *string, page, limit int) ([]domain.StudentFollowUp, int64, error)
	UpdateStage(ctx context.Context, id uint, stage domain.FollowUpStage) error
	GetByStudentID(ctx context.Context, studentID string) ([]domain.StudentFollowUp, error)
}

type ContactLogRepository interface {
	Create(ctx context.Context, log *domain.ContactLog) error
	GetByFollowUpID(ctx context.Context, followUpID uint) ([]domain.ContactLog, error)
}

type MeetingRepository interface {
	Create(ctx context.Context, meeting *domain.Meeting) error
	GetByID(ctx context.Context, id uint) (*domain.Meeting, error)
	GetByFollowUpID(ctx context.Context, followUpID uint) ([]domain.Meeting, error)
	UpdateStatus(ctx context.Context, id uint, status domain.MeetingStatus) error
	AddParticipants(ctx context.Context, meetingID uint, userIDs []string) error
	// GetParticipants removed as it depends on domain.User
}

type MeetingOutcomeRepository interface {
	Create(ctx context.Context, outcome *domain.MeetingOutcome) error
	GetByMeetingID(ctx context.Context, meetingID uint) (*domain.MeetingOutcome, error)
}

type ReminderRepository interface {
	Create(ctx context.Context, reminder *domain.FollowUpReminder) error
	GetPendingReminders(ctx context.Context, before time.Time) ([]domain.FollowUpReminder, error)
	MarkAsSent(ctx context.Context, id uint) error
	GetUpcomingByUserID(ctx context.Context, userID string) ([]domain.FollowUpReminder, error)
}
