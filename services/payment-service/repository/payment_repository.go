package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/streamverse/common-go/database"
	"github.com/streamverse/payment-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PaymentRepository handles payment data operations
type PaymentRepository struct {
	subscriptionCollection *mongo.Collection
	purchaseCollection     *mongo.Collection
	stripeLinkCollection   *mongo.Collection
	webhookEventCollection *mongo.Collection
	plans                  map[string]*models.Plan
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *database.MongoDB) *PaymentRepository {
	// Initialize default plans
	tier1 := &models.Plan{
		ID:         "tier1",
		Name:       "Basic",
		Price:      4.99,
		Currency:   "USD",
		Interval:   "month",
		Features:   []string{"480p", "1 screen", "ad-supported"},
		MaxStreams: 1,
		Quality:    "480p",
	}
	tier2 := &models.Plan{
		ID:         "tier2",
		Name:       "Pro",
		Price:      12.99,
		Currency:   "USD",
		Interval:   "month",
		Features:   []string{"720p", "2 screens", "downloads"},
		MaxStreams: 2,
		Quality:    "720p",
	}
	tier3 := &models.Plan{
		ID:         "tier3",
		Name:       "Premium",
		Price:      19.99,
		Currency:   "USD",
		Interval:   "month",
		Features:   []string{"4K", "4 screens", "downloads", "priority support"},
		MaxStreams: 4,
		Quality:    "4K",
	}

	plans := map[string]*models.Plan{
		"tier1": tier1,
		"tier2": tier2,
		"tier3": tier3,
		// Backward-compatible aliases.
		"basic":    tier1,
		"standard": tier2,
		"pro":      tier2,
		"premium":  tier3,
	}

	subscriptionCollection := db.Collection("subscriptions")
	purchaseCollection := db.Collection("purchases")
	stripeLinkCollection := db.Collection("stripe_customer_links")
	webhookEventCollection := db.Collection("stripe_webhook_events")

	_, _ = webhookEventCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "event_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)

	return &PaymentRepository{
		subscriptionCollection: subscriptionCollection,
		purchaseCollection:     purchaseCollection,
		stripeLinkCollection:   stripeLinkCollection,
		webhookEventCollection: webhookEventCollection,
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

// UpsertSubscriptionByUserID upserts a subscription record keyed by user ID.
func (r *PaymentRepository) UpsertSubscriptionByUserID(ctx context.Context, subscription *models.Subscription) error {
	subscription.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"plan_id":                subscription.PlanID,
			"status":                 subscription.Status,
			"payment_method_id":      subscription.PaymentMethodID,
			"stripe_customer_id":     subscription.StripeCustomerID,
			"stripe_subscription_id": subscription.StripeSubscriptionID,
			"current_period_start":   subscription.CurrentPeriodStart,
			"current_period_end":     subscription.CurrentPeriodEnd,
			"cancel_at_period_end":   subscription.CancelAtPeriodEnd,
			"updated_at":             subscription.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"_id":        subscription.ID,
			"user_id":    subscription.UserID,
			"created_at": subscription.CreatedAt,
		},
	}

	_, err := r.subscriptionCollection.UpdateOne(
		ctx,
		bson.M{"user_id": subscription.UserID},
		update,
		options.Update().SetUpsert(true),
	)
	return err
}

// UpdateSubscriptionStatusByUserID updates subscription status for a user.
func (r *PaymentRepository) UpdateSubscriptionStatusByUserID(ctx context.Context, userID, status string, cancelAtPeriodEnd bool) error {
	update := bson.M{
		"$set": bson.M{
			"status":               status,
			"cancel_at_period_end": cancelAtPeriodEnd,
			"updated_at":           time.Now(),
		},
	}

	_, err := r.subscriptionCollection.UpdateOne(ctx, bson.M{"user_id": userID}, update)
	return err
}

// UpsertStripeLink upserts Stripe-to-user reconciliation IDs.
func (r *PaymentRepository) UpsertStripeLink(ctx context.Context, userID, customerID, subscriptionID string) error {
	if userID == "" || (customerID == "" && subscriptionID == "") {
		return nil
	}

	update := bson.M{
		"$set": bson.M{
			"user_id":         userID,
			"customer_id":     customerID,
			"subscription_id": subscriptionID,
			"updated_at":      time.Now(),
		},
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}

	filter := bson.M{"user_id": userID}
	_, err := r.stripeLinkCollection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

// ResolveUserIDByStripeIDs resolves a user ID from Stripe customer/subscription IDs.
func (r *PaymentRepository) ResolveUserIDByStripeIDs(ctx context.Context, customerID, subscriptionID string) (string, error) {
	if customerID == "" && subscriptionID == "" {
		return "", nil
	}

	filter := bson.M{}
	switch {
	case subscriptionID != "" && customerID != "":
		filter = bson.M{"$or": []bson.M{
			{"subscription_id": subscriptionID},
			{"customer_id": customerID},
		}}
	case subscriptionID != "":
		filter = bson.M{"subscription_id": subscriptionID}
	case customerID != "":
		filter = bson.M{"customer_id": customerID}
	}

	var doc struct {
		UserID string `bson:"user_id"`
	}
	err := r.stripeLinkCollection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil
		}
		return "", err
	}

	return doc.UserID, nil
}

