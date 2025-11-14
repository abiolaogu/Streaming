-- ============================================================================
-- Global Streaming Platform - PostgreSQL Database Schema
-- ============================================================================

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- For fuzzy text search
CREATE EXTENSION IF NOT EXISTS "postgis"; -- For geo-location features

-- ============================================================================
-- USERS & AUTHENTICATION
-- ============================================================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    subscription_tier VARCHAR(50) DEFAULT 'free', -- free, basic, standard, premium, family
    subscription_status VARCHAR(50) DEFAULT 'active', -- active, cancelled, suspended, expired
    subscription_start_date TIMESTAMP,
    subscription_end_date TIMESTAMP,
    payment_method_id VARCHAR(255), -- Stripe payment method ID
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    email_verified BOOLEAN DEFAULT FALSE,
    phone VARCHAR(20),
    country_code CHAR(2), -- ISO 3166-1 alpha-2
    language VARCHAR(10) DEFAULT 'en',
    marketing_opt_in BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_subscription ON users(subscription_tier, subscription_status);

-- Profiles (multi-profile support)
CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    avatar_url VARCHAR(500),
    is_kids BOOLEAN DEFAULT FALSE,
    language VARCHAR(10) DEFAULT 'en',
    preferences JSONB DEFAULT '{}', -- autoplay, audio_language, subtitle_language, etc.
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT max_5_profiles CHECK (
        (SELECT COUNT(*) FROM profiles WHERE user_id = profiles.user_id) <= 5
    )
);

CREATE INDEX idx_profiles_user_id ON profiles(user_id);

-- ============================================================================
-- CONTENT MANAGEMENT
-- ============================================================================

CREATE TABLE content (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type VARCHAR(50) NOT NULL, -- movie, series, live_channel, fast_channel
    title VARCHAR(500) NOT NULL,
    original_title VARCHAR(500),
    description TEXT,
    long_description TEXT,
    duration INT, -- seconds (null for live channels)
    release_year INT,
    genres VARCHAR(50)[] DEFAULT '{}',
    rating VARCHAR(20), -- G, PG, PG-13, R, TV-MA, etc.
    imdb_id VARCHAR(20),
    tmdb_id VARCHAR(20),
    languages VARCHAR(10)[] DEFAULT '{}',
    subtitles VARCHAR(10)[] DEFAULT '{}',
    audio_tracks JSONB DEFAULT '[]', -- [{"language": "en", "type": "stereo|5.1|atmos"}]
    cast JSONB DEFAULT '[]', -- [{"name": "...", "role": "...", "character": "..."}]
    directors VARCHAR(255)[],
    producers VARCHAR(255)[],
    studios VARCHAR(255)[],
    countries VARCHAR(10)[],
    tags VARCHAR(100)[] DEFAULT '{}',
    stream_url VARCHAR(1000), -- HLS master manifest URL
    stream_type VARCHAR(50), -- hls, dash, webrtc
    poster_url VARCHAR(500),
    thumbnail_url VARCHAR(500),
    banner_url VARCHAR(500),
    trailer_url VARCHAR(500),
    metadata JSONB DEFAULT '{}', -- flexible metadata
    availability JSONB DEFAULT '{}', -- {"us": true, "uk": false, ...}
    monetization VARCHAR(50) DEFAULT 'svod', -- avod, svod, tvod, ppv
    rental_price DECIMAL(10, 2),
    purchase_price DECIMAL(10, 2),
    ppv_price DECIMAL(10, 2),
    is_published BOOLEAN DEFAULT FALSE,
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_content_type CHECK (type IN ('movie', 'series', 'episode', 'live_channel', 'fast_channel'))
);

CREATE INDEX idx_content_type ON content(type);
CREATE INDEX idx_content_title_trgm ON content USING gin(title gin_trgm_ops);
CREATE INDEX idx_content_genres ON content USING gin(genres);
CREATE INDEX idx_content_published ON content(is_published, published_at);
CREATE INDEX idx_content_monetization ON content(monetization);

