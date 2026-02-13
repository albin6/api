package port

import (
	"context"

	"github.com/albin6/api/internal/core/domain"
)

type AuthService interface {
	Signup(ctx context.Context, user *domain.User) error
	Login(ctx context.Context, email, password string) (string, string, error)
	Logout(ctx context.Context, userID string, tokenID string) error
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	SearchUsers(ctx context.Context, query string, limit int) ([]domain.User, error)
}

type AdminService interface {
	CreateAdmin(ctx context.Context, admin *domain.Admin) error
}

type StudentService interface {
	GetStudents(ctx context.Context, search string, status *bool, sortBy string, order string, page int, limit int) (map[string]interface{}, error)
	CreateStudent(ctx context.Context, student *domain.Student) error
	SearchStudents(ctx context.Context, query string, limit int) ([]domain.Student, error)
}

type FollowUpService interface {
	CreateFollowUp(ctx context.Context, studentID, assignedTo, createdBy uint) (*domain.StudentFollowUp, error)
	GetFollowUp(ctx context.Context, id, requestingUserID uint) (*domain.StudentFollowUp, error)
	ListFollowUps(ctx context.Context, stage *domain.FollowUpStage, assignedTo *uint, page, limit int) (map[string]interface{}, error)
	AddContactLog(ctx context.Context, followUpID, userID uint, successful bool, notes string) error
	ScheduleMeeting(ctx context.Context, followUpID, userID uint, scheduledAt string, meetingLink string, participantIDs []uint) (*domain.Meeting, error)
	ListMeetings(ctx context.Context, followUpID, userID uint) ([]domain.Meeting, error)
	CompleteMeeting(ctx context.Context, meetingID, userID uint) error
	SubmitOutcome(ctx context.Context, meetingID, userID uint, status domain.OutcomeStatus, remarks, recordingURL string, nextFollowUpAt *string) error
	RestartFollowUp(ctx context.Context, followUpID, userID uint) error
}

type ReminderService interface {
	GetUpcomingReminders(ctx context.Context, userID uint) ([]domain.FollowUpReminder, error)
	ProcessPendingReminders(ctx context.Context) error
}
