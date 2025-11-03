package service

import (
	"context"

	"github.com/streamverse/content-service/models"
	"github.com/streamverse/content-service/repository"
)

// ContentServiceAdvanced extends ContentService with advanced features
type ContentServiceAdvanced struct {
	*ContentService
	collectionsRepo *repository.CollectionsRepository
	fastChannelsRepo *repository.FASTChannelsRepository
}

// NewContentServiceAdvanced creates an advanced content service
func NewContentServiceAdvanced(
	baseService *ContentService,
	collectionsRepo *repository.CollectionsRepository,
	fastChannelsRepo *repository.FASTChannelsRepository,
) *ContentServiceAdvanced {
	return &ContentServiceAdvanced{
		ContentService:   baseService,
		collectionsRepo:  collectionsRepo,
		fastChannelsRepo: fastChannelsRepo,
	}
}

// CreateCollection creates a new collection/playlist
func (s *ContentServiceAdvanced) CreateCollection(ctx context.Context, collection *models.Collection) (*models.Collection, error) {
	return s.collectionsRepo.Create(ctx, collection)
}

// AddContentToCollection adds content to a collection
func (s *ContentServiceAdvanced) AddContentToCollection(ctx context.Context, collectionID, contentID string) error {
	return s.collectionsRepo.AddContent(ctx, collectionID, contentID)
}

// RemoveContentFromCollection removes content from a collection
func (s *ContentServiceAdvanced) RemoveContentFromCollection(ctx context.Context, collectionID, contentID string) error {
	return s.collectionsRepo.RemoveContent(ctx, collectionID, contentID)
}

// CreateFASTChannel creates a new FAST channel
func (s *ContentServiceAdvanced) CreateFASTChannel(ctx context.Context, channel *models.FASTChannel) (*models.FASTChannel, error) {
	return s.fastChannelsRepo.Create(ctx, channel)
}

// UpdateChannelSchedule updates FAST channel schedule
func (s *ContentServiceAdvanced) UpdateChannelSchedule(ctx context.Context, channelID string, schedule []models.ScheduleItem) error {
	return s.fastChannelsRepo.UpdateSchedule(ctx, channelID, schedule)
}

