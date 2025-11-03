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

// CollectionsRepository handles collections/playlists
type CollectionsRepository struct {
	collection *mongo.Collection
}

// NewCollectionsRepository creates a new collections repository
func NewCollectionsRepository(db *database.MongoDB) *CollectionsRepository {
	return &CollectionsRepository{
		collection: db.Collection("collections"),
	}
}

// Create creates a new collection
func (r *CollectionsRepository) Create(ctx context.Context, collection *models.Collection) (*models.Collection, error) {
	collection.ID = primitive.NewObjectID()
	collection.CreatedAt = time.Now()
	collection.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, collection)
	if err != nil {
		return nil, err
	}

	return collection, nil
}

// GetByID retrieves collection by ID
func (r *CollectionsRepository) GetByID(ctx context.Context, id string) (*models.Collection, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid collection ID: %w", err)
	}

	var collection models.Collection
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&collection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("collection not found")
		}
		return nil, err
	}

	return &collection, nil
}

// AddContent adds content to collection
func (r *CollectionsRepository) AddContent(ctx context.Context, collectionID, contentID string) error {
	objectID, err := primitive.ObjectIDFromHex(collectionID)
	if err != nil {
		return fmt.Errorf("invalid collection ID: %w", err)
	}

	update := bson.M{
		"$addToSet": bson.M{"content_ids": contentID},
		"$set":      bson.M{"updated_at": time.Now()},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// RemoveContent removes content from collection
func (r *CollectionsRepository) RemoveContent(ctx context.Context, collectionID, contentID string) error {
	objectID, err := primitive.ObjectIDFromHex(collectionID)
	if err != nil {
		return fmt.Errorf("invalid collection ID: %w", err)
	}

	update := bson.M{
		"$pull": bson.M{"content_ids": contentID},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

