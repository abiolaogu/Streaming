package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/streamverse/common-go/jwt"
	"github.com/streamverse/auth-service/models"
	"github.com/streamverse/auth-service/repository"
	"github.com/streamverse/auth-service/utils"
)

const (
	maxFailedAttempts = 5
	lockoutDuration   = 30 * time.Minute
	tokenExpiration   = 15 * time.Minute
	refreshExpiration = 7 * 24 * time.Hour
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo   *repository.UserRepository
	tokenRepo  *repository.TokenRepository
	deviceRepo *repository.DeviceRepository
	jwtSecret  string
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo *repository.UserRepository,
	tokenRepo *repository.TokenRepository,
	deviceRepo *repository.DeviceRepository,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		deviceRepo: deviceRepo,
		jwtSecret:  jwtSecret,
	}
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, email, password, name string) (*models.User, error) {
	// Validate email
	if !utils.ValidateEmail(email) {
		return nil, fmt.Errorf("invalid email format")
	}

	// Validate password strength
	if err := utils.ValidatePasswordStrength(password); err != nil {
		return nil, err
	}

	// Check if user exists
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Email:        email,
		PasswordHash: passwordHash,
		Name:         name,
	}

	user, err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, email, password, userAgent, ipAddress string) (string, string, *models.User, error) {
	// Get user
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", nil, fmt.Errorf("invalid credentials")
	}

	// Check if account is locked
	if user.AccountLockedUntil != nil && time.Now().Before(*user.AccountLockedUntil) {
		return "", "", nil, fmt.Errorf("account is locked")
	}

	// Verify password
	if !utils.CheckPassword(password, user.PasswordHash) {
		s.userRepo.IncrementFailedLoginAttempts(ctx, user.ID.Hex())
		
		if user.FailedLoginAttempts+1 >= maxFailedAttempts {
			lockUntil := time.Now().Add(lockoutDuration)
			s.userRepo.LockAccount(ctx, user.ID.Hex(), lockUntil)
		}

		return "", "", nil, fmt.Errorf("invalid credentials")
	}

	// Reset failed attempts on successful login
	if user.FailedLoginAttempts > 0 {
		s.userRepo.ResetFailedLoginAttempts(ctx, user.ID.Hex())
	}

	// Generate device fingerprint
	deviceFingerprint := utils.GenerateDeviceFingerprint(userAgent, ipAddress)

	// Create or update device
	device := &models.Device{
		ID:          deviceFingerprint,
		UserID:      user.ID.Hex(),
		Name:        userAgent,
		Type:        "web",
		Fingerprint: deviceFingerprint,
		LastUsedAt:  time.Now(),
		CreatedAt:   time.Now(),
	}
	s.deviceRepo.CreateDevice(ctx, device)

	// Generate tokens
	accessToken, err := jwt.GenerateAccessToken(
		user.ID.Hex(),
		user.Email,
		user.Roles,
		s.jwtSecret,
		tokenExpiration,
	)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := jwt.GenerateRefreshToken(
		user.ID.Hex(),
		s.jwtSecret,
		refreshExpiration,
	)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, user, nil
}

// RefreshToken refreshes an access token
func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenString string) (string, string, error) {
	claims, err := jwt.VerifyToken(refreshTokenString, s.jwtSecret)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return "", "", fmt.Errorf("user not found")
	}

	// Generate new tokens
	accessToken, err := jwt.GenerateAccessToken(
		user.ID.Hex(),
		user.Email,
		user.Roles,
		s.jwtSecret,
		tokenExpiration,
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := jwt.GenerateRefreshToken(
		user.ID.Hex(),
		s.jwtSecret,
		refreshExpiration,
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, newRefreshToken, nil
}

// RequestPasswordReset requests a password reset
func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Don't reveal if user exists
		return "", nil
	}

	// Generate reset token
	token := generateRandomToken()
	resetToken := &models.PasswordResetToken{
		UserID:    user.ID.Hex(),
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := s.tokenRepo.CreatePasswordResetToken(ctx, resetToken); err != nil {
		return "", err
	}

	return token, nil
}

// ResetPassword resets a user's password
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Validate password strength
	if err := utils.ValidatePasswordStrength(newPassword); err != nil {
		return err
	}

	// Get reset token
	resetToken, err := s.tokenRepo.GetPasswordResetToken(ctx, token)
	if err != nil {
		return err
	}

	// Hash new password
	passwordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password
	user, err := s.userRepo.GetByID(ctx, resetToken.UserID)
	if err != nil {
		return err
	}

	user.PasswordHash = passwordHash
	if err := s.userRepo.Update(ctx, resetToken.UserID, user); err != nil {
		return err
	}

	// Mark token as used
	s.tokenRepo.MarkPasswordResetTokenUsed(ctx, token)

	return nil
}

// VerifyEmail verifies a user's email
func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	verifyToken, err := s.tokenRepo.GetEmailVerificationToken(ctx, token)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetByID(ctx, verifyToken.UserID)
	if err != nil {
		return err
	}

	now := time.Now()
	user.EmailVerified = true
	user.EmailVerifiedAt = &now
	if err := s.userRepo.Update(ctx, verifyToken.UserID, user); err != nil {
		return err
	}

	s.tokenRepo.MarkEmailVerificationTokenUsed(ctx, token)
	return nil
}

// EnableMFA enables multi-factor authentication
func (s *AuthService) EnableMFA(ctx context.Context, userID string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "StreamVerse",
		AccountName: userID,
	})
	if err != nil {
		return "", err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	user.MFAEnabled = false // Set to true after verification
	user.MFASecret = key.Secret()
	if err := s.userRepo.Update(ctx, userID, user); err != nil {
		return "", err
	}

	return key.URL(), nil
}

// VerifyMFA verifies MFA code
func (s *AuthService) VerifyMFA(ctx context.Context, userID, code string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if !totp.Validate(code, user.MFASecret) {
		return fmt.Errorf("invalid MFA code")
	}

	user.MFAEnabled = true
	return s.userRepo.Update(ctx, userID, user)
}

// GetDevices retrieves user devices
func (s *AuthService) GetDevices(ctx context.Context, userID string) ([]models.Device, error) {
	return s.deviceRepo.GetDevicesByUserID(ctx, userID)
}

// RevokeDevice revokes a device
func (s *AuthService) RevokeDevice(ctx context.Context, userID, deviceID string) error {
	return s.deviceRepo.DeleteDevice(ctx, deviceID, userID)
}

func generateRandomToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

