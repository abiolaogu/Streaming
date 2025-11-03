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
	"github.com/streamverse/user-service/models"
)

// UserRepository handles user profile data operations
type UserRepository struct {
	profileCollection     *mongo.Collection
	preferencesCollection *mongo.Collection
	profilesCollection    *mongo.Collection
	historyCollection     *mongo.Collection
	watchlistCollection   *mongo.Collection
	devicesCollection     *mongo.Collection
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.MongoDB) *UserRepository {
	profileCollection := db.Collection("user_profiles")
	preferencesCollection := db.Collection("user_preferences")
	profilesCollection := db.Collection("profiles")
	historyCollection := db.Collection("watch_history")
	watchlistCollection := db.Collection("watchlist")
	devicesCollection := db.Collection("devices")

	// Create indexes
	profileCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	historyCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "watched_at", Value: -1}}},
		{Keys: bson.D{{Key: "content_id", Value: 1}}},
	})

	devicesCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{{Key: "device_id", Value: 1}, {Key: "user_id", Value: 1}}},
	})

	return &UserRepository{
		profileCollection:     profileCollection,
		preferencesCollection: preferencesCollection,
		profilesCollection:    profilesCollection,
		historyCollection:     historyCollection,
		watchlistCollection:   watchlistCollection,
		devicesCollection:     devicesCollection,
	}
}

// GetProfile retrieves user profile
func (r *UserRepository) GetProfile(ctx context.Context, userID string) (*models.UserProfile, error) {
	var profile models.UserProfile
	err := r.profileCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Create default profile
			profile = models.UserProfile{
				ID:        primitive.NewObjectID(),
				UserID:    userID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			r.profileCollection.InsertOne(ctx, profile)
			return &profile, nil
		}
		return nil, err
	}
	return &profile, nil
}

// UpdateProfile updates user profile
func (r *UserRepository) UpdateProfile(ctx context.Context, userID string, profile *models.UserProfile) error {
	profile.UpdatedAt = time.Now()
	update := bson.M{"$set": profile}

	opts := options.Update().SetUpsert(true)
	_, err := r.profileCollection.UpdateOne(
		ctx,
		bson.M{"user_id": userID},
		update,
		opts,
	)
	return err
}

// GetPreferences retrieves user preferences
func (r *UserRepository) GetPreferences(ctx context.Context, userID string) (*models.UserPreferences, error) {
	var prefs models.UserPreferences
	err := r.preferencesCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&prefs)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Create default preferences
			prefs = models.UserPreferences{
				ID:            primitive.NewObjectID(),
				UserID:        userID,
				Language:      "en",
				ContentRating: "PG-13",
				Notifications: models.NotificationPrefs{
					Email:     true,
					Push:      true,
					Marketing: false,
				},
				Playback: models.PlaybackPrefs{
					Quality:  "auto",
					Autoplay: true,
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			r.preferencesCollection.InsertOne(ctx, prefs)
			return &prefs, nil
		}
		return nil, err
	}
	return &prefs, nil
}

// UpdatePreferences updates user preferences
func (r *UserRepository) UpdatePreferences(ctx context.Context, userID string, prefs *models.UserPreferences) error {
	prefs.UpdatedAt = time.Now()
	update := bson.M{"$set": prefs}

	opts := options.Update().SetUpsert(true)
	_, err := r.preferencesCollection.UpdateOne(
		ctx,
		bson.M{"user_id": userID},
		update,
		opts,
	)
	return err
}

// CreateProfile creates a sub-profile
func (r *UserRepository) CreateProfile(ctx context.Context, profile *models.Profile) (*models.Profile, error) {
	profile.ID = primitive.NewObjectID()
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	_, err := r.profilesCollection.InsertOne(ctx, profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// GetProfiles retrieves all profiles for a user
func (r *UserRepository) GetProfiles(ctx context.Context, userID string) ([]models.Profile, error) {
	cursor, err := r.profilesCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var profiles []models.Profile
	if err = cursor.All(ctx, &profiles); err != nil {
		return nil, err
	}

	return profiles, nil
}

// GetProfileByID retrieves a profile by ID
func (r *UserRepository) GetProfileByID(ctx context.Context, profileID string) (*models.Profile, error) {
	objectID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return nil, fmt.Errorf("invalid profile ID: %w", err)
	}

	var profile models.Profile
	err = r.profilesCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("profile not found")
		}
		return nil, err
	}

	return &profile, nil
}

