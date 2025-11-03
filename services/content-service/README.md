# Content Service

Microservice for managing video catalog, metadata, and assets.

## Features

- ✅ Content CRUD operations
- ✅ Category-based filtering
- ✅ Search functionality
- ✅ Home content rows
- ✅ MongoDB persistence
- ✅ JWT authentication

## API Endpoints

- `GET /health` - Health check
- `GET /api/v1/content/home` - Get home screen content rows
- `GET /api/v1/content/:id` - Get content by ID
- `GET /api/v1/content/category/:category` - Get content by category
- `POST /api/v1/content` - Create content
- `PUT /api/v1/content/:id` - Update content
- `DELETE /api/v1/content/:id` - Delete content
- `GET /api/v1/content/search?q=query` - Search content

## Environment Variables

- `SERVER_PORT` - Server port (default: 8080)
- `SERVER_HOST` - Server host (default: 0.0.0.0)
- `DATABASE_URI` - MongoDB connection URI
- `DATABASE_NAME` - Database name (default: streamverse)
- `JWT_SECRET_KEY` - JWT secret key
- `LOG_LEVEL` - Log level (default: info)

## Running

```bash
go run main.go
```

## Docker

```bash
docker build -t content-service .
docker run -p 8080:8080 content-service
```

