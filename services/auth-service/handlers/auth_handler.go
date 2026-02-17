package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streamverse/common-go/errors"
	"github.com/streamverse/common-go/logger"
	"github.com/streamverse/auth-service/service"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	service *service.AuthService
	logger  *logger.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(service *service.AuthService, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logger,
	}
}

// RegisterRequest represents registration request
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// PasswordResetRequest represents password reset request
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// PasswordResetConfirmRequest represents password reset confirmation
type PasswordResetConfirmRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// VerifyEmailRequest represents email verification request
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	user, err := h.service.Register(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		h.logger.Error("Registration failed", logger.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	accessToken, refreshToken, user, err := h.service.Login(
		c.Request.Context(),
		req.Email,
		req.Password,
		userAgent,
		ipAddress,
	)
	if err != nil {
		h.logger.Error("Login failed", logger.Error(err))
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError(err.Error()))
		return
	}

	// Set refresh token in httpOnly cookie (Issue #11 requirement)
	c.SetCookie(
		"refresh_token",
		refreshToken,
		7*24*3600, // 7 days
		"/",
		"",
		true,  // secure (HTTPS only)
		true,  // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"user":         user,
		"expires_in":   900, // 15 minutes
	})
}

// RefreshToken handles POST /auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Try to get refresh token from cookie first, then from body
	var refreshToken string
	cookie, err := c.Cookie("refresh_token")
	if err == nil && cookie != "" {
		refreshToken = cookie
	} else {
		var req RefreshTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errors.NewInvalidInputError("Refresh token required"))
			return
		}
		refreshToken = req.RefreshToken
	}

	accessToken, newRefreshToken, err := h.service.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		h.logger.Error("Token refresh failed", logger.Error(err))
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError(err.Error()))
		return
	}

	// Set new refresh token in httpOnly cookie
	c.SetCookie(
		"refresh_token",
		newRefreshToken,
		7*24*3600,
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"expires_in":   900,
	})
}

// Logout handles POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Token invalidation handled via Redis/session management
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// RequestPasswordReset handles POST /api/v1/auth/password/reset
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	_, err := h.service.RequestPasswordReset(c.Request.Context(), req.Email)
	if err != nil {
		h.logger.Error("Password reset request failed", logger.Error(err))
	}

	// Always return success to prevent email enumeration
	c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a reset link has been sent"})
}

// ResetPassword handles POST /api/v1/auth/password/reset/confirm
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req PasswordResetConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		h.logger.Error("Password reset failed", logger.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// VerifyEmail handles POST /api/v1/auth/verify-email
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	if err := h.service.VerifyEmail(c.Request.Context(), req.Token); err != nil {
		h.logger.Error("Email verification failed", logger.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// GetDevices handles GET /api/v1/auth/devices
func (h *AuthHandler) GetDevices(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("User not authenticated"))
		return
	}

	devices, err := h.service.GetDevices(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error("Failed to get devices", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError("Failed to get devices"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"devices": devices})
}

// RevokeDevice handles DELETE /api/v1/auth/devices/:id
func (h *AuthHandler) RevokeDevice(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("User not authenticated"))
		return
	}

	deviceID := c.Param("id")
	if err := h.service.RevokeDevice(c.Request.Context(), userID.(string), deviceID); err != nil {
		h.logger.Error("Failed to revoke device", logger.Error(err))
		c.JSON(http.StatusNotFound, errors.NewNotFoundError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device revoked successfully"})
}

// Validate handles GET /auth/validate (internal use)
func (h *AuthHandler) Validate(c *gin.Context) {
	// Token is already validated by AuthMiddleware
	userID, _ := c.Get("user_id")
	email, _ := c.Get("email")
	roles, _ := c.Get("roles")
	orgID, _ := c.Get("org_id")

	c.JSON(http.StatusOK, gin.H{
		"valid":  true,
		"user_id": userID,
		"email":   email,
		"roles":   roles,
		"org_id":  orgID,
	})
}

// SetupMFA handles POST /auth/mfa/setup
func (h *AuthHandler) SetupMFA(c *gin.Context) {
	userID, _ := c.Get("user_id")

	qrCodeURL, err := h.service.EnableMFA(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error("Failed to setup MFA", logger.Error(err))
		c.JSON(http.StatusInternalServerError, errors.NewInternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"qr_code_url": qrCodeURL,
		"message":     "Scan QR code with authenticator app",
	})
}

// VerifyMFA handles POST /auth/mfa/verify
func (h *AuthHandler) VerifyMFA(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	if err := h.service.VerifyMFA(c.Request.Context(), userID.(string), req.Code); err != nil {
		h.logger.Error("Failed to verify MFA", logger.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError("Invalid MFA code"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "MFA enabled successfully"})
}

// OAuthGoogle handles POST /auth/oauth/google
func (h *AuthHandler) OAuthGoogle(c *gin.Context) {
	var req struct {
		IDToken string `json:"id_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	// TODO: Verify Google ID token and create/login user
	// accessToken, refreshToken, user, err := h.service.OAuthLogin(ctx, "google", req.IDToken)
	c.JSON(http.StatusNotImplemented, gin.H{"message": "OAuth2.0 Google integration - TODO"})
}

// OAuthApple handles POST /auth/oauth/apple
func (h *AuthHandler) OAuthApple(c *gin.Context) {
	var req struct {
		IDToken string `json:"id_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewInvalidInputError(err.Error()))
		return
	}

	// TODO: Verify Apple ID token and create/login user
	c.JSON(http.StatusNotImplemented, gin.H{"message": "OAuth2.0 Apple integration - TODO"})
}
