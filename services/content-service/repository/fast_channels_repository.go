package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/streamverse/common-go/database"
	"github.com/streamverse/content-service/models"
)

// FASTChannelsRepository handles FAST channels
type FASTChannelsRepository struct {
	collection *mongo.Collection
}

// NewFASTChannelsRepository creates a new FAST channels repository
func NewFASTChannelsRepository(db *database.MongoDB) *FASTChannelsRepository {
	return &FASTChannelsRepository{
		collection: db.Collection("fast_channels"),
	}
}

// Create creates a new FAST channel
func (r *FASTChannelsRepository) Create(ctx context.Context, channel *models.FASTChannel) (*models.FASTChannel, error) {
	channel.ID = primitive.NewObjectID()
	channel.CreatedAt = time.Now()
	channel.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, channel)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

// GetByID retrieves channel by ID
func (r *FASTChannelsRepository) GetByID(ctx context.Context, id string) (*models.FASTChannel, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid channel ID: %w", err)
	}

	var channel models.FASTChannel
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&channel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("channel not found")
		}
		return nil, err
	}

	return &channel, nil
}

// UpdateSchedule updates channel schedule
func (r *FASTChannelsRepository) UpdateSchedule(ctx context.Context, id string, schedule []models.ScheduleItem) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid channel ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"schedule":  schedule,
			"updated_at": time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

