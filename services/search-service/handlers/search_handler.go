package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/streamverse/common-go/logger"
	"github.com/streamverse/search-service/service"
)

// SearchHandler handles HTTP requests for search
type SearchHandler struct {
	service *service.SearchService
	logger  *logger.Logger
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(service *service.SearchService, logger *logger.Logger) *SearchHandler {
	return &SearchHandler{
		service: service,
		logger:  logger,
	}
}

// Search handles GET /search - Issue #17
func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// Build filters
	filters := make(map[string]interface{})
	if genre := c.Query("genre"); genre != "" {
		filters["genre"] = genre
	}
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}

	results, total, err := h.service.Search(c.Request.Context(), query, filters, page, pageSize)
	if err != nil {
		h.logger.Error("Search failed", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results":  results,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

// Suggest handles GET /search/suggest - Issue #17
func (h *SearchHandler) Suggest(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	suggestions, err := h.service.Autocomplete(c.Request.Context(), query)
	if err != nil {
		h.logger.Error("Autocomplete failed", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Autocomplete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"suggestions": suggestions})
}

// IndexContent handles POST /api/v1/search/index
func (h *SearchHandler) IndexContent(c *gin.Context) {
	var doc map[string]interface{}
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, ok := doc["id"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document must have 'id' field"})
		return
	}

	if err := h.service.IndexContent(c.Request.Context(), id, doc); err != nil {
		h.logger.Error("Indexing failed", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Indexing failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Content indexed successfully"})
}

// GetFilters handles GET /search/filters - Issue #17
func (h *SearchHandler) GetFilters(c *gin.Context) {
	filters := map[string]interface{}{
		"genres": []string{"action", "comedy", "drama", "horror", "sci-fi", "thriller"},
		"years": []int{2020, 2021, 2022, 2023, 2024},
		"ratings": []string{"G", "PG", "PG-13", "R", "NC-17"},
		"types": []string{"movie", "show", "live"},
	}
	c.JSON(http.StatusOK, filters)
}

// Autocomplete handles GET /api/v1/search/autocomplete (deprecated, use Suggest)
func (h *SearchHandler) Autocomplete(c *gin.Context) {
	h.Suggest(c)
}

