# App Flow — StreamVerse Streaming Platform

## 1. Application Bootstrap Flow

The StreamVerse web application follows a gated initialization sequence defined in `App.tsx`:

```
[App Mount]
    |
    v
[Check Gemini API Key]
    |-- ENV var VITE_GEMINI_API_KEY set? --> Key selected = true
    |-- window.aistudio.hasSelectedApiKey()? --> Key selected = true
    |-- Neither? --> Show ApiKeyPromptView
    |
    v (key selected)
[Authentication Gate]
    |-- isAuthenticated = false? --> Show LoginView
    |-- isAuthenticated = true? --> Show Main App Shell
    |
    v (authenticated)
[Main App Shell]
    |-- Sidebar (navigation)
    |-- Header (logout, theme toggle, breadcrumbs)
    |-- Main content area (renders active view)
    |-- ChatBot (AI assistant overlay)
```

### State Management
- `currentViewState: ViewState` tracks the active view and optional params (e.g., `contentId`)
- `isAuthenticated: boolean` controls the auth gate
- `theme: Theme` persisted to localStorage, toggles dark/light classes on `<html>`
- `isKeySelected / hasCheckedKey` controls the API key initialization gate

---

## 2. View Routing Architecture

StreamVerse uses an in-memory state-based router (no URL-based routing library). Views are organized into three categories:

### Consumer Views
| View | Component | Description |
|------|-----------|-------------|
| `streamverse_home` | `StreamVerseHomeView` | Landing page with featured content, trending rows, continue watching |
| `watch` | `WatchView` | Video player with content ID parameter |

### Creator Views
| View | Component | Description |
|------|-----------|-------------|
| `creator_studio` | `CreatorStudioView` | Upload, manage, and monetize content |

### Admin Views
| View | Component | Description |
|------|-----------|-------------|
| `roadmap` | `RoadmapView` | Project roadmap and issue tracking |
| `overview` | `OverviewView` | Platform health dashboard |
| `gpu_fabric` | `GpuFabricView` | GPU transcoding infrastructure monitoring |
| `cdn` | `CdnView` | CDN analytics, cache ratios, latency |
| `security` | `SecurityView` | Security posture and threat dashboard |
| `media` | `MediaView` | Media processing pipeline status |
| `data` | `DataView` | Data platform metrics |
| `telecom` | `TelecomView` | Telecom service integration status |
| `satellite` | `SatelliteView` | Satellite delivery telemetry |
| `drm` | `DrmView` | DRM license server monitoring |
| `ai_ops` | `AiOpsView` | AI operations and autonomous actions |
| `neural_engine` | `NeuralContentEngineView` | Content analysis ML pipeline |
| `bi` | `BusinessIntelligenceView` | Churn prediction, content ROI |
| `broadcast_ops` | `BroadcastOpsView` | DVB-NIP, DVB-I, DVB-IP status |
| `user_profile` | `UserProfileView` | User settings, devices, preferences |

---

## 3. Authentication Flow

```
[User enters credentials in LoginView]
    |
    v
[POST /auth/login]
    |-- Request: { email, password }
    |-- Auth Service validates against PostgreSQL users table
    |-- Password hash verified (bcrypt)
    |
    v (success)
[JWT Token Pair Generated]
    |-- accessToken (1 hour expiry)
    |-- refreshToken (30 day expiry)
    |-- Response includes user profile and subscription tier
    |
    v
[Client stores tokens]
    |-- Access token in memory (React state)
    |-- Refresh token in httpOnly cookie (or secure storage on mobile)
    |
    v
[Subsequent API calls]
    |-- Authorization: Bearer {accessToken}
    |-- Kong API Gateway validates JWT
    |-- On 401: POST /auth/refresh with refreshToken
    |-- On refresh failure: redirect to LoginView
```

### OAuth 2.0 Social Login
```
[User clicks "Sign in with Google/Apple"]
    |
    v
[Redirect to OAuth provider]
    |-- Google: OAuth 2.0 + OpenID Connect
    |-- Apple: Sign in with Apple
    |
    v
[Provider callback with authorization code]
    |
    v
[POST /auth/oauth/callback]
    |-- Exchange code for tokens
    |-- Create or link user account
    |-- Generate StreamVerse JWT pair
```

---

## 4. Content Discovery Flow

```
[User lands on StreamVerseHomeView]
    |
    v
[GET /content/home?profileId={profileId}]
    |
    v
[Response: personalized home layout]
    |-- featured: Hero banner content
    |-- rows[]:
    |     |-- "Continue Watching" (from watch_history, position > 0, not completed)
    |     |-- "Trending Now" (from content_hotness, sorted by hotness_score)
    |     |-- "Recommended For You" (from recommendation-service hybrid algo)
    |     |-- "New Releases" (content sorted by published_at DESC)
    |     |-- Genre-specific rows (Action, Drama, Comedy, etc.)
    |     |-- "Top 10 in [Country]" (geo-filtered popularity)
    |     |-- "FAST Channels" (live_channels with type=fast_channel)
    |
    v
[User scrolls and interacts]
    |-- Click content card --> GET /content/{contentId} --> Detail overlay
    |-- Click "Play" --> Initiate playback flow
    |-- Click "Add to Watchlist" --> POST /users/me/watchlist/{contentId}
    |-- Search bar --> GET /search/autocomplete?q={query} --> GET /search?q={query}
```

