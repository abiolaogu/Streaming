// ============================================================================
// User Service - Handles user profiles, preferences, watch history
// ============================================================================

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// User model
type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	SubscriptionTier string `json:"subscription_tier"`
	Profiles      []Profile `json:"profiles"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Profile struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Name       string    `json:"name"`
	AvatarURL  string    `json:"avatar_url"`
	IsKids     bool      `json:"is_kids"`
	Language   string    `json:"language"`
	Preferences map[string]interface{} `json:"preferences"`
}

type WatchHistoryItem struct {
	ProfileID   string    `json:"profile_id"`
	ContentID   string    `json:"content_id"`
	ContentType string    `json:"content_type"` // "movie", "episode", "live"
	Position    int       `json:"position"`      // seconds
	Duration    int       `json:"duration"`
	WatchedAt   time.Time `json:"watched_at"`
	Completed   bool      `json:"completed"`
}

// UserService handles all user-related operations
type UserService struct {
	db    *sql.DB
	redis *redis.Client
	
	// Prometheus metrics
	requestCount  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewUserService(db *sql.DB, redisClient *redis.Client) *UserService {
	return &UserService{
		db:    db,
		redis: redisClient,
		requestCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "user_service_requests_total",
				Help: "Total number of requests to user service",
			},
			[]string{"method", "endpoint", "status"},
		),
		requestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "user_service_request_duration_seconds",
				Help: "Request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
	}
}

// GetUser retrieves user by ID with caching
func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
	start := time.Now()
	defer func() {
		s.requestDuration.WithLabelValues("GET", "/users").Observe(time.Since(start).Seconds())
	}()
	
	// Try cache first
	cacheKey := "user:" + userID
	cachedData, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var user User
		if err := json.Unmarshal([]byte(cachedData), &user); err == nil {
			s.requestCount.WithLabelValues("GET", "/users", "200").Inc()
			return &user, nil
		}
	}
	
	// Cache miss - query database
	var user User
	query := `
		SELECT id, email, username, subscription_tier, created_at, updated_at
		FROM users WHERE id = $1
	`
	err = s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Email, &user.Username, &user.SubscriptionTier,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		s.requestCount.WithLabelValues("GET", "/users", "404").Inc()
		return nil, err
	}
	
	// Load profiles
	user.Profiles, err = s.GetProfiles(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	// Cache the result (TTL: 5 minutes)
	userData, _ := json.Marshal(user)
	s.redis.Set(ctx, cacheKey, userData, 5*time.Minute)
	
	s.requestCount.WithLabelValues("GET", "/users", "200").Inc()
	return &user, nil
}

// GetProfiles retrieves all profiles for a user
func (s *UserService) GetProfiles(ctx context.Context, userID string) ([]Profile, error) {
	query := `
		SELECT id, user_id, name, avatar_url, is_kids, language, preferences
		FROM profiles WHERE user_id = $1 ORDER BY created_at
	`
	
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var profiles []Profile
	for rows.Next() {
		var p Profile
		var prefsJSON []byte
		err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.AvatarURL, &p.IsKids, &p.Language, &prefsJSON)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(prefsJSON, &p.Preferences)
		profiles = append(profiles, p)
	}
	
	return profiles, nil
}

// UpdateWatchHistory updates or creates watch history entry
func (s *UserService) UpdateWatchHistory(ctx context.Context, item *WatchHistoryItem) error {
	query := `
		INSERT INTO watch_history (profile_id, content_id, content_type, position, duration, watched_at, completed)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (profile_id, content_id)
		DO UPDATE SET position = $4, duration = $5, watched_at = $6, completed = $7
	`
	
	_, err := s.db.ExecContext(ctx, query,
		item.ProfileID, item.ContentID, item.ContentType,
		item.Position, item.Duration, item.WatchedAt, item.Completed,
	)
	
	if err == nil {
		// Invalidate cache for continue watching
		s.redis.Del(ctx, "continue_watching:"+item.ProfileID)
		
		// Publish event to Kafka for analytics
		// (Implementation depends on your Kafka client)
	}
	
	return err
}

// GetContinueWatching retrieves in-progress content for a profile
func (s *UserService) GetContinueWatching(ctx context.Context, profileID string, limit int) ([]WatchHistoryItem, error) {
	cacheKey := "continue_watching:" + profileID
	
	// Try cache
	cachedData, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var items []WatchHistoryItem
		if err := json.Unmarshal([]byte(cachedData), &items); err == nil {
			return items, nil
		}
	}
	
	// Query database
	query := `
		SELECT profile_id, content_id, content_type, position, duration, watched_at, completed
		FROM watch_history
		WHERE profile_id = $1 AND completed = false AND position > 0
		ORDER BY watched_at DESC
		LIMIT $2
	`
	
	rows, err := s.db.QueryContext(ctx, query, profileID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var items []WatchHistoryItem
	for rows.Next() {
		var item WatchHistoryItem
		err := rows.Scan(&item.ProfileID, &item.ContentID, &item.ContentType,
			&item.Position, &item.Duration, &item.WatchedAt, &item.Completed)
		if err != nil {
			continue
		}
		items = append(items, item)
	}
	
	// Cache results (TTL: 10 minutes)
	itemsData, _ := json.Marshal(items)
	s.redis.Set(ctx, cacheKey, itemsData, 10*time.Minute)
	
	return items, nil
}

// HTTP Handlers
func (s *UserService) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	
	user, err := s.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (s *UserService) UpdateWatchHistoryHandler(w http.ResponseWriter, r *http.Request) {
	var item WatchHistoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	item.WatchedAt = time.Now()
	
	if err := s.UpdateWatchHistory(r.Context(), &item); err != nil {
		http.Error(w, "Failed to update watch history", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

func (s *UserService) GetContinueWatchingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileID := vars["profile_id"]
	
	items, err := s.GetContinueWatching(r.Context(), profileID, 20)
	if err != nil {
		http.Error(w, "Failed to fetch continue watching", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"items": items,
	})
}

// ============================================================================
// Content Service - Handles metadata, catalogs, search
// ============================================================================

type Content struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"` // "movie", "series", "live_channel"
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Duration    int               `json:"duration"` // seconds
	ReleaseYear int               `json:"release_year"`
	Genres      []string          `json:"genres"`
	Rating      string            `json:"rating"` // "G", "PG", "PG-13", "R", "TV-MA"
	Languages   []string          `json:"languages"`
	Subtitles   []string          `json:"subtitles"`
	StreamURL   string            `json:"stream_url"`
	PosterURL   string            `json:"poster_url"`
	ThumbnailURL string           `json:"thumbnail_url"`
	Metadata    map[string]interface{} `json:"metadata"`
	Availability map[string]bool  `json:"availability"` // region -> available
}

