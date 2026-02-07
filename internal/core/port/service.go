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
}

type AdminService interface {
	CreateAdmin(ctx context.Context, admin *domain.Admin) error
}
