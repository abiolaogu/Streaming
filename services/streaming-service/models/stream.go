package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StreamSession represents a playback session
type StreamSession struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          string             `bson:"user_id" json:"userId"`
	ContentID       string             `bson:"content_id" json:"contentId"`
	DeviceID        string             `bson:"device_id" json:"deviceId"`
	SessionToken    string             `bson:"session_token" json:"sessionToken"`
	Position        int64              `bson:"position" json:"position"` // milliseconds
	Duration        int64              `bson:"duration" json:"duration"` // milliseconds
	Quality         string             `bson:"quality" json:"quality"`
	Bandwidth       int64              `bson:"bandwidth" json:"bandwidth"`
	Protocol        string             `bson:"protocol" json:"protocol"` // hls, dash
	ManifestURL     string             `bson:"manifest_url" json:"manifestUrl"`
	CreatedAt       time.Time          `bson:"created_at" json:"createdAt"`
	LastHeartbeat   time.Time          `bson:"last_heartbeat" json:"lastHeartbeat"`
	EndedAt         *time.Time         `bson:"ended_at,omitempty" json:"endedAt,omitempty"`
}

// Manifest represents HLS/DASH manifest
type Manifest struct {
	ContentID    string    `json:"contentId"`
	Protocol     string    `json:"protocol"` // hls, dash
	BaseURL      string    `json:"baseUrl"`
	Variants     []Variant `json:"variants"`
	Subtitles    []Subtitle `json:"subtitles"`
	DRMConfig    *DRMConfig `json:"drmConfig,omitempty"`
}

// Variant represents a quality variant
type Variant struct {
	Bandwidth int    `json:"bandwidth"`
	Resolution string `json:"resolution"`
	Codec     string `json:"codec"`
	URL       string `json:"url"`
}

// Subtitle represents subtitle track
type Subtitle struct {
	Language string `json:"language"`
	Label    string `json:"label"`
	URL      string `json:"url"`
	Default  bool   `json:"default"`
}

// DRMConfig represents DRM configuration
type DRMConfig struct {
	Type        string   `json:"type"` // widevine, fairplay, playready
	LicenseURL  string   `json:"licenseUrl"`
	Certificate string   `json:"certificate,omitempty"`
	KeyIDs      []string `json:"keyIds,omitempty"`
}

// PlaybackEvent represents an analytics event
type PlaybackEvent struct {
	SessionID   string    `json:"sessionId"`
	EventType   string    `json:"eventType"` // play, pause, seek, quality_change, buffering, error
	Timestamp   time.Time `json:"timestamp"`
	Position    int64     `json:"position"`
	Quality     string    `json:"quality,omitempty"`
	Error       string    `json:"error,omitempty"`
}

