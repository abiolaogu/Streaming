# Search Service

Microservice for content search using Elasticsearch.

## Features

- ✅ Full-text search
- ✅ Autocomplete/suggestions
- ✅ Filtering by genre, category, etc.
- ✅ Fuzzy matching
- ✅ Content indexing
- ✅ Result pagination

## API Endpoints

- `GET /health` - Health check
- `GET /api/v1/search?q=query` - Search content
- `GET /api/v1/search/autocomplete?q=query` - Get autocomplete suggestions
- `POST /api/v1/search/index` - Index content

## Environment Variables

- `SERVER_PORT` - Server port (default: 8080)
- `SERVER_HOST` - Server host (default: 0.0.0.0)
- `ELASTICSEARCH_ADDRESSES` - Elasticsearch addresses (comma-separated, default: http://localhost:9200)
- `LOG_LEVEL` - Log level (default: info)

## Running

```bash
go run main.go
```

## Docker

```bash
docker build -t search-service .
docker run -p 8080:8080 search-service
```

