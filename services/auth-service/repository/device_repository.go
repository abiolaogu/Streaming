package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/streamverse/common-go/database"
	"github.com/streamverse/auth-service/models"
)

// DeviceRepository handles device operations
type DeviceRepository struct {
	collection *mongo.Collection
}

// NewDeviceRepository creates a new device repository
func NewDeviceRepository(db *database.MongoDB) *DeviceRepository {
	return &DeviceRepository{
		collection: db.Collection("devices"),
	}
}

// CreateDevice creates a new device
func (r *DeviceRepository) CreateDevice(ctx context.Context, device *models.Device) error {
	_, err := r.collection.InsertOne(ctx, device)
	return err
}

// GetDevicesByUserID retrieves all devices for a user
func (r *DeviceRepository) GetDevicesByUserID(ctx context.Context, userID string) ([]models.Device, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var devices []models.Device
	if err = cursor.All(ctx, &devices); err != nil {
		return nil, err
	}

	return devices, nil
}

// DeleteDevice deletes a device
func (r *DeviceRepository) DeleteDevice(ctx context.Context, deviceID, userID string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"id": deviceID, "user_id": userID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("device not found")
	}

	return nil
}

// UpdateDeviceLastUsed updates device last used timestamp
func (r *DeviceRepository) UpdateDeviceLastUsed(ctx context.Context, deviceID string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"id": deviceID},
		bson.M{"$set": bson.M{"last_used_at": time.Now()}},
	)
	return err
}

