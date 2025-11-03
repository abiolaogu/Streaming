package service

import (
	"context"

	"github.com/streamverse/ad-service/models"
	"github.com/streamverse/ad-service/repository"
)

// AdService handles ad business logic
type AdService struct {
	repo *repository.AdRepository
}

// NewAdService creates a new ad service
func NewAdService(repo *repository.AdRepository) *AdService {
	return &AdService{
		repo: repo,
	}
}

// GetAds retrieves ads for a content request
func (s *AdService) GetAds(ctx context.Context, req *models.AdRequest) (*models.AdResponse, error) {
	// Check if user has ad-free subscription
	if s.isAdFreeUser(ctx, req.UserID) {
		return &models.AdResponse{Ads: []models.Ad{}}, nil
	}

	// Get targeted ads
	ads := s.repo.GetAdsByTargeting(ctx, req)

	return &models.AdResponse{
		Ads:         ads,
		SkipAllowed: req.Position == "pre-roll",
	}, nil
}

// TrackAdEvent tracks an ad event
func (s *AdService) TrackAdEvent(ctx context.Context, tracking *models.AdTracking) error {
	return s.repo.CreateTracking(ctx, tracking)
}

func (s *AdService) isAdFreeUser(ctx context.Context, userID string) bool {
	// TODO: Check subscription status
	return false
}