type ContentService struct {
	db    *sql.DB
	redis *redis.Client
}

func NewContentService(db *sql.DB, redisClient *redis.Client) *ContentService {
	return &ContentService{
		db:    db,
		redis: redisClient,
	}
}

// GetContent retrieves content by ID with caching
func (s *ContentService) GetContent(ctx context.Context, contentID string) (*Content, error) {
	cacheKey := "content:" + contentID
	
	// Try cache
	cachedData, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var content Content
		if err := json.Unmarshal([]byte(cachedData), &content); err == nil {
			return &content, nil
		}
	}
	
	// Query database
	var content Content
	var genresArray pq.StringArray
	var languagesArray pq.StringArray
	var subtitlesArray pq.StringArray
	var availabilityJSON []byte
	var metadataJSON []byte
	
	query := `
		SELECT id, type, title, description, duration, release_year, genres,
		       rating, languages, subtitles, stream_url, poster_url, thumbnail_url,
		       metadata, availability
		FROM content WHERE id = $1
	`
	
	err = s.db.QueryRowContext(ctx, query, contentID).Scan(
		&content.ID, &content.Type, &content.Title, &content.Description,
		&content.Duration, &content.ReleaseYear, &genresArray, &content.Rating,
		&languagesArray, &subtitlesArray, &content.StreamURL, &content.PosterURL,
		&content.ThumbnailURL, &metadataJSON, &availabilityJSON,
	)
	
	if err != nil {
		return nil, err
	}
	
	content.Genres = genresArray
	content.Languages = languagesArray
	content.Subtitles = subtitlesArray
	json.Unmarshal(metadataJSON, &content.Metadata)
	json.Unmarshal(availabilityJSON, &content.Availability)
	
	// Cache result (TTL: 1 hour)
	contentData, _ := json.Marshal(content)
	s.redis.Set(ctx, cacheKey, contentData, 1*time.Hour)
	
	return &content, nil
}

