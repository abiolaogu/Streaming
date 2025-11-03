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

// InitiateUpload handles POST /transcode/uploads - Issue #29
func (h *TranscodingHandler) InitiateUpload(c *gin.Context) {
	var req struct {
		FileName string `json:"file_name" binding:"required"`
		FileSize int64  `json:"file_size" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	uploadID, err := h.service.InitiateUpload(c.Request.Context(), req.FileName, req.FileSize)
	if err != nil {
		h.logger.Error("Failed to initiate upload", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to initiate upload"))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"upload_id": uploadID})
}

// UploadPart handles POST /transcode/uploads/:upload_id/parts - Issue #29
func (h *TranscodingHandler) UploadPart(c *gin.Context) {
	uploadID := c.Param("upload_id")
	partNumber, err := strconv.Atoi(c.Query("part_number"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError("Invalid part number"))
		return
	}

	file, err := c.FormFile("part")
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError("Missing part"))
		return
	}

	etag, err := h.service.UploadPart(c.Request.Context(), uploadID, partNumber, file)
	if err != nil {
		h.logger.Error("Failed to upload part", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to upload part"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"etag": etag})
}

// CompleteUpload handles POST /transcode/uploads/:upload_id/complete - Issue #29
func (h *TranscodingHandler) CompleteUpload(c *gin.Context) {
	uploadID := c.Param("upload_id")

	var req struct {
		Parts []struct {
			ETag       string `json:"etag" binding:"required"`
			PartNumber int    `json:"part_number" binding:"required"`
		} `json:"parts" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	parts := make([]service.UploadPart, len(req.Parts))
	for i, p := range req.Parts {
		parts[i] = service.UploadPart{ETag: p.ETag, PartNumber: p.PartNumber}
	}

	location, err := h.service.CompleteUpload(c.Request.Context(), uploadID, parts)
	if err != nil {
		h.logger.Error("Failed to complete upload", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to complete upload"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"location": location})
}
