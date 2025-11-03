package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user account
type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email             string             `bson:"email" json:"email"`
	PasswordHash      string             `bson:"password_hash" json:"-"`
	Name              string             `bson:"name,omitempty" json:"name,omitempty"`
	EmailVerified     bool               `bson:"email_verified" json:"emailVerified"`
	EmailVerifiedAt   *time.Time         `bson:"email_verified_at,omitempty" json:"emailVerifiedAt,omitempty"`
	MFAEnabled        bool               `bson:"mfa_enabled" json:"mfaEnabled"`
	MFASecret         string             `bson:"mfa_secret,omitempty" json:"-"`
	MFABackupCodes    []string           `bson:"mfa_backup_codes,omitempty" json:"-"`
	FailedLoginAttempts int              `bson:"failed_login_attempts" json:"-"`
	AccountLockedUntil *time.Time        `bson:"account_locked_until,omitempty" json:"-"`
	Roles             []string           `bson:"roles" json:"roles"`
	OAuthProviders    map[string]string  `bson:"oauth_providers,omitempty" json:"-"`
	CreatedAt         time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updatedAt"`
}

// Device represents a device session
type Device struct {
	ID          string    `bson:"id" json:"id"`
	UserID      string    `bson:"user_id" json:"userId"`
	Name        string    `bson:"name" json:"name"`
	Type        string    `bson:"type" json:"type"`
	Fingerprint string    `bson:"fingerprint" json:"fingerprint"`
	LastUsedAt  time.Time `bson:"last_used_at" json:"lastUsedAt"`
	CreatedAt   time.Time `bson:"created_at" json:"createdAt"`
}

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"userId"`
	Token     string             `bson:"token" json:"-"`
	ExpiresAt time.Time          `bson:"expires_at" json:"expiresAt"`
	Used      bool               `bson:"used" json:"used"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
}

// EmailVerificationToken represents an email verification token
type EmailVerificationToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"userId"`
	Token     string             `bson:"token" json:"-"`
	ExpiresAt time.Time          `bson:"expires_at" json:"expiresAt"`
	Used      bool               `bson:"used" json:"used"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
}

// OAuthState represents OAuth state for CSRF protection
type OAuthState struct {
	ID        string    `bson:"id" json:"id"`
	Provider  string    `bson:"provider" json:"provider"`
	State     string    `bson:"state" json:"state"`
	ExpiresAt time.Time `bson:"expires_at" json:"expiresAt"`
}

