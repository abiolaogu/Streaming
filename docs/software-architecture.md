# Backend Architecture — StreamVerse Streaming Platform

## 1. Microservices Overview

StreamVerse employs a domain-driven microservices architecture with 16 distinct services organized into four service groups: Core, Media, Data, and Support. Each service owns its domain logic, exposes both HTTP (Gin) and gRPC interfaces, and communicates asynchronously via Apache Kafka.

### Service Registry

| Service | Language | HTTP Port | gRPC Port | Domain |
|---------|----------|-----------|-----------|--------|
| auth-service | Go | 8081 | 9081 | Authentication, OAuth, MFA, JWT |
| user-service | Go | 8082 | 9082 | Profiles, preferences, watch history |
| content-service | Go | 8083 | 9083 | Metadata, categories, series, episodes |
| streaming-service | Go | 8084 | 9084 | Playback URLs, ABR, DRM, sessions |
| payment-service | Go | 8085 | 9085 | Subscriptions, billing, Stripe |
| transcoding-service | Go | 8086 | 9086 | FFmpeg, GPU transcoding, ABR ladder |
| ad-service | Go | 8087 | 9087 | Ad inventory, targeting, VAST/VMAP |
| ad-compositing-service | Go | 8088 | 9088 | SSAI, scene detection, ad stitching |
| analytics-service | Go/TS | 8089 | 9089 | User behavior, A/B testing, dashboards |
| search-service | Go | 8090 | 9090 | Elasticsearch, autocomplete, facets |
| recommendation-service | Python | 8091 | 9091 | NCF, collaborative filtering, ML |
| notification-service | Node.js | 8092 | 9092 | FCM, APNs, SendGrid, Twilio |
| websocket-service | Node.js | 8093 | 9093 | Real-time sync, watch parties, chat |
| scheduler-service | Go | 8094 | 9094 | Cron jobs, content publishing |
| admin-service | Go | 8095 | 9095 | CMS, moderation, user management |
| training-bot-service | Go | 8096 | -- | AI assistant (Gemini integration) |

---

## 2. Service Communication Patterns

### Synchronous (Request-Response)
- **gRPC**: Inter-service calls for real-time data needs
  - Streaming service calls content-service via gRPC to verify entitlements
  - Streaming service calls payment-service via gRPC to check subscription status
  - Protobuf definitions in `packages/proto/` with generated Go code in `packages/proto/gen/go/`

- **HTTP REST**: External client-facing APIs through Kong API Gateway
  - All services expose `/health` and `/ready` endpoints
  - `/metrics` endpoint for Prometheus scraping

### Asynchronous (Event-Driven)
- **Apache Kafka**: Primary event bus
  - Topics: `content.uploaded`, `content.transcoded`, `playback.events`, `subscription.activated`, `user.registered`, `notification.send`
  - Exactly-once semantics where required
  - Partitioning by user_id or content_id for ordering guarantees

- **gNATS**: Lightweight real-time messaging
  - Subject-based routing for notification fan-out
  - Low-latency pub/sub for WebSocket event distribution

---

## 3. Shared Libraries (packages/common-go/)

The `packages/common-go/` directory contains shared Go packages used across all Go services:

### cache/
Redis client wrapper with connection pooling, serialization helpers, and TTL management. Used for session caching, rate limiting counters, and recommendation result caching.

### config/
Viper-based configuration loader supporting environment variables, YAML files, and HashiCorp Vault secrets. Implements the 12-factor app methodology.

### database/
PostgreSQL connection pool manager using pgx driver. Provides transaction helpers, migration runner interface, and connection health checks.

### errors/
Standardized error types with error codes mapping to HTTP status codes and gRPC status codes. Includes error wrapping for distributed tracing context.

### i18n/
Internationalization support with message catalogs and locale resolution. Supports 50+ languages with fallback chains.

### jwt/
JWT token generation, validation, and refresh logic. Supports RS256 signing with key rotation via Vault.

### logger/
Structured logging with zerolog. Includes request ID propagation, correlation ID injection, and log level configuration per environment.

### middleware/
Common HTTP middleware: authentication, authorization, rate limiting, CORS, request logging, panic recovery, and request ID injection.

### tenant/
Multi-tenancy context propagation. Extracts tenant ID from JWT claims or headers and injects into database query scopes.

---

## 4. Auth Service Architecture

Located in `services/auth-service/`, this service handles all authentication concerns:

```
auth-service/
  main.go              # Service entry point, server initialization
  handlers/
    auth_handler.go    # HTTP route handlers (register, login, refresh, logout)
  models/
    user.go            # User model, validation rules
  repository/
    auth_repository.go # PostgreSQL queries for user CRUD
  service/
    auth_service.go    # Business logic (password hashing, token generation)
  utils/
    jwt.go             # JWT helper functions
```

### Key Implementation Details
- Password hashing: bcrypt with cost factor 12
- JWT access tokens: 1-hour expiry, RS256 signed
- JWT refresh tokens: 30-day expiry, stored in Redis for revocation
- OAuth 2.0: Google and Apple providers (configurable via env vars)
- MFA: TOTP-based (Google Authenticator compatible)
- Session tracking: Redis-backed with concurrent session limits per subscription tier

---

## 5. Streaming Service Architecture

Located in `services/streaming-service/`, this is the critical path for video delivery:

```
streaming-service/
  main.go
  config/              # Service-specific configuration
  handlers/
    streaming_handler.go  # Play, heartbeat, complete endpoints
  models/
    stream.go            # Stream session, quality, DRM models
    streaming.go         # Additional streaming types
  repository/
    session_repository.go    # Playback session persistence
    streaming_repository.go  # Stream URL and manifest queries
  service/
    streaming_service.go     # Orchestration: entitlement + DRM + CDN
  utils/
    cdn.go              # CDN signed URL generation
    manifest.go         # HLS/DASH manifest generation
  internal/clients/
    content/
      content_client.go   # gRPC client to content-service
    payment/
      payment_client.go   # gRPC client to payment-service
```

### Playback Session Management
1. Session created on `/play` request with unique `sessionId`
2. Heartbeat received every 30 seconds with position, quality, buffer health
3. Sessions tracked in Redis (fast reads) and PostgreSQL (durability)
4. Concurrent stream enforcement: check active sessions against subscription tier limit
5. Session cleanup: expired sessions purged by scheduler-service

---

## 6. Transcoding Service Architecture

Located in `services/transcoding-service/`:

```
transcoding-service/
  main.go
  handlers/
    transcoding_handler.go  # Job submission and status endpoints
  models/
    job.go                  # Transcoding job model (source, profile, status, progress)
  repository/
    transcoding_repository.go  # Job persistence
  service/
    transcoding_service.go     # FFmpeg pipeline orchestration
```

### Transcoding Pipeline
1. Job received from Kafka (`content.uploaded` topic)
2. Source file downloaded from S3/MinIO to local SSD
3. FFprobe analysis: codec, resolution, bitrate, audio channels, duration
4. ABR ladder computed based on source quality:
   - 480p @ 1.5 Mbps (H.264)
   - 720p @ 3 Mbps (H.264)
   - 1080p @ 6 Mbps (H.264 + HEVC)
   - 4K @ 25 Mbps (HEVC)
5. GPU-accelerated transcoding via NVIDIA NVENC
6. HLS packaging: m3u8 master + variant playlists, .ts segments (6s each)
7. DASH packaging: MPD manifest, .m4s segments
8. Thumbnail extraction: poster image + timeline thumbnails every 10 seconds
9. Upload transcoded assets to S3/MinIO
10. Publish `content.transcoded` event to Kafka

### Hybrid GPU Scaling (streaming-saas layer)
- Local GPUs handle baseline queue (up to `LOCAL_GPU_COUNT * 10` jobs/hour)
- When queue depth exceeds threshold, Runpod.io serverless pods auto-provision
- Runpod.io client (`streaming-saas/transcoding-service/src/runpod_client.rs`) manages pod lifecycle
- Cost optimization: RTX 4090 for standard ($0.39/hr), A100 for AI enhancement ($1.89/hr), H100 for 8K ($4.50/hr)

---

## 7. Content Service Architecture

Located in `services/content-service/`:

```
content-service/
  main.go
  handlers/
    content_handler.go    # CRUD, categories, search, home page
  models/
    content.go            # Content, Season, Episode models
  repository/
    content_repository.go # PostgreSQL queries with pagination
  service/
    content_service.go    # Business logic, cache management
```

### Content Data Model
- `content` table: movies, series, episodes, live channels, FAST channels
- `seasons` table: season metadata linked to series
- `episodes` table: episode metadata linked to season and content
- Monetization field: `avod`, `svod`, `tvod`, `ppv`
- Availability field (JSONB): geo-restriction per country
- Flexible metadata (JSONB): extensible without schema changes

---

## 8. Database Architecture

