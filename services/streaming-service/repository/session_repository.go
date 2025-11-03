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

// SessionRepository handles stream session operations
type SessionRepository struct {
	collection *mongo.Collection
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *database.MongoDB) *SessionRepository {
	return &SessionRepository{
		collection: db.Collection("stream_sessions"),
	}
}

// CreateSession creates a new stream session
func (r *SessionRepository) CreateSession(ctx context.Context, session *models.StreamSession) (*models.StreamSession, error) {
	session.ID = primitive.NewObjectID()
	session.CreatedAt = time.Now()
	session.LastHeartbeat = time.Now()

	_, err := r.collection.InsertOne(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession retrieves a session by ID
func (r *SessionRepository) GetSession(ctx context.Context, sessionID string) (*models.StreamSession, error) {
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	var session models.StreamSession
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("session not found")
		}
		return nil, err
	}

	return &session, nil
}

// UpdateSession updates a session
func (r *SessionRepository) UpdateSession(ctx context.Context, sessionID string, session *models.StreamSession) error {
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"position":       session.Position,
			"quality":        session.Quality,
			"bandwidth":      session.Bandwidth,
			"last_heartbeat": time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// UpdateHeartbeat updates session heartbeat
func (r *SessionRepository) UpdateHeartbeat(ctx context.Context, sessionID string) error {
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"last_heartbeat": time.Now()}},
	)
	return err
}

// EndSession ends a session
func (r *SessionRepository) EndSession(ctx context.Context, sessionID string) error {
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}

	now := time.Now()
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"ended_at": now}},
	)
	return err
}

// GetActiveSessionsByUser retrieves active sessions for a user
func (r *SessionRepository) GetActiveSessionsByUser(ctx context.Context, userID string) ([]models.StreamSession, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"user_id": userID,
		"ended_at": nil,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []models.StreamSession
	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

