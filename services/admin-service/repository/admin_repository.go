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
	"github.com/streamverse/admin-service/models"
)

// AdminRepository handles admin data operations
type AdminRepository struct {
	auditCollection   *mongo.Collection
	settingsCollection *mongo.Collection
	usersCollection   *mongo.Collection
	contentCollection *mongo.Collection
}

// NewAdminRepository creates a new admin repository
func NewAdminRepository(db *database.MongoDB) *AdminRepository {
	auditCollection := db.Collection("audit_logs")
	settingsCollection := db.Collection("system_settings")
	usersCollection := db.Collection("users")
	contentCollection := db.Collection("contents")

	// Create indexes
	auditCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "resource", Value: 1}, {Key: "resource_id", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	})

	settingsCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	return &AdminRepository{
		auditCollection:    auditCollection,
		settingsCollection: settingsCollection,
		usersCollection:    usersCollection,
		contentCollection:  contentCollection,
	}
}

// CreateAuditLog creates an audit log entry
func (r *AdminRepository) CreateAuditLog(ctx context.Context, log *models.AuditLog) error {
	_, err := r.auditCollection.InsertOne(ctx, log)
	return err
}

// GetAuditLogs retrieves audit logs with filters and pagination
func (r *AdminRepository) GetAuditLogs(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*models.AuditLog, int64, error) {
	filter := bson.M{}
	if userID, ok := filters["user_id"]; ok {
		filter["user_id"] = userID
	}
	if resource, ok := filters["resource"]; ok {
		filter["resource"] = resource
	}
	if action, ok := filters["action"]; ok {
		filter["action"] = action
	}

	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	opts := options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort:  bson.D{{Key: "created_at", Value: -1}},
	}

	cursor, err := r.auditCollection.Find(ctx, filter, &opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var logs []*models.AuditLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, 0, err
	}

	total, err := r.auditCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetSystemSettings retrieves system settings
func (r *AdminRepository) GetSystemSettings(ctx context.Context) (*models.SystemSettings, error) {
	var settings models.SystemSettings
	err := r.settingsCollection.FindOne(ctx, bson.M{}).Decode(&settings)
	if err == mongo.ErrNoDocuments {
		// Return default settings
		return &models.SystemSettings{
			ID:              primitive.NewObjectID(),
			FeatureFlags:    make(map[string]bool),
			MaxUploadSize:   5 * 1024 * 1024 * 1024, // 5GB
			MaintenanceMode: false,
			CDNBaseURL:      "https://cdn.streamverse.io",
			UpdatedAt:       time.Now(),
		}, nil
	}
	return &settings, err
}

// UpdateSystemSettings updates system settings
func (r *AdminRepository) UpdateSystemSettings(ctx context.Context, settings *models.SystemSettings) error {
	settings.UpdatedAt = time.Now()
	opts := options.Update().SetUpsert(true)
	_, err := r.settingsCollection.UpdateOne(
		ctx,
		bson.M{},
		bson.M{"$set": settings},
		opts,
	)
	return err
}

// ListUsers lists users with filters and pagination
func (r *AdminRepository) ListUsers(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]map[string]interface{}, int64, error) {
	filter := bson.M{}
	if status, ok := filters["status"]; ok {
		filter["status"] = status
	}
	if role, ok := filters["role"]; ok {
		if roleStr, isString := role.(string); isString && roleStr != "" {
			filter["roles"] = bson.M{"$in": []string{roleStr}}
		}
	}
	if email, ok := filters["email"]; ok {
		filter["email"] = bson.M{"$regex": email, "$options": "i"}
	}

	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	opts := options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort:  bson.D{{Key: "created_at", Value: -1}},
	}

	cursor, err := r.usersCollection.Find(ctx, filter, &opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []map[string]interface{}
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	total, err := r.usersCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUser updates a user
func (r *AdminRepository) UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	updates["updated_at"] = time.Now()
	_, err = r.usersCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updates},
	)
	return err
}

// SoftDeleteUser soft deletes a user (sets status to "deleted")
func (r *AdminRepository) SoftDeleteUser(ctx context.Context, userID string) error {
	return r.UpdateUser(ctx, userID, map[string]interface{}{
		"status": "deleted",
		"deleted_at": time.Now(),
	})
}

// ListContent lists content with filters and pagination
func (r *AdminRepository) ListContent(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]map[string]interface{}, int64, error) {
	filter := bson.M{}
	if status, ok := filters["status"]; ok {
		filter["status"] = status
	}
	if category, ok := filters["category"]; ok {
		filter["category"] = category
	}
	if genre, ok := filters["genre"]; ok {
		filter["genre"] = genre
	}

	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	opts := options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort:  bson.D{{Key: "created_at", Value: -1}},
	}

	cursor, err := r.contentCollection.Find(ctx, filter, &opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var content []map[string]interface{}
	if err = cursor.All(ctx, &content); err != nil {
		return nil, 0, err
	}

	total, err := r.contentCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return content, total, nil
}

// UpdateContent updates content metadata
func (r *AdminRepository) UpdateContent(ctx context.Context, contentID string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(contentID)
	if err != nil {
		return fmt.Errorf("invalid content ID: %w", err)
	}

	updates["updated_at"] = time.Now()
	_, err = r.contentCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updates},
	)
	return err
}

// DeleteContent deletes content
func (r *AdminRepository) DeleteContent(ctx context.Context, contentID string) error {
	objectID, err := primitive.ObjectIDFromHex(contentID)
	if err != nil {
		return fmt.Errorf("invalid content ID: %w", err)
	}

	_, err = r.contentCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