---

## 5. Video Playback Flow

This is the critical path for the streaming platform:

```
[User clicks Play on content]
    |
    v
[GET /streaming/{contentId}/play?profileId={pid}&quality=auto]
    |
    v
[Streaming Service orchestration]
    |-- 1. Content Service: Verify content exists and is published
    |-- 2. Payment Service: Verify subscription entitlement
    |       |-- SVOD: Check subscription_tier >= content requirement
    |       |-- TVOD: Check transaction exists for this content
    |       |-- AVOD: Always allowed (ads will be inserted)
    |       |-- PPV: Check PPV transaction
    |-- 3. DRM Service: Generate license based on platform
    |       |-- Web (Chrome/Firefox): Widevine
    |       |-- Web (Safari): FairPlay
    |       |-- Android: Widevine
    |       |-- iOS/tvOS: FairPlay
    |       |-- Roku/Xbox: PlayReady
    |-- 4. CDN: Generate signed manifest URL (time-limited)
    |
    v
[Response to client]
    {
      streamUrl: "https://cdn.streamverse.io/streams/{id}/master.m3u8",
      protocol: "hls",
      qualities: [
        { quality: "720p", bitrate: 3000000, url: "...720p.m3u8" },
        { quality: "1080p", bitrate: 6000000, url: "...1080p.m3u8" },
        { quality: "4K", bitrate: 25000000, url: "...4k.m3u8" }
      ],
      drm: { type: "widevine", licenseUrl, certificateUrl },
      subtitles: [{ language, url }],
      sessionId: "sess_abc123"
    }
    |
    v
[Client-side player initialization]
    |-- HLS.js (web) or ExoPlayer (Android) or AVPlayer (iOS)
    |-- Initialize DRM module with license URL
    |-- Load master manifest (m3u8 or mpd)
    |-- ABR engine selects initial quality based on bandwidth estimate
    |-- Begin buffering and playback
    |
    v
[During playback]
    |-- Every 30s: POST /streaming/sessions/{sessionId}/heartbeat
    |     { currentTime, duration, quality, bufferHealth }
    |-- Heartbeat events published to Kafka topic: playback.events
    |-- Analytics Service consumes events for real-time dashboards
    |-- Watch history updated (position, duration)
    |
    v
[Playback completion]
    |-- POST /streaming/sessions/{sessionId}/complete
    |-- Watch history marked as completed
    |-- Recommendation engine weight adjusted
    |-- If series: auto-play next episode triggers
```

---

## 6. Content Upload and Transcoding Flow

```
[Creator uploads video via Creator Studio]
    |
    v
[POST /content/upload (multipart)]
    |-- Content Service validates metadata
    |-- Raw video uploaded to S3/MinIO (object storage)
    |-- Content record created in PostgreSQL (is_published=false)
    |
    v
[Kafka event: content.uploaded]
    |-- Transcoding Service consumes event
    |
    v
[Transcoding Pipeline]
    |-- 1. Probe: FFprobe analyzes source (codec, resolution, bitrate, audio)
    |-- 2. ABR Ladder: Generate quality variants
    |       720p @ 3 Mbps
    |       1080p @ 6 Mbps
    |       4K @ 25 Mbps
    |       Audio: AAC stereo + 5.1 surround
    |-- 3. GPU transcoding (NVIDIA NVENC on local or Runpod.io)
    |       H.264 for compatibility
    |       HEVC for efficiency
    |       AV1 for next-gen (optional)
    |-- 4. HLS packaging: Generate m3u8 manifests + .ts segments
    |-- 5. DASH packaging: Generate mpd manifests + .m4s segments
    |-- 6. Thumbnail extraction: Poster, thumbnails at intervals
    |-- 7. Subtitle processing: VTT conversion
    |
    v
[Upload transcoded assets to CDN origin (S3/MinIO)]
    |
    v
[Kafka event: content.transcoded]
    |-- Content Service updates stream_url, poster_url, thumbnail_url
    |-- Search Service indexes content in Elasticsearch
    |-- CDN cache warming for expected popular content
    |
    v
[Admin review (optional)]
    |-- Admin Service flags content for moderation
    |-- On approval: is_published = true, published_at = NOW()
    |
    v
[Content available to users]
```

---

## 7. Recommendation Engine Flow

