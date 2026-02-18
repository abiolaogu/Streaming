package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// EntitlementProvider defines how content service retrieves user entitlement records.
type EntitlementProvider interface {
	GetUserEntitlements(ctx context.Context, userID, authHeader string) ([]map[string]interface{}, error)
}

// PaymentEntitlementsClient fetches user entitlements from payment-service boundary.
type PaymentEntitlementsClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewPaymentEntitlementsClient creates a payment-service entitlement client.
func NewPaymentEntitlementsClient(baseURL string) *PaymentEntitlementsClient {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	return &PaymentEntitlementsClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GetUserEntitlements retrieves entitlement records for a user.
func (c *PaymentEntitlementsClient) GetUserEntitlements(ctx context.Context, userID, authHeader string) ([]map[string]interface{}, error) {
	if userID == "" {
		return nil, fmt.Errorf("user id is required")
	}
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header is required")
	}

	endpoint := fmt.Sprintf("%s/payments/entitlements/%s", c.baseURL, url.PathEscape(userID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", authHeader)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("payment service returned status %d", res.StatusCode)
	}

	var payload struct {
		Entitlements []map[string]interface{} `json:"entitlements"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Entitlements, nil
}
