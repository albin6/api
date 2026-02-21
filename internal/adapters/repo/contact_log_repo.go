package repo

import (
	"context"

	"github.com/albin6/api/internal/core/domain"
	"gorm.io/gorm"
)

type PostgresContactLogRepo struct {
	db *gorm.DB
}

func NewPostgresContactLogRepo(db *gorm.DB) *PostgresContactLogRepo {
	return &PostgresContactLogRepo{db: db}
}

func (r *PostgresContactLogRepo) Create(ctx context.Context, log *domain.ContactLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *PostgresContactLogRepo) GetByFollowUpID(ctx context.Context, followUpID uint) ([]domain.ContactLog, error) {
	var logs []domain.ContactLog
	err := r.db.WithContext(ctx).
		Where("follow_up_id = ?", followUpID).
		Order("contacted_at DESC").
		Find(&logs).Error
	return logs, err
}
