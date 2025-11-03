package repository

import (
	"context"

	"github.com/streamverse/common-go/database"
	"go.mongodb.org/mongo-driver/mongo"
)

// AdCompositingRepository handles ad compositing data operations
type AdCompositingRepository struct {
	collection *mongo.Collection
}

// NewAdCompositingRepository creates a new ad compositing repository
func NewAdCompositingRepository(db *database.MongoDB) *AdCompositingRepository {
	return &AdCompositingRepository{
		collection: db.Collection("composited_videos"),
	}
}

// SaveCompositedVideo saves the composited video data
func (r *AdCompositingRepository) SaveCompositedVideo(ctx context.Context, videoData map[string]interface{}) (map[string]interface{}, error) {
	_, err := r.collection.InsertOne(ctx, videoData)
	if err != nil {
		return nil, err
	}

	return videoData, nil
}
