package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Content represents a video content item
type Content struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title         string              `bson:"title" json:"title"`
	Description   string              `bson:"description" json:"description"`
	Genre         string              `bson:"genre" json:"genre"`
	Category      string              `bson:"category" json:"category"` // "movie", "show", "live"
	PosterURL     string              `bson:"poster_url" json:"posterUrl"`
	BackdropURL   string              `bson:"backdrop_url" json:"backdropUrl"`
	StreamURL     string              `bson:"stream_url" json:"streamUrl"`
	Duration      int64               `bson:"duration" json:"duration"` // milliseconds
	ReleaseYear   int                 `bson:"release_year" json:"releaseYear"`
	Rating        float64             `bson:"rating" json:"rating"`
	IsDRMProtected bool               `bson:"is_drm_protected" json:"isDrmProtected"`
	DRMType       string              `bson:"drm_type,omitempty" json:"drmType,omitempty"`
	ThumbnailURL  string              `bson:"thumbnail_url,omitempty" json:"thumbnailUrl,omitempty"`
	Cast          []string            `bson:"cast" json:"cast"`
	Directors     []string            `bson:"directors" json:"directors"`
	Tags          []string            `bson:"tags" json:"tags"`
	Status        string              `bson:"status" json:"status"` // "draft", "published", "archived"
	CreatedAt     time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time           `bson:"updated_at" json:"updatedAt"`
}

// ContentRow represents a row of content for home screen
type ContentRow struct {
	ID    string    `json:"id"`
	Title string    `json:"title"`
	Items []Content `json:"items"`
}

// Series represents a TV series
type Series struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Seasons     []Season           `bson:"seasons" json:"seasons"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

// Season represents a season in a series
type Season struct {
	Number    int       `bson:"number" json:"number"`
	Episodes  []Content `bson:"episodes" json:"episodes"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
}

// Collection represents a curated collection or playlist
type Collection struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string              `bson:"title" json:"title"`
	Description string              `bson:"description" json:"description"`
	ContentIDs  []string            `bson:"content_ids" json:"contentIds"`
	Type        string              `bson:"type" json:"type"` // "curated", "user"
	UserID      string              `bson:"user_id,omitempty" json:"userId,omitempty"`
	CreatedAt   time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time           `bson:"updated_at" json:"updatedAt"`
}

// FASTChannel represents a 24/7 programmed channel
type FASTChannel struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string              `bson:"name" json:"name"`
	Description string              `bson:"description" json:"description"`
	EPG         []EPGEntry          `bson:"epg" json:"epg"`
	Schedule    []ScheduleItem      `bson:"schedule" json:"schedule"`
	CreatedAt   time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time           `bson:"updated_at" json:"updatedAt"`
}

// EPGEntry represents an Electronic Program Guide entry
type EPGEntry struct {
	Title       string    `bson:"title" json:"title"`
	StartTime   time.Time `bson:"start_time" json:"startTime"`
	EndTime     time.Time `bson:"end_time" json:"endTime"`
	Description string    `bson:"description" json:"description"`
}

// ScheduleItem represents a scheduled content item
type ScheduleItem struct {
	ContentID  string    `bson:"content_id" json:"contentId"`
	StartTime time.Time `bson:"start_time" json:"startTime"`
	EndTime   time.Time `bson:"end_time" json:"endTime"`
}

// LiveEvent represents a live event
type LiveEvent struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	StreamURL   string             `bson:"stream_url" json:"streamUrl"`
	StartTime   time.Time          `bson:"start_time" json:"startTime"`
	EndTime     time.Time          `bson:"end_time" json:"endTime"`
	IsPPV       bool               `bson:"is_ppv" json:"isPpv"`
	TicketURL   string             `bson:"ticket_url,omitempty" json:"ticketUrl,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

// Rating represents a user rating for content
type Rating struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ContentID string             `bson:"content_id" json:"contentId"`
	UserID    string             `bson:"user_id" json:"userId"`
	Stars     int                `bson:"stars" json:"stars"` // 1-5
	Comment   string             `bson:"comment,omitempty" json:"comment,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

// RatingAggregate represents aggregated rating statistics
type RatingAggregate struct {
	ContentID    string  `json:"contentId"`
	AverageStars float64 `json:"averageStars"`
	TotalRatings int     `json:"totalRatings"`
	Distribution map[int]int `json:"distribution"` // stars -> count
}

// Entitlement represents user's right to access content
type Entitlement struct {
	ContentID   string    `json:"contentId"`
	UserID      string    `json:"userId"`
	HasAccess   bool      `json:"hasAccess"`
	Reason      string    `json:"reason,omitempty"` // "subscription", "purchased", "free", "geo_blocked", "expired"
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
	DRMLevel    string    `json:"drmLevel,omitempty"` // "1" (4K), "2" (1080p), "3" (SD)
	LicenseURL  string    `json:"licenseUrl,omitempty"`
}

// Category represents a content category with count
type Category struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