// SearchContent performs full-text search on content
func (s *ContentService) SearchContent(ctx context.Context, query string, limit int) ([]Content, error) {
	// In production, use Elasticsearch for better performance
	sqlQuery := `
		SELECT id, type, title, description, poster_url, rating
		FROM content
		WHERE to_tsvector('english', title || ' ' || description) @@ plainto_tsquery('english', $1)
		ORDER BY ts_rank(to_tsvector('english', title || ' ' || description), plainto_tsquery('english', $1)) DESC
		LIMIT $2
	`
	
	rows, err := s.db.QueryContext(ctx, sqlQuery, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var results []Content
	for rows.Next() {
		var c Content
		rows.Scan(&c.ID, &c.Type, &c.Title, &c.Description, &c.PosterURL, &c.Rating)
		results = append(results, c)
	}
	
	return results, nil
}

// ============================================================================
// Stream Service - Handles manifest generation and playback
// ============================================================================

type StreamService struct {
	contentService *ContentService
	redis          *redis.Client
	cdnBaseURL     string
}

func NewStreamService(contentService *ContentService, redisClient *redis.Client, cdnBaseURL string) *StreamService {
	return &StreamService{
		contentService: contentService,
		redis:          redisClient,
		cdnBaseURL:     cdnBaseURL,
	}
}

// GeneratePlaybackURL generates signed playback URL with CDN selection
func (s *StreamService) GeneratePlaybackURL(ctx context.Context, contentID, userID, profileID, region string) (string, error) {
	// 1. Verify entitlement (subscription/purchase)
	// (Call to subscription service)
	
	// 2. Get content metadata
	content, err := s.contentService.GetContent(ctx, contentID)
	if err != nil {
		return "", err
	}
	
	// 3. Check regional availability
	if available, ok := content.Availability[region]; !ok || !available {
		return "", &ErrContentNotAvailable{Region: region}
	}
	
	// 4. Select optimal CDN edge (ML-based in production)
	edge := s.selectOptimalEdge(ctx, region, userID)
	
	// 5. Generate signed URL (token valid for 6 hours)
	token := s.generatePlaybackToken(userID, profileID, contentID, 6*time.Hour)
	
	// HLS manifest URL
	playbackURL := fmt.Sprintf("%s/%s/master.m3u8?token=%s", edge, contentID, token)
	
	// 6. Log playback start event (async to Kafka)
	go s.logPlaybackEvent(contentID, userID, profileID, "play_start")
	
	return playbackURL, nil
}

func (s *StreamService) selectOptimalEdge(ctx context.Context, region, userID string) string {
	// Simplified edge selection (in production: ML model predicts best edge)
	edgeMap := map[string]string{
		"eu-west":      "https://eu-west-edge.cdn.example.com",
		"eu-central":   "https://eu-central-edge.cdn.example.com",
		"us-east":      "https://us-east-edge.cdn.example.com",
		"us-west":      "https://us-west-edge.cdn.example.com",
		"af-west":      "https://af-west-edge.cdn.example.com",
		"ap-southeast": "https://ap-southeast-edge.cdn.example.com",
		"sa-east":      "https://sa-east-edge.cdn.example.com",
	}
	
	if edge, ok := edgeMap[region]; ok {
		return edge
	}
	return s.cdnBaseURL // Fallback to default
}

func (s *StreamService) generatePlaybackToken(userID, profileID, contentID string, validity time.Duration) string {
	// In production: Use HMAC-SHA256 with secret key
	claims := map[string]interface{}{
		"user_id":    userID,
		"profile_id": profileID,
		"content_id": contentID,
		"exp":        time.Now().Add(validity).Unix(),
	}
	
	// Simplified token generation (use JWT in production)
	tokenData, _ := json.Marshal(claims)
	return base64.URLEncoding.EncodeToString(tokenData)
}

func (s *StreamService) logPlaybackEvent(contentID, userID, profileID, eventType string) {
	// Publish to Kafka topic "playback.events"
	event := map[string]interface{}{
		"event_type":  eventType,
		"content_id":  contentID,
		"user_id":     userID,
		"profile_id":  profileID,
		"timestamp":   time.Now().Unix(),
	}
	
	// Kafka publish implementation
	log.Printf("Playback event: %+v", event)
}

// ============================================================================
// Main Application Setup
// ============================================================================

func main() {
	// Database connection
	db, err := sql.Open("postgres", "postgres://user:pass@localhost/streaming?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	// Redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	
	// Initialize services
	userService := NewUserService(db, redisClient)
	contentService := NewContentService(db, redisClient)
	streamService := NewStreamService(contentService, redisClient, "https://cdn.example.com")
	
	// Setup router
	r := mux.NewRouter()
	
	// User endpoints
	r.HandleFunc("/api/v1/users/{id}", userService.GetUserHandler).Methods("GET")
	r.HandleFunc("/api/v1/watch-history", userService.UpdateWatchHistoryHandler).Methods("POST")
	r.HandleFunc("/api/v1/profiles/{profile_id}/continue-watching", userService.GetContinueWatchingHandler).Methods("GET")
	
	// Content endpoints
	r.HandleFunc("/api/v1/content/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		content, err := contentService.GetContent(r.Context(), vars["id"])
		if err != nil {
			http.Error(w, "Content not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(content)
	}).Methods("GET")
	
	// Stream endpoints
	r.HandleFunc("/api/v1/stream/{content_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		// Extract from auth token (simplified)
		userID := r.Header.Get("X-User-ID")
		profileID := r.Header.Get("X-Profile-ID")
		region := r.Header.Get("X-Region")
		
		url, err := streamService.GeneratePlaybackURL(r.Context(), vars["content_id"], userID, profileID, region)
		if err != nil {
			http.Error(w, "Failed to generate playback URL", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"playback_url": url})
	}).Methods("GET")
	
	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
	
	// Prometheus metrics endpoint
	r.Handle("/metrics", promhttp.Handler())
	
	// Start server
	log.Println("User Service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Custom errors
type ErrContentNotAvailable struct {
	Region string
}

func (e *ErrContentNotAvailable) Error() string {
	return fmt.Sprintf("Content not available in region: %s", e.Region)
}
