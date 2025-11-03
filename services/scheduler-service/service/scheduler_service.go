package service

import (
	"context"
	"fmt"
	"time"

	"github.com/streamverse/scheduler-service/models"
	"github.com/streamverse/scheduler-service/repository"
)

// SchedulerService handles scheduler business logic
type SchedulerService struct {
	repo *repository.SchedulerRepository
	cdnBaseURL string // TODO: Load from config
}

// NewSchedulerService creates a new scheduler service
func NewSchedulerService(repo *repository.SchedulerRepository, cdnBaseURL string) *SchedulerService {
	return &SchedulerService{
		repo: repo,
		cdnBaseURL: cdnBaseURL,
	}
}

// ListChannels lists all channels
func (s *SchedulerService) ListChannels(ctx context.Context, status string) ([]*models.Channel, error) {
	return s.repo.ListChannels(ctx, status)
}

// GetChannelByID retrieves a channel by ID
func (s *SchedulerService) GetChannelByID(ctx context.Context, channelID string) (*models.Channel, error) {
	return s.repo.GetChannelByID(ctx, channelID)
}

// CreateChannel creates a new channel
func (s *SchedulerService) CreateChannel(ctx context.Context, channel *models.Channel) error {
	return s.repo.CreateChannel(ctx, channel)
}

// UpdateChannel updates a channel
func (s *SchedulerService) UpdateChannel(ctx context.Context, channelID string, updates map[string]interface{}) error {
	return s.repo.UpdateChannel(ctx, channelID, updates)
}

// GetChannelEPG generates EPG for a channel (next 7 days)
func (s *SchedulerService) GetChannelEPG(ctx context.Context, channelID string) (*models.EPG, error) {
	channel, err := s.repo.GetChannelByID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("channel not found: %w", err)
	}

	startTime := time.Now()
	endTime := startTime.Add(7 * 24 * time.Hour) // Next 7 days

	entries, err := s.repo.GetScheduleEntries(ctx, channelID, &startTime, &endTime, 0)
	if err != nil {
		return nil, err
	}

	epgEntries := make([]models.EPGEntry, 0, len(entries))
	for _, entry := range entries {
		epgEntries = append(epgEntries, models.EPGEntry{
			Title:       entry.Title,
			StartTime:   entry.StartTime,
			Duration:    entry.Duration,
			Description: entry.Description,
			Poster:      entry.Poster,
			ContentID:   entry.ContentID,
		})
	}

	return &models.EPG{
		ChannelID:   channel.ChannelID,
		ChannelName: channel.Name,
		Schedule:   epgEntries,
		GeneratedAt: time.Now(),
	}, nil
}

// GetChannelManifest generates streaming manifest URL for a channel
func (s *SchedulerService) GetChannelManifest(ctx context.Context, channelID string) (*models.ChannelManifest, error) {
	channel, err := s.repo.GetChannelByID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("channel not found: %w", err)
	}

	if channel.Type == "live" && channel.ManifestURL != "" {
		return &models.ChannelManifest{
			ChannelID:   channelID,
			ManifestURL: channel.ManifestURL,
			Type:        "hls",
		}, nil
	}

	// For FAST channels, generate manifest URL based on current schedule
	// TODO: Generate actual HLS manifest with current content
	manifestURL := fmt.Sprintf("%s/channels/%s/manifest.m3u8", s.cdnBaseURL, channelID)

	return &models.ChannelManifest{
		ChannelID:   channelID,
		ManifestURL: manifestURL,
		Type:        "hls",
	}, nil
}

// CreateScheduleEntry creates a new schedule entry
func (s *SchedulerService) CreateScheduleEntry(ctx context.Context, entry *models.ScheduleEntry) error {
	// Validate entry
	if entry.StartTime.After(entry.EndTime) {
		return fmt.Errorf("start time must be before end time")
	}

	entry.Duration = int(entry.EndTime.Sub(entry.StartTime).Seconds())
	return s.repo.CreateScheduleEntry(ctx, entry)
}

// UpdateScheduleEntry updates a schedule entry
func (s *SchedulerService) UpdateScheduleEntry(ctx context.Context, entryID string, updates map[string]interface{}) error {
	return s.repo.UpdateScheduleEntry(ctx, entryID, updates)
}

// DeleteScheduleEntry deletes a schedule entry
func (s *SchedulerService) DeleteScheduleEntry(ctx context.Context, entryID string) error {
	return s.repo.DeleteScheduleEntry(ctx, entryID)
}

// GetCurrentScheduleEntry gets the currently playing entry for a channel
func (s *SchedulerService) GetCurrentScheduleEntry(ctx context.Context, channelID string) (*models.ScheduleEntry, error) {
	return s.repo.GetCurrentScheduleEntry(ctx, channelID)
}

