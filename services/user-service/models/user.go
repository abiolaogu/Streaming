package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserProfile represents a user profile
type UserProfile struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"userId"`
	DisplayName string             `bson:"display_name,omitempty" json:"displayName,omitempty"`
	Bio         string             `bson:"bio,omitempty" json:"bio,omitempty"`
	AvatarURL   string             `bson:"avatar_url,omitempty" json:"avatarUrl,omitempty"`
	DateOfBirth *time.Time         `bson:"date_of_birth,omitempty" json:"dateOfBirth,omitempty"`
	Location    string             `bson:"location,omitempty" json:"location,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

// UserPreferences represents user preferences
type UserPreferences struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID            string             `bson:"user_id" json:"userId"`
	Language          string             `bson:"language" json:"language"`
	SubtitleLanguage  string             `bson:"subtitle_language,omitempty" json:"subtitleLanguage,omitempty"`
	ContentRating     string             `bson:"content_rating" json:"contentRating"` // G, PG, PG-13, R
	Notifications     NotificationPrefs  `bson:"notifications" json:"notifications"`
	Playback          PlaybackPrefs      `bson:"playback" json:"playback"`
	CreatedAt         time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updatedAt"`
}

// NotificationPrefs represents notification preferences
type NotificationPrefs struct {
	Email    bool `bson:"email" json:"email"`
	Push     bool `bson:"push" json:"push"`
	Marketing bool `bson:"marketing" json:"marketing"`
}

// PlaybackPrefs represents playback preferences
type PlaybackPrefs struct {
	Quality      string `bson:"quality" json:"quality"`           // auto, 1080p, 720p, 480p
	Autoplay     bool   `bson:"autoplay" json:"autoplay"`
	SkipIntro    bool   `bson:"skip_intro" json:"skipIntro"`
	SkipCredits  bool   `bson:"skip_credits" json:"skipCredits"`
}

// Profile represents a sub-profile (for family accounts)
type Profile struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"userId"`
	Name        string             `bson:"name" json:"name"`
	AvatarURL   string             `bson:"avatar_url,omitempty" json:"avatarUrl,omitempty"`
	IsKids      bool               `bson:"is_kids" json:"isKids"`
	PIN         string             `bson:"pin,omitempty" json:"-"`
	ContentRating string           `bson:"content_rating" json:"contentRating"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

// WatchHistory represents a watch history entry
type WatchHistory struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"userId"`
	ProfileID   string             `bson:"profile_id,omitempty" json:"profileId,omitempty"`
	ContentID   string             `bson:"content_id" json:"contentId"`
	Position    int64              `bson:"position" json:"position"` // milliseconds
	Duration    int64              `bson:"duration" json:"duration"` // milliseconds
	WatchedAt   time.Time          `bson:"watched_at" json:"watchedAt"`
	Completed   bool               `bson:"completed" json:"completed"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

// Watchlist represents a user's watchlist/favorites
type Watchlist struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"userId"`
	ContentID   string             `bson:"content_id" json:"contentId"`
	AddedAt     time.Time          `bson:"added_at" json:"addedAt"`
}

// Device represents a user device
type Device struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"userId"`
	DeviceID    string             `bson:"device_id" json:"deviceId"`
	DeviceName  string             `bson:"device_name" json:"deviceName"`
	DeviceType  string             `bson:"device_type" json:"deviceType"` // mobile, tv, tablet, web
	OS          string             `bson:"os" json:"os"`
	OSVersion   string             `bson:"os_version,omitempty" json:"osVersion,omitempty"`
	AppVersion  string             `bson:"app_version,omitempty" json:"appVersion,omitempty"`
	LastUsedAt  time.Time          `bson:"last_used_at" json:"lastUsedAt"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
}

// UserDataExport represents GDPR data export
type UserDataExport struct {
	UserID      string             `json:"userId"`
	Profile     *UserProfile       `json:"profile"`
	Preferences *UserPreferences   `json:"preferences"`
	Profiles    []Profile          `json:"profiles"`
	WatchHistory []WatchHistory    `json:"watchHistory"`
	Watchlist   []Watchlist        `json:"watchlist"`
	Devices     []Device           `json:"devices"`
	ExportedAt  time.Time          `json:"exportedAt"`
}

