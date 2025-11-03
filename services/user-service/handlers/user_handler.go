package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/streamverse/common-go/errors"
	"github.com/streamverse/common-go/logger"
	"github.com/streamverse/user-service/models"
	"github.com/streamverse/user-service/service"
)

// UserHandler handles HTTP requests for user profiles
type UserHandler struct {
	service *service.UserService
	logger  *logger.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(service *service.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// GetProfile handles GET /users/{id}
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.Param("id")
	currentUserID, _ := c.Get("user_id")

	// Users can only access their own profile unless admin
	if userID != currentUserID.(string) {
		roles, _ := c.Get("roles")
		rolesList, ok := roles.([]string)
		if !ok || !contains(rolesList, "admin") {
			c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Access denied"))
			return
		}
	}

	profile, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get profile", logger.Error(err))
		c.JSON(http.StatusNotFound, errors.NewNotFoundError("User not found"))
		return
	}

	c.JSON(http.StatusOK, profile)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// UpdateProfile handles PUT /users/{id}
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.Param("id")
	currentUserID, _ := c.Get("user_id")

	if userID != currentUserID.(string) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Access denied"))
		return
	}

	var profile models.UserProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	if err := h.service.UpdateProfile(c.Request.Context(), userID, &profile); err != nil {
		h.logger.Error("Failed to update profile", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to update profile"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// GetPreferences handles GET /users/{id}/preferences
func (h *UserHandler) GetPreferences(c *gin.Context) {
	userID := c.Param("id")
	currentUserID, _ := c.Get("user_id")

	if userID != currentUserID.(string) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Access denied"))
		return
	}

	prefs, err := h.service.GetPreferences(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get preferences", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to get preferences"))
		return
	}

	c.JSON(http.StatusOK, prefs)
}

// UpdatePreferences handles PUT /users/{id}/preferences
func (h *UserHandler) UpdatePreferences(c *gin.Context) {
	userID := c.Param("id")
	currentUserID, _ := c.Get("user_id")

	if userID != currentUserID.(string) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Access denied"))
		return
	}

	var prefs models.UserPreferences
	if err := c.ShouldBindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	if err := h.service.UpdatePreferences(c.Request.Context(), userID, &prefs); err != nil {
		h.logger.Error("Failed to update preferences", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to update preferences"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Preferences updated successfully"})
}

// CreateProfile handles POST /api/v1/users/me/profiles
func (h *UserHandler) CreateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("User not authenticated"))
		return
	}

	var profile models.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	created, err := h.service.CreateProfile(c.Request.Context(), userID.(string), &profile)
	if err != nil {
		h.logger.Error("Failed to create profile", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to create profile"))
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetProfiles handles GET /api/v1/users/me/profiles
func (h *UserHandler) GetProfiles(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("User not authenticated"))
		return
	}

	profiles, err := h.service.GetProfiles(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error("Failed to get profiles", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to get profiles"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"profiles": profiles})
}

// UpdateWatchHistory handles PUT /api/v1/users/me/watch-history
func (h *UserHandler) UpdateWatchHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("User not authenticated"))
		return
	}

	var history models.WatchHistory
	if err := c.ShouldBindJSON(&history); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	history.UserID = userID.(string)
	if err := h.service.UpdateWatchHistory(c.Request.Context(), &history); err != nil {
		h.logger.Error("Failed to update watch history", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to update watch history"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Watch history updated"})
}

// GetWatchHistory handles GET /users/{id}/watch-history
func (h *UserHandler) GetWatchHistory(c *gin.Context) {
	userID := c.Param("id")
	currentUserID, _ := c.Get("user_id")

	if userID != currentUserID.(string) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Access denied"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	if pageSize > 1000 {
		pageSize = 1000 // Max 1000 entries as per requirement
	}

	history, err := h.service.GetWatchHistory(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		h.logger.Error("Failed to get watch history", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to get watch history"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}

// GetContinueWatching handles GET /api/v1/users/me/continue-watching
func (h *UserHandler) GetContinueWatching(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("User not authenticated"))
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	history, err := h.service.GetContinueWatching(c.Request.Context(), userID.(string), limit)
	if err != nil {
		h.logger.Error("Failed to get continue watching", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to get continue watching"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": history})
}

// ClearWatchHistory handles DELETE /api/v1/users/me/watch-history
func (h *UserHandler) ClearWatchHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("User not authenticated"))
		return
	}

	if err := h.service.ClearWatchHistory(c.Request.Context(), userID.(string)); err != nil {
		h.logger.Error("Failed to clear watch history", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to clear watch history"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Watch history cleared"})
}

