package repo

import (
	"context"
	"time"

	"github.com/albin6/api/internal/core/domain"
	"gorm.io/gorm"
)

type PostgresReminderRepo struct {
	db *gorm.DB
}

func NewPostgresReminderRepo(db *gorm.DB) *PostgresReminderRepo {
	return &PostgresReminderRepo{db: db}
}

func (r *PostgresReminderRepo) Create(ctx context.Context, reminder *domain.FollowUpReminder) error {
	return r.db.WithContext(ctx).Create(reminder).Error
}

func (r *PostgresReminderRepo) GetPendingReminders(ctx context.Context, before time.Time) ([]domain.FollowUpReminder, error) {
	var reminders []domain.FollowUpReminder
	err := r.db.WithContext(ctx).
		Where("remind_at <= ? AND sent = ?", before, false).
		Preload("FollowUp").
		Find(&reminders).Error
	return reminders, err
}

func (r *PostgresReminderRepo) MarkAsSent(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.FollowUpReminder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"sent":    true,
			"sent_at": &now,
		}).Error
}

func (r *PostgresReminderRepo) GetUpcomingByUserID(ctx context.Context, userID string) ([]domain.FollowUpReminder, error) {
	var reminders []domain.FollowUpReminder
	err := r.db.WithContext(ctx).
		Joins("JOIN student_follow_ups ON student_follow_ups.id = follow_up_reminders.follow_up_id").
		Where("student_follow_ups.assigned_to = ? AND follow_up_reminders.sent = ?", userID, false).
		Preload("FollowUp").
		Order("follow_up_reminders.remind_at ASC").
		Find(&reminders).Error
	return reminders, err
}
