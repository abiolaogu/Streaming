package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/streamverse/payment-service/models"
	"github.com/streamverse/payment-service/repository"
)

// PaymentService handles payment business logic
type PaymentService struct {
	repo *repository.PaymentRepository
}

// NewPaymentService creates a new payment service
func NewPaymentService(repo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{
		repo: repo,
	}
}

// Subscribe subscribes user to a plan
func (s *PaymentService) Subscribe(ctx context.Context, userID, planID, paymentMethodID string) (*models.Subscription, error) {
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
		ID:                 primitive.NewObjectID(),
		UserID:             userID,
		PlanID:            planID,
		Status:            "active",
		PaymentMethodID:   paymentMethodID,
		CurrentPeriodStart: now,
		CurrentPeriodEnd:  periodEnd,
		CancelAtPeriodEnd: false,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	return s.repo.CreateSubscription(ctx, subscription)
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
			"type":        "subscription",
			"plan_id":     subscription.PlanID,
			"status":      subscription.Status,
			"expires_at":  subscription.CurrentPeriodEnd,
		})
	}

	// TODO: Add PPV purchases to entitlements
	// purchases, _ := s.repo.GetUserPurchases(ctx, userID)

	return entitlements, nil
}

