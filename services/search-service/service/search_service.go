package service

import (
	"context"

	"github.com/streamverse/common-go/logger"
	"github.com/streamverse/search-service/repository"
)

// SearchService handles search business logic
type SearchService struct {
	repo   *repository.SearchRepository
	logger *logger.Logger
}

// NewSearchService creates a new search service
func NewSearchService(repo *repository.SearchRepository, logger *logger.Logger) *SearchService {
	return &SearchService{
		repo:   repo,
		logger: logger,
	}
}

// Search searches for content
func (s *SearchService) Search(ctx context.Context, query string, filters map[string]interface{}, page, pageSize int) ([]map[string]interface{}, int64, error) {
	return s.repo.Search(ctx, query, filters, page, pageSize)
}

// Autocomplete provides autocomplete suggestions
func (s *SearchService) Autocomplete(ctx context.Context, query string) ([]string, error) {
	return s.repo.Autocomplete(ctx, query)
}

// IndexContent indexes content in Elasticsearch
func (s *SearchService) IndexContent(ctx context.Context, id string, doc interface{}) error {
	return s.repo.IndexContent(ctx, id, doc)
}

// DeleteContent removes content from index
func (s *SearchService) DeleteContent(ctx context.Context, id string) error {
	return s.repo.DeleteContent(ctx, id)
}

