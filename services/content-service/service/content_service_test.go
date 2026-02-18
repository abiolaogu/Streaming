package service

import (
	"testing"
	"time"

	"github.com/streamverse/content-service/models"
)

func TestContractEntitlementsPrefersPurchaseForMatchingContent(t *testing.T) {
	records := []map[string]interface{}{
		{
			"type":    "subscription",
			"plan_id": "premium",
			"status":  "active",
		},
		{
			"type":       "purchase",
			"content_id": "content-1",
			"status":     "completed",
			"expires_at": "2030-01-01T00:00:00Z",
		},
	}

	planID, hasSubscription, hasPurchase, expiresAt := evaluateEntitlements(records, "content-1")
	if planID != "premium" {
		t.Fatalf("expected plan premium, got %q", planID)
	}
	if !hasSubscription {
		t.Fatalf("expected active subscription")
	}
	if !hasPurchase {
		t.Fatalf("expected purchase entitlement for content")
	}
	if expiresAt == nil {
		t.Fatalf("expected purchase expiry to be parsed")
	}
}

func TestEvaluateEntitlementsIgnoresPurchaseForDifferentContent(t *testing.T) {
	records := []map[string]interface{}{
		{
			"type":       "purchase",
			"content_id": "another-content",
			"status":     "completed",
		},
	}

	_, _, hasPurchase, _ := evaluateEntitlements(records, "content-1")
	if hasPurchase {
		t.Fatalf("purchase should not apply to different content id")
	}
}

func TestContractGeoBlockingReadsConfiguredCountries(t *testing.T) {
	t.Setenv("GEO_BLOCKED_COUNTRIES", "KP,IR,SY")
	if !isGeoBlocked("ir") {
		t.Fatalf("expected IR to be blocked")
	}
	if isGeoBlocked("US") {
		t.Fatalf("US should not be blocked")
	}
}

func TestToTimeParsesRFC3339(t *testing.T) {
	parsed := toTime("2029-05-01T10:30:00Z")
	if parsed == nil {
		t.Fatalf("expected time to parse")
	}
	if parsed.UTC().Format(time.RFC3339) != "2029-05-01T10:30:00Z" {
		t.Fatalf("unexpected parsed time: %s", parsed.UTC().Format(time.RFC3339))
	}
}

func TestContractDRMLevelForPlanMappings(t *testing.T) {
	if got := drmLevelForPlan("tier3"); got != "1" {
		t.Fatalf("expected tier3 DRM level 1, got %s", got)
	}
	if got := drmLevelForPlan("tier2"); got != "2" {
		t.Fatalf("expected tier2 DRM level 2, got %s", got)
	}
	if got := drmLevelForPlan("unknown"); got != "3" {
		t.Fatalf("expected unknown plan DRM level 3, got %s", got)
	}
}

func TestContractDefaultLicenseURLHonorsEnvironment(t *testing.T) {
	t.Setenv("DRM_LICENSE_URL", "https://drm.example.com/license")

	if got := defaultLicenseURL(); got != "https://drm.example.com/license" {
		t.Fatalf("expected DRM license URL from env, got %s", got)
	}
}

func TestToTimeParsesBSONDateMap(t *testing.T) {
	parsed := toTime(map[string]interface{}{"$date": "2030-06-30T12:00:00Z"})
	if parsed == nil {
		t.Fatalf("expected bson date map to parse")
	}
	if parsed.UTC().Format(time.RFC3339) != "2030-06-30T12:00:00Z" {
		t.Fatalf("unexpected parsed bson date: %s", parsed.UTC().Format(time.RFC3339))
	}
}

func TestContractLocalEntitlementEvaluationForFreeGeoBlockedDRM(t *testing.T) {
	t.Setenv("GEO_BLOCKED_COUNTRIES", "KP")
	t.Setenv("DRM_LICENSE_URL", "https://license.test/v1")

	content := &models.Content{
		Category:       "free",
		IsDRMProtected: true,
	}

	entitlement := evaluateLocalEntitlement(content, "content-1", "user-1", "US", nil)
	if !entitlement.HasAccess || entitlement.Reason != "free" {
		t.Fatalf("expected free access decision, got %+v", entitlement)
	}
	if entitlement.LicenseURL != "https://license.test/v1" {
		t.Fatalf("expected drm license url from env")
	}

	blocked := evaluateLocalEntitlement(content, "content-1", "user-1", "kp", nil)
	if blocked.HasAccess {
		t.Fatalf("expected geo blocked access denial")
	}
	if blocked.Reason != "geo_blocked" {
		t.Fatalf("expected geo_blocked reason, got %s", blocked.Reason)
	}
}

func TestToStringAndToTimeHandleUnexpectedInput(t *testing.T) {
	if got := toString(123); got != "" {
		t.Fatalf("expected non-string to return empty string, got %q", got)
	}
	if got := toTime(map[string]interface{}{"$date": "invalid"}); got != nil {
		t.Fatalf("expected invalid bson date to return nil")
	}
	if got := toTime("not-a-date"); got != nil {
		t.Fatalf("expected invalid date string to return nil")
	}
}
