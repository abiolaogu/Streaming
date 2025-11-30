package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/streamverse/common-go/logger"
	"github.com/streamverse/content-service/models"
	"github.com/streamverse/content-service/service"
	"go.uber.org/zap"
)

// ContentHandler handles HTTP requests for content
type ContentHandler struct {
	service *service.ContentService
	logger  *logger.Logger
}

// NewContentHandler creates a new content handler
func NewContentHandler(service *service.ContentService, logger *logger.Logger) *ContentHandler {
	return &ContentHandler{
		service: service,
		logger:  logger,
	}
}

// GetContentByID handles GET /api/v1/content/:id
func (h *ContentHandler) GetContentByID(c *gin.Context) {
	id := c.Param("id")

	content, err := h.service.GetContentByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get content", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
		return
	}

	c.JSON(http.StatusOK, content)
}

// GetContentByCategory handles GET /api/v1/content/category/:category
func (h *ContentHandler) GetContentByCategory(c *gin.Context) {
	category := c.Param("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	contents, total, err := h.service.ListContent(c.Request.Context(), category, page, pageSize)
	if err != nil {
		h.logger.Error("Failed to list content", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     contents,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateContent handles POST /api/v1/content
func (h *ContentHandler) CreateContent(c *gin.Context) {
	var content models.Content
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.service.CreateContent(c.Request.Context(), &content)
	if err != nil {
		h.logger.Error("Failed to create content", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// UpdateContent handles PUT /api/v1/content/:id
func (h *ContentHandler) UpdateContent(c *gin.Context) {
	id := c.Param("id")

	var content models.Content
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateContent(c.Request.Context(), id, &content); err != nil {
		h.logger.Error("Failed to update content", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Content updated successfully"})
}

// DeleteContent handles DELETE /api/v1/content/:id
func (h *ContentHandler) DeleteContent(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteContent(c.Request.Context(), id); err != nil {
		h.logger.Error("Failed to delete content", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Content deleted successfully"})
}

// SearchContent handles GET /api/v1/content/search
func (h *ContentHandler) SearchContent(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	contents, total, err := h.service.SearchContent(c.Request.Context(), query, page, pageSize)
	if err != nil {
		h.logger.Error("Failed to search content", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results":   contents,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetHomeContent handles GET /api/v1/content/home
func (h *ContentHandler) GetHomeContent(c *gin.Context) {
	rows, err := h.service.GetHomeContent(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get home content", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch home content"})
		return
	}

	c.JSON(http.StatusOK, rows)
}

// GetCategories handles GET /content/categories - Issue #13
func (h *ContentHandler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get categories", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// GetTrending handles GET /content/trending - Issue #13
func (h *ContentHandler) GetTrending(c *gin.Context) {
	region := c.DefaultQuery("region", "US")
	deviceType := c.DefaultQuery("device_type", "")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	trending, err := h.service.GetTrending(c.Request.Context(), region, deviceType, limit)
	if err != nil {
		h.logger.Error("Failed to get trending", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trending content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": trending})
}

// RateContent handles POST /content/{id}/ratings - Issue #13
func (h *ContentHandler) RateContent(c *gin.Context) {
	contentID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		Stars   int    `json:"stars" binding:"required,min=1,max=5"`
		Comment string `json:"comment,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.RateContent(c.Request.Context(), contentID, userID.(string), req.Stars, req.Comment)
	if err != nil {
		h.logger.Error("Failed to rate content", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit rating"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rating submitted successfully"})
}

// GetRatings handles GET /content/{id}/ratings - Issue #13
func (h *ContentHandler) GetRatings(c *gin.Context) {
	contentID := c.Param("id")

	aggregate, err := h.service.GetRatings(c.Request.Context(), contentID)
	if err != nil {
		h.logger.Error("Failed to get ratings", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ratings"})
		return
	}

	c.JSON(http.StatusOK, aggregate)
}

// GetSimilar handles GET /content/{id}/similar - Issue #13
func (h *ContentHandler) GetSimilar(c *gin.Context) {
	contentID := c.Param("id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	similar, err := h.service.GetSimilar(c.Request.Context(), contentID, limit)
	if err != nil {
		h.logger.Error("Failed to get similar content", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch similar content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": similar})
}

// GetEntitlements handles GET /content/{id}/entitlements - Issue #13
func (h *ContentHandler) GetEntitlements(c *gin.Context) {
	contentID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	entitlement, err := h.service.GetEntitlements(c.Request.Context(), contentID, userID.(string))
	if err != nil {
		h.logger.Error("Failed to get entitlements", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check entitlements"})
		return
	}

	c.JSON(http.StatusOK, entitlement)
}
