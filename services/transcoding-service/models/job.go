package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TranscodingJob represents a transcoding job
type TranscodingJob struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ContentID   string             `bson:"content_id" json:"contentId"`
	InputURL    string             `bson:"input_url" json:"inputUrl"`
	OutputURL   string             `bson:"output_url,omitempty" json:"outputUrl,omitempty"`
	Status      string             `bson:"status" json:"status"` // "pending", "processing", "completed", "failed"
	Progress    float64            `bson:"progress" json:"progress"` // 0-100
	Priority    int                `bson:"priority" json:"priority"` // 1-10
	QualityLevels []string          `bson:"quality_levels" json:"qualityLevels"` // ["1080p", "720p", "480p"]
	Error       string             `bson:"error,omitempty" json:"error,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
	CompletedAt *time.Time         `bson:"completed_at,omitempty" json:"completedAt,omitempty"`
}

// ThumbnailJob represents a thumbnail generation job
type ThumbnailJob struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ContentID   string             `bson:"content_id" json:"contentId"`
	VideoURL    string             `bson:"video_url" json:"videoUrl"`
	OutputURL   string             `bson:"output_url,omitempty" json:"outputUrl,omitempty"`
	Status      string             `bson:"status" json:"status"`
	Progress    float64            `bson:"progress" json:"progress"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}
