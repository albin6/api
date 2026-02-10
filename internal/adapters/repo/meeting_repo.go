package repo

import (
	"context"

	"github.com/albin6/api/internal/core/domain"
	"gorm.io/gorm"
)

type PostgresMeetingRepo struct {
	db *gorm.DB
}

func NewPostgresMeetingRepo(db *gorm.DB) *PostgresMeetingRepo {
	return &PostgresMeetingRepo{db: db}
}

func (r *PostgresMeetingRepo) Create(ctx context.Context, meeting *domain.Meeting) error {
	return r.db.WithContext(ctx).Create(meeting).Error
}

func (r *PostgresMeetingRepo) GetByID(ctx context.Context, id uint) (*domain.Meeting, error) {
	var meeting domain.Meeting
	err := r.db.WithContext(ctx).
		Preload("FollowUp").
		Preload("FollowUp.Student").
		Preload("Creator").
		First(&meeting, id).Error
	if err != nil {
		return nil, err
	}
	return &meeting, nil
}

func (r *PostgresMeetingRepo) GetByFollowUpID(ctx context.Context, followUpID uint) ([]domain.Meeting, error) {
	var meetings []domain.Meeting
	err := r.db.WithContext(ctx).
		Where("follow_up_id = ?", followUpID).
		Preload("Creator").
		Order("scheduled_at DESC").
		Find(&meetings).Error
	return meetings, err
}

func (r *PostgresMeetingRepo) UpdateStatus(ctx context.Context, id uint, status domain.MeetingStatus) error {
	return r.db.WithContext(ctx).
		Model(&domain.Meeting{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *PostgresMeetingRepo) AddParticipants(ctx context.Context, meetingID uint, userIDs []uint) error {
	participants := make([]domain.MeetingParticipant, len(userIDs))
	for i, userID := range userIDs {
		participants[i] = domain.MeetingParticipant{
			MeetingID: meetingID,
			UserID:    userID,
		}
	}
	return r.db.WithContext(ctx).Create(&participants).Error
}

func (r *PostgresMeetingRepo) GetParticipants(ctx context.Context, meetingID uint) ([]domain.User, error) {
	var participants []domain.MeetingParticipant
	err := r.db.WithContext(ctx).
		Where("meeting_id = ?", meetingID).
		Preload("User").
		Find(&participants).Error
	if err != nil {
		return nil, err
	}

	users := make([]domain.User, len(participants))
	for i, p := range participants {
		users[i] = p.User
	}
	return users, nil
}
