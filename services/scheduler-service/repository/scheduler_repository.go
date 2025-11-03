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
	"github.com/streamverse/scheduler-service/models"
)

// SchedulerRepository handles scheduler data operations
type SchedulerRepository struct {
	channelCollection  *mongo.Collection
	scheduleCollection *mongo.Collection
}

// NewSchedulerRepository creates a new scheduler repository
func NewSchedulerRepository(db *database.MongoDB) *SchedulerRepository {
	channelCollection := db.Collection("channels")
	scheduleCollection := db.Collection("schedule")

	// Create indexes
	channelCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "channel_id", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "type", Value: 1}}},
	})

	scheduleCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "channel_id", Value: 1}, {Key: "start_time", Value: 1}}},
		{Keys: bson.D{{Key: "channel_id", Value: 1}, {Key: "end_time", Value: 1}}},
		{Keys: bson.D{{Key: "content_id", Value: 1}}},
	})

	return &SchedulerRepository{
		channelCollection:  channelCollection,
		scheduleCollection: scheduleCollection,
	}
}

// ListChannels lists all channels
func (r *SchedulerRepository) ListChannels(ctx context.Context, status string) ([]*models.Channel, error) {
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}

	cursor, err := r.channelCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var channels []*models.Channel
	if err = cursor.All(ctx, &channels); err != nil {
		return nil, err
	}

	return channels, nil
}

// GetChannelByID retrieves a channel by ID
func (r *SchedulerRepository) GetChannelByID(ctx context.Context, channelID string) (*models.Channel, error) {
	var channel models.Channel
	err := r.channelCollection.FindOne(ctx, bson.M{"channel_id": channelID}).Decode(&channel)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// CreateChannel creates a new channel
func (r *SchedulerRepository) CreateChannel(ctx context.Context, channel *models.Channel) error {
	channel.CreatedAt = time.Now()
	channel.UpdatedAt = time.Now()
	_, err := r.channelCollection.InsertOne(ctx, channel)
	return err
}

// UpdateChannel updates a channel
func (r *SchedulerRepository) UpdateChannel(ctx context.Context, channelID string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	_, err := r.channelCollection.UpdateOne(
		ctx,
		bson.M{"channel_id": channelID},
		bson.M{"$set": updates},
	)
	return err
}

// GetScheduleEntries retrieves schedule entries for a channel
func (r *SchedulerRepository) GetScheduleEntries(ctx context.Context, channelID string, startTime, endTime *time.Time, limit int) ([]*models.ScheduleEntry, error) {
	filter := bson.M{"channel_id": channelID}

	if startTime != nil && endTime != nil {
		filter["$or"] = []bson.M{
			{"start_time": bson.M{"$gte": startTime, "$lte": endTime}},
			{"end_time": bson.M{"$gte": startTime, "$lte": endTime}},
			{"start_time": bson.M{"$lte": startTime}, "end_time": bson.M{"$gte": endTime}},
		}
	}

	opts := options.Find()
	if limit > 0 {
		limit64 := int64(limit)
		opts.SetLimit(limit64)
	}
	opts.SetSort(bson.D{{Key: "start_time", Value: 1}})

	cursor, err := r.scheduleCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var entries []*models.ScheduleEntry
	if err = cursor.All(ctx, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// CreateScheduleEntry creates a new schedule entry
func (r *SchedulerRepository) CreateScheduleEntry(ctx context.Context, entry *models.ScheduleEntry) error {
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()
	_, err := r.scheduleCollection.InsertOne(ctx, entry)
	return err
}

// UpdateScheduleEntry updates a schedule entry
func (r *SchedulerRepository) UpdateScheduleEntry(ctx context.Context, entryID string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(entryID)
	if err != nil {
		return fmt.Errorf("invalid entry ID: %w", err)
	}

	updates["updated_at"] = time.Now()
	_, err = r.scheduleCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updates},
	)
	return err
}

// DeleteScheduleEntry deletes a schedule entry
func (r *SchedulerRepository) DeleteScheduleEntry(ctx context.Context, entryID string) error {
	objectID, err := primitive.ObjectIDFromHex(entryID)
	if err != nil {
		return fmt.Errorf("invalid entry ID: %w", err)
	}

	_, err = r.scheduleCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// GetCurrentScheduleEntry gets the currently playing schedule entry for a channel
func (r *SchedulerRepository) GetCurrentScheduleEntry(ctx context.Context, channelID string) (*models.ScheduleEntry, error) {
	now := time.Now()
	filter := bson.M{
		"channel_id": channelID,
		"start_time": bson.M{"$lte": now},
		"end_time":   bson.M{"$gte": now},
	}

	var entry models.ScheduleEntry
	err := r.scheduleCollection.FindOne(ctx, filter).Decode(&entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

