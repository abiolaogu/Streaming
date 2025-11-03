package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streamverse/ad-compositing-service/service"
	"github.com/streamverse/common-go/logger"
)

// AdCompositingHandler handles HTTP requests for ad compositing
type AdCompositingHandler struct {
	service *service.AdCompositingService
	logger  *logger.Logger
}

// NewAdCompositingHandler creates a new ad compositing handler
func NewAdCompositingHandler(service *service.AdCompositingService, logger *logger.Logger) *AdCompositingHandler {
	return &AdCompositingHandler{
		service: service,
		logger:  logger,
	}
}

// CompositeAds handles POST /ad-compositing/composite
func (h *AdCompositingHandler) CompositeAds(c *gin.Context) {
	var req struct {
		VideoID     string                 `json:"video_id" binding:"required"`
		UserProfile map[string]interface{} `json:"user_profile" binding:"required"`
		SceneData   map[string]interface{} `json:"scene_data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	compositedVideo, err := h.service.CompositeAds(c.Request.Context(), req.VideoID, req.UserProfile, req.SceneData)
	if err != nil {
		h.logger.Error("Failed to composite ads", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to composite ads"})
		return
	}

	c.JSON(http.StatusOK, compositedVideo)
}
