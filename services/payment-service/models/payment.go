package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Subscription represents a user subscription
type Subscription struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"userId"`
	PlanID      string             `bson:"plan_id" json:"planId"`
	Status      string             `bson:"status" json:"status"` // "active", "cancelled", "expired", "paused"
	PaymentMethodID string         `bson:"payment_method_id" json:"paymentMethodId"`
	CurrentPeriodStart time.Time    `bson:"current_period_start" json:"currentPeriodStart"`
	CurrentPeriodEnd   time.Time    `bson:"current_period_end" json:"currentPeriodEnd"`
	CancelAtPeriodEnd  bool         `bson:"cancel_at_period_end" json:"cancelAtPeriodEnd"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

// Plan represents a subscription plan
type Plan struct {
	ID          string             `bson:"id" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Price       float64            `bson:"price" json:"price"`
	Currency    string             `bson:"currency" json:"currency"`
	Interval    string             `bson:"interval" json:"interval"` // "month", "year"
	Features    []string           `bson:"features" json:"features"`
	MaxStreams  int                `bson:"max_streams" json:"maxStreams"`
	Quality     string             `bson:"quality" json:"quality"`
}

// Payment represents a payment transaction
type Payment struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"userId"`
	Amount      float64            `bson:"amount" json:"amount"`
	Currency    string             `bson:"currency" json:"currency"`
	Status      string             `bson:"status" json:"status"` // "pending", "completed", "failed", "refunded"
	PaymentMethod string           `bson:"payment_method" json:"paymentMethod"`
	Gateway     string             `bson:"gateway" json:"gateway"` // "stripe", "paypal", etc.
	GatewayTransactionID string    `bson:"gateway_transaction_id" json:"gatewayTransactionId"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

// Invoice represents an invoice
type Invoice struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"userId"`
	SubscriptionID string          `bson:"subscription_id,omitempty" json:"subscriptionId,omitempty"`
	Amount      float64            `bson:"amount" json:"amount"`
	Currency    string             `bson:"currency" json:"currency"`
	Status      string             `bson:"status" json:"status"`
	PDFURL      string             `bson:"pdf_url,omitempty" json:"pdfUrl,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
}

// Purchase represents a TVOD/PPV purchase
type Purchase struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"userId"`
	ContentID   string             `bson:"content_id" json:"contentId"`
	Type        string             `bson:"type" json:"type"` // "rent", "buy", "ppv"
	Amount      float64            `bson:"amount" json:"amount"`
	Currency    string             `bson:"currency" json:"currency"`
	ExpiresAt   *time.Time         `bson:"expires_at,omitempty" json:"expiresAt,omitempty"`
	Status      string             `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
}

