package utils

import (
	"fmt"
	"github.com/albin6/api/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenClaims struct {
	Sub  string `json:"sub"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateTokens(userID string, role string, cfg *config.Config) (string, string, error) {
	accessClaims := TokenClaims{
		Sub:  userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccessToken, err := accessToken.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", "", err
	}

	refreshClaims := TokenClaims{
		Sub:  userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}

func ValidateToken(tokenString string, cfg *config.Config) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
