package service

import (
	"context"
	"testing"
	"time"

	"github.com/streamverse/common-go/logger"
	"github.com/streamverse/payment-service/models"
)

func TestContractReplayFailedWebhookEventsProcessesFailedBatch(t *testing.T) {
	repo := &replayRepo{
		failedEvents: []models.WebhookEvent{
			{
				EventID:     "evt_failed_1",
				EventType:   "invoice.payment_succeeded",
				PayloadHash: "hash1",
				EventObject: map[string]interface{}{
					"customer":     "cus_1",
					"subscription": "sub_1",
					"metadata": map[string]interface{}{
						"user_id": "user-1",
					},
				},
				Status: "failed",
			},
		},
	}

	svc := NewPaymentService(repo)
	replayed, err := svc.ReplayFailedWebhookEvents(context.Background(), 10)
	if err != nil {
		t.Fatalf("expected replay to succeed, got %v", err)
	}
	if replayed != 1 {
		t.Fatalf("expected one replayed event, got %d", replayed)
	}
	if repo.processedCount != 1 {
		t.Fatalf("expected one processed marker, got %d", repo.processedCount)
	}
	if repo.statusUpdateCount != 1 {
		t.Fatalf("expected one status update, got %d", repo.statusUpdateCount)
	}
}

func TestWebhookReconciliationConfigRespectsResidencyStrictMode(t *testing.T) {
	cfg := WebhookReconciliationConfig{
		Enabled:       true,
		Region:        "us-east-1",
		ResidencyMode: "strict",
		AllowedRegions: map[string]struct{}{
			"eu-west-1": {},
		},
	}

	if cfg.residencyAllowsRegion() {
		t.Fatalf("expected strict residency config to block disallowed region")
	}
}

func TestWebhookReconciliationWorkerRunOnce(t *testing.T) {
	repo := &replayRepo{
		failedEvents: []models.WebhookEvent{
			{
				EventID:     "evt_failed_2",
				EventType:   "invoice.payment_succeeded",
				PayloadHash: "hash2",
				EventObject: map[string]interface{}{
					"metadata": map[string]interface{}{
						"user_id": "user-2",
					},
				},
				Status: "failed",
			},
		},
	}

	svc := NewPaymentService(repo)
	log := logger.NewDefault()
	worker := NewWebhookReconciliationWorker(svc, log, WebhookReconciliationConfig{
		Enabled:   true,
		Interval:  10 * time.Millisecond,
		BatchSize: 10,
		Region:    "us-east-1",
	})

	worker.runOnce(context.Background())
	if repo.processedCount != 1 {
		t.Fatalf("expected worker to process one failed event")
	}
}

type replayRepo struct {
	failedEvents        []models.WebhookEvent
	processedCount      int
	failedCount         int
	statusUpdateCount   int
	beginCalls          int
	createdSubscription *models.Subscription
}

func (r *replayRepo) GetPlan(planID string) (*models.Plan, error) {
	return &models.Plan{ID: planID, Interval: "month"}, nil
}

func (r *replayRepo) CreateSubscription(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {
	r.createdSubscription = subscription
	return subscription, nil
}

func (r *replayRepo) GetSubscriptionByUserID(ctx context.Context, userID string) (*models.Subscription, error) {
	return nil, nil
}

func (r *replayRepo) CancelSubscription(ctx context.Context, userID, subscriptionID string) error {
	return nil
}

func (r *replayRepo) CreatePurchase(ctx context.Context, purchase *models.Purchase) (*models.Purchase, error) {
	return purchase, nil
}

func (r *replayRepo) UpsertSubscriptionByUserID(ctx context.Context, subscription *models.Subscription) error {
	return nil
}

func (r *replayRepo) UpdateSubscriptionStatusByUserID(ctx context.Context, userID, status string, cancelAtPeriodEnd bool) error {
	r.statusUpdateCount++
	return nil
}

func (r *replayRepo) UpsertStripeLink(ctx context.Context, userID, customerID, subscriptionID string) error {
	return nil
}

func (r *replayRepo) ResolveUserIDByStripeIDs(ctx context.Context, customerID, subscriptionID string) (string, error) {
	return "", nil
}

func (r *replayRepo) GetActivePurchasesByUserID(ctx context.Context, userID string) ([]models.Purchase, error) {
	return nil, nil
}

func (r *replayRepo) BeginWebhookEvent(ctx context.Context, eventID, eventType, payloadHash string, eventObject map[string]interface{}) (bool, error) {
	r.beginCalls++
	return true, nil
}

func (r *replayRepo) MarkWebhookEventProcessed(ctx context.Context, eventID string) error {
	r.processedCount++
	return nil
}

func (r *replayRepo) MarkWebhookEventFailed(ctx context.Context, eventID, lastError string) error {
	r.failedCount++
	return nil
}

func (r *replayRepo) ListFailedWebhookEvents(ctx context.Context, limit int) ([]models.WebhookEvent, error) {
	if limit <= 0 || limit >= len(r.failedEvents) {
		return r.failedEvents, nil
	}
	return r.failedEvents[:limit], nil
}
