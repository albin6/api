package service

import (
	"context"
	"github.com/albin6/api/internal/core/domain"
	"github.com/albin6/api/internal/core/port"
	"github.com/albin6/api/pkg/utils"
)

type AdminService struct {
	adminRepo port.AdminRepository
}

func NewAdminService(adminRepo port.AdminRepository) *AdminService {
	return &AdminService{adminRepo: adminRepo}
}

func (s *AdminService) CreateAdmin(ctx context.Context, admin *domain.Admin) error {
	hashedPassword, err := utils.HashPassword(admin.Password)
	if err != nil {
		return err
	}
	admin.Password = hashedPassword
	return s.adminRepo.Create(ctx, admin)
}
