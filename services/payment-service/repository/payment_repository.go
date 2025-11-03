package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/streamverse/common-go/database"
	"github.com/streamverse/payment-service/models"
)

// PaymentRepository handles payment data operations
type PaymentRepository struct {
	subscriptionCollection *mongo.Collection
	purchaseCollection     *mongo.Collection
	plans                  map[string]*models.Plan
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *database.MongoDB) *PaymentRepository {
	// Initialize default plans
	plans := map[string]*models.Plan{
		"basic": {
			ID:         "basic",
			Name:       "Basic",
			Price:      9.99,
			Currency:   "USD",
			Interval:   "month",
			Features:   []string{"1080p", "2 screens"},
			MaxStreams: 2,
			Quality:    "1080p",
		},
		"premium": {
			ID:         "premium",
			Name:       "Premium",
			Price:      14.99,
			Currency:   "USD",
			Interval:   "month",
			Features:   []string{"4K", "4 screens", "offline downloads"},
			MaxStreams: 4,
			Quality:    "4K",
		},
	}

	return &PaymentRepository{
		subscriptionCollection: db.Collection("subscriptions"),
		purchaseCollection:     db.Collection("purchases"),
		plans:                  plans,
	}
}

// GetPlan retrieves a plan by ID
func (r *PaymentRepository) GetPlan(planID string) (*models.Plan, error) {
	plan, ok := r.plans[planID]
	if !ok {
		return nil, fmt.Errorf("plan not found")
	}
	return plan, nil
}

// CreateSubscription creates a subscription
func (r *PaymentRepository) CreateSubscription(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {
	_, err := r.subscriptionCollection.InsertOne(ctx, subscription)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

// GetSubscriptionByUserID retrieves user subscription
func (r *PaymentRepository) GetSubscriptionByUserID(ctx context.Context, userID string) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.subscriptionCollection.FindOne(ctx, bson.M{"user_id": userID, "status": "active"}).Decode(&subscription)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// CancelSubscription cancels a subscription - Issue #16
func (r *PaymentRepository) CancelSubscription(ctx context.Context, userID, subscriptionID string) error {
	objectID, err := primitive.ObjectIDFromHex(subscriptionID)
	if err != nil {
		return fmt.Errorf("invalid subscription ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"status":               "canceled",
			"cancel_at_period_end": true,
			"updated_at":           time.Now(),
		},
	}

	_, err = r.subscriptionCollection.UpdateOne(ctx, bson.M{"_id": objectID, "user_id": userID}, update)
	return err
}

// CreatePurchase creates a purchase
func (r *PaymentRepository) CreatePurchase(ctx context.Context, purchase *models.Purchase) (*models.Purchase, error) {
	_, err := r.purchaseCollection.InsertOne(ctx, purchase)
	if err != nil {
		return nil, err
	}
	return purchase, nil
}

