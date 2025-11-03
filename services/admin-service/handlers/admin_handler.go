package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/streamverse/common-go/errors"
	"github.com/streamverse/common-go/logger"
	"github.com/streamverse/admin-service/models"
	"github.com/streamverse/admin-service/service"
)

// AdminHandler handles HTTP requests for admin operations
type AdminHandler struct {
	service *service.AdminService
	logger  *logger.Logger
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(service *service.AdminService, logger *logger.Logger) *AdminHandler {
	return &AdminHandler{
		service: service,
		logger:  logger,
	}
}

// checkRole checks if user has required role
func (h *AdminHandler) checkRole(c *gin.Context, requiredRoles []string) bool {
	roles, exists := c.Get("roles")
	if !exists {
		return false
	}

	rolesList, ok := roles.([]string)
	if !ok {
		return false
	}

	for _, required := range requiredRoles {
		for _, userRole := range rolesList {
			if userRole == required {
				return true
			}
		}
	}
	return false
}

// logAudit logs an audit event
func (h *AdminHandler) logAudit(c *gin.Context, action, resource, resourceID string, changes map[string]interface{}) {
	userID, _ := c.Get("user_id")
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	go func() {
		_ = h.service.LogAuditEvent(
			c.Request.Context(),
			userID.(string),
			action,
			resource,
			resourceID,
			changes,
			ipAddress,
			userAgent,
		)
	}()
}

// ListUsers handles GET /admin/users - Issue #21
func (h *AdminHandler) ListUsers(c *gin.Context) {
	if !h.checkRole(c, []string{"superadmin", "admin"}) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Insufficient permissions"))
		return
	}

	var filters models.UserListFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		// Optional filters, continue with empty filters
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	users, total, err := h.service.ListUsers(c.Request.Context(), &filters, page, pageSize)
	if err != nil {
		h.logger.Error("Failed to list users", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to list users"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":     users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetUser handles GET /admin/users/{id} - Issue #21
func (h *AdminHandler) GetUser(c *gin.Context) {
	if !h.checkRole(c, []string{"superadmin", "admin"}) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Insufficient permissions"))
		return
	}

	userID := c.Param("id")
	user, err := h.service.GetUser(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user", logger.Error(err))
		c.JSON(http.StatusNotFound, errors.NewNotFoundError("User not found"))
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser handles PUT /admin/users/{id} - Issue #21
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	if !h.checkRole(c, []string{"superadmin", "admin"}) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Insufficient permissions"))
		return
	}

	userID := c.Param("id")
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	oldUser, _ := h.service.GetUser(c.Request.Context(), userID)
	if err := h.service.UpdateUser(c.Request.Context(), userID, updates); err != nil {
		h.logger.Error("Failed to update user", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to update user"))
		return
	}

	h.logAudit(c, "update", "user", userID, map[string]interface{}{"old": oldUser, "new": updates})
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser handles DELETE /admin/users/{id} - Issue #21
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	if !h.checkRole(c, []string{"superadmin"}) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Only superadmin can delete users"))
		return
	}

	userID := c.Param("id")
	if err := h.service.DeleteUser(c.Request.Context(), userID); err != nil {
		h.logger.Error("Failed to delete user", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to delete user"))
		return
	}

	h.logAudit(c, "delete", "user", userID, nil)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ListContent handles GET /admin/content - Issue #21
func (h *AdminHandler) ListContent(c *gin.Context) {
	if !h.checkRole(c, []string{"superadmin", "admin", "editor"}) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Insufficient permissions"))
		return
	}

	var filters models.ContentListFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		// Optional filters
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	content, total, err := h.service.ListContent(c.Request.Context(), &filters, page, pageSize)
	if err != nil {
		h.logger.Error("Failed to list content", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to list content"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"content":   content,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// BulkImportContent handles POST /admin/content - Issue #21
func (h *AdminHandler) BulkImportContent(c *gin.Context) {
	if !h.checkRole(c, []string{"superadmin", "admin"}) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Insufficient permissions"))
		return
	}

	format := c.DefaultQuery("format", "json") // "csv" or "json"
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError("File is required"))
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError("Failed to open file"))
		return
	}
	defer f.Close()

	result, err := h.service.BulkImportContent(c.Request.Context(), f, format)
	if err != nil {
		h.logger.Error("Failed to import content", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to import content"))
		return
	}

	h.logAudit(c, "bulk_import", "content", "", map[string]interface{}{"format": format, "result": result})
	c.JSON(http.StatusOK, result)
}

// UpdateContent handles PUT /admin/content/{id} - Issue #21
func (h *AdminHandler) UpdateContent(c *gin.Context) {
	if !h.checkRole(c, []string{"superadmin", "admin", "editor"}) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Insufficient permissions"))
		return
	}

	contentID := c.Param("id")
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	if err := h.service.UpdateContent(c.Request.Context(), contentID, updates); err != nil {
		h.logger.Error("Failed to update content", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to update content"))
		return
	}

	h.logAudit(c, "update", "content", contentID, updates)
	c.JSON(http.StatusOK, gin.H{"message": "Content updated successfully"})
}

// DeleteContent handles DELETE /admin/content/{id} - Issue #21
func (h *AdminHandler) DeleteContent(c *gin.Context) {
	if !h.checkRole(c, []string{"superadmin", "admin"}) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Insufficient permissions"))
		return
	}

	contentID := c.Param("id")
	if err := h.service.DeleteContent(c.Request.Context(), contentID); err != nil {
		h.logger.Error("Failed to delete content", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to delete content"))
		return
	}

	h.logAudit(c, "delete", "content", contentID, nil)
	c.JSON(http.StatusOK, gin.H{"message": "Content deleted successfully"})
}

// GetDashboardMetrics handles GET /admin/analytics - Issue #21
func (h *AdminHandler) GetDashboardMetrics(c *gin.Context) {
	if !h.checkRole(c, []string{"superadmin", "admin"}) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Insufficient permissions"))
		return
	}

	// TODO: Fetch metrics from Analytics Service
	c.JSON(http.StatusOK, gin.H{
		"concurrent_viewers": 1500,
		"total_users":        10000,
		"total_content":      5000,
		"video_starts_today": 5000,
		"error_rate":         0.01,
	})
}

// GetSystemSettings handles GET /admin/settings - Issue #21
func (h *AdminHandler) GetSystemSettings(c *gin.Context) {
	if !h.checkRole(c, []string{"superadmin", "admin"}) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Insufficient permissions"))
		return
	}

	settings, err := h.service.GetSystemSettings(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get settings", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to get settings"))
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateSystemSettings handles PUT /admin/settings - Issue #21
func (h *AdminHandler) UpdateSystemSettings(c *gin.Context) {
	if !h.checkRole(c, []string{"superadmin"}) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Only superadmin can update settings"))
		return
	}

	var settings models.SystemSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	userID, _ := c.Get("user_id")
	if err := h.service.UpdateSystemSettings(c.Request.Context(), &settings, userID.(string)); err != nil {
		h.logger.Error("Failed to update settings", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to update settings"))
		return
	}

	h.logAudit(c, "update", "settings", "", map[string]interface{}{"settings": settings})
	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}

// GetAuditLogs handles GET /admin/audit-logs - Issue #21
func (h *AdminHandler) GetAuditLogs(c *gin.Context) {
	if !h.checkRole(c, []string{"superadmin", "admin"}) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Insufficient permissions"))
		return
	}

	filters := make(map[string]interface{})
	if resource := c.Query("resource"); resource != "" {
		filters["resource"] = resource
	}
	if action := c.Query("action"); action != "" {
		filters["action"] = action
	}
	if userID := c.Query("user_id"); userID != "" {
		filters["user_id"] = userID
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	logs, total, err := h.service.GetAuditLogs(c.Request.Context(), filters, page, pageSize)
	if err != nil {
		h.logger.Error("Failed to get audit logs", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to get audit logs"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