// UpdateWatchHistory updates watch history
func (r *UserRepository) UpdateWatchHistory(ctx context.Context, history *models.WatchHistory) error {
	history.UpdatedAt = time.Now()
	filter := bson.M{
		"user_id":    history.UserID,
		"content_id": history.ContentID,
	}

	if history.ProfileID != "" {
		filter["profile_id"] = history.ProfileID
	}

	update := bson.M{
		"$set": history,
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.historyCollection.UpdateOne(ctx, filter, update, opts)
	return err
}

// GetWatchHistory retrieves watch history with pagination (Issue #12: max 1000 entries per request)
func (r *UserRepository) GetWatchHistory(ctx context.Context, userID string, page, pageSize int) ([]models.WatchHistory, error) {
	if pageSize > 1000 {
		pageSize = 1000 // Max 1000 as per requirement
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.D{{Key: "watched_at", Value: -1}})

	cursor, err := r.historyCollection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var history []models.WatchHistory
	if err = cursor.All(ctx, &history); err != nil {
		return nil, err
	}

	return history, nil
}

// GetContinueWatching retrieves continue watching items
func (r *UserRepository) GetContinueWatching(ctx context.Context, userID string, limit int) ([]models.WatchHistory, error) {
	filter := bson.M{
		"user_id":  userID,
		"completed": false,
	}
	opts := options.Find().SetLimit(int64(limit)).SetSort(bson.D{{Key: "updated_at", Value: -1}})
	cursor, err := r.historyCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var history []models.WatchHistory
	if err = cursor.All(ctx, &history); err != nil {
		return nil, err
	}

	return history, nil
}

// ClearWatchHistory clears watch history
func (r *UserRepository) ClearWatchHistory(ctx context.Context, userID string) error {
	_, err := r.historyCollection.DeleteMany(ctx, bson.M{"user_id": userID})
	return err
}

// AddToWatchlist adds content to watchlist
func (r *UserRepository) AddToWatchlist(ctx context.Context, userID, contentID string) error {
	watchlist := models.Watchlist{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		ContentID: contentID,
		AddedAt:   time.Now(),
	}

	filter := bson.M{
		"user_id":    userID,
		"content_id": contentID,
	}

	update := bson.M{"$setOnInsert": watchlist}
	opts := options.Update().SetUpsert(true)

	_, err := r.watchlistCollection.UpdateOne(ctx, filter, update, opts)
	return err
}

// RemoveFromWatchlist removes content from watchlist
func (r *UserRepository) RemoveFromWatchlist(ctx context.Context, userID, contentID string) error {
	_, err := r.watchlistCollection.DeleteOne(ctx, bson.M{"user_id": userID, "content_id": contentID})
	return err
}

// GetWatchlist retrieves watchlist
func (r *UserRepository) GetWatchlist(ctx context.Context, userID string) ([]models.Watchlist, error) {
	cursor, err := r.watchlistCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var watchlist []models.Watchlist
	if err = cursor.All(ctx, &watchlist); err != nil {
		return nil, err
	}

	return watchlist, nil
}

// GetDevices retrieves all devices for a user
func (r *UserRepository) GetDevices(ctx context.Context, userID string) ([]models.Device, error) {
	cursor, err := r.devicesCollection.Find(ctx, bson.M{"user_id": userID})
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

// RegisterDevice registers a new device
func (r *UserRepository) RegisterDevice(ctx context.Context, device *models.Device) error {
	device.ID = primitive.NewObjectID()
	device.CreatedAt = time.Now()
	device.LastUsedAt = time.Now()

	filter := bson.M{
		"user_id":   device.UserID,
		"device_id": device.DeviceID,
	}

	update := bson.M{
		"$set":         bson.M{"last_used_at": time.Now()},
		"$setOnInsert": device,
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.devicesCollection.UpdateOne(ctx, filter, update, opts)
	return err
}

// DeregisterDevice removes a device
func (r *UserRepository) DeregisterDevice(ctx context.Context, userID, deviceID string) error {
	result, err := r.devicesCollection.DeleteOne(ctx, bson.M{
		"user_id":   userID,
		"device_id": deviceID,
	})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("device not found")
	}
	return nil
}

// ExportUserData exports all user data for GDPR compliance
func (r *UserRepository) ExportUserData(ctx context.Context, userID string) (*models.UserDataExport, error) {
	export := &models.UserDataExport{
		UserID:     userID,
		ExportedAt: time.Now(),
	}

	// Get profile
	profile, err := r.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}
	export.Profile = profile

	// Get preferences
	prefs, err := r.GetPreferences(ctx, userID)
	if err != nil {
		return nil, err
	}
	export.Preferences = prefs

	// Get profiles
	profiles, err := r.GetProfiles(ctx, userID)
	if err != nil {
		return nil, err
	}
	export.Profiles = profiles

	// Get watch history (all, no pagination for export)
	cursor, err := r.historyCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &export.WatchHistory); err != nil {
		return nil, err
	}

	// Get watchlist
	watchlist, err := r.GetWatchlist(ctx, userID)
	if err != nil {
		return nil, err
	}
	export.Watchlist = watchlist

	// Get devices
	devices, err := r.GetDevices(ctx, userID)
	if err != nil {
		return nil, err
	}
	export.Devices = devices

	return export, nil
}

// DeleteUser deletes all user data (GDPR compliance)
func (r *UserRepository) DeleteUser(ctx context.Context, userID string) error {
	// Delete all user-related data
	r.profileCollection.DeleteMany(ctx, bson.M{"user_id": userID})
	r.preferencesCollection.DeleteMany(ctx, bson.M{"user_id": userID})
	r.profilesCollection.DeleteMany(ctx, bson.M{"user_id": userID})
	r.devicesCollection.DeleteMany(ctx, bson.M{"user_id": userID})
	r.historyCollection.DeleteMany(ctx, bson.M{"user_id": userID})
	r.watchlistCollection.DeleteMany(ctx, bson.M{"user_id": userID})
	return nil
}

