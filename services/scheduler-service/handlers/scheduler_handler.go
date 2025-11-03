package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/streamverse/common-go/errors"
	"github.com/streamverse/common-go/logger"
	"github.com/streamverse/scheduler-service/models"
	"github.com/streamverse/scheduler-service/service"
)

// SchedulerHandler handles HTTP requests for scheduler operations
type SchedulerHandler struct {
	service *service.SchedulerService
	logger  *logger.Logger
}

// NewSchedulerHandler creates a new scheduler handler
func NewSchedulerHandler(service *service.SchedulerService, logger *logger.Logger) *SchedulerHandler {
	return &SchedulerHandler{
		service: service,
		logger:  logger,
	}
}

// ListChannels handles GET /scheduler/channels - Issue #27
func (h *SchedulerHandler) ListChannels(c *gin.Context) {
	status := c.Query("status") // "active" or "inactive"

	channels, err := h.service.ListChannels(c.Request.Context(), status)
	if err != nil {
		h.logger.Error("Failed to list channels", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to list channels"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"channels": channels})
}

// GetChannelEPG handles GET /scheduler/channels/{channel_id}/epg - Issue #27
func (h *SchedulerHandler) GetChannelEPG(c *gin.Context) {
	channelID := c.Param("channel_id")

	epg, err := h.service.GetChannelEPG(c.Request.Context(), channelID)
	if err != nil {
		h.logger.Error("Failed to get EPG", logger.Error(err))
		c.JSON(http.StatusNotFound, errors.NewNotFoundError("Channel not found or EPG unavailable"))
		return
	}

	c.JSON(http.StatusOK, epg)
}

// GetChannelManifest handles GET /scheduler/channels/{channel_id}/manifest - Issue #27
func (h *SchedulerHandler) GetChannelManifest(c *gin.Context) {
	channelID := c.Param("channel_id")

	manifest, err := h.service.GetChannelManifest(c.Request.Context(), channelID)
	if err != nil {
		h.logger.Error("Failed to get manifest", logger.Error(err))
		c.JSON(http.StatusNotFound, errors.NewNotFoundError("Channel not found"))
		return
	}

	c.JSON(http.StatusOK, manifest)
}

// CreateScheduleEntry handles POST /scheduler/schedule - Issue #27
func (h *SchedulerHandler) CreateScheduleEntry(c *gin.Context) {
	var entry models.ScheduleEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	if err := h.service.CreateScheduleEntry(c.Request.Context(), &entry); err != nil {
		h.logger.Error("Failed to create schedule entry", logger.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, entry)
}

// UpdateScheduleEntry handles PUT /scheduler/schedule/{id} - Issue #27
func (h *SchedulerHandler) UpdateScheduleEntry(c *gin.Context) {
	entryID := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	if err := h.service.UpdateScheduleEntry(c.Request.Context(), entryID, updates); err != nil {
		h.logger.Error("Failed to update schedule entry", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to update schedule entry"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule entry updated successfully"})
}

// DeleteScheduleEntry handles DELETE /scheduler/schedule/{id} - Issue #27
func (h *SchedulerHandler) DeleteScheduleEntry(c *gin.Context) {
	entryID := c.Param("id")

	if err := h.service.DeleteScheduleEntry(c.Request.Context(), entryID); err != nil {
		h.logger.Error("Failed to delete schedule entry", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to delete schedule entry"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule entry deleted successfully"})
}

// GetCurrentScheduleEntry handles GET /scheduler/channels/{channel_id}/now - Issue #27
func (h *SchedulerHandler) GetCurrentScheduleEntry(c *gin.Context) {
	channelID := c.Param("channel_id")

	entry, err := h.service.GetCurrentScheduleEntry(c.Request.Context(), channelID)
	if err != nil {
		h.logger.Error("Failed to get current schedule", logger.Error(err))
		c.JSON(http.StatusNotFound, errors.NewNotFoundError("No current schedule entry"))
		return
	}

	c.JSON(http.StatusOK, entry)
}

