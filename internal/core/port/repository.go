package port

import (
	"context"
	"github.com/albin6/api/internal/core/domain"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id uint) (*domain.User, error)
}

type AdminRepository interface {
	Create(ctx context.Context, admin *domain.Admin) error
	GetByEmail(ctx context.Context, email string) (*domain.Admin, error)
	GetByID(ctx context.Context, id uint) (*domain.Admin, error)
}

type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error
	ValidateRefreshToken(ctx context.Context, userID string, tokenID string) (bool, error)
}