// AddToWatchlist handles POST /users/{id}/watchlist
func (h *UserHandler) AddToWatchlist(c *gin.Context) {
	userID := c.Param("id")
	currentUserID, _ := c.Get("user_id")

	if userID != currentUserID.(string) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Access denied"))
		return
	}

	var req struct {
		ContentID string `json:"content_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	if err := h.service.AddToWatchlist(c.Request.Context(), userID, req.ContentID); err != nil {
		h.logger.Error("Failed to add to watchlist", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to add to watchlist"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Added to watchlist"})
}

// RemoveFromWatchlist handles DELETE /users/{id}/watchlist/{content_id}
func (h *UserHandler) RemoveFromWatchlist(c *gin.Context) {
	userID := c.Param("id")
	currentUserID, _ := c.Get("user_id")

	if userID != currentUserID.(string) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Access denied"))
		return
	}

	contentID := c.Param("content_id")
	if err := h.service.RemoveFromWatchlist(c.Request.Context(), userID, contentID); err != nil {
		h.logger.Error("Failed to remove from watchlist", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to remove from watchlist"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed from watchlist"})
}

// GetWatchlist handles GET /users/{id}/watchlist
func (h *UserHandler) GetWatchlist(c *gin.Context) {
	userID := c.Param("id")
	currentUserID, _ := c.Get("user_id")

	if userID != currentUserID.(string) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Access denied"))
		return
	}

	watchlist, err := h.service.GetWatchlist(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get watchlist", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to get watchlist"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"watchlist": watchlist})
}

// GetDevices handles GET /users/{id}/devices
func (h *UserHandler) GetDevices(c *gin.Context) {
	userID := c.Param("id")
	currentUserID, _ := c.Get("user_id")

	if userID != currentUserID.(string) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Access denied"))
		return
	}

	devices, err := h.service.GetDevices(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get devices", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to get devices"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"devices": devices})
}

// RegisterDevice handles POST /users/{id}/devices
func (h *UserHandler) RegisterDevice(c *gin.Context) {
	userID := c.Param("id")
	currentUserID, _ := c.Get("user_id")

	if userID != currentUserID.(string) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Access denied"))
		return
	}

	var device models.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	device.UserID = userID
	if err := h.service.RegisterDevice(c.Request.Context(), &device); err != nil {
		h.logger.Error("Failed to register device", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to register device"))
		return
	}

	c.JSON(http.StatusCreated, device)
}

// DeregisterDevice handles DELETE /users/{id}/devices/{device_id}
func (h *UserHandler) DeregisterDevice(c *gin.Context) {
	userID := c.Param("id")
	deviceID := c.Param("device_id")
	currentUserID, _ := c.Get("user_id")

	if userID != currentUserID.(string) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Access denied"))
		return
	}

	if err := h.service.DeregisterDevice(c.Request.Context(), userID, deviceID); err != nil {
		h.logger.Error("Failed to deregister device", logger.Error(err))
		c.JSON(http.StatusNotFound, errors.NewNotFoundError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device deregistered"})
}

// ExportUserData handles GET /users/{id}/export (GDPR)
func (h *UserHandler) ExportUserData(c *gin.Context) {
	userID := c.Param("id")
	currentUserID, _ := c.Get("user_id")

	if userID != currentUserID.(string) {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Access denied"))
		return
	}

	data, err := h.service.ExportUserData(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to export user data", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to export data"))
		return
	}

	c.JSON(http.StatusOK, data)
}