```
[Triggered by: user login, playback complete, daily cron]
    |
    v
[Recommendation Service (Python/FastAPI)]
    |
    v
[Hybrid Strategy with weighted scoring]
    |-- NCF (Neural Collaborative Filtering) - weight: 0.4
    |     |-- PyTorch model trained on watch_history
    |     |-- User embedding + Item embedding --> MLP --> prediction score
    |-- Collaborative Filtering - weight: 0.3
    |     |-- Item-based co-occurrence from watch_history
    |     |-- "Users who watched X also watched Y"
    |-- Content-Based Filtering - weight: 0.2
    |     |-- Genre preference extraction from user history
    |     |-- Match unwatched content by genre overlap
    |-- Trending/Popular - weight: 0.1
    |     |-- content_hotness table (view_count_24h, 7d, 30d)
    |
    v
[Merge and rank by composite score]
    |
    v
[Cache results in Redis (15 min TTL)]
    |
    v
[Return top-K recommendations with content metadata]
```

---

## 8. Live Streaming Flow

```
[Broadcaster starts live stream]
    |
    v
[Ingest via RTMP/SRT/WebRTC to Ingestion Service]
    |-- Rust-based multi-protocol ingestion (streaming-saas/ingestion-service)
    |-- 10,000+ concurrent streams per node
    |
    v
[Real-time transcoding (GPU)]
    |-- ABR ladder generated in real-time
    |-- LL-HLS (Low-Latency HLS) packaging
    |-- WebRTC for sub-second latency path
    |
    v
[CDN distribution]
    |-- Edge caching on CloudFlare PoPs
    |-- P2P mesh augmentation (WebRTC) for bandwidth savings
    |
    v
[SSAI (Server-Side Ad Insertion)]
    |-- Ad-compositing-service analyzes scene boundaries
    |-- Dynamic ad stitching at natural break points
    |-- VAST/VMAP compliance
    |
    v
[EPG (Electronic Program Guide)]
    |-- epg_events table tracks schedule
    |-- Live channel metadata updated in real-time
    |-- Viewer count tracked via WebSocket heartbeats
```

---

## 9. Payment and Subscription Flow

```
[User selects subscription plan]
    |
    v
[Client collects payment method via Stripe Elements]
    |
    v
[POST /payments/subscribe]
    { planId: "plan_premium", paymentMethodId: "pm_card_visa" }
    |
    v
[Payment Service]
    |-- Create Stripe subscription
    |-- Create transaction record (status: pending)
    |-- On Stripe webhook (invoice.paid):
    |     |-- Update transaction status to completed
    |     |-- Update user subscription_tier and subscription_end_date
    |     |-- Publish Kafka event: subscription.activated
    |
    v
[Subscription active]
    |-- User can access content up to their tier's quality limit
    |-- Concurrent stream limit enforced by streaming-service
    |
    v
[Renewal / Cancellation]
    |-- Stripe handles automatic renewal
    |-- PUT /payments/subscriptions/{id} for plan changes
    |-- Cancellation schedules end at period end
```

---

## 10. Notification Delivery Flow

```
[Trigger event (new content, subscription, system alert)]
    |
    v
[Kafka event consumed by Notification Service]
    |
    v
[Fan-out to delivery channels]
    |-- Push Notification:
    |     |-- FCM (Firebase Cloud Messaging) for Android/Web
    |     |-- APNs (Apple Push Notification service) for iOS/tvOS
    |-- Email: SendGrid API
    |-- SMS: Twilio API (for critical alerts only)
    |-- In-App: WebSocket broadcast via websocket-service
    |
    v
[Notification record stored in PostgreSQL]
    |-- GET /notifications returns user's notification feed
    |-- PUT /notifications/{id}/read marks as read
```

---

## 11. Search Flow

```
[User types in search bar]
    |
    v
[GET /search/autocomplete?q={partial}]
    |-- Search Service queries Elasticsearch
    |-- Prefix matching on title, cast, director fields
    |-- Returns top 5 suggestions
    |
    v
[User submits full search]
    |
    v
[GET /search?q={query}&type=all&limit=20]
    |-- Elasticsearch full-text search
    |-- Trigram matching via pg_trgm (PostgreSQL fallback)
    |-- Faceted results: by type (movie/series), by genre
    |-- Relevance scoring with boost on title match
    |
    v
[Results rendered in grid/list view]
    |-- Each result links to content detail or playback
```

---

## 12. WebSocket Real-Time Flow

```
[Client connects: wss://api.streamverse.io/v1/ws?token={jwt}]
    |
    v
[WebSocket Service authenticates JWT]
    |
    v
[Client subscribes to channels]
    |-- content_updates: New content notifications
    |-- watch_party:{partyId}: Synchronized playback
    |-- notifications:{userId}: Personal notification feed
    |
    v
[Server pushes events]
    |-- Watch party sync: { currentTime, state: "playing"|"paused" }
    |-- New content: { contentId, title, type }
    |-- Live viewer count updates
    |-- Chat messages (watch party)
```

---

**Document Version**: 2.0
**Last Updated**: 2026-02-17