-- Series and episodes
CREATE TABLE seasons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    series_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    season_number INT NOT NULL,
    title VARCHAR(255),
    description TEXT,
    poster_url VARCHAR(500),
    release_year INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (series_id, season_number)
);

CREATE INDEX idx_seasons_series ON seasons(series_id);

CREATE TABLE episodes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    season_id UUID NOT NULL REFERENCES seasons(id) ON DELETE CASCADE,
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    episode_number INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duration INT NOT NULL, -- seconds
    air_date DATE,
    thumbnail_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (season_id, episode_number)
);

CREATE INDEX idx_episodes_season ON episodes(season_id);
CREATE INDEX idx_episodes_content ON episodes(content_id);

-- ============================================================================
-- WATCH HISTORY & ENGAGEMENT
-- ============================================================================

CREATE TABLE watch_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    content_type VARCHAR(50) NOT NULL,
    position INT DEFAULT 0, -- playback position in seconds
    duration INT DEFAULT 0, -- total duration in seconds
    watched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed BOOLEAN DEFAULT FALSE,
    UNIQUE (profile_id, content_id)
);

CREATE INDEX idx_watch_history_profile ON watch_history(profile_id, watched_at DESC);
CREATE INDEX idx_watch_history_content ON watch_history(content_id);
CREATE INDEX idx_watch_history_completed ON watch_history(profile_id, completed);

-- Watchlist (favorites)
CREATE TABLE watchlist (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (profile_id, content_id)
);

CREATE INDEX idx_watchlist_profile ON watchlist(profile_id, added_at DESC);

-- Ratings and reviews
CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    review TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (profile_id, content_id)
);

CREATE INDEX idx_ratings_content ON ratings(content_id);

-- ============================================================================
-- SUBSCRIPTIONS & BILLING
-- ============================================================================

CREATE TABLE subscription_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE, -- free, basic, standard, premium, family
    price DECIMAL(10, 2) NOT NULL,
    currency CHAR(3) DEFAULT 'USD',
    billing_period VARCHAR(50) NOT NULL, -- monthly, annual
    features JSONB DEFAULT '{}', -- max_streams, quality, downloads, ads, etc.
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO subscription_plans (name, price, billing_period, features) VALUES
('free', 0.00, 'monthly', '{"max_streams": 1, "quality": "720p", "ads": true, "downloads": false}'),
('basic', 4.99, 'monthly', '{"max_streams": 1, "quality": "1080p", "ads": true, "downloads": false}'),
('standard', 9.99, 'monthly', '{"max_streams": 2, "quality": "1080p", "ads": false, "downloads": true}'),
('premium', 14.99, 'monthly', '{"max_streams": 4, "quality": "4k", "ads": false, "downloads": true}'),
('family', 19.99, 'monthly', '{"max_streams": 6, "quality": "4k", "ads": false, "downloads": true, "profiles": 5}');

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- subscription, rental, purchase, ppv
    amount DECIMAL(10, 2) NOT NULL,
    currency CHAR(3) DEFAULT 'USD',
    payment_method VARCHAR(50), -- stripe, paypal, apple_pay, google_pay
    payment_provider_id VARCHAR(255), -- Stripe transaction ID
    status VARCHAR(50) DEFAULT 'pending', -- pending, completed, failed, refunded
    content_id UUID REFERENCES content(id), -- for TVOD/PPV
    subscription_plan_id UUID REFERENCES subscription_plans(id),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transactions_user ON transactions(user_id, created_at DESC);
CREATE INDEX idx_transactions_status ON transactions(status);

-- ============================================================================
-- LIVE CHANNELS & SCHEDULING
-- ============================================================================

CREATE TABLE live_channels (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    channel_number INT UNIQUE,
    stream_url VARCHAR(1000) NOT NULL,
    ingest_url VARCHAR(1000),
    is_live BOOLEAN DEFAULT FALSE,
    viewer_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_live_channels_live ON live_channels(is_live);

-- EPG (Electronic Program Guide) for live channels
CREATE TABLE epg_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    channel_id UUID NOT NULL REFERENCES live_channels(id) ON DELETE CASCADE,
    content_id UUID REFERENCES content(id), -- linked VOD content (for FAST channels)
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    duration INT NOT NULL, -- seconds
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_time_range CHECK (end_time > start_time)
);

