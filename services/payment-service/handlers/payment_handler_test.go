package handlers

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streamverse/common-go/logger"
	"github.com/streamverse/payment-service/models"
	"github.com/streamverse/payment-service/service"
)

func TestVerifyStripeSignatureAcceptsValidSignature(t *testing.T) {
	payload := []byte(`{"id":"evt_test","type":"invoice.payment_succeeded"}`)
	secret := "whsec_test_secret"
	timestamp := time.Now().Unix()

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("%d.%s", timestamp, payload)))
	signature := hex.EncodeToString(mac.Sum(nil))
	header := fmt.Sprintf("t=%d,v1=%s", timestamp, signature)

	if err := verifyStripeSignature(payload, header, secret, 5*time.Minute); err != nil {
		t.Fatalf("expected valid signature, got error: %v", err)
	}
}

func TestVerifyStripeSignatureRejectsInvalidSignature(t *testing.T) {
	payload := []byte(`{"id":"evt_test","type":"invoice.payment_succeeded"}`)
	secret := "whsec_test_secret"
	header := fmt.Sprintf("t=%d,v1=%s", time.Now().Unix(), "invalid")

	if err := verifyStripeSignature(payload, header, secret, 5*time.Minute); err == nil {
		t.Fatalf("expected invalid signature to be rejected")
	}
}

func TestContractStripeWebhookReplayIsIdempotent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "whsec_test_secret"
	t.Setenv("STRIPE_WEBHOOK_SECRET", secret)

	repo := newTestWebhookRepo()
	handler := newWebhookTestHandler(t, repo)
	payload := []byte(`{"id":"evt_replay_1","type":"invoice.payment_succeeded","data":{"object":{"customer":"cus_123","subscription":"sub_123","metadata":{"user_id":"user_123"}}}}`)

	first := performWebhookRequest(handler, payload, secret)
	if first.Code != http.StatusOK {
		t.Fatalf("expected first delivery 200, got %d body=%s", first.Code, first.Body.String())
	}

	second := performWebhookRequest(handler, payload, secret)
	if second.Code != http.StatusOK {
		t.Fatalf("expected replay delivery 200, got %d body=%s", second.Code, second.Body.String())
	}

	if got := repo.statusUpdateCount(); got != 1 {
		t.Fatalf("expected business mutation to run once, got %d", got)
	}

	event := repo.webhookEvent("evt_replay_1")
	if event == nil {
		t.Fatalf("expected webhook event to be persisted")
	}
	if event.Status != "processed" {
		t.Fatalf("expected webhook status processed, got %q", event.Status)
	}
	if event.Attempts != 1 {
		t.Fatalf("expected one processing attempt, got %d", event.Attempts)
	}
}

func TestContractStripeWebhookRetriesFailedEvent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "whsec_test_secret"
	t.Setenv("STRIPE_WEBHOOK_SECRET", secret)

	repo := newTestWebhookRepo()
	handler := newWebhookTestHandler(t, repo)
	payload := []byte(`{"id":"evt_retry_1","type":"invoice.payment_succeeded","data":{"object":{"customer":"cus_retry","subscription":"sub_retry"}}}`)

	first := performWebhookRequest(handler, payload, secret)
	if first.Code != http.StatusInternalServerError {
		t.Fatalf("expected first attempt to fail with 500, got %d body=%s", first.Code, first.Body.String())
	}

	eventAfterFirst := repo.webhookEvent("evt_retry_1")
	if eventAfterFirst == nil || eventAfterFirst.Status != "failed" {
		t.Fatalf("expected webhook event to be marked failed after first attempt")
	}

	repo.setStripeMapping("cus_retry", "sub_retry", "user_retry")

	second := performWebhookRequest(handler, payload, secret)
	if second.Code != http.StatusOK {
		t.Fatalf("expected second attempt to succeed, got %d body=%s", second.Code, second.Body.String())
	}

	if got := repo.statusUpdateCount(); got != 1 {
		t.Fatalf("expected successful mutation after retry, got %d", got)
	}

	eventAfterSecond := repo.webhookEvent("evt_retry_1")
	if eventAfterSecond == nil {
		t.Fatalf("expected webhook event after retry")
	}
	if eventAfterSecond.Status != "processed" {
		t.Fatalf("expected status processed after retry, got %q", eventAfterSecond.Status)
	}
	if eventAfterSecond.Attempts != 2 {
		t.Fatalf("expected two attempts after retry, got %d", eventAfterSecond.Attempts)
	}
}

func performWebhookRequest(handler *PaymentHandler, payload []byte, secret string) *httptest.ResponseRecorder {
	router := gin.New()
	router.POST("/payments/webhook", handler.HandleStripeWebhook)

	req := httptest.NewRequest(http.MethodPost, "/payments/webhook", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Stripe-Signature", stripeSignatureHeader(payload, secret, time.Now().Unix()))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

func stripeSignatureHeader(payload []byte, secret string, timestamp int64) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("%d.%s", timestamp, payload)))
	signature := hex.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("t=%d,v1=%s", timestamp, signature)
}

func newWebhookTestHandler(t *testing.T, repo *testWebhookRepo) *PaymentHandler {
	t.Helper()

	log := logger.NewDefault()
	paymentService := service.NewPaymentService(repo)
	return NewPaymentHandler(paymentService, log)
}

