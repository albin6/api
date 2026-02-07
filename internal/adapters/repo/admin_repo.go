package repo

import (
	"context"
	"github.com/albin6/api/internal/core/domain"
	"gorm.io/gorm"
)

type PostgresAdminRepo struct {
	db *gorm.DB
}

func NewPostgresAdminRepo(db *gorm.DB) *PostgresAdminRepo {
	return &PostgresAdminRepo{db: db}
}

func (r *PostgresAdminRepo) Create(ctx context.Context, admin *domain.Admin) error {
	return r.db.WithContext(ctx).Create(admin).Error
}

func (r *PostgresAdminRepo) GetByEmail(ctx context.Context, email string) (*domain.Admin, error) {
	var admin domain.Admin
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *PostgresAdminRepo) GetByID(ctx context.Context, id uint) (*domain.Admin, error) {
	var admin domain.Admin
	err := r.db.WithContext(ctx).First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}
