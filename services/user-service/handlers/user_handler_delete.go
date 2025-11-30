package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streamverse/common-go/errors"
	"go.uber.org/zap"
)

// DeleteUser handles DELETE /api/v1/users/me
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("User not authenticated"))
		return
	}

	if err := h.service.DeleteUser(c.Request.Context(), userID.(string)); err != nil {
		h.logger.Error("Failed to delete user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to delete user"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