CREATE INDEX idx_epg_events_channel_time ON epg_events(channel_id, start_time, end_time);
CREATE INDEX idx_epg_events_start ON epg_events(start_time);

-- ============================================================================
-- ADVERTISING
-- ============================================================================

CREATE TABLE ad_campaigns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    advertiser VARCHAR(255) NOT NULL,
    video_url VARCHAR(500) NOT NULL,
    click_through_url VARCHAR(500),
    duration INT NOT NULL, -- seconds
    targeting JSONB DEFAULT '{}', -- {"age_range": [18, 35], "genres": ["action"], ...}
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    budget DECIMAL(10, 2),
    cost_per_impression DECIMAL(10, 4),
    impressions_count INT DEFAULT 0,
    clicks_count INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ad_campaigns_active ON ad_campaigns(is_active, start_date, end_date);

CREATE TABLE ad_impressions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID NOT NULL REFERENCES ad_campaigns(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id),
    profile_id UUID REFERENCES profiles(id),
    content_id UUID REFERENCES content(id),
    impression_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed BOOLEAN DEFAULT FALSE,
    clicked BOOLEAN DEFAULT FALSE,
    metadata JSONB DEFAULT '{}'
);

CREATE INDEX idx_ad_impressions_campaign ON ad_impressions(campaign_id, impression_time);

-- ============================================================================
-- ANALYTICS & METRICS (Write to Cassandra in production)
-- ============================================================================

CREATE TABLE playback_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    profile_id UUID REFERENCES profiles(id),
    content_id UUID NOT NULL REFERENCES content(id),
    event_type VARCHAR(50) NOT NULL, -- play, pause, seek, buffering, error, completed
    position INT, -- playback position in seconds
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    session_id UUID,
    device_type VARCHAR(50), -- web, mobile, tv
    platform VARCHAR(50), -- ios, android, roku, etc.
    quality VARCHAR(20), -- 480p, 720p, 1080p, 4k
    bitrate INT, -- kbps
    cdn_edge VARCHAR(100),
    buffer_duration INT, -- ms
    error_code VARCHAR(50),
    metadata JSONB DEFAULT '{}'
);

CREATE INDEX idx_playback_events_content ON playback_events(content_id, timestamp);
CREATE INDEX idx_playback_events_user ON playback_events(user_id, timestamp);
CREATE INDEX idx_playback_events_session ON playback_events(session_id);

-- Partitioning strategy for time-series data (monthly partitions)
-- In production: Use TimescaleDB or move to Cassandra

-- ============================================================================
-- CDN & CACHING (Sync with Aerospike)
-- ============================================================================

CREATE TABLE content_hotness (
    content_id UUID PRIMARY KEY REFERENCES content(id) ON DELETE CASCADE,
    hotness_score FLOAT DEFAULT 0.0, -- 0.0 to 1.0
    view_count_24h INT DEFAULT 0,
    view_count_7d INT DEFAULT 0,
    view_count_30d INT DEFAULT 0,
    last_viewed_at TIMESTAMP,
    ttl_hint INT DEFAULT 3600, -- seconds, for CDN cache
    replication_tier INT DEFAULT 3, -- 1=global, 2=regional, 3=on-demand
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_content_hotness_score ON content_hotness(hotness_score DESC);
CREATE INDEX idx_content_hotness_tier ON content_hotness(replication_tier);

-- ============================================================================
-- NOTIFICATIONS & MESSAGING
-- ============================================================================

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- new_content, subscription, payment, system
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    link VARCHAR(500),
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user ON notifications(user_id, created_at DESC);
CREATE INDEX idx_notifications_unread ON notifications(user_id, is_read);

-- ============================================================================
-- ADMIN & CMS
-- ============================================================================

CREATE TABLE admin_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL, -- super_admin, admin, editor, viewer
    permissions JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    admin_user_id UUID REFERENCES admin_users(id),
    action VARCHAR(100) NOT NULL, -- create, update, delete, publish
    resource_type VARCHAR(50) NOT NULL, -- content, user, subscription
    resource_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_admin ON audit_logs(admin_user_id, created_at DESC);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);

