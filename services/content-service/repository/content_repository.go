package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/streamverse/common-go/database"
	"github.com/streamverse/content-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ContentRepository handles content data operations
type ContentRepository struct {
	collection        *mongo.Collection
	ratingsCollection *mongo.Collection
}

// NewContentRepository creates a new content repository
func NewContentRepository(db *database.MongoDB) *ContentRepository {
	collection := db.Collection("contents")
	ratingsCollection := db.Collection("ratings")

	// Create indexes
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "title", Value: 1}}},
		{Keys: bson.D{{Key: "category", Value: 1}}},
		{Keys: bson.D{{Key: "genre", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "rating", Value: -1}}}, // For trending
	}
	collection.Indexes().CreateMany(context.Background(), indexes)

	// Ratings indexes
	ratingsIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "content_id", Value: 1}, {Key: "user_id", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "content_id", Value: 1}}},
	}
	ratingsCollection.Indexes().CreateMany(context.Background(), ratingsIndexes)

	return &ContentRepository{
		collection:        collection,
		ratingsCollection: ratingsCollection,
	}
}

// Create creates a new content item
func (r *ContentRepository) Create(ctx context.Context, content *models.Content) (*models.Content, error) {
	content.ID = primitive.NewObjectID()
	content.CreatedAt = time.Now()
	content.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, content)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// GetByID retrieves content by ID
func (r *ContentRepository) GetByID(ctx context.Context, id string) (*models.Content, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid content ID: %w", err)
	}

	var content models.Content
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&content)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("content not found")
		}
		return nil, err
	}

	return &content, nil
}

// List retrieves content with pagination and filters
func (r *ContentRepository) List(ctx context.Context, filter bson.M, page, pageSize int) ([]models.Content, int64, error) {
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var contents []models.Content
	if err = cursor.All(ctx, &contents); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return contents, total, nil
}

// Update updates content
func (r *ContentRepository) Update(ctx context.Context, id string, content *models.Content) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid content ID: %w", err)
	}

	content.UpdatedAt = time.Now()
	update := bson.M{
		"$set": content,
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("content not found")
	}

	return nil
}

// Delete soft deletes content
func (r *ContentRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid content ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"status":     "archived",
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("content not found")
	}

	return nil
}

// GetByCategory retrieves content by category
func (r *ContentRepository) GetByCategory(ctx context.Context, category string, page, pageSize int) ([]models.Content, int64, error) {
	filter := bson.M{
		"category": category,
		"status":   "published",
	}
	return r.List(ctx, filter, page, pageSize)
}

// Search searches content by query
func (r *ContentRepository) Search(ctx context.Context, query string, page, pageSize int) ([]models.Content, int64, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
			{"cast": bson.M{"$in": []string{query}}},
		},
		"status": "published",
	}
	return r.List(ctx, filter, page, pageSize)
}

// GetHomeContent retrieves content rows for home screen
func (r *ContentRepository) GetHomeContent(ctx context.Context) ([]models.ContentRow, error) {
	// Get different categories
	categories := []string{"trending", "new-releases", "popular", "featured"}
	rows := make([]models.ContentRow, 0)

	for i, category := range categories {
		contents, _, err := r.GetByCategory(ctx, category, 1, 20)
		if err != nil {
			continue
		}

		rows = append(rows, models.ContentRow{
			ID:    fmt.Sprintf("%d", i+1),
			Title: category,
			Items: contents,
		})
	}

	return rows, nil
}

// GetCategories gets all categories with counts - Issue #13
func (r *ContentRepository) GetCategories(ctx context.Context) ([]models.Category, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"status": "published"}},
		{"$group": bson.M{
			"_id":   "$category",
			"count": bson.M{"$sum": 1},
		}},
		{"$project": bson.M{
			"name":  "$_id",
			"count": 1,
			"_id":   0,
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []models.Category
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetTrending gets trending content - Issue #13
// TODO: Should integrate with Analytics Service for real-time trending based on playback events
func (r *ContentRepository) GetTrending(ctx context.Context, region, deviceType string, limit int) ([]models.Content, error) {
	filter := bson.M{
		"status": "published",
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "rating", Value: -1}, {Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var contents []models.Content
	if err = cursor.All(ctx, &contents); err != nil {
		return nil, err
	}

	return contents, nil
}

// RateContent submits a rating - Issue #13
func (r *ContentRepository) RateContent(ctx context.Context, contentID, userID string, stars int, comment string) error {
	rating := models.Rating{
		ID:        primitive.NewObjectID(),
		ContentID: contentID,
		UserID:    userID,
		Stars:     stars,
		Comment:   comment,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	filter := bson.M{
		"content_id": contentID,
		"user_id":    userID,
	}

	update := bson.M{
		"$set":         bson.M{"stars": stars, "comment": comment, "updated_at": time.Now()},
		"$setOnInsert": rating,
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.ratingsCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	// Update content rating aggregate (simplified - should be done via aggregation)
	// TODO: Use aggregation pipeline to calculate average rating
	return nil
}

// GetRatings gets aggregated ratings - Issue #13
func (r *ContentRepository) GetRatings(ctx context.Context, contentID string) (*models.RatingAggregate, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"content_id": contentID}},
		{"$group": bson.M{
			"_id":          nil,
			"averageStars": bson.M{"$avg": "$stars"},
			"totalRatings": bson.M{"$sum": 1},
			"distribution": bson.M{"$push": "$stars"},
		}},
	}

	cursor, err := r.ratingsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result struct {
		AverageStars float64 `bson:"averageStars"`
		TotalRatings int     `bson:"totalRatings"`
		Distribution []int   `bson:"distribution"`
	}

	if !cursor.Next(ctx) {
		// No ratings yet
		return &models.RatingAggregate{
			ContentID:    contentID,
			AverageStars: 0,
			TotalRatings: 0,
			Distribution: make(map[int]int),
		}, nil
	}

	if err = cursor.Decode(&result); err != nil {
		return nil, err
	}

	// Calculate distribution
	dist := make(map[int]int)
	for _, stars := range result.Distribution {
		dist[stars]++
	}

	return &models.RatingAggregate{
		ContentID:    contentID,
		AverageStars: result.AverageStars,
		TotalRatings: result.TotalRatings,
		Distribution: dist,
	}, nil
}

// GetSimilar gets similar content - Issue #13
func (r *ContentRepository) GetSimilar(ctx context.Context, contentID string, limit int) ([]models.Content, error) {
	// Get the content to find similar items
	content, err := r.GetByID(ctx, contentID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":    bson.M{"$ne": content.ID},
		"genre":  content.Genre,
		"status": "published",
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "rating", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var contents []models.Content
	if err = cursor.All(ctx, &contents); err != nil {
		return nil, err
	}

	return contents, nil
}
