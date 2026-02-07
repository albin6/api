package middleware

import (
	"github.com/albin6/api/config"
	"github.com/albin6/api/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenString, cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("userID", claims.Sub)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AdminKeyMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("x-admin-auth-key")
		if apiKey == "" || apiKey != cfg.AdminSecretKey {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid or missing admin key"})
			return
		}
		c.Next()
	}
}
