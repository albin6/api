package port

import (
	"context"

	"github.com/albin6/api/internal/core/domain"
)

type FollowUpService interface {
	CreateFollowUp(ctx context.Context, studentID string, assignedTo, createdBy string) (*domain.StudentFollowUp, error)
	GetFollowUp(ctx context.Context, id uint, requestingUserID string) (*domain.StudentFollowUp, error)
	ListFollowUps(ctx context.Context, stage *domain.FollowUpStage, assignedTo *string, page, limit int) (map[string]interface{}, error)
	AddContactLog(ctx context.Context, followUpID uint, userID string, successful bool, notes string) error
	ScheduleMeeting(ctx context.Context, followUpID uint, userID string, scheduledAt string, meetingLink string, participantIDs []string) (*domain.Meeting, error)
	ListMeetings(ctx context.Context, followUpID uint, userID string) ([]domain.Meeting, error)
	CompleteMeeting(ctx context.Context, meetingID uint, userID string) error
	SubmitOutcome(ctx context.Context, meetingID uint, userID string, status domain.OutcomeStatus, remarks, recordingURL string, nextFollowUpAt *string) error
	RestartFollowUp(ctx context.Context, followUpID uint, userID string) error
}

type ReminderService interface {
	GetUpcomingReminders(ctx context.Context, userID string) ([]domain.FollowUpReminder, error)
	ProcessPendingReminders(ctx context.Context) error
}
