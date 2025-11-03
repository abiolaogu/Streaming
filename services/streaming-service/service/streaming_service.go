package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/golang-jwt/jwt/v5"
	"github.com/streamverse/streaming-service/models"
	"github.com/streamverse/streaming-service/repository"
)

// StreamingService handles streaming business logic
type StreamingService struct {
	repo             *repository.StreamingRepository
	contentRepo      ContentRepository
	subscriptionRepo SubscriptionRepository
	jwtSecret        string // TODO: Load from config
}

// ContentRepository interface for content service
type ContentRepository interface {
	GetContentByID(ctx context.Context, id string) (map[string]interface{}, error)
}

// SubscriptionRepository interface for subscription checks
type SubscriptionRepository interface {
	GetUserSubscription(ctx context.Context, userID string) (map[string]interface{}, error)
	CheckConcurrentStreamLimit(ctx context.Context, userID string) (int, error)
}

// NewStreamingService creates a new streaming service
func NewStreamingService(
	repo *repository.StreamingRepository,
	contentRepo ContentRepository,
	subscriptionRepo SubscriptionRepository,
	jwtSecret string,
) *StreamingService {
	return &StreamingService{
		repo:             repo,
		contentRepo:      contentRepo,
		subscriptionRepo: subscriptionRepo,
		jwtSecret:        jwtSecret,
	}
}

// GenerateToken generates a JWT token for manifest access - Issue #14
func (s *StreamingService) GenerateToken(ctx context.Context, contentID, userID, ip, deviceID string) (*models.StreamingToken, error) {
	now := time.Now()
	expiresIn := 3600 // 1 hour

	claims := jwt.MapClaims{
		"content_id": contentID,
		"user_id":    userID,
		"ip":         ip,
		"device_id":  deviceID,
		"exp":        jwt.NewNumericDate(now.Add(time.Duration(expiresIn) * time.Second)),
		"nbf":        jwt.NewNumericDate(now),
		"iat":        jwt.NewNumericDate(now),
		"aud":        "cdn.streamverse.io",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.StreamingToken{
		Token:     tokenString,
		ExpiresIn: expiresIn,
		ContentID: contentID,
	}, nil
}

// ValidateToken validates a token and returns user_id - Issue #14
func (s *StreamingService) ValidateToken(ctx context.Context, tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", fmt.Errorf("invalid token claims")
		}
		return userID, nil
	}

	return "", fmt.Errorf("invalid token")
}

// GenerateHLSManifest generates an HLS manifest - Issue #14
func (s *StreamingService) GenerateHLSManifest(ctx context.Context, contentID, userID string) (string, error) {
	// Get content metadata
	content, err := s.contentRepo.GetContentByID(ctx, contentID)
	if err != nil {
		return "", fmt.Errorf("content not found: %w", err)
	}

	// Select ABR profile based on device/network
	deviceType := s.detectDeviceType(ctx, userID) // TODO: Get from session or device info
	selectedProfile := s.SelectABRProfile(ctx, userID, deviceType)

	// Generate HLS manifest
	manifest := fmt.Sprintf(`#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:6
#EXTINF:6.0,
%s/%s/segment_%d.ts
#EXT-X-ENDLIST
`, s.getCDNBaseURL(), selectedProfile, 0)

	return manifest, nil
}

