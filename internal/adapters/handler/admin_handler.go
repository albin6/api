package handler

import (
	"github.com/albin6/api/internal/core/domain"
	"github.com/albin6/api/internal/core/port"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminHandler struct {
	adminService port.AdminService
}

func NewAdminHandler(adminService port.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var admin domain.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.adminService.CreateAdmin(c.Request.Context(), &admin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Admin created successfully"})
}
