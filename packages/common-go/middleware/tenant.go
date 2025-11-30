package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streamverse/common-go/errors"
)

// TenantContext holds tenant information
type TenantContext struct {
	TenantID string
	UserID   string
	Roles    []string
}

const TenantContextKey = "tenant_context"

// TenantMiddleware extracts tenant ID from X-Tenant-ID header
// Issue #29: Multi-Tenancy & White-Label Support
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")

		// If no tenant ID in header, try to get from JWT token
		if tenantID == "" {
			if tenantIDClaim, exists := c.Get("tenant_id"); exists {
				if tid, ok := tenantIDClaim.(string); ok {
					tenantID = tid
				}
			}
		}

		// Set tenant ID in context
		if tenantID != "" {
			c.Set("tenant_id", tenantID)

			// Build tenant context
			roles, _ := c.Get("roles")
			rolesList := []string{}
			if roles != nil {
				if r, ok := roles.([]string); ok {
					rolesList = r
				}
			}

			ctx := &TenantContext{
				TenantID: tenantID,
				UserID:   getStringFromContext(c, "user_id"),
				Roles:    rolesList,
			}
			c.Set(TenantContextKey, ctx)
		}

		c.Next()
	}
}

// RequireTenant ensures tenant ID is present
func RequireTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, exists := c.Get("tenant_id")
		if !exists || tenantID == "" {
			c.JSON(http.StatusBadRequest, errors.NewInvalidInputError("X-Tenant-ID header is required"))
			c.Abort()
			return
		}
		c.Next()
	}
}

func getStringFromContext(c *gin.Context, key string) string {
	value, exists := c.Get(key)
	if !exists {
		return ""
	}
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}
