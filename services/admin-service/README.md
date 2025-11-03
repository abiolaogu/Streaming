# Admin Service

Admin Service provides system management and monitoring capabilities with RBAC (Role-Based Access Control).

## Features

- **User Management**: List, view, update, and soft-delete users
- **Content Management**: List, update, delete, and bulk import content
- **Analytics Dashboard**: System-wide metrics and KPIs
- **System Settings**: Feature flags, maintenance mode, CDN configuration
- **Audit Logging**: Complete audit trail of all admin actions

## RBAC Roles

- **superadmin**: All operations including user deletion and system settings
- **admin**: Manage users and content (except billing operations)
- **editor**: Edit content metadata only

## Endpoints

### User Management
- `GET /admin/users` - List users (filter, pagination, export)
- `GET /admin/users/{id}` - Get user details
- `PUT /admin/users/{id}` - Update user
- `DELETE /admin/users/{id}` - Delete user (soft delete)

### Content Management
- `GET /admin/content` - List content
- `POST /admin/content` - Bulk import content (CSV/JSON)
- `PUT /admin/content/{id}` - Edit content metadata
- `DELETE /admin/content/{id}` - Remove content

### Analytics
- `GET /admin/analytics` - Dashboard metrics

### Settings
- `GET /admin/settings` - System configuration
- `PUT /admin/settings` - Update settings (superadmin only)

### Audit Logs
- `GET /admin/audit-logs` - Audit trail (filterable by resource, action, user)

## Usage

All endpoints require JWT authentication and appropriate RBAC role.

### Example: List Users

```bash
curl -X GET "http://localhost:8080/admin/users?page=1&page_size=20&status=active" \
  -H "Authorization: Bearer <token>"
```

### Example: Bulk Import Content

```bash
curl -X POST "http://localhost:8080/admin/content?format=json" \
  -H "Authorization: Bearer <token>" \
  -F "file=@content.json"
```

## Environment Variables

```bash
DATABASE_URI=mongodb://localhost:27017
DATABASE_NAME=streamverse
JWT_SECRET_KEY=your-secret-key
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
```

## Development

```bash
go mod download
go run main.go
```

## Testing

```bash
go test ./...
```

## Docker

```bash
docker build -t admin-service .
docker run -p 8080:8080 admin-service
```