type testWebhookEvent struct {
	Status      string
	PayloadHash string
	Attempts    int
	LastError   string
}

type testWebhookRepo struct {
	mu                 sync.Mutex
	events             map[string]*testWebhookEvent
	customerToUser     map[string]string
	subscriptionToUser map[string]string
	statusUpdates      int
}

func newTestWebhookRepo() *testWebhookRepo {
	return &testWebhookRepo{
		events:             make(map[string]*testWebhookEvent),
		customerToUser:     make(map[string]string),
		subscriptionToUser: make(map[string]string),
	}
}

func (r *testWebhookRepo) setStripeMapping(customerID, subscriptionID, userID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if customerID != "" {
		r.customerToUser[customerID] = userID
	}
	if subscriptionID != "" {
		r.subscriptionToUser[subscriptionID] = userID
	}
}

func (r *testWebhookRepo) statusUpdateCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.statusUpdates
}

func (r *testWebhookRepo) webhookEvent(eventID string) *testWebhookEvent {
	r.mu.Lock()
	defer r.mu.Unlock()

	event, ok := r.events[eventID]
	if !ok {
		return nil
	}
	copy := *event
	return &copy
}

func (r *testWebhookRepo) GetPlan(planID string) (*models.Plan, error) {
	return &models.Plan{ID: planID, Interval: "month"}, nil
}

func (r *testWebhookRepo) CreateSubscription(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {
	return subscription, nil
}

func (r *testWebhookRepo) GetSubscriptionByUserID(ctx context.Context, userID string) (*models.Subscription, error) {
	return nil, fmt.Errorf("not found")
}

func (r *testWebhookRepo) CancelSubscription(ctx context.Context, userID, subscriptionID string) error {
	return nil
}

func (r *testWebhookRepo) CreatePurchase(ctx context.Context, purchase *models.Purchase) (*models.Purchase, error) {
	return purchase, nil
}

func (r *testWebhookRepo) UpsertSubscriptionByUserID(ctx context.Context, subscription *models.Subscription) error {
	return nil
}

func (r *testWebhookRepo) UpdateSubscriptionStatusByUserID(ctx context.Context, userID, status string, cancelAtPeriodEnd bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.statusUpdates++
	return nil
}

func (r *testWebhookRepo) UpsertStripeLink(ctx context.Context, userID, customerID, subscriptionID string) error {
	r.setStripeMapping(customerID, subscriptionID, userID)
	return nil
}

func (r *testWebhookRepo) ResolveUserIDByStripeIDs(ctx context.Context, customerID, subscriptionID string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if subscriptionID != "" {
		if userID, ok := r.subscriptionToUser[subscriptionID]; ok {
			return userID, nil
		}
	}
	if customerID != "" {
		if userID, ok := r.customerToUser[customerID]; ok {
			return userID, nil
		}
	}
	return "", nil
}

func (r *testWebhookRepo) GetActivePurchasesByUserID(ctx context.Context, userID string) ([]models.Purchase, error) {
	return nil, nil
}

func (r *testWebhookRepo) BeginWebhookEvent(ctx context.Context, eventID, eventType, payloadHash string, eventObject map[string]interface{}) (bool, error) {
	if eventID == "" {
		return false, fmt.Errorf("missing event id")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	event, exists := r.events[eventID]
	if !exists {
		r.events[eventID] = &testWebhookEvent{
			Status:      "processing",
			PayloadHash: payloadHash,
			Attempts:    1,
		}
		return true, nil
	}

	if event.PayloadHash != "" && payloadHash != "" && event.PayloadHash != payloadHash {
		return false, fmt.Errorf("payload mismatch for event %s", eventID)
	}

	switch event.Status {
	case "processed", "processing":
		return false, nil
	case "failed":
		event.Status = "processing"
		event.LastError = ""
		event.Attempts++
		event.PayloadHash = payloadHash
		return true, nil
	default:
		return false, nil
	}
}

func (r *testWebhookRepo) ListFailedWebhookEvents(ctx context.Context, limit int) ([]models.WebhookEvent, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var events []models.WebhookEvent
	for eventID, event := range r.events {
		if event.Status != "failed" {
			continue
		}
		events = append(events, models.WebhookEvent{
			EventID:     eventID,
			EventType:   "invoice.payment_succeeded",
			PayloadHash: event.PayloadHash,
			Status:      event.Status,
			Attempts:    event.Attempts,
			LastError:   event.LastError,
		})
		if limit > 0 && len(events) >= limit {
			break
		}
	}
	return events, nil
}

func (r *testWebhookRepo) MarkWebhookEventProcessed(ctx context.Context, eventID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	event, exists := r.events[eventID]
	if !exists {
		return fmt.Errorf("event not found")
	}
	event.Status = "processed"
	event.LastError = ""
	return nil
}

func (r *testWebhookRepo) MarkWebhookEventFailed(ctx context.Context, eventID, lastError string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	event, exists := r.events[eventID]
	if !exists {
		return fmt.Errorf("event not found")
	}
	event.Status = "failed"
	event.LastError = lastError
	return nil
}
