package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CompositedVideo represents a video with composited ads
type CompositedVideo struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	VideoID        string             `bson:"video_id" json:"videoId"`
	CompositedURL  string             `bson:"composited_url" json:"compositedUrl"`
	TrackingPixels []string           `bson:"tracking_pixels" json:"trackingPixels"`
}
