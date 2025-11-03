package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/streamverse/common-go/errors"
	"github.com/streamverse/common-go/logger"
	"github.com/streamverse/transcoding-service/service"
)

// TranscodingHandler handles HTTP requests for transcoding
type TranscodingHandler struct {
	service *service.TranscodingService
	logger  *logger.Logger
}

// NewTranscodingHandler creates a new transcoding handler
func NewTranscodingHandler(service *service.TranscodingService, logger *logger.Logger) *TranscodingHandler {
	return &TranscodingHandler{
		service: service,
		logger:  logger,
	}
}

// SubmitTranscodeJob handles POST /transcode/jobs - Issue #15
func (h *TranscodingHandler) SubmitTranscodeJob(c *gin.Context) {
	var req struct {
		ContentID    string   `json:"content_id" binding:"required"`
		InputURL     string   `json:"input_url" binding:"required"`
		QualityLevels []string `json:"quality_levels"`
		Priority     int      `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	if req.QualityLevels == nil {
		req.QualityLevels = []string{"1080p", "720p", "480p"}
	}
	if req.Priority == 0 {
		req.Priority = 5
	}

	job, err := h.service.CreateJob(c.Request.Context(), req.ContentID, req.InputURL, req.QualityLevels, req.Priority)
	if err != nil {
		h.logger.Error("Failed to create job", logger.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, job)
}

// GetTranscodeJobStatus handles GET /transcode/jobs/{job_id} - Issue #15
func (h *TranscodingHandler) GetTranscodeJobStatus(c *gin.Context) {
	jobID := c.Param("job_id")

	job, err := h.service.GetJob(c.Request.Context(), jobID)
	if err != nil {
		h.logger.Error("Failed to get job", logger.Error(err))
		c.JSON(http.StatusNotFound, errors.NewNotFoundError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, job)
}

// ListTranscodeJobs handles GET /transcode/jobs - Issue #15
func (h *TranscodingHandler) ListTranscodeJobs(c *gin.Context) {
	status := c.Query("status") // filter by status: queued, processing, done, failed
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	jobs, total, err := h.service.ListJobs(c.Request.Context(), status, page, pageSize)
	if err != nil {
		h.logger.Error("Failed to list jobs", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to list jobs"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs":      jobs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ListProfiles handles GET /transcode/profiles - Issue #15
func (h *TranscodingHandler) ListProfiles(c *gin.Context) {
	profiles, err := h.service.ListProfiles(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to list profiles", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to list profiles"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"profiles": profiles})
}

// CreateProfile handles POST /transcode/profiles - Issue #15
func (h *TranscodingHandler) CreateProfile(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		Codec     string `json:"codec" binding:"required"` // "h264" or "h265"
		Bitrate   int    `json:"bitrate" binding:"required"`
		Resolution string `json:"resolution" binding:"required"` // "360p", "480p", "720p", "1080p", "2160p"
		FPS       int    `json:"fps"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	if req.FPS == 0 {
		req.FPS = 24
	}

	profile, err := h.service.CreateProfile(c.Request.Context(), req.Name, req.Codec, req.Bitrate, req.Resolution, req.FPS)
	if err != nil {
		h.logger.Error("Failed to create profile", logger.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, profile)
}

