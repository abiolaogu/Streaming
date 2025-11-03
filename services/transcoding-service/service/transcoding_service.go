package service

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/streamverse/transcoding-service/models"
	"github.com/streamverse/transcoding-service/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TranscodingService handles transcoding business logic
type TranscodingService struct {
	repo *repository.TranscodingRepository
}

// NewTranscodingService creates a new transcoding service
func NewTranscodingService(repo *repository.TranscodingRepository) *TranscodingService {
	return &TranscodingService{
		repo: repo,
	}
}

// CreateJob creates a new transcoding job
func (s *TranscodingService) CreateJob(ctx context.Context, contentID, inputURL string, qualityLevels []string, priority int) (*models.TranscodingJob, error) {
	job := &models.TranscodingJob{
		ID:            primitive.NewObjectID(),
		ContentID:     contentID,
		InputURL:      inputURL,
		Status:        "pending",
		Progress:      0,
		Priority:      priority,
		QualityLevels: qualityLevels,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return s.repo.CreateJob(ctx, job)
}

// GetJob retrieves a job by ID
func (s *TranscodingService) GetJob(ctx context.Context, jobID string) (*models.TranscodingJob, error) {
	return s.repo.GetJob(ctx, jobID)
}

// UpdateJobProgress updates job progress
func (s *TranscodingService) UpdateJobProgress(ctx context.Context, jobID string, progress float64) error {
	return s.repo.UpdateProgress(ctx, jobID, progress)
}

// CompleteJob marks job as completed
func (s *TranscodingService) CompleteJob(ctx context.Context, jobID, outputURL string) error {
	now := time.Now()
	return s.repo.UpdateStatus(ctx, jobID, "completed", &now, outputURL)
}

// FailJob marks job as failed
func (s *TranscodingService) FailJob(ctx context.Context, jobID, errorMsg string) error {
	return s.repo.FailJob(ctx, jobID, errorMsg)
}

// CreateThumbnailJob creates a thumbnail generation job
func (s *TranscodingService) CreateThumbnailJob(ctx context.Context, contentID, videoURL string) (*models.ThumbnailJob, error) {
	job := &models.ThumbnailJob{
		ID:        primitive.NewObjectID(),
		ContentID: contentID,
		VideoURL:  videoURL,
		Status:    "pending",
		Progress:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repo.CreateThumbnailJob(ctx, job)
}

// ListJobs lists transcoding jobs with filters - Issue #15
func (s *TranscodingService) ListJobs(ctx context.Context, status string, page, pageSize int) ([]*models.TranscodingJob, int64, error) {
	return s.repo.ListJobs(ctx, status, page, pageSize)
}

// ListProfiles lists available transcoding profiles - Issue #15
func (s *TranscodingService) ListProfiles(ctx context.Context) ([]map[string]interface{}, error) {
	// Default profiles as per Issue #15
	profiles := []map[string]interface{}{
		{"name": "baseline", "codec": "h264", "bitrate": 800000, "resolution": "360p", "fps": 24},
		{"name": "low", "codec": "h264", "bitrate": 2500000, "resolution": "480p", "fps": 24},
		{"name": "medium", "codec": "h264", "bitrate": 5000000, "resolution": "720p", "fps": 30},
		{"name": "high", "codec": "h265", "bitrate": 8000000, "resolution": "1080p", "fps": 30},
		{"name": "uhd", "codec": "h265", "bitrate": 15000000, "resolution": "2160p", "fps": 60},
		{"name": "hdr", "codec": "h265", "bitrate": 20000000, "resolution": "2160p", "fps": 60, "hdr": true},
	}
	// TODO: Load from database if custom profiles exist
	return profiles, nil
}

// CreateProfile creates a new transcoding profile - Issue #15
func (s *TranscodingService) CreateProfile(ctx context.Context, name, codec string, bitrate int, resolution string, fps int) (map[string]interface{}, error) {
	// TODO: Store profile in database
	profile := map[string]interface{}{
		"name":       name,
		"codec":      codec,
		"bitrate":    bitrate,
		"resolution": resolution,
		"fps":        fps,
	}
	return profile, nil
}

// UploadPart defines a part of a multipart upload - Issue #29
type UploadPart struct {
	ETag       string
	PartNumber int
}

// InitiateUpload initiates a multipart upload - Issue #29
func (s *TranscodingService) InitiateUpload(ctx context.Context, fileName string, fileSize int64) (string, error) {
	return s.repo.InitiateUpload(ctx, fileName, fileSize)
}

// UploadPart uploads a part of a multipart upload - Issue #29
func (s *TranscodingService) UploadPart(ctx context.Context, uploadID string, partNumber int, file *multipart.FileHeader) (string, error) {
	return s.repo.UploadPart(ctx, uploadID, partNumber, file)
}

// CompleteUpload completes a multipart upload - Issue #29
func (s *TranscodingService) CompleteUpload(ctx context.Context, uploadID string, parts []UploadPart) (string, error) {
	completedParts := make([]repository.UploadPart, len(parts))
	for i, p := range parts {
		completedParts[i] = repository.UploadPart(p)
	}
	return s.repo.CompleteUpload(ctx, uploadID, completedParts)
}