### PostgreSQL (Primary Datastore)
Tables defined in `database-schema.sql`:
- `users` (UUID PK, email, password_hash, subscription_tier/status)
- `profiles` (UUID PK, user_id FK, name, avatar, is_kids, preferences JSONB)
- `content` (UUID PK, type, title, genres[], stream_url, monetization, availability JSONB)
- `seasons` (UUID PK, series_id FK, season_number)
- `episodes` (UUID PK, season_id FK, content_id FK, episode_number)
- `watch_history` (profile_id FK, content_id FK, position, completed)
- `watchlist` (profile_id FK, content_id FK)
- `ratings` (profile_id FK, content_id FK, 1-5 rating, review text)
- `subscription_plans` (name, price, features JSONB)
- `transactions` (user_id FK, type, amount, payment_provider_id, status)
- `live_channels` (content_id FK, stream_url, is_live, viewer_count)
- `epg_events` (channel_id FK, title, start_time, end_time)
- `ad_campaigns` (targeting JSONB, budget, impressions_count)
- `ad_impressions` (campaign_id FK, user_id FK, completed, clicked)
- `playback_events` (content_id FK, event_type, position, quality, bitrate)
- `content_hotness` (hotness_score, view_count_24h/7d/30d, replication_tier)
- `notifications` (user_id FK, type, title, message, is_read)
- `admin_users` (email, role, permissions JSONB)
- `audit_logs` (admin_user_id FK, action, resource_type, old/new values JSONB)

### Extensions
- `uuid-ossp`: UUID generation
- `pg_trgm`: Trigram fuzzy text search
- `postgis`: Geo-location features

### Indexes
- B-tree on primary keys and foreign keys
- GIN indexes on array columns (genres, tags)
- Trigram GIN index on content title for fuzzy search
- Partial indexes on active subscriptions and published content
- Composite indexes on time-series queries (profile + timestamp DESC)

### Materialized Views
- `content_popularity`: Pre-aggregated view counts, completion rates, average ratings

### Triggers
- `update_updated_at_column()`: Auto-update `updated_at` on row modification

### Functions
- `update_content_hotness()`: Recalculate content hotness scores (called by cron)

---

## 9. Caching Strategy

### Redis Layers

| Cache Key Pattern | TTL | Purpose |
|------------------|-----|---------|
| `session:{sessionId}` | 24h | Playback session state |
| `recs:hybrid:{profileId}:{k}` | 15min | Hybrid recommendation results |
| `recs:ncf:{profileId}:{k}` | 1h | NCF model recommendations |
| `recs:cf:{profileId}:{k}` | 30min | Collaborative filtering results |
| `popular:{region}:{k}` | 10min | Trending content by region |
| `content:home:{profileId}` | 5min | Personalized home page |
| `rate_limit:{ip}` | 1min | API rate limiting counter |
| `auth:refresh:{token}` | 30d | Refresh token validity |

### CDN Edge Caching
- Video segments (.ts, .m4s): 1 year (immutable content-addressed)
- HLS/DASH manifests: 6 seconds (live), 1 hour (VoD)
- Thumbnails and posters: 1 year (content-addressed)
- API responses: Not cached at CDN (dynamic content)

---

## 10. Error Handling and Resilience

### Circuit Breaker Pattern
- Inter-service gRPC calls wrapped with circuit breaker (Hystrix/Resilience4j pattern)
- States: closed (normal) --> open (failing) --> half-open (testing recovery)
- Thresholds: 5 failures in 30 seconds opens circuit, 60 second recovery window

### Retry Strategy
- gRPC calls: 3 retries with exponential backoff (100ms, 200ms, 400ms)
- Kafka consumer: infinite retry with dead-letter queue after 5 failures
- External API calls (Stripe, SendGrid): 3 retries with jitter

### Graceful Degradation
- Recommendation service down: fall back to trending/popular content
- Search service down: fall back to PostgreSQL trigram search
- DRM service down: allow playback with watermark-only protection
- Analytics service down: buffer events in Kafka (no data loss)

---

## 11. API Gateway Configuration (Kong)

Kong sits at the ingress point and handles:
- JWT validation on all authenticated routes
- Rate limiting: 100 req/min (anonymous), 1000 req/min (authenticated)
- Request routing to upstream services
- SSL/TLS termination
- Request/response transformation
- CORS headers
- Request logging and analytics
- Plugin ecosystem: OAuth, JWT, CORS, rate-limiting, prometheus

---

## 12. Protobuf Service Definitions

Located in `packages/proto/`:
- `auth/v1/`: Auth service proto definitions
- `content/v1/`: Content service proto definitions
- `streaming/`: Streaming service proto definitions
- `payment/`: Payment service proto definitions
- Generated Go code in `packages/proto/gen/go/`

---

**Document Version**: 2.0
**Last Updated**: 2026-02-17
