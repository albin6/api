package repo

import (
	"context"
	"github.com/albin6/api/internal/core/domain"
	"gorm.io/gorm"
	"strings"
)

type PostgresStudentRepo struct {
	db *gorm.DB
}

func NewPostgresStudentRepo(db *gorm.DB) *PostgresStudentRepo {
	return &PostgresStudentRepo{db: db}
}

func (r *PostgresStudentRepo) GetAll(ctx context.Context, search string, status *bool, sortBy string, order string, page int, limit int) ([]domain.Student, int64, error) {
	var students []domain.Student
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Student{})

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("full_name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	if status != nil {
		query = query.Where("program_status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	validSortColumns := map[string]bool{"full_name": true, "email": true, "created_at": true}
	if !validSortColumns[sortBy] {
		sortBy = "created_at"
	}
	
	if strings.ToUpper(order) != "ASC" {
		order = "DESC"
	}

	offset := (page - 1) * limit
	err := query.Order(sortBy + " " + order).
		Limit(limit).
		Offset(offset).
		Find(&students).Error

	return students, total, err
}

func (r *PostgresStudentRepo) Create(ctx context.Context, student *domain.Student) error {
	return r.db.WithContext(ctx).Create(student).Error
}
