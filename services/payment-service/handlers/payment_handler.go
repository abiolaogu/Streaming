package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streamverse/common-go/errors"
	"github.com/streamverse/common-go/logger"
	"github.com/streamverse/payment-service/service"
)

// PaymentHandler handles HTTP requests for payments
type PaymentHandler struct {
	service *service.PaymentService
	logger  *logger.Logger
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(service *service.PaymentService, logger *logger.Logger) *PaymentHandler {
	return &PaymentHandler{
		service: service,
		logger:  logger,
	}
}

// CreateSubscription handles POST /payments/subscribe - Issue #16
func (h *PaymentHandler) CreateSubscription(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("User not authenticated"))
		return
	}

	var req struct {
		PlanID               string `json:"plan_id" binding:"required"`
		PaymentMethodID      string `json:"payment_method_id" binding:"required"`
		StripeCustomerID     string `json:"stripe_customer_id,omitempty"`
		StripeSubscriptionID string `json:"stripe_subscription_id,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	subscription, err := h.service.Subscribe(
		c.Request.Context(),
		userID,
		req.PlanID,
		req.PaymentMethodID,
		req.StripeCustomerID,
		req.StripeSubscriptionID,
	)
	if err != nil {
		h.logger.Error("Failed to subscribe", logger.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

// GetSubscription handles GET /api/v1/payments/subscription
func (h *PaymentHandler) GetSubscription(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("User not authenticated"))
		return
	}

	subscription, err := h.service.GetSubscription(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.NewNotFoundError("No active subscription"))
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// CancelSubscription handles POST /payments/subscribe/{subscription_id}/cancel - Issue #16
func (h *PaymentHandler) CancelSubscription(c *gin.Context) {
	subscriptionID := c.Param("subscription_id")
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("User not authenticated"))
		return
	}

	if err := h.service.CancelSubscription(c.Request.Context(), userID, subscriptionID); err != nil {
		h.logger.Error("Failed to cancel subscription", logger.Error(err))
		c.JSON(http.StatusNotFound, errors.NewNotFoundError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription cancelled"})
}

// PurchaseContent handles POST /payments/purchase - Issue #16
func (h *PaymentHandler) PurchaseContent(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("User not authenticated"))
		return
	}

	var req struct {
		ContentID string  `json:"content_id" binding:"required"`
		Type      string  `json:"type" binding:"required"` // "rent" or "buy"
		Amount    float64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	purchase, err := h.service.CreatePurchase(c.Request.Context(), userID, req.ContentID, req.Type, req.Amount)
	if err != nil {
		h.logger.Error("Failed to create purchase", logger.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, purchase)
}

// GetUserEntitlements handles GET /payments/entitlements/{user_id} - Issue #16
func (h *PaymentHandler) GetUserEntitlements(c *gin.Context) {
	userID := c.Param("user_id")
	currentUser, ok := currentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("User not authenticated"))
		return
	}

	// Users can only access their own entitlements unless admin
	if userID != currentUser && !hasRole(c, "admin") {
		c.JSON(http.StatusForbidden, errors.NewUnauthorizedError("Access denied"))
		return
	}

	entitlements, err := h.service.GetUserEntitlements(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get entitlements", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to get entitlements"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"entitlements": entitlements})
}

// ListPlans handles GET /payments/plans - Issue #16
func (h *PaymentHandler) ListPlans(c *gin.Context) {
	plans := []map[string]interface{}{
		{
			"id":       "tier1",
			"name":     "Basic",
			"price":    4.99,
			"currency": "USD",
			"interval": "month",
			"quality":  "480p",
			"screens":  1,
			"ads":      true,
			"download": false,
		},
		{
			"id":       "tier2",
			"name":     "Pro",
			"price":    12.99,
			"currency": "USD",
			"interval": "month",
			"quality":  "720p",
			"screens":  2,
			"ads":      false,
			"download": true,
		},
		{
			"id":               "tier3",
			"name":             "Premium",
			"price":            19.99,
			"currency":         "USD",
			"interval":         "month",
			"quality":          "4K",
			"screens":          4,
			"ads":              false,
			"download":         true,
			"priority_support": true,
		},
	}
	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

func currentUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	id, ok := userID.(string)
	if !ok || strings.TrimSpace(id) == "" {
		return "", false
	}
	return id, true
}

func hasRole(c *gin.Context, requiredRole string) bool {
	rolesRaw, exists := c.Get("roles")
	if !exists {
		return false
	}
	roles, ok := rolesRaw.([]string)
	if !ok {
		return false
	}
	for _, role := range roles {
		if strings.EqualFold(strings.TrimSpace(role), requiredRole) {
			return true
		}
	}
	return false
}

// HandleStripeWebhook handles POST /payments/webhook - Issue #16
func (h *PaymentHandler) HandleStripeWebhook(c *gin.Context) {
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if webhookSecret == "" {
		c.JSON(http.StatusServiceUnavailable, errors.NewInternalError("Stripe webhook is not configured"))
		return
	}

	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError("Failed to read webhook payload"))
		return
	}

	signatureHeader := c.GetHeader("Stripe-Signature")
	if err := verifyStripeSignature(payload, signatureHeader, webhookSecret, 5*time.Minute); err != nil {
		h.logger.Error("Invalid Stripe signature", logger.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError("Invalid Stripe signature"))
		return
	}

	var event struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Data struct {
			Object map[string]interface{} `json:"object"`
		} `json:"data"`
	}

	if err := json.Unmarshal(payload, &event); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError("Invalid webhook payload"))
		return
	}
	if strings.TrimSpace(event.ID) == "" || strings.TrimSpace(event.Type) == "" {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError("Invalid webhook payload"))
		return
	}

	if err := h.service.ProcessStripeWebhook(
		c.Request.Context(),
		event.ID,
		event.Type,
		event.Data.Object,
		payloadSHA256(payload),
	); err != nil {
		h.logger.Error("Failed to process Stripe webhook", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to process webhook"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Webhook received",
		"id":      event.ID,
		"event":   event.Type,
	})
}

func payloadSHA256(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

func verifyStripeSignature(payload []byte, header, secret string, tolerance time.Duration) error {
	if header == "" {
		return fmt.Errorf("missing Stripe-Signature header")
	}

	parts := strings.Split(header, ",")
	var timestamp string
	var signatures []string

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "t=") {
			timestamp = strings.TrimPrefix(part, "t=")
		}
		if strings.HasPrefix(part, "v1=") {
			signatures = append(signatures, strings.TrimPrefix(part, "v1="))
		}
	}

	if timestamp == "" || len(signatures) == 0 {
		return fmt.Errorf("invalid Stripe-Signature format")
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid Stripe timestamp: %w", err)
	}

	if delta := time.Since(time.Unix(ts, 0)); delta > tolerance || delta < -tolerance {
		return fmt.Errorf("stripe signature timestamp outside tolerance")
	}

	signedPayload := []byte(timestamp + "." + string(payload))
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(signedPayload)
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	for _, sig := range signatures {
		if hmac.Equal([]byte(sig), []byte(expectedSig)) {
			return nil
		}
	}

	return fmt.Errorf("no matching Stripe signature")
}
