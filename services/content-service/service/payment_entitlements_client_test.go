package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContractPaymentEntitlementsClientFetchesEntitlements(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET method, got %s", r.Method)
		}
		if r.URL.Path != "/payments/entitlements/user-123" {
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer token-123" {
			t.Fatalf("authorization header not forwarded")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"entitlements":[{"type":"subscription","plan_id":"tier3","status":"active"}]}`))
	}))
	defer server.Close()

	client := NewPaymentEntitlementsClient(server.URL)
	entitlements, err := client.GetUserEntitlements(context.Background(), "user-123", "Bearer token-123")
	if err != nil {
		t.Fatalf("expected successful entitlement fetch, got error: %v", err)
	}

	if len(entitlements) != 1 {
		t.Fatalf("expected one entitlement record, got %d", len(entitlements))
	}
	if entitlements[0]["plan_id"] != "tier3" {
		t.Fatalf("expected plan_id tier3, got %v", entitlements[0]["plan_id"])
	}
}

func TestPaymentEntitlementsClientRejectsMissingAuthorization(t *testing.T) {
	client := NewPaymentEntitlementsClient("http://localhost:9999")
	if _, err := client.GetUserEntitlements(context.Background(), "user-123", ""); err == nil {
		t.Fatalf("expected missing authorization to fail")
	}
}

func TestPaymentEntitlementsClientReturnsErrorOnNonSuccessStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "upstream error", http.StatusBadGateway)
	}))
	defer server.Close()

	client := NewPaymentEntitlementsClient(server.URL)
	if _, err := client.GetUserEntitlements(context.Background(), "user-123", "Bearer token-123"); err == nil {
		t.Fatalf("expected non-success upstream status to fail")
	}
}
