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
	"github.com/streamverse/auth-service/models"
)

// UserRepository handles user data operations
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.MongoDB) *UserRepository {
	collection := db.Collection("users")
	
	// Create indexes
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	}
	collection.Indexes().CreateMany(context.Background(), indexes)

	return &UserRepository{
		collection: collection,
	}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.EmailVerified = false
	user.MFAEnabled = false
	user.FailedLoginAttempts = 0
	user.Roles = []string{"user"}

	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, fmt.Errorf("email already exists")
		}
		return nil, err
	}

	return user, nil
}

// GetByEmail retrieves user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetByID retrieves user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// Update updates user
func (r *UserRepository) Update(ctx context.Context, id string, user *models.User) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	user.UpdatedAt = time.Now()
	update := bson.M{"$set": user}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// IncrementFailedLoginAttempts increments failed login attempts
func (r *UserRepository) IncrementFailedLoginAttempts(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	update := bson.M{
		"$inc": bson.M{"failed_login_attempts": 1},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// ResetFailedLoginAttempts resets failed login attempts
func (r *UserRepository) ResetFailedLoginAttempts(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"failed_login_attempts": 0,
			"account_locked_until":   nil,
			"updated_at":            time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// LockAccount locks the account until the specified time
func (r *UserRepository) LockAccount(ctx context.Context, id string, until time.Time) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"account_locked_until": until,
			"updated_at":          time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

