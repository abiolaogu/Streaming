package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/streamverse/payment-service/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type paymentRepository interface {
	GetPlan(planID string) (*models.Plan, error)
	CreateSubscription(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error)
	GetSubscriptionByUserID(ctx context.Context, userID string) (*models.Subscription, error)
	CancelSubscription(ctx context.Context, userID, subscriptionID string) error
	CreatePurchase(ctx context.Context, purchase *models.Purchase) (*models.Purchase, error)
	UpsertSubscriptionByUserID(ctx context.Context, subscription *models.Subscription) error
	UpdateSubscriptionStatusByUserID(ctx context.Context, userID, status string, cancelAtPeriodEnd bool) error
	UpsertStripeLink(ctx context.Context, userID, customerID, subscriptionID string) error
	ResolveUserIDByStripeIDs(ctx context.Context, customerID, subscriptionID string) (string, error)
	GetActivePurchasesByUserID(ctx context.Context, userID string) ([]models.Purchase, error)
	BeginWebhookEvent(ctx context.Context, eventID, eventType, payloadHash string) (bool, error)
	MarkWebhookEventProcessed(ctx context.Context, eventID string) error
	MarkWebhookEventFailed(ctx context.Context, eventID, lastError string) error
}

// PaymentService handles payment business logic
type PaymentService struct {
	repo paymentRepository
}

// NewPaymentService creates a new payment service
func NewPaymentService(repo paymentRepository) *PaymentService {
	return &PaymentService{
		repo: repo,
	}
}

// Subscribe subscribes user to a plan
func (s *PaymentService) Subscribe(ctx context.Context, userID, planID, paymentMethodID, stripeCustomerID, stripeSubscriptionID string) (*models.Subscription, error) {
	plan, err := s.repo.GetPlan(planID)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}

	now := time.Now()
	periodEnd := now.AddDate(0, 1, 0) // 1 month
	if plan.Interval == "year" {
		periodEnd = now.AddDate(1, 0, 0)
	}

	subscription := &models.Subscription{
		ID:                   primitive.NewObjectID(),
		UserID:               userID,
		PlanID:               planID,
		Status:               "active",
		PaymentMethodID:      paymentMethodID,
		StripeCustomerID:     stripeCustomerID,
		StripeSubscriptionID: stripeSubscriptionID,
		CurrentPeriodStart:   now,
		CurrentPeriodEnd:     periodEnd,
		CancelAtPeriodEnd:    false,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	created, err := s.repo.CreateSubscription(ctx, subscription)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpsertStripeLink(ctx, userID, stripeCustomerID, stripeSubscriptionID); err != nil {
		return nil, err
	}

	return created, nil
}

// GetSubscription retrieves user subscription
func (s *PaymentService) GetSubscription(ctx context.Context, userID string) (*models.Subscription, error) {
	return s.repo.GetSubscriptionByUserID(ctx, userID)
}

// CancelSubscription cancels a subscription - Issue #16
func (s *PaymentService) CancelSubscription(ctx context.Context, userID, subscriptionID string) error {
	return s.repo.CancelSubscription(ctx, userID, subscriptionID)
}

// CreatePurchase creates a TVOD/PPV purchase
func (s *PaymentService) CreatePurchase(ctx context.Context, userID, contentID, purchaseType string, amount float64) (*models.Purchase, error) {
	purchase := &models.Purchase{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		ContentID: contentID,
		Type:      purchaseType,
		Amount:    amount,
		Currency:  "USD",
		Status:    "completed",
		CreatedAt: time.Now(),
	}

	if purchaseType == "rent" {
		expiresAt := time.Now().Add(48 * time.Hour)
		purchase.ExpiresAt = &expiresAt
	}

	return s.repo.CreatePurchase(ctx, purchase)
}

// GetUserEntitlements gets user entitlements (subscriptions and PPV) - Issue #16
func (s *PaymentService) GetUserEntitlements(ctx context.Context, userID string) ([]map[string]interface{}, error) {
	subscription, err := s.repo.GetSubscriptionByUserID(ctx, userID)
	if err != nil {
		// No subscription, return empty entitlements
		subscription = nil
	}

	entitlements := []map[string]interface{}{}

	if subscription != nil && subscription.Status == "active" {
		entitlements = append(entitlements, map[string]interface{}{
			"type":       "subscription",
			"plan_id":    subscription.PlanID,
			"status":     subscription.Status,
			"expires_at": subscription.CurrentPeriodEnd,
		})
	}

	purchases, err := s.repo.GetActivePurchasesByUserID(ctx, userID)
	if err == nil {
		for _, purchase := range purchases {
			entitlement := map[string]interface{}{
				"type":       "purchase",
				"content_id": purchase.ContentID,
				"status":     purchase.Status,
			}
			if purchase.ExpiresAt != nil {
				entitlement["expires_at"] = *purchase.ExpiresAt
			}
			entitlements = append(entitlements, entitlement)
		}
	}

	return entitlements, nil
}

