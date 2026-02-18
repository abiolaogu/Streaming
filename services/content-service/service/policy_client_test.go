package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContractEntitlementPolicyClientEvaluatesDecision(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST method, got %s", r.Method)
		}
		if r.URL.Path != "/policy/v1/entitlements/evaluate" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer token-123" {
			t.Fatalf("expected forwarded auth header, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"contract_version":"v1","decision":{"has_access":true,"reason":"subscription","drm_level":"2","license_url":"https://drm.test/license"}}`))
	}))
	defer server.Close()

	client := NewEntitlementPolicyClient(server.URL)
	entitlement, err := client.EvaluateEntitlement(context.Background(), PolicyEntitlementInput{
		ContentID:       "content-1",
		UserID:          "user-1",
		ContentCategory: "movie",
		IsDRMProtected:  true,
		Entitlements: []map[string]interface{}{
			{"type": "subscription", "plan_id": "tier2", "status": "active"},
		},
	}, "Bearer token-123")
	if err != nil {
		t.Fatalf("expected policy evaluation success, got %v", err)
	}

	if entitlement == nil || !entitlement.HasAccess {
		t.Fatalf("expected entitlement access to be granted")
	}
	if entitlement.Reason != "subscription" {
		t.Fatalf("expected subscription reason, got %q", entitlement.Reason)
	}
	if entitlement.DRMLevel != "2" {
		t.Fatalf("expected drm level 2, got %q", entitlement.DRMLevel)
	}
}

func TestPolicyClientRejectsUnexpectedContractVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"contract_version":"v2","decision":{"has_access":false,"reason":"subscription_required"}}`))
	}))
	defer server.Close()

	client := NewEntitlementPolicyClient(server.URL)
	_, err := client.EvaluateEntitlement(context.Background(), PolicyEntitlementInput{
		ContentID: "content-1",
		UserID:    "user-1",
	}, "Bearer token-123")
	if err == nil {
		t.Fatalf("expected unsupported contract version error")
	}
}

func TestPolicyClientRequiresAuthorizationHeader(t *testing.T) {
	client := NewEntitlementPolicyClient("http://example.test")
	_, err := client.EvaluateEntitlement(context.Background(), PolicyEntitlementInput{
		ContentID: "content-1",
		UserID:    "user-1",
	}, "")
	if err == nil {
		t.Fatalf("expected missing authorization header error")
	}
}
