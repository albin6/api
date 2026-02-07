package service

import (
	"context"
	"errors"
	"github.com/albin6/api/config"
	"github.com/albin6/api/internal/core/domain"
	"github.com/albin6/api/internal/core/port"
	"github.com/albin6/api/pkg/utils"
	"github.com/google/uuid"
	"strconv"
)

type AuthService struct {
	userRepo  port.UserRepository
	tokenRepo port.TokenRepository
	cfg       *config.Config
}

func NewAuthService(userRepo port.UserRepository, tokenRepo port.TokenRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		cfg:       cfg,
	}
}

func (s *AuthService) Signup(ctx context.Context, user *domain.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	if user.Role == "" {
		user.Role = domain.RoleMember
	}

	return s.userRepo.Create(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("invalid credentials")
	}

	return s.generateAndStoreTokens(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, userID string, tokenID string) error {
	return s.tokenRepo.DeleteRefreshToken(ctx, userID, tokenID)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := utils.ValidateToken(refreshToken, s.cfg)
	if err != nil {
		return "", "", err
	}

	uid, _ := strconv.Atoi(claims.Sub)
	user, err := s.userRepo.GetByID(ctx, uint(uid))
	if err != nil {
		return "", "", err
	}

	return s.generateAndStoreTokens(ctx, user)
}

func (s *AuthService) generateAndStoreTokens(ctx context.Context, user *domain.User) (string, string, error) {
	strUserID := strconv.Itoa(int(user.ID))
	accessToken, refreshToken, err := utils.GenerateTokens(strUserID, user.Role.String(), s.cfg)
	if err != nil {
		return "", "", err
	}

	tokenID := uuid.New().String()

	err = s.tokenRepo.SetRefreshToken(ctx, strUserID, tokenID, s.cfg.RefreshTokenExpiry)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