-- ============================================================================
-- FUNCTIONS & TRIGGERS
-- ============================================================================

-- Automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply trigger to all relevant tables
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_profiles_updated_at BEFORE UPDATE ON profiles
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_content_updated_at BEFORE UPDATE ON content
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Update hotness scores (call this periodically via cron job)
CREATE OR REPLACE FUNCTION update_content_hotness()
RETURNS void AS $$
BEGIN
    UPDATE content_hotness ch
    SET
        view_count_24h = (
            SELECT COUNT(*) FROM playback_events
            WHERE content_id = ch.content_id
            AND event_type = 'play'
            AND timestamp > CURRENT_TIMESTAMP - INTERVAL '24 hours'
        ),
        view_count_7d = (
            SELECT COUNT(*) FROM playback_events
            WHERE content_id = ch.content_id
            AND event_type = 'play'
            AND timestamp > CURRENT_TIMESTAMP - INTERVAL '7 days'
        ),
        view_count_30d = (
            SELECT COUNT(*) FROM playback_events
            WHERE content_id = ch.content_id
            AND event_type = 'play'
            AND timestamp > CURRENT_TIMESTAMP - INTERVAL '30 days'
        ),
        hotness_score = LEAST(1.0,
            (view_count_24h * 0.5 + view_count_7d * 0.3 + view_count_30d * 0.2) / 10000.0
        ),
        replication_tier = CASE
            WHEN hotness_score >= 0.8 THEN 1  -- Top 1% - global replication
            WHEN hotness_score >= 0.5 THEN 2  -- Top 10% - regional replication
            ELSE 3                             -- Long tail - on-demand
        END,
        updated_at = CURRENT_TIMESTAMP;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- MATERIALIZED VIEWS (for analytics)
-- ============================================================================

CREATE MATERIALIZED VIEW content_popularity AS
SELECT
    c.id,
    c.title,
    c.type,
    c.genres,
    COUNT(DISTINCT wh.profile_id) AS unique_viewers,
    COUNT(wh.id) AS total_views,
    AVG(CASE WHEN wh.completed THEN 1.0 ELSE 0.0 END) AS completion_rate,
    AVG(r.rating) AS avg_rating,
    COUNT(wl.id) AS watchlist_count
FROM content c
LEFT JOIN watch_history wh ON c.id = wh.content_id
LEFT JOIN ratings r ON c.id = r.content_id
LEFT JOIN watchlist wl ON c.id = wl.content_id
WHERE c.is_published = TRUE
GROUP BY c.id, c.title, c.type, c.genres;

CREATE UNIQUE INDEX idx_content_popularity_id ON content_popularity(id);

-- Refresh daily
-- REFRESH MATERIALIZED VIEW CONCURRENTLY content_popularity;

-- ============================================================================
-- INITIAL DATA SEEDING (Optional)
-- ============================================================================

-- Create default admin user (password: admin123 - change in production!)
INSERT INTO admin_users (email, password_hash, role)
VALUES ('admin@example.com', '$2a$10$XtYN8.wZ7Y.jQJmZqWVd5uKZ3j8pQZ3l3Zjm5Yl0hYZm5Yl0hYZm5Y', 'super_admin');

-- ============================================================================
-- PERFORMANCE OPTIMIZATION QUERIES
-- ============================================================================

-- Analyze tables for query optimization
ANALYZE users;
ANALYZE profiles;
ANALYZE content;
ANALYZE watch_history;
ANALYZE playback_events;

-- Vacuum to reclaim storage
VACUUM ANALYZE users;
VACUUM ANALYZE content;
VACUUM ANALYZE watch_history;
