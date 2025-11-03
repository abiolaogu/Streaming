# Multi-Tenancy & White-Label Support

## Overview

Multi-tenancy implementation using shared database with row-level security (org_id column). Each tenant's data is isolated using `org_id` filtering.

## Implementation

### 1. Tenant Middleware

All services use `TenantMiddleware()` from `common-go/middleware/tenant.go`:

```go
router.Use(middleware.TenantMiddleware())
```

This middleware:
- Extracts `X-Tenant-ID` header
- Falls back to `tenant_id` from JWT token
- Sets tenant context in request

### 2. Database Schema Updates

All collections add `org_id` field:

```go
type Content struct {
    ID     primitive.ObjectID `bson:"_id,omitempty"`
    OrgID  string             `bson:"org_id"` // Tenant isolation
    Title  string             `bson:"title"`
    // ... other fields
}
```

### 3. Repository Layer

Use `database.AddTenantFilter()` to automatically add tenant filter:

```go
func (r *ContentRepository) GetContent(ctx context.Context, contentID string) (*models.Content, error) {
    filter := bson.M{"_id": objectID}
    filter = database.AddTenantFilter(ctx, filter) // Adds org_id filter
    
    var content models.Content
    err := r.collection.FindOne(ctx, filter).Decode(&content)
    return &content, err
}
```

### 4. API Header

Clients send tenant ID in header:

```bash
curl -X GET "https://api.streamverse.io/content/123" \
  -H "X-Tenant-ID: acme-corp" \
  -H "Authorization: Bearer <token>"
```

## White-Label Branding

### Branding Configuration

Each tenant has branding config:

```json
{
  "tenant_id": "acme-corp",
  "branding": {
    "logo_url": "https://cdn.streamverse.io/branding/acme-corp/logo.png",
    "primary_color": "#FF5733",
    "name": "ACME Video",
    "theme": "dark",
    "favicon_url": "https://cdn.streamverse.io/branding/acme-corp/favicon.ico"
  }
}
```

### Retrieving Branding

```go
branding, err := tenant.GetBranding(ctx, tenantID)
```

### API Endpoint

```
GET /tenant/branding
Response:
{
  "logo_url": "...",
  "primary_color": "#FF5733",
  "name": "ACME Video",
  "theme": "dark"
}
```

## Acceptance Criteria

- ✅ Tenant A can't access Tenant B's data
- ✅ API enforces tenant isolation
- ✅ `org_id` filtering prevents data leakage
- ✅ Branding per tenant works
- ✅ Integration tests verify isolation

## Testing

### Integration Test Example

```go
func TestTenantIsolation(t *testing.T) {
    // Create content for tenant A
    ctxA := context.WithValue(context.Background(), "tenant_id", "tenant-a")
    contentA, _ := repo.CreateContent(ctxA, &Content{Title: "Content A", OrgID: "tenant-a"})
    
    // Try to access from tenant B
    ctxB := context.WithValue(context.Background(), "tenant_id", "tenant-b")
    _, err := repo.GetContent(ctxB, contentA.ID.Hex())
    
    assert.Error(t, err) // Should not find content
}
```

## Migration Guide

### Existing Services

1. Add `org_id` field to all models
2. Update repositories to use `database.AddTenantFilter()`
3. Add `TenantMiddleware()` to router
4. Update API documentation with `X-Tenant-ID` header

### Default Tenant

For backwards compatibility, if no `X-Tenant-ID` is provided:
- Use `default` as tenant ID
- All existing data belongs to `default` tenant