// GetActivePurchasesByUserID retrieves active purchases for a user.
func (r *PaymentRepository) GetActivePurchasesByUserID(ctx context.Context, userID string) ([]models.Purchase, error) {
	now := time.Now()
	filter := bson.M{
		"user_id": userID,
		"status":  "completed",
		"$or": []bson.M{
			{"expires_at": bson.M{"$exists": false}},
			{"expires_at": bson.M{"$gt": now}},
		},
	}

	cursor, err := r.purchaseCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var purchases []models.Purchase
	if err := cursor.All(ctx, &purchases); err != nil {
		return nil, err
	}

	return purchases, nil
}

// BeginWebhookEvent starts webhook processing if the event is new or failed previously.
// Returns false when the event was already processed or is currently in-flight.
func (r *PaymentRepository) BeginWebhookEvent(
	ctx context.Context,
	eventID,
	eventType,
	payloadHash string,
	eventObject map[string]interface{},
) (bool, error) {
	if eventID == "" {
		return false, fmt.Errorf("webhook event id is required")
	}

	now := time.Now().UTC()
	_, err := r.webhookEventCollection.InsertOne(ctx, bson.M{
		"event_id":     eventID,
		"event_type":   eventType,
		"payload_hash": payloadHash,
		"event_object": eventObject,
		"status":       "processing",
		"attempts":     1,
		"created_at":   now,
		"updated_at":   now,
	})
	if err == nil {
		return true, nil
	}
	if !mongo.IsDuplicateKeyError(err) {
		return false, err
	}

	var existing struct {
		Status      string `bson:"status"`
		PayloadHash string `bson:"payload_hash"`
	}
	if err := r.webhookEventCollection.FindOne(ctx, bson.M{"event_id": eventID}).Decode(&existing); err != nil {
		return false, err
	}

	if existing.PayloadHash != "" && payloadHash != "" && existing.PayloadHash != payloadHash {
		return false, fmt.Errorf("payload mismatch for existing webhook event: %s", eventID)
	}

	switch existing.Status {
	case "processed", "processing":
		return false, nil
	case "failed":
		update := bson.M{
			"$set": bson.M{
				"status":       "processing",
				"event_type":   eventType,
				"payload_hash": payloadHash,
				"event_object": eventObject,
				"last_error":   "",
				"updated_at":   now,
			},
			"$inc": bson.M{
				"attempts": 1,
			},
		}
		result, err := r.webhookEventCollection.UpdateOne(ctx, bson.M{"event_id": eventID, "status": "failed"}, update)
		if err != nil {
			return false, err
		}
		return result.ModifiedCount == 1, nil
	default:
		return false, nil
	}
}

// MarkWebhookEventProcessed marks an event as successfully processed.
func (r *PaymentRepository) MarkWebhookEventProcessed(ctx context.Context, eventID string) error {
	if eventID == "" {
		return fmt.Errorf("webhook event id is required")
	}

	now := time.Now().UTC()
	result, err := r.webhookEventCollection.UpdateOne(ctx, bson.M{"event_id": eventID}, bson.M{
		"$set": bson.M{
			"status":       "processed",
			"processed_at": now,
			"updated_at":   now,
		},
		"$unset": bson.M{
			"last_error": "",
		},
	})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("webhook event not found: %s", eventID)
	}
	return nil
}

// MarkWebhookEventFailed marks an event as failed and stores the latest error.
func (r *PaymentRepository) MarkWebhookEventFailed(ctx context.Context, eventID, lastError string) error {
	if eventID == "" {
		return fmt.Errorf("webhook event id is required")
	}

	result, err := r.webhookEventCollection.UpdateOne(ctx, bson.M{"event_id": eventID}, bson.M{
		"$set": bson.M{
			"status":     "failed",
			"last_error": lastError,
			"updated_at": time.Now().UTC(),
		},
	})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("webhook event not found: %s", eventID)
	}
	return nil
}

// ListFailedWebhookEvents returns failed webhook events for reconciliation replay.
func (r *PaymentRepository) ListFailedWebhookEvents(ctx context.Context, limit int) ([]models.WebhookEvent, error) {
	if limit <= 0 {
		limit = 50
	}

	cursor, err := r.webhookEventCollection.Find(
		ctx,
		bson.M{"status": "failed"},
		options.Find().SetSort(bson.D{{Key: "updated_at", Value: 1}}).SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []models.WebhookEvent
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}

	return events, nil
}
