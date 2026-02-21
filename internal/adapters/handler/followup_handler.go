package handler

import (
	"net/http"
	"strconv"

	"github.com/albin6/api/internal/core/domain"
	"github.com/albin6/api/internal/core/port"
	"github.com/gin-gonic/gin"
)

type FollowUpHandler struct {
	service port.FollowUpService
}

func NewFollowUpHandler(service port.FollowUpService) *FollowUpHandler {
	return &FollowUpHandler{service: service}
}

func (h *FollowUpHandler) CreateFollowUp(c *gin.Context) {
	var req struct {
		StudentID  string `json:"student_id" binding:"required"`
		AssignedTo string `json:"assigned_to" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"errors":  []gin.H{{"field": "request_body", "error": err.Error()}},
		})
		return
	}

	createdBy := getUserIDFromContext(c)

	followUp, err := h.service.CreateFollowUp(c.Request.Context(), req.StudentID, req.AssignedTo, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Follow-up created successfully",
		"data":    followUp,
	})
}

func (h *FollowUpHandler) GetFollowUp(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid follow-up ID",
		})
		return
	}

	userID := getUserIDFromContext(c)
	followUp, err := h.service.GetFollowUp(c.Request.Context(), uint(id), userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Follow-up retrieved successfully",
		"data":    followUp,
	})
}

func (h *FollowUpHandler) ListFollowUps(c *gin.Context) {
	stageStr := c.Query("stage")
	assignedToStr := c.Query("assigned_to")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	var stage *domain.FollowUpStage
	if stageStr != "" {
		s := domain.FollowUpStage(stageStr)
		if s.IsValid() {
			stage = &s
		}
	}

	var assignedTo *string
	if assignedToStr != "" {
		assignedTo = &assignedToStr
	}

	result, err := h.service.ListFollowUps(c.Request.Context(), stage, assignedTo, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "Follow-ups retrieved successfully",
		"data":       result["data"],
		"pagination": result["pagination"],
	})
}

func (h *FollowUpHandler) AddContactLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid follow-up ID",
		})
		return
	}

	var req struct {
		Successful bool   `json:"successful" binding:"required"`
		Notes      string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"errors":  []gin.H{{"field": "request_body", "error": err.Error()}},
		})
		return
	}

	if len(req.Notes) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Notes cannot exceed 500 characters",
		})
		return
	}

	userID := getUserIDFromContext(c)
	err = h.service.AddContactLog(c.Request.Context(), uint(id), userID, req.Successful, req.Notes)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "unauthorized: only the assigned member can add contact logs" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Contact log added successfully",
	})
}

func (h *FollowUpHandler) GetContactLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid follow-up ID",
		})
		return
	}

	userID := getUserIDFromContext(c)
	_, err = h.service.GetFollowUp(c.Request.Context(), uint(id), userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Contact logs retrieved successfully",
		"data":    []interface{}{},
	})
}

func (h *FollowUpHandler) ScheduleMeeting(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid follow-up ID",
		})
		return
	}

	var req struct {
		ScheduledAt    string   `json:"scheduled_at" binding:"required"`
		MeetingLink    string   `json:"meeting_link" binding:"required"`
		ParticipantIDs []string `json:"participant_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"errors":  []gin.H{{"field": "request_body", "error": err.Error()}},
		})
		return
	}

	userID := getUserIDFromContext(c)
	meeting, err := h.service.ScheduleMeeting(c.Request.Context(), uint(id), userID, req.ScheduledAt, req.MeetingLink, req.ParticipantIDs)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "unauthorized: only the assigned member can schedule meetings" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Meeting scheduled successfully",
		"data":    meeting,
	})
}

func (h *FollowUpHandler) ListMeetings(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid follow-up ID",
		})
		return
	}

	userID := getUserIDFromContext(c)
	meetings, err := h.service.ListMeetings(c.Request.Context(), uint(id), userID)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "unauthorized: only the assigned member can view meetings" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Meetings retrieved successfully",
		"data":    meetings,
	})
}

func (h *FollowUpHandler) CompleteMeeting(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid meeting ID",
		})
		return
	}

	userID := getUserIDFromContext(c)
	err = h.service.CompleteMeeting(c.Request.Context(), uint(id), userID)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "unauthorized: only the assigned member can mark meetings as complete" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Meeting marked as completed",
	})
}

func (h *FollowUpHandler) SubmitOutcome(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid meeting ID",
		})
		return
	}

	var req struct {
		Status         string  `json:"status" binding:"required"`
		Remarks        string  `json:"remarks" binding:"required"`
		RecordingURL   string  `json:"recording_url" binding:"required"`
		NextFollowUpAt *string `json:"next_follow_up_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"errors":  []gin.H{{"field": "request_body", "error": err.Error()}},
		})
		return
	}

	if len(req.Remarks) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Remarks cannot exceed 1000 characters",
		})
		return
	}

	status := domain.OutcomeStatus(req.Status)
	userID := getUserIDFromContext(c)
	err = h.service.SubmitOutcome(c.Request.Context(), uint(id), userID, status, req.Remarks, req.RecordingURL, req.NextFollowUpAt)
	if err != nil {
		httpStatus := http.StatusBadRequest
		if err.Error() == "unauthorized: only the assigned member can submit outcomes" {
			httpStatus = http.StatusForbidden
		}
		c.JSON(httpStatus, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Outcome submitted successfully",
	})
}

func (h *FollowUpHandler) RestartFollowUp(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid follow-up ID",
		})
		return
	}

	userID := getUserIDFromContext(c)
	err = h.service.RestartFollowUp(c.Request.Context(), uint(id), userID)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "unauthorized: only the assigned member can restart follow-ups" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Follow-up restarted successfully",
	})
}

func getUserIDFromContext(c *gin.Context) string {
	userIDStr, exists := c.Get("userID")
	if !exists {
		return ""
	}

	userID, ok := userIDStr.(string)
	if !ok {
		return ""
	}

	return userID
}
