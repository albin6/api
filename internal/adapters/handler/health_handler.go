package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
)

type HealthHandler struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewHealthHandler(db *gorm.DB, redis *redis.Client) *HealthHandler {
	return &HealthHandler{db: db, redis: redis}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	status := gin.H{
		"postgres": "up",
		"redis":    "up",
	}

	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		status["postgres"] = "down"
	}

	if h.redis.Ping(context.Background()).Err() != nil {
		status["redis"] = "down"
	}

	if status["postgres"] == "down" || status["redis"] == "down" {
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	c.JSON(http.StatusOK, status)
}
