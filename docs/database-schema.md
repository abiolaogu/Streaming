# Database Schema — StreamVerse Streaming Platform

## 1. Database Strategy

StreamVerse uses a polyglot persistence strategy with PostgreSQL as the primary relational datastore:

| Database | Purpose | Location |
|----------|---------|----------|
| PostgreSQL 14+ | Primary relational data (users, content, billing) | `database-schema.sql` |
| Redis 7+ | Caching, sessions, rate limiting | In-memory |
| Elasticsearch | Full-text search, content discovery | Search indices |
| ScyllaDB | Time-series playback events (planned) | Analytics pipeline |
| MongoDB | Session data, user preferences (planned) | Document store |

---

## 2. PostgreSQL Extensions

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";   -- UUID v4 generation
CREATE EXTENSION IF NOT EXISTS "pg_trgm";     -- Trigram fuzzy text search
CREATE EXTENSION IF NOT EXISTS "postgis";     -- Geo-location features
```

---

## 3. Core Tables

### 3.1 Users and Authentication

#### users
Primary user account table. Each user may have multiple profiles and one subscription.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PK, DEFAULT uuid_generate_v4() | Unique user identifier |
| email | VARCHAR(255) | UNIQUE NOT NULL | Login email |
| username | VARCHAR(100) | UNIQUE NOT NULL | Display username |
| password_hash | VARCHAR(255) | NOT NULL | bcrypt hashed password |
| subscription_tier | VARCHAR(50) | DEFAULT 'free' | free, basic, standard, premium, family |
| subscription_status | VARCHAR(50) | DEFAULT 'active' | active, cancelled, suspended, expired |
| subscription_start_date | TIMESTAMP | | Subscription period start |
| subscription_end_date | TIMESTAMP | | Subscription period end |
| payment_method_id | VARCHAR(255) | | Stripe payment method reference |
| country_code | CHAR(2) | | ISO 3166-1 alpha-2 country code |
| language | VARCHAR(10) | DEFAULT 'en' | Preferred language |
| email_verified | BOOLEAN | DEFAULT FALSE | Email verification status |
| is_active | BOOLEAN | DEFAULT TRUE | Account active flag |
| marketing_opt_in | BOOLEAN | DEFAULT FALSE | Marketing consent |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Account creation |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last modification |
| last_login_at | TIMESTAMP | | Most recent login |

**Indexes**: `idx_users_email` (email), `idx_users_subscription` (tier, status)

#### profiles
Multi-profile support (up to 5 per account). Each profile has independent preferences and watch history.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PK | Profile identifier |
| user_id | UUID | FK users(id) CASCADE | Parent account |
| name | VARCHAR(100) | NOT NULL | Profile display name |
| avatar_url | VARCHAR(500) | | Avatar image URL |
| is_kids | BOOLEAN | DEFAULT FALSE | Kids-safe mode filter |
| language | VARCHAR(10) | DEFAULT 'en' | Profile language |
| preferences | JSONB | DEFAULT '{}' | autoplay, audio_language, subtitle_language |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |

**Constraint**: Maximum 5 profiles per user_id
**Index**: `idx_profiles_user_id` (user_id)

---

### 3.2 Content Management

#### content
Central content metadata table supporting movies, series, episodes, and live channels.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PK | Content identifier |
| type | VARCHAR(50) | NOT NULL, CHECK | movie, series, episode, live_channel, fast_channel |
| title | VARCHAR(500) | NOT NULL | Display title |
| original_title | VARCHAR(500) | | Original language title |
| description | TEXT | | Short synopsis |
| long_description | TEXT | | Full description |
| duration | INT | | Duration in seconds (null for live) |
| release_year | INT | | Year of release |
| genres | VARCHAR(50)[] | DEFAULT '{}' | Array of genre tags |
| rating | VARCHAR(20) | | Age rating (G, PG, PG-13, R, TV-MA) |
| imdb_id | VARCHAR(20) | | IMDB cross-reference |
| tmdb_id | VARCHAR(20) | | TMDB cross-reference |
| languages | VARCHAR(10)[] | DEFAULT '{}' | Available audio languages |
| subtitles | VARCHAR(10)[] | DEFAULT '{}' | Available subtitle languages |
| audio_tracks | JSONB | DEFAULT '[]' | Audio track details (language, stereo/5.1/atmos) |
| cast | JSONB | DEFAULT '[]' | Cast list (name, role, character) |
| directors | VARCHAR(255)[] | | Director names |
| stream_url | VARCHAR(1000) | | HLS master manifest URL |
| stream_type | VARCHAR(50) | | hls, dash, webrtc |
| poster_url | VARCHAR(500) | | Poster image URL |
| thumbnail_url | VARCHAR(500) | | Thumbnail image URL |
| banner_url | VARCHAR(500) | | Wide banner image URL |
| trailer_url | VARCHAR(500) | | Trailer manifest URL |
| metadata | JSONB | DEFAULT '{}' | Extensible metadata |
| availability | JSONB | DEFAULT '{}' | Geo-availability map {"us": true, "uk": false} |
| monetization | VARCHAR(50) | DEFAULT 'svod' | avod, svod, tvod, ppv |
| rental_price | DECIMAL(10,2) | | TVOD rental price |
| purchase_price | DECIMAL(10,2) | | TVOD purchase price |
| ppv_price | DECIMAL(10,2) | | PPV event price |
| is_published | BOOLEAN | DEFAULT FALSE | Publication status |
| published_at | TIMESTAMP | | Publication timestamp |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | |

**Indexes**:
- `idx_content_type` (type)
- `idx_content_title_trgm` GIN (title gin_trgm_ops) for fuzzy search
- `idx_content_genres` GIN (genres) for array contains queries
- `idx_content_published` (is_published, published_at)
- `idx_content_monetization` (monetization)

#### seasons

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PK | Season identifier |
| series_id | UUID | FK content(id) CASCADE | Parent series |
| season_number | INT | NOT NULL | Season ordinal |
| title | VARCHAR(255) | | Season title |
| description | TEXT | | Season description |
| poster_url | VARCHAR(500) | | Season poster |
| release_year | INT | | Season release year |

**Unique**: (series_id, season_number)

#### episodes

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PK | Episode identifier |
| season_id | UUID | FK seasons(id) CASCADE | Parent season |
| content_id | UUID | FK content(id) CASCADE | Linked content record |
| episode_number | INT | NOT NULL | Episode ordinal |
| title | VARCHAR(255) | NOT NULL | Episode title |
| description | TEXT | | Episode synopsis |
| duration | INT | NOT NULL | Duration in seconds |
| air_date | DATE | | Original air date |
| thumbnail_url | VARCHAR(500) | | Episode thumbnail |

**Unique**: (season_id, episode_number)

---

### 3.3 User Engagement

#### watch_history
Tracks per-profile viewing progress for continue-watching and completion tracking.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PK | |
| profile_id | UUID | FK profiles(id) CASCADE | Viewing profile |
| content_id | UUID | FK content(id) CASCADE | Watched content |
| content_type | VARCHAR(50) | NOT NULL | Content type at time of watch |
| position | INT | DEFAULT 0 | Playback position in seconds |
| duration | INT | DEFAULT 0 | Total duration in seconds |
| watched_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Last watch timestamp |
| completed | BOOLEAN | DEFAULT FALSE | Completion flag |

**Unique**: (profile_id, content_id) -- upsert on re-watch
**Indexes**: Profile+time DESC, content_id, profile+completed

#### watchlist

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | PK |
| profile_id | UUID | FK profiles(id) |
| content_id | UUID | FK content(id) |
| added_at | TIMESTAMP | When added |

**Unique**: (profile_id, content_id)

#### ratings

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PK | |
| profile_id | UUID | FK profiles(id) | Rating profile |
| content_id | UUID | FK content(id) | Rated content |
| rating | INT | CHECK (1-5) | Star rating |
| review | TEXT | | Optional text review |

**Unique**: (profile_id, content_id)

---

### 3.4 Subscriptions and Billing

#### subscription_plans

| Plan | Price | Features |
|------|-------|----------|
| free | $0.00 | 1 stream, 720p, ads, no downloads |
| basic | $4.99 | 1 stream, 1080p, ads, no downloads |
| standard | $9.99 | 2 streams, 1080p, no ads, downloads |
| premium | $14.99 | 4 streams, 4K, no ads, downloads |
| family | $19.99 | 6 streams, 4K, no ads, downloads, 5 profiles |

#### transactions
Records all financial transactions (subscription, rental, purchase, PPV).

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | PK |
| user_id | UUID | FK users(id) |
| type | VARCHAR(50) | subscription, rental, purchase, ppv |
| amount | DECIMAL(10,2) | Transaction amount |
| currency | CHAR(3) | Currency code (USD) |
| payment_method | VARCHAR(50) | stripe, paypal, apple_pay, google_pay |
| payment_provider_id | VARCHAR(255) | Stripe transaction ID |
| status | VARCHAR(50) | pending, completed, failed, refunded |
| content_id | UUID | FK content(id) for TVOD/PPV |
| subscription_plan_id | UUID | FK subscription_plans(id) |
| metadata | JSONB | Additional transaction data |

---

### 3.5 Live Channels and EPG

#### live_channels

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | PK |
| content_id | UUID | FK content(id) |
| channel_number | INT | UNIQUE channel number |
| stream_url | VARCHAR(1000) | Live stream URL |
| ingest_url | VARCHAR(1000) | Ingest endpoint |
| is_live | BOOLEAN | Currently broadcasting |
| viewer_count | INT | Real-time viewer count |

#### epg_events (Electronic Program Guide)

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | PK |
| channel_id | UUID | FK live_channels(id) |
| content_id | UUID | FK content(id) optional (for FAST) |
| title | VARCHAR(255) | Program title |
| start_time | TIMESTAMP | Program start |
| end_time | TIMESTAMP | Program end |
| duration | INT | Duration in seconds |

**Constraint**: end_time > start_time
**Index**: (channel_id, start_time, end_time) for schedule queries

---

### 3.6 Advertising

#### ad_campaigns
Manages ad inventory with targeting, budgets, and performance metrics.

#### ad_impressions
Tracks individual ad views with completion and click-through data.

---

### 3.7 Analytics

#### playback_events
High-volume event log for player telemetry.

| Column | Type | Description |
|--------|------|-------------|
| event_type | VARCHAR(50) | play, pause, seek, buffering, error, completed |
| quality | VARCHAR(20) | 480p, 720p, 1080p, 4k |
| bitrate | INT | Current bitrate (kbps) |
| device_type | VARCHAR(50) | web, mobile, tv |
| platform | VARCHAR(50) | ios, android, roku, tizen, etc. |
| cdn_edge | VARCHAR(100) | CDN edge server identifier |
| buffer_duration | INT | Buffer time in milliseconds |
| error_code | VARCHAR(50) | Player error code |

**Note**: In production, this table should be partitioned by month or migrated to ScyllaDB for time-series performance.

#### content_hotness
Pre-computed content popularity scores used for CDN cache tiering and trending content.

| Column | Type | Description |
|--------|------|-------------|
| content_id | UUID | PK, FK content(id) |
| hotness_score | FLOAT | 0.0 to 1.0 normalized score |
| view_count_24h | INT | Views in last 24 hours |
| view_count_7d | INT | Views in last 7 days |
| view_count_30d | INT | Views in last 30 days |
| replication_tier | INT | 1=global, 2=regional, 3=on-demand |
| ttl_hint | INT | Suggested CDN cache TTL in seconds |

**Score calculation**: `hotness = min(1.0, (24h*0.5 + 7d*0.3 + 30d*0.2) / 10000)`
**Tiering**: score >= 0.8 = global (top 1%), >= 0.5 = regional (top 10%), else on-demand

---

### 3.8 Administration

#### admin_users
Separate admin user table with role-based access.

#### audit_logs
Complete audit trail for all administrative actions with before/after JSONB snapshots.

---

## 4. Materialized Views

### content_popularity
Pre-aggregated analytics for content performance dashboards:
- Unique viewers per content
- Total view count
- Completion rate (percentage who finish)
- Average rating
- Watchlist count

Refreshed daily via `REFRESH MATERIALIZED VIEW CONCURRENTLY content_popularity`.

---

## 5. Database Functions and Triggers

### update_updated_at_column()
Trigger function automatically setting `updated_at = CURRENT_TIMESTAMP` on any row update. Applied to: users, profiles, content.

### update_content_hotness()
Batch function recalculating hotness scores and CDN replication tiers. Called periodically by scheduler-service cron job.

---

## 6. Migration Strategy

Schema versioning should use golang-migrate with numbered SQL migration files:
```
migrations/
  000001_create_users.up.sql
  000001_create_users.down.sql
  000002_create_content.up.sql
  ...
```

---

**Document Version**: 2.0
**Last Updated**: 2026-02-17