// ProcessStripeWebhook processes relevant Stripe webhook events.
func (s *PaymentService) ProcessStripeWebhook(ctx context.Context, eventID, eventType string, object map[string]interface{}, payloadHash string) error {
	if eventID == "" {
		return fmt.Errorf("stripe webhook event id is required")
	}

	shouldProcess, err := s.repo.BeginWebhookEvent(ctx, eventID, eventType, payloadHash)
	if err != nil {
		return err
	}
	if !shouldProcess {
		return nil
	}

	if err := s.applyStripeWebhook(ctx, eventType, object); err != nil {
		markErr := s.repo.MarkWebhookEventFailed(ctx, eventID, trimWebhookError(err.Error()))
		if markErr != nil {
			return fmt.Errorf("webhook apply error: %v (mark failed: %w)", err, markErr)
		}
		return err
	}

	return s.repo.MarkWebhookEventProcessed(ctx, eventID)
}

func (s *PaymentService) applyStripeWebhook(ctx context.Context, eventType string, object map[string]interface{}) error {
	customerID := stripeCustomerID(object)
	subscriptionID := stripeSubscriptionID(eventType, object)
	userID := nestedString(object, "metadata", "user_id")

	if userID == "" {
		resolvedUserID, err := s.repo.ResolveUserIDByStripeIDs(ctx, customerID, subscriptionID)
		if err != nil {
			return err
		}
		userID = resolvedUserID
	}

	if userID == "" {
		return fmt.Errorf(
			"unable to resolve user for stripe event %s (customer=%s subscription=%s)",
			eventType,
			customerID,
			subscriptionID,
		)
	}

	if err := s.repo.UpsertStripeLink(ctx, userID, customerID, subscriptionID); err != nil {
		return err
	}

	switch eventType {
	case "customer.subscription.created", "customer.subscription.updated":
		subscription := &models.Subscription{
			ID:                   primitive.NewObjectID(),
			UserID:               userID,
			PlanID:               subscriptionPlanID(object),
			Status:               normalizeStripeStatus(stringValue(object["status"])),
			PaymentMethodID:      stringValue(object["default_payment_method"]),
			StripeCustomerID:     customerID,
			StripeSubscriptionID: subscriptionID,
			CurrentPeriodStart:   unixToTime(object["current_period_start"]),
			CurrentPeriodEnd:     unixToTime(object["current_period_end"]),
			CancelAtPeriodEnd:    boolValue(object["cancel_at_period_end"]),
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		return s.repo.UpsertSubscriptionByUserID(ctx, subscription)

	case "customer.subscription.deleted":
		return s.repo.UpdateSubscriptionStatusByUserID(ctx, userID, "canceled", true)

	case "invoice.payment_failed":
		return s.repo.UpdateSubscriptionStatusByUserID(ctx, userID, "paused", false)

	case "invoice.payment_succeeded":
		return s.repo.UpdateSubscriptionStatusByUserID(ctx, userID, "active", false)

	default:
		return nil
	}
}

func trimWebhookError(message string) string {
	const maxLen = 1000
	if len(message) <= maxLen {
		return message
	}
	return message[:maxLen]
}

func stripeCustomerID(object map[string]interface{}) string {
	return stringValue(object["customer"])
}

func stripeSubscriptionID(eventType string, object map[string]interface{}) string {
	switch eventType {
	case "customer.subscription.created", "customer.subscription.updated", "customer.subscription.deleted":
		return stringValue(object["id"])
	default:
		return stringValue(object["subscription"])
	}
}

func normalizeStripeStatus(status string) string {
	switch strings.ToLower(status) {
	case "active", "trialing":
		return "active"
	case "past_due", "incomplete", "incomplete_expired", "unpaid":
		return "paused"
	case "canceled":
		return "canceled"
	default:
		return "active"
	}
}

func nestedString(data map[string]interface{}, path ...string) string {
	var current interface{} = data
	for _, key := range path {
		asMap, ok := current.(map[string]interface{})
		if !ok {
			return ""
		}
		current = asMap[key]
	}
	return stringValue(current)
}

func stringValue(v interface{}) string {
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}

func boolValue(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}

func int64Value(v interface{}) int64 {
	switch val := v.(type) {
	case float64:
		return int64(val)
	case int64:
		return val
	case int:
		return int64(val)
	default:
		return 0
	}
}

func unixToTime(v interface{}) time.Time {
	timestamp := int64Value(v)
	if timestamp <= 0 {
		return time.Now()
	}
	return time.Unix(timestamp, 0).UTC()
}

func subscriptionPlanID(object map[string]interface{}) string {
	items, ok := object["items"].(map[string]interface{})
	if !ok {
		return "premium"
	}
	data, ok := items["data"].([]interface{})
	if !ok || len(data) == 0 {
		return "premium"
	}
	first, ok := data[0].(map[string]interface{})
	if !ok {
		return "premium"
	}
	if plan, ok := first["plan"].(map[string]interface{}); ok {
		if id := stringValue(plan["id"]); id != "" {
			return id
		}
	}
	if price, ok := first["price"].(map[string]interface{}); ok {
		if id := stringValue(price["id"]); id != "" {
			return id
		}
	}
	return "premium"
}
