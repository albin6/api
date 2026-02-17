package handler

import (
	"net/http"

	"github.com/albin6/api/internal/core/port"
	"github.com/gin-gonic/gin"
)

type ToolHandler struct {
	toolService port.ToolService
}

func NewToolHandler(toolService port.ToolService) *ToolHandler {
	return &ToolHandler{
		toolService: toolService,
	}
}

func (h *ToolHandler) GetStudents(c *gin.Context) {
	// Extract all query parameters
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	resp, err := h.toolService.GetStudents(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ToolHandler) GetPageFilters(c *gin.Context) {
	pageName := c.Query("pageName")
	if pageName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pageName is required"})
		return
	}

	resp, err := h.toolService.GetPageFilters(c.Request.Context(), pageName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ToolHandler) GetBatches(c *gin.Context) {
	resp, err := h.toolService.GetBatches(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ToolHandler) GetCourses(c *gin.Context) {
	resp, err := h.toolService.GetCourses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ToolHandler) GetDomains(c *gin.Context) {
	resp, err := h.toolService.GetDomains(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ToolHandler) GetEmployees(c *gin.Context) {
	roles := c.Query("roles")
	if roles == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "roles is required"})
		return
	}

	resp, err := h.toolService.GetEmployees(c.Request.Context(), roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ToolHandler) GetStatusOptions(c *gin.Context) {
	category := c.Query("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category is required"})
		return
	}

	resp, err := h.toolService.GetStatusOptions(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
