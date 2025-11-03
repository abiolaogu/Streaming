# Kong API Gateway Configuration

Kong API Gateway setup for StreamVerse platform.

## Features

- ✅ Service routing
- ✅ JWT authentication
- ✅ Rate limiting
- ✅ CORS configuration
- ✅ Request/response transformation
- ✅ Load balancing
- ✅ Request ID and correlation ID

## Setup

### Using Docker Compose

```bash
cd infrastructure/kong
docker-compose up -d
```

### Manual Setup

```bash
# Start Kong database
docker run -d --name kong-database \
  -p 5432:5432 \
  -e POSTGRES_USER=kong \
  -e POSTGRES_PASSWORD=kong \
  -e POSTGRES_DB=kong \
  postgres:15

# Run migrations
docker run --rm \
  --network host \
  -e KONG_DATABASE=postgres \
  -e KONG_PG_HOST=localhost \
  -e KONG_PG_USER=kong \
  -e KONG_PG_PASSWORD=kong \
  kong:3.4 kong migrations bootstrap

# Start Kong
docker run -d --name kong \
  --network host \
  -e KONG_DATABASE=postgres \
  -e KONG_PG_HOST=localhost \
  -e KONG_PG_USER=kong \
  -e KONG_PG_PASSWORD=kong \
  -v $(pwd)/kong.yml:/kong/kong.yml:ro \
  kong:3.4
```

## Configuration

Edit `kong.yml` to add/remove services and configure plugins.

## API Endpoints

- Proxy: `http://localhost:8000`
- Admin API: `http://localhost:8001`

## Services Configured

- `/api/v1/auth` → Auth Service
- `/api/v1/content` → Content Service
- `/api/v1/search` → Search Service
- `/api/v1/users` → User Service

