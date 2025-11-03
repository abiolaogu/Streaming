package service

import (
	"context"
	"fmt"

	"github.com/streamverse/content-service/models"
	"github.com/streamverse/content-service/repository"
)

// ContentService handles content business logic
type ContentService struct {
	repo *repository.ContentRepository
}

// NewContentService creates a new content service
func NewContentService(repo *repository.ContentRepository) *ContentService {
	return &ContentService{
		repo: repo,
	}
}

// GetContentByID retrieves content by ID
func (s *ContentService) GetContentByID(ctx context.Context, id string) (*models.Content, error) {
	return s.repo.GetByID(ctx, id)
}

// ListContent lists content with filters
func (s *ContentService) ListContent(ctx context.Context, category string, page, pageSize int) ([]models.Content, int64, error) {
	if category != "" {
		return s.repo.GetByCategory(ctx, category, page, pageSize)
	}
	return s.repo.List(ctx, map[string]interface{}{"status": "published"}, page, pageSize)
}

// CreateContent creates new content
func (s *ContentService) CreateContent(ctx context.Context, content *models.Content) (*models.Content, error) {
	// Validation
	if content.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if content.Category == "" {
		return nil, fmt.Errorf("category is required")
	}

	// Set default status if not set
	if content.Status == "" {
		content.Status = "draft"
	}

	return s.repo.Create(ctx, content)
}

// UpdateContent updates content
func (s *ContentService) UpdateContent(ctx context.Context, id string, content *models.Content) error {
	// Verify content exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, id, content)
}

// DeleteContent deletes content
func (s *ContentService) DeleteContent(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// SearchContent searches content
func (s *ContentService) SearchContent(ctx context.Context, query string, page, pageSize int) ([]models.Content, int64, error) {
	return s.repo.Search(ctx, query, page, pageSize)
}

// GetHomeContent gets home screen content rows
func (s *ContentService) GetHomeContent(ctx context.Context) ([]models.ContentRow, error) {
	return s.repo.GetHomeContent(ctx)
}

// GetCategories gets all categories with counts - Issue #13
func (s *ContentService) GetCategories(ctx context.Context) ([]models.Category, error) {
	return s.repo.GetCategories(ctx)
}

// GetTrending gets trending content by region/device - Issue #13
func (s *ContentService) GetTrending(ctx context.Context, region, deviceType string, limit int) ([]models.Content, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.GetTrending(ctx, region, deviceType, limit)
}

// RateContent submits a rating for content - Issue #13
func (s *ContentService) RateContent(ctx context.Context, contentID, userID string, stars int, comment string) error {
	return s.repo.RateContent(ctx, contentID, userID, stars, comment)
}

// GetRatings gets aggregated ratings for content - Issue #13
func (s *ContentService) GetRatings(ctx context.Context, contentID string) (*models.RatingAggregate, error) {
	return s.repo.GetRatings(ctx, contentID)
}

// GetSimilar gets similar content based on genre/tags - Issue #13
func (s *ContentService) GetSimilar(ctx context.Context, contentID string, limit int) ([]models.Content, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetSimilar(ctx, contentID, limit)
}

// GetEntitlements checks if user can access content - Issue #13
func (s *ContentService) GetEntitlements(ctx context.Context, contentID, userID string) (*models.Entitlement, error) {
	// TODO: Integrate with Payment Service to check subscription
	// TODO: Check geo-blocking based on user IP
	// TODO: Check DRM level based on subscription tier
	// For now, return a basic entitlement check
	return s.repo.GetEntitlements(ctx, contentID, userID)
}

