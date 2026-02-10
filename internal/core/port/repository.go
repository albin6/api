package port

import (
	"context"
	"time"

	"github.com/albin6/api/internal/core/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id uint) (*domain.User, error)
}

type AdminRepository interface {
	Create(ctx context.Context, admin *domain.Admin) error
	GetByEmail(ctx context.Context, email string) (*domain.Admin, error)
	GetByID(ctx context.Context, id uint) (*domain.Admin, error)
}

type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error
	ValidateRefreshToken(ctx context.Context, userID string, tokenID string) (bool, error)
}

type StudentRepository interface {
	GetAll(ctx context.Context, search string, status *bool, sortBy string, order string, page int, limit int) ([]domain.Student, int64, error)
	Create(ctx context.Context, student *domain.Student) error
}

type FollowUpRepository interface {
	Create(ctx context.Context, followUp *domain.StudentFollowUp) error
	GetByID(ctx context.Context, id uint) (*domain.StudentFollowUp, error)
	GetAll(ctx context.Context, stage *domain.FollowUpStage, assignedTo *uint, page, limit int) ([]domain.StudentFollowUp, int64, error)
	UpdateStage(ctx context.Context, id uint, stage domain.FollowUpStage) error
	GetByStudentID(ctx context.Context, studentID uint) ([]domain.StudentFollowUp, error)
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
	AddParticipants(ctx context.Context, meetingID uint, userIDs []uint) error
	GetParticipants(ctx context.Context, meetingID uint) ([]domain.User, error)
}

type MeetingOutcomeRepository interface {
	Create(ctx context.Context, outcome *domain.MeetingOutcome) error
	GetByMeetingID(ctx context.Context, meetingID uint) (*domain.MeetingOutcome, error)
}

type ReminderRepository interface {
	Create(ctx context.Context, reminder *domain.FollowUpReminder) error
	GetPendingReminders(ctx context.Context, before time.Time) ([]domain.FollowUpReminder, error)
	MarkAsSent(ctx context.Context, id uint) error
	GetUpcomingByUserID(ctx context.Context, userID uint) ([]domain.FollowUpReminder, error)
}
