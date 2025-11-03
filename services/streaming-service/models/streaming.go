package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PlaybackSession represents a video playback session
type PlaybackSession struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          string             `bson:"user_id" json:"userId"`
	ContentID       string             `bson:"content_id" json:"contentId"`
	Position        int64              `bson:"position" json:"position"` // milliseconds
	Duration        int64              `bson:"duration" json:"duration"` // milliseconds
	Quality         string             `bson:"quality" json:"quality"`
	Bandwidth       int64              `bson:"bandwidth" json:"bandwidth"` // bps
	LastHeartbeat   time.Time          `bson:"last_heartbeat" json:"lastHeartbeat"`
	DeviceID        string             `bson:"device_id" json:"deviceId"`
	StreamURL       string             `bson:"stream_url" json:"streamUrl"`
	DRMType         string             `bson:"drm_type,omitempty" json:"drmType,omitempty"`
	DRMLicenseURL   string             `bson:"drm_license_url,omitempty" json:"drmLicenseUrl,omitempty"`
	CreatedAt       time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updatedAt"`
}

// StreamManifest represents HLS or DASH manifest
type StreamManifest struct {
	ContentID       string             `json:"contentId"`
	ManifestURL     string             `json:"manifestUrl"`
	Type            string             `json:"type"` // "hls" or "dash"
	Qualities       []QualityLevel     `json:"qualities"`
	Subtitles       []SubtitleTrack    `json:"subtitles"`
	DRMInfo         *DRMInfo           `json:"drmInfo,omitempty"`
}

// QualityLevel represents a quality level
type QualityLevel struct {
	ID          string `json:"id"`
	Resolution  string `json:"resolution"`
	Bitrate     int    `json:"bitrate"`
	Codec       string `json:"codec"`
	URL         string `json:"url"`
}

// SubtitleTrack represents a subtitle track
type SubtitleTrack struct {
	Language string `json:"language"`
	Label    string `json:"label"`
	URL      string `json:"url"`
	Format   string `json:"format"` // "vtt", "srt"
}

// DRMInfo represents DRM information
type DRMInfo struct {
	Type         string `json:"type"` // "widevine", "fairplay", "playready"
	LicenseURL   string `json:"licenseUrl"`
	CertificateURL string `json:"certificateUrl,omitempty"`
}

// StreamingRequest represents a streaming request
type StreamingRequest struct {
	ContentID string `json:"contentId" binding:"required"`
	Quality   string `json:"quality,omitempty"` // "auto", "1080p", "720p", etc.
}

// StreamingToken represents a token for accessing manifests - Issue #14
type StreamingToken struct {
	Token      string `json:"token"`
	ExpiresIn  int    `json:"expiresIn"` // seconds
	ContentID  string `json:"contentId"`
}

// QoEEvent represents a Quality of Experience event - Issue #14
type QoEEvent struct {
	UserID         string    `json:"userId" binding:"required"`
	SessionID      string    `json:"sessionId"`
	ContentID      string    `json:"contentId" binding:"required"`
	Event          string    `json:"event" binding:"required"` // "play|pause|seek|buffering|error|ended"
	Bitrate        int       `json:"bitrate"`
	BufferDuration float64   `json:"bufferDuration"`
	Timestamp      time.Time `json:"timestamp"`
}

