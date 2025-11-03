package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditLog represents an audit trail entry
type AuditLog struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string            `bson:"user_id" json:"userId"`
	Action      string            `bson:"action" json:"action"` // "create", "update", "delete"
	Resource    string            `bson:"resource" json:"resource"` // "user", "content", "settings"
	ResourceID  string            `bson:"resource_id" json:"resourceId"`
	Changes     map[string]interface{} `bson:"changes,omitempty" json:"changes,omitempty"`
	IPAddress   string            `bson:"ip_address,omitempty" json:"ipAddress,omitempty"`
	UserAgent   string            `bson:"user_agent,omitempty" json:"userAgent,omitempty"`
	CreatedAt   time.Time         `bson:"created_at" json:"createdAt"`
}

// SystemSettings represents system configuration
type SystemSettings struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FeatureFlags    map[string]bool    `bson:"feature_flags" json:"featureFlags"`
	MaxUploadSize   int64              `bson:"max_upload_size" json:"maxUploadSize"` // bytes
	MaintenanceMode bool                `bson:"maintenance_mode" json:"maintenanceMode"`
	CDNBaseURL      string              `bson:"cdn_base_url" json:"cdnBaseUrl"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updatedAt"`
	UpdatedBy       string              `bson:"updated_by" json:"updatedBy"`
}

// BulkImportResult represents bulk import operation result
type BulkImportResult struct {
	Total    int `json:"total"`
	Success  int `json:"success"`
	Failed   int `json:"failed"`
	Errors   []string `json:"errors,omitempty"`
}

// UserListFilters represents filters for listing users
type UserListFilters struct {
	Status   string `json:"status,omitempty"` // "active", "suspended", "deleted"
	Role     string `json:"role,omitempty"`
	Email    string `json:"email,omitempty"`
	CreatedAfter *time.Time `json:"createdAfter,omitempty"`
	CreatedBefore *time.Time `json:"createdBefore,omitempty"`
}

// ContentListFilters represents filters for listing content
type ContentListFilters struct {
	Status     string `json:"status,omitempty"` // "published", "draft", "archived"
	Category   string `json:"category,omitempty"`
	Genre      string `json:"genre,omitempty"`
	CreatedAfter *time.Time `json:"createdAfter,omitempty"`
	CreatedBefore *time.Time `json:"createdBefore,omitempty"`
}

