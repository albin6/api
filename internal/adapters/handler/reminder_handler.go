package handler

import (
	"net/http"

	"github.com/albin6/api/internal/core/port"
	"github.com/gin-gonic/gin"
)

type ReminderHandler struct {
	service port.ReminderService
}

func NewReminderHandler(service port.ReminderService) *ReminderHandler {
	return &ReminderHandler{service: service}
}


func (h *ReminderHandler) GetUpcomingReminders(c *gin.Context) {
	userID := getUserIDFromContext(c)

	reminders, err := h.service.GetUpcomingReminders(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Upcoming reminders retrieved successfully",
		"data":    reminders,
	})
}
