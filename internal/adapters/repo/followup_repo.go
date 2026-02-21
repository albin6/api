package repo

import (
	"context"

	"github.com/albin6/api/internal/core/domain"
	"gorm.io/gorm"
)

type PostgresFollowUpRepo struct {
	db *gorm.DB
}

func NewPostgresFollowUpRepo(db *gorm.DB) *PostgresFollowUpRepo {
	return &PostgresFollowUpRepo{db: db}
}

func (r *PostgresFollowUpRepo) Create(ctx context.Context, followUp *domain.StudentFollowUp) error {
	return r.db.WithContext(ctx).Create(followUp).Error
}

func (r *PostgresFollowUpRepo) GetByID(ctx context.Context, id uint) (*domain.StudentFollowUp, error) {
	var followUp domain.StudentFollowUp
	err := r.db.WithContext(ctx).
		First(&followUp, id).Error
	if err != nil {
		return nil, err
	}
	return &followUp, nil
}

func (r *PostgresFollowUpRepo) GetAll(ctx context.Context, stage *domain.FollowUpStage, assignedTo *string, page, limit int) ([]domain.StudentFollowUp, int64, error) {
	var followUps []domain.StudentFollowUp
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.StudentFollowUp{})

	if stage != nil {
		query = query.Where("stage = ?", *stage)
	}
	if assignedTo != nil {
		query = query.Where("assigned_to = ?", *assignedTo) // assignedTo is now string, DB column should be varchar
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&followUps).Error

	return followUps, total, err
}

func (r *PostgresFollowUpRepo) UpdateStage(ctx context.Context, id uint, stage domain.FollowUpStage) error {
	return r.db.WithContext(ctx).
		Model(&domain.StudentFollowUp{}).
		Where("id = ?", id).
		Update("stage", stage).Error
}

func (r *PostgresFollowUpRepo) GetByStudentID(ctx context.Context, studentID string) ([]domain.StudentFollowUp, error) {
	var followUps []domain.StudentFollowUp
	err := r.db.WithContext(ctx).
		Where("student_id = ?", studentID).
		Order("created_at DESC").
		Find(&followUps).Error
	return followUps, err
}
