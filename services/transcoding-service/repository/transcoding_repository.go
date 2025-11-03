package repository

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/streamverse/common-go/database"
	"github.com/streamverse/transcoding-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TranscodingRepository handles transcoding job operations
type TranscodingRepository struct {
	jobCollection       *mongo.Collection
	thumbnailCollection *mongo.Collection
	uploadCollection    *mongo.Collection
	s3Client            *s3.S3
}

// NewTranscodingRepository creates a new transcoding repository
func NewTranscodingRepository(db *database.MongoDB, s3Client *s3.S3) *TranscodingRepository {
	return &TranscodingRepository{
		jobCollection:       db.Collection("transcoding_jobs"),
		thumbnailCollection: db.Collection("thumbnail_jobs"),
		uploadCollection:    db.Collection("multipart_uploads"),
		s3Client:            s3Client,
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

// UploadPart defines a part of a multipart upload for the service layer.
type UploadPart struct {
	ETag       string
	PartNumber int
}

// multipartUpload represents a multipart upload in the database.
type multipartUpload struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UploadID   string             `bson:"upload_id"`
	FileName   string             `bson:"file_name"`
	FileSize   int64              `bson:"file_size"`
	Parts      []s3.Part          `bson:"parts"`
	CreatedAt  time.Time          `bson:"created_at"`
}

// InitiateUpload initiates a multipart upload.
func (r *TranscodingRepository) InitiateUpload(ctx context.Context, fileName string, fileSize int64) (string, error) {
	// Mock S3 CreateMultipartUpload
	uploadID := uuid.New().String()
	// Create a record in the database
	upload := &multipartUpload{
		UploadID:  uploadID,
		FileName:  fileName,
		FileSize:  fileSize,
		Parts:     []s3.Part{},
		CreatedAt: time.Now(),
	}
	_, err := r.uploadCollection.InsertOne(ctx, upload)
	if err != nil {
		return "", err
	}
	return uploadID, nil
}

// UploadPart uploads a part of a multipart upload.
func (r *TranscodingRepository) UploadPart(ctx context.Context, uploadID string, partNumber int, file *multipart.FileHeader) (string, error) {
	// Mock S3 UploadPart
	etag := uuid.New().String()
	part := s3.Part{
		ETag:       &etag,
		PartNumber: aws.Int64(int64(partNumber)),
	}

	filter := bson.M{"upload_id": uploadID}
	update := bson.M{"$push": bson.M{"parts": part}}
	_, err := r.uploadCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	return etag, nil
}

// CompleteUpload completes a multipart upload.
func (r *TranscodingRepository) CompleteUpload(ctx context.Context, uploadID string, parts []UploadPart) (string, error) {
	// Find the upload record
	var upload multipartUpload
	err := r.uploadCollection.FindOne(ctx, bson.M{"upload_id": uploadID}).Decode(&upload)
	if err != nil {
		return "", err
	}

	// Mock S3 CompleteMultipartUpload
	location := fmt.Sprintf("s3://mock-bucket/%s", upload.FileName)

	// Here you would typically delete the multipart upload record from the database
	// _, err = r.uploadCollection.DeleteOne(ctx, bson.M{"upload_id": uploadID})
	// if err != nil {
	// 	return "", err
	// }

	return location, nil
}
