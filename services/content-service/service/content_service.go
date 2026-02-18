package service

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/streamverse/common-go/cache"
	"github.com/streamverse/content-service/models"
	"github.com/streamverse/content-service/repository"
)

// ContentService handles content business logic
type ContentService struct {
	repo                *repository.ContentRepository
	cache               *cache.RedisClient
	entitlementProvider EntitlementProvider
	policyProvider      EntitlementPolicyProvider
}

// NewContentService creates a new content service
func NewContentService(
	repo *repository.ContentRepository,
	cache *cache.RedisClient,
	entitlementProvider EntitlementProvider,
	policyProvider EntitlementPolicyProvider,
) *ContentService {
	return &ContentService{
		repo:                repo,
		cache:               cache,
		entitlementProvider: entitlementProvider,
		policyProvider:      policyProvider,
	}
}

// GetContentByID retrieves content by ID
func (s *ContentService) GetContentByID(ctx context.Context, id string) (*models.Content, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("content:%s", id)
	var content models.Content
	if err := s.cache.Get(ctx, cacheKey, &content); err == nil {
		return &content, nil
	}

	// Cache miss - query DB
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache result (TTL: 1 hour)
	if err := s.cache.Set(ctx, cacheKey, c, time.Hour); err != nil {
		// Log error
	}

	return c, nil
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
	// Try cache first
	cacheKey := "content:categories"
	var categories []models.Category
	if err := s.cache.Get(ctx, cacheKey, &categories); err == nil {
		return categories, nil
	}

	// Cache miss - query DB
	cats, err := s.repo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	// Cache result (TTL: 1 hour)
	if err := s.cache.Set(ctx, cacheKey, cats, time.Hour); err != nil {
		// Log error
	}

	return cats, nil
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
func (s *ContentService) GetEntitlements(ctx context.Context, contentID, userID, countryCode, authHeader string) (*models.Entitlement, error) {
	content, err := s.repo.GetByID(ctx, contentID)
	if err != nil {
		return nil, err
	}

	var entitlementRecords []map[string]interface{}
	if s.entitlementProvider != nil {
		entitlementRecords, err = s.entitlementProvider.GetUserEntitlements(ctx, userID, authHeader)
		if err != nil {
			return nil, err
		}
	}

	if s.policyProvider != nil {
		return s.policyProvider.EvaluateEntitlement(ctx, PolicyEntitlementInput{
			ContentID:       contentID,
			UserID:          userID,
			CountryCode:     countryCode,
			ContentCategory: content.Category,
			IsDRMProtected:  content.IsDRMProtected,
			Entitlements:    entitlementRecords,
		}, authHeader)
	}

	return evaluateLocalEntitlement(content, contentID, userID, countryCode, entitlementRecords), nil
}

func evaluateLocalEntitlement(content *models.Content, contentID, userID, countryCode string, entitlementRecords []map[string]interface{}) *models.Entitlement {
	planID, hasSubscription, hasPurchase, purchaseExpiresAt := evaluateEntitlements(entitlementRecords, contentID)

	entitlement := &models.Entitlement{
		ContentID: contentID,
		UserID:    userID,
		HasAccess: false,
		Reason:    "subscription_required",
		DRMLevel:  drmLevelForPlan(planID),
	}

	if hasPurchase {
		entitlement.HasAccess = true
		entitlement.Reason = "purchased"
		entitlement.ExpiresAt = purchaseExpiresAt
	} else if hasSubscription {
		entitlement.HasAccess = true
		entitlement.Reason = "subscription"
	} else if content.Category == "free" || content.Category == "avod" {
		entitlement.HasAccess = true
		entitlement.Reason = "free"
	}

	if content.IsDRMProtected {
		entitlement.LicenseURL = defaultLicenseURL()
	}

	if isGeoBlocked(countryCode) {
		entitlement.HasAccess = false
		entitlement.Reason = "geo_blocked"
		entitlement.ExpiresAt = nil
	}

	return entitlement
}

func isGeoBlocked(countryCode string) bool {
	countryCode = strings.ToUpper(strings.TrimSpace(countryCode))
	if countryCode == "" {
		return false
	}

	blocked := strings.Split(os.Getenv("GEO_BLOCKED_COUNTRIES"), ",")
	for _, country := range blocked {
		if strings.ToUpper(strings.TrimSpace(country)) == countryCode {
			return true
		}
	}
	return false
}

func evaluateEntitlements(records []map[string]interface{}, contentID string) (planID string, hasSubscription bool, hasPurchase bool, purchaseExpiresAt *time.Time) {
	for _, record := range records {
		recordType := strings.ToLower(toString(record["type"]))
		switch recordType {
		case "subscription":
			status := strings.ToLower(toString(record["status"]))
			if status == "active" || status == "trialing" {
				hasSubscription = true
				if planID == "" {
					planID = toString(record["plan_id"])
				}
			}
		case "purchase":
			if toString(record["content_id"]) != contentID {
				continue
			}
			status := strings.ToLower(toString(record["status"]))
			if status != "" && status != "completed" {
				continue
			}
			hasPurchase = true
			if parsed := toTime(record["expires_at"]); parsed != nil {
				purchaseExpiresAt = parsed
			}
		}
	}
	return planID, hasSubscription, hasPurchase, purchaseExpiresAt
}

func toString(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}

func toTime(value interface{}) *time.Time {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case string:
		if v == "" {
			return nil
		}
		if parsed, err := time.Parse(time.RFC3339, v); err == nil {
			return &parsed
		}
	case map[string]interface{}:
		if dateStr, ok := v["$date"].(string); ok {
			if parsed, err := time.Parse(time.RFC3339, dateStr); err == nil {
				return &parsed
			}
		}
	}
	return nil
}

func drmLevelForPlan(planID string) string {
	switch planID {
	case "premium", "tier3":
		return "1"
	case "pro", "standard", "tier2":
		return "2"
	default:
		return "3"
	}
}

func defaultLicenseURL() string {
	if url := os.Getenv("DRM_LICENSE_URL"); url != "" {
		return url
	}
	return "https://license.widevine.com/license"
}
