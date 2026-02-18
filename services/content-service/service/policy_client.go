package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/streamverse/content-service/models"
)

const policyContractVersionV1 = "v1"

// PolicyEntitlementInput captures the data required by policy-service.
type PolicyEntitlementInput struct {
	ContentID       string
	UserID          string
	CountryCode     string
	ContentCategory string
	IsDRMProtected  bool
	Entitlements    []map[string]interface{}
}

// EntitlementPolicyProvider evaluates entitlement access via policy-service boundary.
type EntitlementPolicyProvider interface {
	EvaluateEntitlement(ctx context.Context, input PolicyEntitlementInput, authHeader string) (*models.Entitlement, error)
}

// EntitlementPolicyClient calls policy-service.
type EntitlementPolicyClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewEntitlementPolicyClient creates a policy-service client.
func NewEntitlementPolicyClient(baseURL string) *EntitlementPolicyClient {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = "http://localhost:8090"
	}

	return &EntitlementPolicyClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// EvaluateEntitlement delegates access evaluation to policy-service.
func (c *EntitlementPolicyClient) EvaluateEntitlement(ctx context.Context, input PolicyEntitlementInput, authHeader string) (*models.Entitlement, error) {
	if strings.TrimSpace(input.ContentID) == "" {
		return nil, fmt.Errorf("content id is required")
	}
	if strings.TrimSpace(input.UserID) == "" {
		return nil, fmt.Errorf("user id is required")
	}
	if strings.TrimSpace(authHeader) == "" {
		return nil, fmt.Errorf("authorization header is required")
	}

	requestBody := map[string]interface{}{
		"contract_version": policyContractVersionV1,
		"content_id":       input.ContentID,
		"user_id":          input.UserID,
		"country_code":     input.CountryCode,
		"content_category": input.ContentCategory,
		"is_drm_protected": input.IsDRMProtected,
		"entitlements":     input.Entitlements,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/policy/v1/entitlements/evaluate", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("policy service returned status %d", res.StatusCode)
	}

	var payload struct {
		ContractVersion string `json:"contract_version"`
		Decision        struct {
			HasAccess  bool       `json:"has_access"`
			Reason     string     `json:"reason"`
			ExpiresAt  *time.Time `json:"expires_at"`
			DRMLevel   string     `json:"drm_level"`
			LicenseURL string     `json:"license_url"`
		} `json:"decision"`
	}

	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if payload.ContractVersion != policyContractVersionV1 {
		return nil, fmt.Errorf("unsupported policy contract version: %s", payload.ContractVersion)
	}

	return &models.Entitlement{
		ContentID:  input.ContentID,
		UserID:     input.UserID,
		HasAccess:  payload.Decision.HasAccess,
		Reason:     payload.Decision.Reason,
		ExpiresAt:  payload.Decision.ExpiresAt,
		DRMLevel:   payload.Decision.DRMLevel,
		LicenseURL: payload.Decision.LicenseURL,
	}, nil
}