// GenerateDASHManifest generates a DASH manifest - Issue #14
func (s *StreamingService) GenerateDASHManifest(ctx context.Context, contentID, userID string) (string, error) {
	// Get content metadata
	content, err := s.contentRepo.GetContentByID(ctx, contentID)
	if err != nil {
		return "", fmt.Errorf("content not found: %w", err)
	}

	// Select ABR profile
	deviceType := s.detectDeviceType(ctx, userID)
	selectedProfile := s.SelectABRProfile(ctx, userID, deviceType)

	// Generate DASH manifest XML (simplified)
	manifest := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<MPD xmlns="urn:mpeg:dash:schema:mpd:2011" type="static" mediaPresentationDuration="PT0H0M0S">
  <Period>
    <AdaptationSet>
      <Representation id="%s" bandwidth="5000000" width="1920" height="1080">
        <BaseURL>%s/%s/</BaseURL>
      </Representation>
    </AdaptationSet>
  </Period>
</MPD>`, selectedProfile, s.getCDNBaseURL(), selectedProfile)

	return manifest, nil
}

// SelectABRProfile selects ABR profile based on device and network - Issue #14
func (s *StreamingService) SelectABRProfile(ctx context.Context, userID, deviceType string) string {
	// Bitrate ladder: 240p (512k), 360p (1.5M), 480p (2.5M), 720p (5M), 1080p (8M), 4K (15M)
	
	// Detect device type and select initial profile
	switch deviceType {
	case "mobile":
		return "360p" // 1.5M
	case "tablet":
		return "720p" // 5M
	case "desktop":
		return "1080p" // 8M
	case "tv":
		return "4K" // 15M
	default:
		return "480p" // 2.5M default
	}
	
	// TODO: Estimate bandwidth from previous QoE events
	// TODO: Adapt based on buffer level and rebuffer events
}

// SubmitQoE submits QoE metrics - Issue #14
func (s *StreamingService) SubmitQoE(ctx context.Context, event *models.QoEEvent) error {
	// TODO: Send to Kafka topic "qoe-events" for Analytics Service
	// For now, just log it
	fmt.Printf("QoE Event: %+v\n", event)
	return nil
}

// Helper methods for Issue #14
func (s *StreamingService) detectDeviceType(ctx context.Context, userID string) string {
	// TODO: Get device type from device registry or session
	return "desktop" // Default
}

func (s *StreamingService) getCDNBaseURL() string {
	return "https://cdn.streamverse.com/videos"
}

// CreateSession creates a new playback session
func (s *StreamingService) CreateSession(ctx context.Context, userID, contentID, deviceID string) (*models.PlaybackSession, error) {
	// Check subscription and concurrent stream limits
	activeStreams, err := s.subscriptionRepo.CheckConcurrentStreamLimit(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check stream limits: %w", err)
	}

	if activeStreams >= 4 { // Assuming max 4 concurrent streams
		return nil, fmt.Errorf("concurrent stream limit reached")
	}

	// Get content metadata
	content, err := s.contentRepo.GetContentByID(ctx, contentID)
	if err != nil {
		return nil, fmt.Errorf("content not found: %w", err)
	}

	// Generate stream URL (simplified - would integrate with CDN)
	streamURL := s.generateStreamURL(contentID)

	session := &models.PlaybackSession{
		ID:            primitive.NewObjectID(),
		UserID:        userID,
		ContentID:     contentID,
		Position:      0,
		Duration:      getContentDuration(content),
		Quality:       "auto",
		LastHeartbeat: time.Now(),
		DeviceID:      deviceID,
		StreamURL:     streamURL,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return s.repo.CreateSession(ctx, session)
}

// GetManifest generates HLS or DASH manifest
func (s *StreamingService) GetManifest(ctx context.Context, contentID, format string, userID string) (*models.StreamManifest, error) {
	// Get content metadata
	content, err := s.contentRepo.GetContentByID(ctx, contentID)
	if err != nil {
		return nil, fmt.Errorf("content not found: %w", err)
	}

	// Check geo-restrictions
	if err := s.checkGeoRestrictions(ctx, contentID, userID); err != nil {
		return nil, err
	}

	// Generate manifest URL
	manifestURL := s.generateManifestURL(contentID, format)

	// Get quality levels
	qualities := s.getQualityLevels(contentID)

	// Get subtitles
	subtitles := s.getSubtitles(contentID)

	// Get DRM info if protected
	var drmInfo *models.DRMInfo
	if isDRMProtected(content) {
		drmInfo = s.getDRMInfo(content, format)
	}

	return &models.StreamManifest{
		ContentID:   contentID,
		ManifestURL: manifestURL,
		Type:        format,
		Qualities:   qualities,
		Subtitles:   subtitles,
		DRMInfo:     drmInfo,
	}, nil
}

// UpdateSessionPosition updates playback position
func (s *StreamingService) UpdateSessionPosition(ctx context.Context, sessionID string, position int64) error {
	return s.repo.UpdatePosition(ctx, sessionID, position)
}

// SendHeartbeat updates session heartbeat
func (s *StreamingService) SendHeartbeat(ctx context.Context, sessionID string) error {
	return s.repo.UpdateHeartbeat(ctx, sessionID)
}

// EndSession ends a playback session
func (s *StreamingService) EndSession(ctx context.Context, sessionID string) error {
	return s.repo.DeleteSession(ctx, sessionID)
}

// Helper methods
func (s *StreamingService) generateStreamURL(contentID string) string {
	// This would integrate with CDN and generate signed URLs
	return fmt.Sprintf("https://cdn.streamverse.com/videos/%s/master.m3u8", contentID)
}

func (s *StreamingService) generateManifestURL(contentID, format string) string {
	if format == "dash" {
		return fmt.Sprintf("https://cdn.streamverse.com/videos/%s/manifest.mpd", contentID)
	}
	return fmt.Sprintf("https://cdn.streamverse.com/videos/%s/master.m3u8", contentID)
}

func (s *StreamingService) getQualityLevels(contentID string) []models.QualityLevel {
	return []models.QualityLevel{
		{ID: "1080p", Resolution: "1920x1080", Bitrate: 5000000, Codec: "h264", URL: fmt.Sprintf("https://cdn.streamverse.com/videos/%s/1080p.m3u8", contentID)},
		{ID: "720p", Resolution: "1280x720", Bitrate: 3000000, Codec: "h264", URL: fmt.Sprintf("https://cdn.streamverse.com/videos/%s/720p.m3u8", contentID)},
		{ID: "480p", Resolution: "854x480", Bitrate: 1500000, Codec: "h264", URL: fmt.Sprintf("https://cdn.streamverse.com/videos/%s/480p.m3u8", contentID)},
	}
}

func (s *StreamingService) getSubtitles(contentID string) []models.SubtitleTrack {
	return []models.SubtitleTrack{
		{Language: "en", Label: "English", URL: fmt.Sprintf("https://cdn.streamverse.com/subtitles/%s/en.vtt", contentID), Format: "vtt"},
	}
}

func (s *StreamingService) checkGeoRestrictions(ctx context.Context, contentID, userID string) error {
	// TODO: Implement IP geolocation check
	return nil
}

func (s *StreamingService) getDRMInfo(content map[string]interface{}, format string) *models.DRMInfo {
	drmType := getDRMType(content, format)
	return &models.DRMInfo{
		Type:        drmType,
		LicenseURL:  fmt.Sprintf("https://drm.streamverse.com/license/%s", drmType),
		CertificateURL: getCertificateURL(drmType),
	}
}

func getContentDuration(content map[string]interface{}) int64 {
	if duration, ok := content["duration"].(int64); ok {
		return duration
	}
	return 0
}

func isDRMProtected(content map[string]interface{}) bool {
	if protected, ok := content["isDrmProtected"].(bool); ok {
		return protected
	}
	return false
}

func getDRMType(content map[string]interface{}, format string) string {
	if format == "dash" {
		return "widevine"
	}
	if drmType, ok := content["drmType"].(string); ok {
		return drmType
	}
	return "widevine"
}

func getCertificateURL(drmType string) string {
	if drmType == "fairplay" {
		return "https://drm.streamverse.com/certificates/fairplay.cer"
	}
	return ""
}
