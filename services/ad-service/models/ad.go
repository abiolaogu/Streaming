package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdRequest represents an ad request
type AdRequest struct {
	ContentID   string `json:"contentId"`
	UserID      string `json:"userId"`
	DeviceType  string `json:"deviceType"`
	Position    string `json:"position"` // "pre-roll", "mid-roll", "post-roll"
	CuePoint    int64  `json:"cuePoint,omitempty"` // For mid-roll
}

// AdResponse represents ad response
type AdResponse struct {
	Ads         []Ad   `json:"ads"`
	AdPodURL    string `json:"adPodUrl,omitempty"` // For SSAI
	SkipAllowed bool   `json:"skipAllowed"`
}

// Ad represents an ad
type Ad struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	VASTURL     string `json:"vastUrl"`
	Duration    int    `json:"duration"`
	ClickURL    string `json:"clickUrl,omitempty"`
	ImpressURL  string `json:"impressUrl,omitempty"`
}

// AdTracking represents ad tracking event
type AdTracking struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AdID        string             `bson:"ad_id" json:"adId"`
	UserID      string             `bson:"user_id" json:"userId"`
	EventType   string             `bson:"event_type" json:"eventType"` // "impression", "click", "complete", "skip"
	ContentID   string             `bson:"content_id" json:"contentId"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
}

