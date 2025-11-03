package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Channel represents a FAST channel or live TV channel
type Channel struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ChannelID   string             `bson:"channel_id" json:"channelId"` // e.g., "pluto-drama"
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Type        string             `bson:"type" json:"type"` // "fast" or "live"
	ManifestURL string             `bson:"manifest_url,omitempty" json:"manifestUrl,omitempty"`
	IngestURL   string             `bson:"ingest_url,omitempty" json:"ingestUrl,omitempty"` // For live channels
	Status      string             `bson:"status" json:"status"` // "active", "inactive"
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

// ScheduleEntry represents a scheduled content item
type ScheduleEntry struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ChannelID   string             `bson:"channel_id" json:"channelId"`
	ContentID   string             `bson:"content_id" json:"contentId"`
	StartTime   time.Time          `bson:"start_time" json:"startTime"`
	EndTime     time.Time          `bson:"end_time" json:"endTime"`
	Duration    int                `bson:"duration" json:"duration"` // seconds
	Title       string             `bson:"title,omitempty" json:"title,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Poster      string             `bson:"poster,omitempty" json:"poster,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

// EPG represents an Electronic Program Guide for a channel
type EPG struct {
	ChannelID string          `json:"channelId"`
	ChannelName string         `json:"channelName"`
	Schedule   []EPGEntry      `json:"schedule"`
	GeneratedAt time.Time      `json:"generatedAt"`
}

// EPGEntry represents a single EPG entry
type EPGEntry struct {
	Title       string    `json:"title"`
	StartTime   time.Time `json:"startTime"`
	Duration    int       `json:"duration"` // seconds
	Description string    `json:"description,omitempty"`
	Poster      string    `json:"poster,omitempty"`
	ContentID   string    `json:"contentId,omitempty"`
}

// ChannelManifest represents a streaming manifest for a channel
type ChannelManifest struct {
	ChannelID string `json:"channelId"`
	ManifestURL string `json:"manifestUrl"`
	Type       string `json:"type"` // "hls" or "dash"
}

