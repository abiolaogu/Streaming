package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/streamverse/admin-service/models"
	"github.com/streamverse/admin-service/repository"
)

// AdminService handles admin business logic
type AdminService struct {
	repo *repository.AdminRepository
}

// NewAdminService creates a new admin service
func NewAdminService(repo *repository.AdminRepository) *AdminService {
	return &AdminService{
		repo: repo,
	}
}

// LogAuditEvent logs an audit event
func (s *AdminService) LogAuditEvent(ctx context.Context, userID, action, resource, resourceID string, changes map[string]interface{}, ipAddress, userAgent string) error {
	log := &models.AuditLog{
		ID:         primitive.NewObjectID(),
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Changes:    changes,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		CreatedAt:  time.Now(),
	}
	return s.repo.CreateAuditLog(ctx, log)
}

// GetAuditLogs retrieves audit logs
func (s *AdminService) GetAuditLogs(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*models.AuditLog, int64, error) {
	return s.repo.GetAuditLogs(ctx, filters, page, pageSize)
}

// GetSystemSettings retrieves system settings
func (s *AdminService) GetSystemSettings(ctx context.Context) (*models.SystemSettings, error) {
	return s.repo.GetSystemSettings(ctx)
}

// UpdateSystemSettings updates system settings
func (s *AdminService) UpdateSystemSettings(ctx context.Context, settings *models.SystemSettings, updatedBy string) error {
	settings.UpdatedBy = updatedBy
	return s.repo.UpdateSystemSettings(ctx, settings)
}

// ListUsers lists users with filters
func (s *AdminService) ListUsers(ctx context.Context, filters *models.UserListFilters, page, pageSize int) ([]map[string]interface{}, int64, error) {
	filterMap := make(map[string]interface{})
	if filters != nil {
		if filters.Status != "" {
			filterMap["status"] = filters.Status
		}
		if filters.Role != "" {
			filterMap["role"] = filters.Role
		}
		if filters.Email != "" {
			filterMap["email"] = filters.Email
		}
	}
	return s.repo.ListUsers(ctx, filterMap, page, pageSize)
}

// GetUser retrieves a user by ID
func (s *AdminService) GetUser(ctx context.Context, userID string) (map[string]interface{}, error) {
	users, _, err := s.repo.ListUsers(ctx, map[string]interface{}{"id": userID}, 1, 1)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return users[0], nil
}

// UpdateUser updates a user
func (s *AdminService) UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error {
	return s.repo.UpdateUser(ctx, userID, updates)
}

// DeleteUser soft deletes a user
func (s *AdminService) DeleteUser(ctx context.Context, userID string) error {
	return s.repo.SoftDeleteUser(ctx, userID)
}

// ListContent lists content with filters
func (s *AdminService) ListContent(ctx context.Context, filters *models.ContentListFilters, page, pageSize int) ([]map[string]interface{}, int64, error) {
	filterMap := make(map[string]interface{})
	if filters != nil {
		if filters.Status != "" {
			filterMap["status"] = filters.Status
		}
		if filters.Category != "" {
			filterMap["category"] = filters.Category
		}
		if filters.Genre != "" {
			filterMap["genre"] = filters.Genre
		}
	}
	return s.repo.ListContent(ctx, filterMap, page, pageSize)
}

// UpdateContent updates content metadata
func (s *AdminService) UpdateContent(ctx context.Context, contentID string, updates map[string]interface{}) error {
	return s.repo.UpdateContent(ctx, contentID, updates)
}

// DeleteContent deletes content
func (s *AdminService) DeleteContent(ctx context.Context, contentID string) error {
	return s.repo.DeleteContent(ctx, contentID)
}

// BulkImportContent imports content from CSV/JSON
func (s *AdminService) BulkImportContent(ctx context.Context, reader io.Reader, format string) (*models.BulkImportResult, error) {
	result := &models.BulkImportResult{
		Errors: []string{},
	}

	// TODO: Parse CSV/JSON and bulk insert
	// For now, return a placeholder
	if format == "csv" {
		csvReader := csv.NewReader(reader)
		records, err := csvReader.ReadAll()
		if err != nil {
			return nil, fmt.Errorf("failed to parse CSV: %w", err)
		}
		result.Total = len(records) - 1 // Exclude header
		result.Success = result.Total
		// TODO: Insert records into database
	} else if format == "json" {
		var items []map[string]interface{}
		decoder := json.NewDecoder(reader)
		if err := decoder.Decode(&items); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
		result.Total = len(items)
		result.Success = result.Total
		// TODO: Insert items into database
	}

	return result, nil
}

