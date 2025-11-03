package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streamverse/common-go/errors"
	"github.com/streamverse/common-go/logger"
	"github.com/streamverse/ad-service/models"
	"github.com/streamverse/ad-service/service"
)

// AdHandler handles HTTP requests for ads
type AdHandler struct {
	service *service.AdService
	logger  *logger.Logger
}

// NewAdHandler creates a new ad handler
func NewAdHandler(service *service.AdService, logger *logger.Logger) *AdHandler {
	return &AdHandler{
		service: service,
		logger:  logger,
	}
}

// GetAds handles POST /api/v1/ads/request
func (h *AdHandler) GetAds(c *gin.Context) {
	var req models.AdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	userID, _ := c.Get("user_id")
	req.UserID = userID.(string)

	response, err := h.service.GetAds(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to get ads", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to get ads"))
		return
	}

	c.JSON(http.StatusOK, response)
}

// TrackAdEvent handles POST /api/v1/ads/track
func (h *AdHandler) TrackAdEvent(c *gin.Context) {
	var tracking models.AdTracking
	if err := c.ShouldBindJSON(&tracking); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	userID, _ := c.Get("user_id")
	tracking.UserID = userID.(string)

	if err := h.service.TrackAdEvent(c.Request.Context(), &tracking); err != nil {
		h.logger.Error("Failed to track ad event", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to track event"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event tracked"})
}

