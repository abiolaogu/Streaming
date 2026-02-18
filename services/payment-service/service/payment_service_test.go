package service

import (
	"strings"
	"testing"
	"time"
)

func TestStripeSubscriptionIDParsing(t *testing.T) {
	subscriptionObject := map[string]interface{}{
		"id": "sub_123",
	}
	invoiceObject := map[string]interface{}{
		"subscription": "sub_456",
	}

	if got := stripeSubscriptionID("customer.subscription.updated", subscriptionObject); got != "sub_123" {
		t.Fatalf("expected subscription event to use object id, got %q", got)
	}

	if got := stripeSubscriptionID("invoice.payment_succeeded", invoiceObject); got != "sub_456" {
		t.Fatalf("expected invoice event to use subscription field, got %q", got)
	}
}

func TestNormalizeStripeStatus(t *testing.T) {
	if got := normalizeStripeStatus("past_due"); got != "paused" {
		t.Fatalf("expected paused, got %q", got)
	}
	if got := normalizeStripeStatus("active"); got != "active" {
		t.Fatalf("expected active, got %q", got)
	}
	if got := normalizeStripeStatus("canceled"); got != "canceled" {
		t.Fatalf("expected canceled, got %q", got)
	}
}

func TestWebhookHelperParsing(t *testing.T) {
	metadata := map[string]interface{}{
		"metadata": map[string]interface{}{
			"user_id": "user_123",
		},
	}
	if got := nestedString(metadata, "metadata", "user_id"); got != "user_123" {
		t.Fatalf("expected nested user id, got %q", got)
	}
	if got := nestedString(metadata, "metadata", "missing"); got != "" {
		t.Fatalf("expected missing path to return empty string, got %q", got)
	}

	if !boolValue(true) {
		t.Fatalf("expected boolValue(true) to be true")
	}
	if boolValue("true") {
		t.Fatalf("expected non-bool input to return false")
	}

	if got := int64Value(float64(42)); got != 42 {
		t.Fatalf("expected float64 conversion to 42, got %d", got)
	}
	if got := int64Value(int(7)); got != 7 {
		t.Fatalf("expected int conversion to 7, got %d", got)
	}
}

func TestWebhookTimeAndPlanHelpers(t *testing.T) {
	if got := unixToTime(int64(0)); got.IsZero() {
		t.Fatalf("expected unixToTime fallback to current time")
	}

	ts := int64(1730000000)
	got := unixToTime(ts)
	if got.Unix() != ts {
		t.Fatalf("expected unix timestamp %d, got %d", ts, got.Unix())
	}
	if got.Location() != time.UTC {
		t.Fatalf("expected UTC timestamp")
	}

	object := map[string]interface{}{
		"items": map[string]interface{}{
			"data": []interface{}{
				map[string]interface{}{
					"price": map[string]interface{}{
						"id": "price_tier3",
					},
				},
			},
		},
	}
	if planID := subscriptionPlanID(object); planID != "price_tier3" {
		t.Fatalf("expected price-based plan id, got %s", planID)
	}
}

func TestTrimWebhookError(t *testing.T) {
	short := "short error"
	if got := trimWebhookError(short); got != short {
		t.Fatalf("expected short error to stay unchanged")
	}

	long := strings.Repeat("x", 1200)
	trimmed := trimWebhookError(long)
	if len(trimmed) != 1000 {
		t.Fatalf("expected trimmed error length 1000, got %d", len(trimmed))
	}
}
