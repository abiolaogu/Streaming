package handlers

import (
	"net/http"

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
	userID, _ := c.Get("user_id")

	var req struct {
		PlanID          string `json:"plan_id" binding:"required"`
		PaymentMethodID string `json:"payment_method_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	subscription, err := h.service.Subscribe(c.Request.Context(), userID.(string), req.PlanID, req.PaymentMethodID)
	if err != nil {
		h.logger.Error("Failed to subscribe", logger.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

// GetSubscription handles GET /api/v1/payments/subscription
func (h *PaymentHandler) GetSubscription(c *gin.Context) {
	userID, _ := c.Get("user_id")

	subscription, err := h.service.GetSubscription(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, errors.NewNotFoundError("No active subscription"))
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// CancelSubscription handles POST /payments/subscribe/{subscription_id}/cancel - Issue #16
func (h *PaymentHandler) CancelSubscription(c *gin.Context) {
	subscriptionID := c.Param("subscription_id")
	userID, _ := c.Get("user_id")

	if err := h.service.CancelSubscription(c.Request.Context(), userID.(string), subscriptionID); err != nil {
		h.logger.Error("Failed to cancel subscription", logger.Error(err))
		c.JSON(http.StatusNotFound, errors.NewNotFoundError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription cancelled"})
}

// PurchaseContent handles POST /payments/purchase - Issue #16
func (h *PaymentHandler) PurchaseContent(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		ContentID string  `json:"content_id" binding:"required"`
		Type      string  `json:"type" binding:"required"` // "rent" or "buy"
		Amount    float64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	purchase, err := h.service.CreatePurchase(c.Request.Context(), userID.(string), req.ContentID, req.Type, req.Amount)
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
	currentUserID, _ := c.Get("user_id")

	// Users can only access their own entitlements unless admin
	if userID != currentUserID.(string) {
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
			"id":          "tier1",
			"name":        "Basic",
			"price":       4.99,
			"currency":    "USD",
			"interval":    "month",
			"quality":     "480p",
			"screens":     1,
			"ads":         true,
			"download":    false,
		},
		{
			"id":          "tier2",
			"name":        "Pro",
			"price":       12.99,
			"currency":    "USD",
			"interval":    "month",
			"quality":     "720p",
			"screens":     2,
			"ads":         false,
			"download":    true,
		},
		{
			"id":          "tier3",
			"name":        "Premium",
			"price":       19.99,
			"currency":    "USD",
			"interval":    "month",
			"quality":     "4K",
			"screens":     4,
			"ads":         false,
			"download":    true,
			"priority_support": true,
		},
	}
	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

// HandleStripeWebhook handles POST /payments/webhook - Issue #16
func (h *PaymentHandler) HandleStripeWebhook(c *gin.Context) {
	// TODO: Verify Stripe signature
	// TODO: Parse webhook event
	// TODO: Update subscription status based on event type

	c.JSON(http.StatusOK, gin.H{"message": "Webhook received"})
}

