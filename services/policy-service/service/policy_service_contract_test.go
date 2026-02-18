package service

import (
	"context"
	"testing"
	"time"
)

func TestContractPolicyV1AllowsActiveSubscription(t *testing.T) {
	svc := NewPolicyService()

	resp, err := svc.EvaluateEntitlement(context.Background(), EntitlementEvaluationRequest{
		ContractVersion: ContractVersionV1,
		ContentID:       "content-1",
		UserID:          "user-1",
		ContentCategory: "movie",
		IsDRMProtected:  true,
		Entitlements: []map[string]interface{}{
			{
				"type":    "subscription",
				"plan_id": "tier3",
				"status":  "active",
			},
		},
	})
	if err != nil {
		t.Fatalf("expected policy evaluation success, got %v", err)
	}
	if resp.ContractVersion != ContractVersionV1 {
		t.Fatalf("expected contract version v1, got %q", resp.ContractVersion)
	}
	if !resp.Decision.HasAccess {
		t.Fatalf("expected access from active subscription")
	}
	if resp.Decision.Reason != "subscription" {
		t.Fatalf("expected subscription reason, got %q", resp.Decision.Reason)
	}
	if resp.Decision.DRMLevel != "1" {
		t.Fatalf("expected tier3 drm level 1, got %q", resp.Decision.DRMLevel)
	}
	if resp.Decision.LicenseURL == "" {
		t.Fatalf("expected license URL when drm is enabled")
	}
}

func TestContractPolicyV1PrefersPurchaseForMatchingContent(t *testing.T) {
	svc := NewPolicyService()

	resp, err := svc.EvaluateEntitlement(context.Background(), EntitlementEvaluationRequest{
		ContractVersion: ContractVersionV1,
		ContentID:       "content-1",
		UserID:          "user-1",
		ContentCategory: "movie",
		Entitlements: []map[string]interface{}{
			{
				"type":    "subscription",
				"plan_id": "tier1",
				"status":  "active",
			},
			{
				"type":       "purchase",
				"content_id": "content-1",
				"status":     "completed",
				"expires_at": "2032-01-02T03:04:05Z",
			},
		},
	})
	if err != nil {
		t.Fatalf("expected policy evaluation success, got %v", err)
	}
	if !resp.Decision.HasAccess || resp.Decision.Reason != "purchased" {
		t.Fatalf("expected purchased access decision, got %+v", resp.Decision)
	}
	if resp.Decision.ExpiresAt == nil {
		t.Fatalf("expected expires_at for purchase decision")
	}
	if got := resp.Decision.ExpiresAt.UTC().Format(time.RFC3339); got != "2032-01-02T03:04:05Z" {
		t.Fatalf("unexpected expires_at value: %s", got)
	}
}

func TestContractPolicyV1DeniesGeoBlockedCountry(t *testing.T) {
	t.Setenv("GEO_BLOCKED_COUNTRIES", "IR,KP")
	svc := NewPolicyService()

	resp, err := svc.EvaluateEntitlement(context.Background(), EntitlementEvaluationRequest{
		ContractVersion: ContractVersionV1,
		ContentID:       "content-1",
		UserID:          "user-1",
		CountryCode:     "ir",
		ContentCategory: "free",
		Entitlements:    []map[string]interface{}{},
	})
	if err != nil {
		t.Fatalf("expected policy evaluation success, got %v", err)
	}
	if resp.Decision.HasAccess {
		t.Fatalf("expected geo-blocked access denial")
	}
	if resp.Decision.Reason != "geo_blocked" {
		t.Fatalf("expected geo_blocked reason, got %q", resp.Decision.Reason)
	}
}

func TestContractPolicyV1RejectsUnsupportedVersion(t *testing.T) {
	svc := NewPolicyService()

	_, err := svc.EvaluateEntitlement(context.Background(), EntitlementEvaluationRequest{
		ContractVersion: "v2",
		ContentID:       "content-1",
		UserID:          "user-1",
	})
	if err == nil {
		t.Fatalf("expected unsupported contract version error")
	}
	if err.Error() != "unsupported contract version: v2" {
		t.Fatalf("unexpected error: %v", err)
	}
}
