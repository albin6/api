package repo

import (
	"context"

	"github.com/albin6/api/internal/core/domain"
	"gorm.io/gorm"
)

type PostgresMeetingOutcomeRepo struct {
	db *gorm.DB
}

func NewPostgresMeetingOutcomeRepo(db *gorm.DB) *PostgresMeetingOutcomeRepo {
	return &PostgresMeetingOutcomeRepo{db: db}
}

func (r *PostgresMeetingOutcomeRepo) Create(ctx context.Context, outcome *domain.MeetingOutcome) error {
	return r.db.WithContext(ctx).Create(outcome).Error
}

func (r *PostgresMeetingOutcomeRepo) GetByMeetingID(ctx context.Context, meetingID uint) (*domain.MeetingOutcome, error) {
	var outcome domain.MeetingOutcome
	err := r.db.WithContext(ctx).
		Where("meeting_id = ?", meetingID).
		Preload("Meeting").
		First(&outcome).Error
	if err != nil {
		return nil, err
	}
	return &outcome, nil
}
