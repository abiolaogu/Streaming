package service

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	// ContractVersionV1 is the current stable policy contract.
	ContractVersionV1 = "v1"
)

// EntitlementEvaluationRequest captures versioned policy input.
type EntitlementEvaluationRequest struct {
	ContractVersion string                   `json:"contract_version"`
	ContentID       string                   `json:"content_id"`
	UserID          string                   `json:"user_id"`
	CountryCode     string                   `json:"country_code"`
	ContentCategory string                   `json:"content_category"`
	IsDRMProtected  bool                     `json:"is_drm_protected"`
	Entitlements    []map[string]interface{} `json:"entitlements"`
}

// EntitlementDecision captures normalized policy output.
type EntitlementDecision struct {
	HasAccess  bool       `json:"has_access"`
	Reason     string     `json:"reason"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	DRMLevel   string     `json:"drm_level,omitempty"`
	LicenseURL string     `json:"license_url,omitempty"`
}

// EntitlementEvaluationResponse wraps a versioned decision payload.
type EntitlementEvaluationResponse struct {
	ContractVersion string              `json:"contract_version"`
	Decision        EntitlementDecision `json:"decision"`
}

// PolicyService evaluates entitlement decisions behind a dedicated service boundary.
type PolicyService struct{}

// NewPolicyService creates a policy service instance.
func NewPolicyService() *PolicyService {
	return &PolicyService{}
}

// EvaluateEntitlement applies versioned entitlement policy evaluation.
func (s *PolicyService) EvaluateEntitlement(ctx context.Context, req EntitlementEvaluationRequest) (*EntitlementEvaluationResponse, error) {
	_ = ctx

	if strings.TrimSpace(req.ContractVersion) == "" {
		req.ContractVersion = ContractVersionV1
	}
	if req.ContractVersion != ContractVersionV1 {
		return nil, fmt.Errorf("unsupported contract version: %s", req.ContractVersion)
	}
	if strings.TrimSpace(req.ContentID) == "" {
		return nil, fmt.Errorf("content_id is required")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	planID, hasSubscription, hasPurchase, purchaseExpiresAt := evaluateEntitlements(req.Entitlements, req.ContentID)
	decision := EntitlementDecision{
		HasAccess: false,
		Reason:    "subscription_required",
		DRMLevel:  drmLevelForPlan(planID),
	}

	if hasPurchase {
		decision.HasAccess = true
		decision.Reason = "purchased"
		decision.ExpiresAt = purchaseExpiresAt
	} else if hasSubscription {
		decision.HasAccess = true
		decision.Reason = "subscription"
	} else if isFreeCategory(req.ContentCategory) {
		decision.HasAccess = true
		decision.Reason = "free"
	}

	if req.IsDRMProtected {
		decision.LicenseURL = defaultLicenseURL()
	}

	if isGeoBlocked(req.CountryCode) {
		decision.HasAccess = false
		decision.Reason = "geo_blocked"
		decision.ExpiresAt = nil
	}

	return &EntitlementEvaluationResponse{
		ContractVersion: ContractVersionV1,
		Decision:        decision,
	}, nil
}

func isFreeCategory(category string) bool {
	switch strings.ToLower(strings.TrimSpace(category)) {
	case "free", "avod":
		return true
	default:
		return false
	}
}

func isGeoBlocked(countryCode string) bool {
	countryCode = strings.ToUpper(strings.TrimSpace(countryCode))
	if countryCode == "" {
		return false
	}

	blocked := strings.Split(os.Getenv("GEO_BLOCKED_COUNTRIES"), ",")
	for _, country := range blocked {
		if strings.ToUpper(strings.TrimSpace(country)) == countryCode {
			return true
		}
	}
	return false
}

func evaluateEntitlements(records []map[string]interface{}, contentID string) (planID string, hasSubscription bool, hasPurchase bool, purchaseExpiresAt *time.Time) {
	for _, record := range records {
		recordType := strings.ToLower(toString(record["type"]))
		switch recordType {
		case "subscription":
			status := strings.ToLower(toString(record["status"]))
			if status == "active" || status == "trialing" {
				hasSubscription = true
				if planID == "" {
					planID = toString(record["plan_id"])
				}
			}
		case "purchase":
			if toString(record["content_id"]) != contentID {
				continue
			}
			status := strings.ToLower(toString(record["status"]))
			if status != "" && status != "completed" {
				continue
			}
			hasPurchase = true
			if parsed := toTime(record["expires_at"]); parsed != nil {
				purchaseExpiresAt = parsed
			}
		}
	}
	return planID, hasSubscription, hasPurchase, purchaseExpiresAt
}

func toString(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}

func toTime(value interface{}) *time.Time {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case string:
		if v == "" {
			return nil
		}
		if parsed, err := time.Parse(time.RFC3339, v); err == nil {
			return &parsed
		}
	case map[string]interface{}:
		if dateStr, ok := v["$date"].(string); ok {
			if parsed, err := time.Parse(time.RFC3339, dateStr); err == nil {
				return &parsed
			}
		}
	}
	return nil
}

func drmLevelForPlan(planID string) string {
	switch planID {
	case "premium", "tier3":
		return "1"
	case "standard", "pro", "tier2":
		return "2"
	default:
		return "3"
	}
}

func defaultLicenseURL() string {
	if value := strings.TrimSpace(os.Getenv("DRM_LICENSE_URL")); value != "" {
		return value
	}
	return "https://drm.streamverse.io/license"
}
