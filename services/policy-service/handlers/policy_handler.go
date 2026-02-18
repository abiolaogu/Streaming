package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streamverse/policy-service/service"
)

// PolicyHandler handles policy evaluation requests.
type PolicyHandler struct {
	service *service.PolicyService
}

// NewPolicyHandler creates a policy handler.
func NewPolicyHandler(policyService *service.PolicyService) *PolicyHandler {
	return &PolicyHandler{service: policyService}
}

// EvaluateEntitlement handles POST /policy/v1/entitlements/evaluate.
func (h *PolicyHandler) EvaluateEntitlement(c *gin.Context) {
	var req service.EntitlementEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.service.EvaluateEntitlement(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
