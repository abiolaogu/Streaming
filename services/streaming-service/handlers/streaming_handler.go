package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streamverse/common-go/errors"
	"github.com/streamverse/common-go/logger"
	"github.com/streamverse/streaming-service/models"
	"github.com/streamverse/streaming-service/service"
)

// StreamingHandler handles HTTP requests for streaming
type StreamingHandler struct {
	service *service.StreamingService
	logger  *logger.Logger
}

// NewStreamingHandler creates a new streaming handler
func NewStreamingHandler(service *service.StreamingService, logger *logger.Logger) *StreamingHandler {
	return &StreamingHandler{
		service: service,
		logger:  logger,
	}
}

// GetHLSManifest handles GET /streaming/manifest/:content_id/:token.m3u8 - Issue #14
func (h *StreamingHandler) GetHLSManifest(c *gin.Context) {
	contentID := c.Param("content_id")
	token := c.Param("token")

	// Validate token
	userID, err := h.service.ValidateToken(c.Request.Context(), token)
	if err != nil {
		h.logger.Error("Invalid token", logger.Error(err))
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("Invalid token"))
		return
	}

	// Generate HLS manifest
	manifest, err := h.service.GenerateHLSManifest(c.Request.Context(), contentID, userID)
	if err != nil {
		h.logger.Error("Failed to get manifest", logger.Error(err))
		c.JSON(http.StatusNotFound, errors.NewNotFoundError(err.Error()))
		return
	}

	// Return HLS manifest as text
	c.Header("Content-Type", "application/vnd.apple.mpegurl")
	c.String(http.StatusOK, manifest)
}

// GetDASHManifest handles GET /streaming/manifest/:content_id/:token.mpd - Issue #14
func (h *StreamingHandler) GetDASHManifest(c *gin.Context) {
	contentID := c.Param("content_id")
	token := c.Param("token")

	// Validate token
	userID, err := h.service.ValidateToken(c.Request.Context(), token)
	if err != nil {
		h.logger.Error("Invalid token", logger.Error(err))
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("Invalid token"))
		return
	}

	// Generate DASH manifest
	manifest, err := h.service.GenerateDASHManifest(c.Request.Context(), contentID, userID)
	if err != nil {
		h.logger.Error("Failed to get manifest", logger.Error(err))
		c.JSON(http.StatusNotFound, errors.NewNotFoundError(err.Error()))
		return
	}

	// Return DASH manifest as XML
	c.Header("Content-Type", "application/dash+xml")
	c.String(http.StatusOK, manifest)
}

// GenerateToken handles POST /streaming/token - Issue #14
func (h *StreamingHandler) GenerateToken(c *gin.Context) {
	var req struct {
		ContentID string `json:"content_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	userID, _ := c.Get("user_id")
	deviceID := c.GetHeader("X-Device-ID")
	ip := c.ClientIP()

	token, err := h.service.GenerateToken(c.Request.Context(), req.ContentID, userID.(string), ip, deviceID)
	if err != nil {
		h.logger.Error("Failed to generate token", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to generate token"))
		return
	}

	c.JSON(http.StatusOK, token)
}

// SubmitQoE handles POST /streaming/qoe - Issue #14
func (h *StreamingHandler) SubmitQoE(c *gin.Context) {
	var event models.QoEEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	userID, _ := c.Get("user_id")
	event.UserID = userID.(string)
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	if err := h.service.SubmitQoE(c.Request.Context(), &event); err != nil {
		h.logger.Error("Failed to submit QoE", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to submit QoE"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "QoE event recorded"})
}

// GetManifest handles GET /api/v1/streaming/:contentId/manifest (deprecated, kept for backward compatibility)
func (h *StreamingHandler) GetManifest(c *gin.Context) {
	contentID := c.Param("contentId")
	format := c.DefaultQuery("format", "hls")
	userID, _ := c.Get("user_id")

	manifest, err := h.service.GetManifest(c.Request.Context(), contentID, format, userID.(string))
	if err != nil {
		h.logger.Error("Failed to get manifest", logger.Error(err))
		c.JSON(http.StatusNotFound, errors.NewNotFoundError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, manifest)
}

// CreateSession handles POST /api/v1/streaming/sessions
func (h *StreamingHandler) CreateSession(c *gin.Context) {
	var req models.StreamingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	userID, _ := c.Get("user_id")
	deviceID := c.GetHeader("X-Device-ID")

	session, err := h.service.CreateSession(c.Request.Context(), userID.(string), req.ContentID, deviceID)
	if err != nil {
		h.logger.Error("Failed to create session", logger.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, session)
}

// UpdatePosition handles PUT /api/v1/streaming/sessions/:sessionId/position
func (h *StreamingHandler) UpdatePosition(c *gin.Context) {
	sessionID := c.Param("sessionId")

	var req struct {
		Position int64 `json:"position" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	if err := h.service.UpdateSessionPosition(c.Request.Context(), sessionID, req.Position); err != nil {
		h.logger.Error("Failed to update position", logger.Error(err))
		c.JSON(http.StatusNotFound, errors.NewNotFoundError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Position updated"})
}

// Heartbeat handles POST /api/v1/streaming/sessions/:sessionId/heartbeat
func (h *StreamingHandler) Heartbeat(c *gin.Context) {
	sessionID := c.Param("sessionId")

	if err := h.service.SendHeartbeat(c.Request.Context(), sessionID); err != nil {
		h.logger.Error("Failed to send heartbeat", logger.Error(err))
		c.JSON(http.StatusNotFound, errors.NewNotFoundError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Heartbeat received"})
}

// EndSession handles DELETE /api/v1/streaming/sessions/:sessionId
func (h *StreamingHandler) EndSession(c *gin.Context) {
	sessionID := c.Param("sessionId")

	if err := h.service.EndSession(c.Request.Context(), sessionID); err != nil {
		h.logger.Error("Failed to end session", logger.Error(err))
		c.JSON(http.StatusNotFound, errors.NewNotFoundError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session ended"})
}
