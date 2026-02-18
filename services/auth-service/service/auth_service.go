package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/streamverse/auth-service/models"
	"github.com/streamverse/auth-service/utils"
	"github.com/streamverse/common-go/jwt"
)

const (
	maxFailedAttempts = 5
	lockoutDuration   = 30 * time.Minute
	tokenExpiration   = 15 * time.Minute
	refreshExpiration = 7 * 24 * time.Hour
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo   authUserRepository
	tokenRepo  authTokenRepository
	deviceRepo authDeviceRepository
	jwtSecret  string
	oauth      oauthVerifier
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo authUserRepository,
	tokenRepo authTokenRepository,
	deviceRepo authDeviceRepository,
	jwtSecret string,
	oauth oauthVerifier,
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		deviceRepo: deviceRepo,
		jwtSecret:  jwtSecret,
		oauth:      oauth,
	}
}

type authUserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByOAuthProvider(ctx context.Context, provider, subject string) (*models.User, error)
	Update(ctx context.Context, id string, user *models.User) error
	IncrementFailedLoginAttempts(ctx context.Context, id string) error
	ResetFailedLoginAttempts(ctx context.Context, id string) error
	LockAccount(ctx context.Context, id string, until time.Time) error
}

type authTokenRepository interface {
	CreatePasswordResetToken(ctx context.Context, token *models.PasswordResetToken) error
	GetPasswordResetToken(ctx context.Context, token string) (*models.PasswordResetToken, error)
	MarkPasswordResetTokenUsed(ctx context.Context, token string) error
	GetEmailVerificationToken(ctx context.Context, token string) (*models.EmailVerificationToken, error)
	MarkEmailVerificationTokenUsed(ctx context.Context, token string) error
}

type authDeviceRepository interface {
	CreateDevice(ctx context.Context, device *models.Device) error
	GetDevicesByUserID(ctx context.Context, userID string) ([]models.Device, error)
	DeleteDevice(ctx context.Context, deviceID, userID string) error
}

type oauthVerifier interface {
	Verify(ctx context.Context, provider, token string) (*OAuthIdentity, error)
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
	accessToken, refreshToken, err := s.issueTokens(user)
	if err != nil {
		return "", "", nil, err
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
	accessToken, newRefreshToken, err := s.issueTokens(user)
	if err != nil {
		return "", "", err
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

// OAuthLogin verifies OAuth id_token and performs login/signup for provider users.
func (s *AuthService) OAuthLogin(ctx context.Context, provider, idToken, userAgent, ipAddress string) (string, string, *models.User, error) {
	identity, err := s.oauth.Verify(ctx, provider, idToken)
	if err != nil {
		return "", "", nil, err
	}

	if identity.Subject == "" {
		return "", "", nil, fmt.Errorf("oauth subject is missing")
	}

	user, err := s.userRepo.GetByOAuthProvider(ctx, strings.ToLower(provider), identity.Subject)
	if err != nil {
		if identity.Email != "" {
			existingUser, emailErr := s.userRepo.GetByEmail(ctx, identity.Email)
			if emailErr == nil {
				user = existingUser
			}
		}
	}

	if user == nil {
		email := identity.Email
		if email == "" {
			email = syntheticOAuthEmail(provider, identity.Subject)
		}

		user = &models.User{
			Email:         email,
			PasswordHash:  "",
			Name:          identity.Name,
			EmailVerified: identity.EmailVerified,
			Roles:         []string{"user"},
			OAuthProviders: map[string]string{
				strings.ToLower(provider): identity.Subject,
			},
		}

		user, err = s.userRepo.Create(ctx, user)
		if err != nil {
			return "", "", nil, err
		}
	} else {
		updated := false
		if user.OAuthProviders == nil {
			user.OAuthProviders = make(map[string]string)
		}
		if user.OAuthProviders[strings.ToLower(provider)] != identity.Subject {
			user.OAuthProviders[strings.ToLower(provider)] = identity.Subject
			updated = true
		}
		if user.Name == "" && identity.Name != "" {
			user.Name = identity.Name
			updated = true
		}
		if user.Email == "" && identity.Email != "" {
			user.Email = identity.Email
			updated = true
		}
		if identity.EmailVerified && !user.EmailVerified {
			now := time.Now()
			user.EmailVerified = true
			user.EmailVerifiedAt = &now
			updated = true
		}

		if updated {
			if err := s.userRepo.Update(ctx, user.ID.Hex(), user); err != nil {
				return "", "", nil, err
			}
		}
	}

	deviceFingerprint := utils.GenerateDeviceFingerprint(userAgent, ipAddress)
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

	accessToken, refreshToken, err := s.issueTokens(user)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

func (s *AuthService) issueTokens(user *models.User) (string, string, error) {
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

	refreshToken, err := jwt.GenerateRefreshToken(
		user.ID.Hex(),
		s.jwtSecret,
		refreshExpiration,
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func syntheticOAuthEmail(provider, subject string) string {
	cleanProvider := strings.ToLower(strings.TrimSpace(provider))
	cleanSubject := strings.ToLower(strings.TrimSpace(subject))
	cleanSubject = strings.ReplaceAll(cleanSubject, "@", "_")
	return fmt.Sprintf("%s+%s@oauth.streamverse.local", cleanProvider, cleanSubject)
}

func generateRandomToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
