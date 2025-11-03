package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/streamverse/common-go/database"
	"github.com/streamverse/streaming-service/models"
)

// StreamingRepository handles playback session operations
type StreamingRepository struct {
	collection *mongo.Collection
}

// NewStreamingRepository creates a new streaming repository
func NewStreamingRepository(db *database.MongoDB) *StreamingRepository {
	return &StreamingRepository{
		collection: db.Collection("playback_sessions"),
	}
}

// CreateSession creates a new playback session
func (r *StreamingRepository) CreateSession(ctx context.Context, session *models.PlaybackSession) (*models.PlaybackSession, error) {
	_, err := r.collection.InsertOne(ctx, session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// GetSession retrieves a session by ID
func (r *StreamingRepository) GetSession(ctx context.Context, sessionID string) (*models.PlaybackSession, error) {
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	var session models.PlaybackSession
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// UpdatePosition updates playback position
func (r *StreamingRepository) UpdatePosition(ctx context.Context, sessionID string, position int64) error {
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"position":   position,
			"updated_at": time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// UpdateHeartbeat updates session heartbeat
func (r *StreamingRepository) UpdateHeartbeat(ctx context.Context, sessionID string) error {
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"last_heartbeat": time.Now(),
			"updated_at":     time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// DeleteSession deletes a session
func (r *StreamingRepository) DeleteSession(ctx context.Context, sessionID string) error {
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// GetActiveSessions retrieves active sessions for a user
func (r *StreamingRepository) GetActiveSessions(ctx context.Context, userID string) ([]models.PlaybackSession, error) {
	filter := bson.M{
		"user_id":        userID,
		"last_heartbeat": bson.M{"$gte": time.Now().Add(-5 * time.Minute)},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []models.PlaybackSession
	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

