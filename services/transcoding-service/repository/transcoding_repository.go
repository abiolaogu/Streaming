package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/streamverse/common-go/database"
	"github.com/streamverse/transcoding-service/models"
)

// TranscodingRepository handles transcoding job operations
type TranscodingRepository struct {
	jobCollection       *mongo.Collection
	thumbnailCollection *mongo.Collection
}

// NewTranscodingRepository creates a new transcoding repository
func NewTranscodingRepository(db *database.MongoDB) *TranscodingRepository {
	return &TranscodingRepository{
		jobCollection:       db.Collection("transcoding_jobs"),
		thumbnailCollection: db.Collection("thumbnail_jobs"),
	}
}

// CreateJob creates a new transcoding job
func (r *TranscodingRepository) CreateJob(ctx context.Context, job *models.TranscodingJob) (*models.TranscodingJob, error) {
	_, err := r.jobCollection.InsertOne(ctx, job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

// GetJob retrieves a job by ID
func (r *TranscodingRepository) GetJob(ctx context.Context, jobID string) (*models.TranscodingJob, error) {
	objectID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return nil, fmt.Errorf("invalid job ID: %w", err)
	}

	var job models.TranscodingJob
	err = r.jobCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// UpdateProgress updates job progress
func (r *TranscodingRepository) UpdateProgress(ctx context.Context, jobID string, progress float64) error {
	objectID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return fmt.Errorf("invalid job ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"progress":   progress,
			"updated_at": time.Now(),
		},
	}

	_, err = r.jobCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// UpdateStatus updates job status
func (r *TranscodingRepository) UpdateStatus(ctx context.Context, jobID, status string, completedAt *time.Time, outputURL string) error {
	objectID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return fmt.Errorf("invalid job ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	if completedAt != nil {
		update["$set"].(bson.M)["completed_at"] = completedAt
	}

	if outputURL != "" {
		update["$set"].(bson.M)["output_url"] = outputURL
	}

	_, err = r.jobCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// FailJob marks job as failed
func (r *TranscodingRepository) FailJob(ctx context.Context, jobID, errorMsg string) error {
	objectID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return fmt.Errorf("invalid job ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"status":     "failed",
			"error":      errorMsg,
			"updated_at": time.Now(),
		},
	}

	_, err = r.jobCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// CreateThumbnailJob creates a thumbnail job
func (r *TranscodingRepository) CreateThumbnailJob(ctx context.Context, job *models.ThumbnailJob) (*models.ThumbnailJob, error) {
	_, err := r.thumbnailCollection.InsertOne(ctx, job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

// ListJobs lists transcoding jobs with filters - Issue #15
func (r *TranscodingRepository) ListJobs(ctx context.Context, status string, page, pageSize int) ([]*models.TranscodingJob, int64, error) {
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}

	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	opts := mongo.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort:  bson.D{{Key: "created_at", Value: -1}},
	}

	cursor, err := r.jobCollection.Find(ctx, filter, &opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var jobs []*models.TranscodingJob
	if err = cursor.All(ctx, &jobs); err != nil {
		return nil, 0, err
	}

	total, err := r.jobCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

