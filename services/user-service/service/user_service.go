package service

import (
	"context"
	"fmt"

	"github.com/streamverse/user-service/models"
	"github.com/streamverse/user-service/repository"
)

// UserService handles user profile business logic
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// GetProfile retrieves user profile
func (s *UserService) GetProfile(ctx context.Context, userID string) (*models.UserProfile, error) {
	return s.repo.GetProfile(ctx, userID)
}

// UpdateProfile updates user profile
func (s *UserService) UpdateProfile(ctx context.Context, userID string, profile *models.UserProfile) error {
	profile.UserID = userID
	return s.repo.UpdateProfile(ctx, userID, profile)
}

// GetPreferences retrieves user preferences
func (s *UserService) GetPreferences(ctx context.Context, userID string) (*models.UserPreferences, error) {
	return s.repo.GetPreferences(ctx, userID)
}

// UpdatePreferences updates user preferences
func (s *UserService) UpdatePreferences(ctx context.Context, userID string, prefs *models.UserPreferences) error {
	prefs.UserID = userID
	return s.repo.UpdatePreferences(ctx, userID, prefs)
}

// CreateProfile creates a sub-profile
func (s *UserService) CreateProfile(ctx context.Context, userID string, profile *models.Profile) (*models.Profile, error) {
	profile.UserID = userID
	return s.repo.CreateProfile(ctx, profile)
}

// GetProfiles retrieves all profiles for a user
func (s *UserService) GetProfiles(ctx context.Context, userID string) ([]models.Profile, error) {
	return s.repo.GetProfiles(ctx, userID)
}

// UpdateWatchHistory updates watch history
func (s *UserService) UpdateWatchHistory(ctx context.Context, history *models.WatchHistory) error {
	return s.repo.UpdateWatchHistory(ctx, history)
}

// GetWatchHistory retrieves watch history with pagination (Issue #12: max 1000 entries)
func (s *UserService) GetWatchHistory(ctx context.Context, userID string, page, pageSize int) ([]models.WatchHistory, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}
	if pageSize > 1000 {
		pageSize = 1000 // Max 1000 as per requirement
	}
	return s.repo.GetWatchHistory(ctx, userID, page, pageSize)
}

// GetContinueWatching retrieves continue watching items
func (s *UserService) GetContinueWatching(ctx context.Context, userID string, limit int) ([]models.WatchHistory, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetContinueWatching(ctx, userID, limit)
}

// ClearWatchHistory clears watch history
func (s *UserService) ClearWatchHistory(ctx context.Context, userID string) error {
	return s.repo.ClearWatchHistory(ctx, userID)
}

// AddToWatchlist adds content to watchlist
func (s *UserService) AddToWatchlist(ctx context.Context, userID, contentID string) error {
	return s.repo.AddToWatchlist(ctx, userID, contentID)
}

// RemoveFromWatchlist removes content from watchlist
func (s *UserService) RemoveFromWatchlist(ctx context.Context, userID, contentID string) error {
	return s.repo.RemoveFromWatchlist(ctx, userID, contentID)
}

// GetWatchlist retrieves watchlist
func (s *UserService) GetWatchlist(ctx context.Context, userID string) ([]models.Watchlist, error) {
	return s.repo.GetWatchlist(ctx, userID)
}

// DeleteUser deletes all user data (GDPR compliance)
func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	return s.repo.DeleteUser(ctx, userID)
}

// GetDevices retrieves all devices for a user
func (s *UserService) GetDevices(ctx context.Context, userID string) ([]models.Device, error) {
	return s.repo.GetDevices(ctx, userID)
}

// RegisterDevice registers a new device
func (s *UserService) RegisterDevice(ctx context.Context, device *models.Device) error {
	return s.repo.RegisterDevice(ctx, device)
}

// DeregisterDevice removes a device
func (s *UserService) DeregisterDevice(ctx context.Context, userID, deviceID string) error {
	return s.repo.DeregisterDevice(ctx, userID, deviceID)
}

// ExportUserData exports all user data for GDPR compliance
func (s *UserService) ExportUserData(ctx context.Context, userID string) (*models.UserDataExport, error) {
	return s.repo.ExportUserData(ctx, userID)
}

