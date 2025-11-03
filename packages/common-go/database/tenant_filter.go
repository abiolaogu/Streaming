package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

// AddTenantFilter adds tenant_id filter to MongoDB query
// Issue #29: Multi-Tenancy & White-Label Support
func AddTenantFilter(ctx context.Context, filter bson.M) bson.M {
	// Get tenant_id from context
	tenantID, exists := ctx.Value("tenant_id").(string)
	if exists && tenantID != "" {
		filter["org_id"] = tenantID
	}
	return filter
}

// GetTenantIDFromContext extracts tenant ID from context
func GetTenantIDFromContext(ctx context.Context) string {
	if tenantID, ok := ctx.Value("tenant_id").(string); ok {
		return tenantID
	}
	return ""
}

