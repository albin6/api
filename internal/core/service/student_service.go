package service

import (
	"context"
	"errors"
	"github.com/albin6/api/internal/core/domain"
	"github.com/albin6/api/internal/core/port"
	"math"
	"strings"
)

type StudentService struct {
	repo port.StudentRepository
}

func NewStudentService(repo port.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

func (s *StudentService) GetStudents(ctx context.Context, search string, status *bool, sortBy string, order string, page int, limit int) (map[string]interface{}, error) {
	const maxLimit = 100
	if limit > maxLimit {
		limit = maxLimit
	}
	if limit < 1 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}

	students, total, err := s.repo.GetAll(ctx, search, status, sortBy, order, page, limit)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return map[string]interface{}{
		"data": students,
		"pagination": map[string]interface{}{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	}, nil
}

func (s *StudentService) CreateStudent(ctx context.Context, student *domain.Student) error {
	// 3. Validation
	if student.FullName == "" || student.Email == "" || student.Phone == "" {
		return errors.New("full_name, email, and phone are required")
	}
	// Basic email validation
	if !strings.Contains(student.Email, "@") {
		return errors.New("invalid email format")
	}
	// Basic phone validation (e.g. at least 10 digits)
	if len(student.Phone) < 10 {
		return errors.New("phone must be at least 10 digits")
	}

	return s.repo.Create(ctx, student)
}

func (s *StudentService) SearchStudents(ctx context.Context, query string, limit int) ([]domain.Student, error) {
	if query == "" {
		return []domain.Student{}, nil
	}
	return s.repo.Search(ctx, query, limit)
}
