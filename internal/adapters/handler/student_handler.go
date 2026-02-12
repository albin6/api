package handler

import (
	"github.com/albin6/api/internal/core/domain"
	"github.com/albin6/api/internal/core/port"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type StudentHandler struct {
	service port.StudentService
}

func NewStudentHandler(service port.StudentService) *StudentHandler {
	return &StudentHandler{service: service}
}

func (h *StudentHandler) GetStudents(c *gin.Context) {
	search := c.Query("search")
	statusStr := c.Query("program_status")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	var status *bool
	if statusStr != "" {
		val, err := strconv.ParseBool(statusStr)
		if err == nil {
			status = &val
		}
	}

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	result, err := h.service.GetStudents(c.Request.Context(), search, status, sortBy, order, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var student domain.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateStudent(c.Request.Context(), &student); err != nil {
		if strings.Contains(err.Error(), "required") || strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, student)
}

func (h *StudentHandler) SearchStudents(c *gin.Context) {
	query := c.Query("q")
	limitStr := c.DefaultQuery("limit", "10")
	
	limit, _ := strconv.Atoi(limitStr)
	
	students, err := h.service.SearchStudents(c.Request.Context(), query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"students": students})
}
