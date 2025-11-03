package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/streamverse/common-go/database"
	"github.com/streamverse/auth-service/models"
)

// TokenRepository handles token operations
type TokenRepository struct {
	passwordResetCollection *mongo.Collection
	emailVerifCollection    *mongo.Collection
	oauthStateCollection    *mongo.Collection
}

// NewTokenRepository creates a new token repository
func NewTokenRepository(db *database.MongoDB) *TokenRepository {
	return &TokenRepository{
		passwordResetCollection: db.Collection("password_reset_tokens"),
		emailVerifCollection:    db.Collection("email_verification_tokens"),
		oauthStateCollection:    db.Collection("oauth_states"),
	}
}

// CreatePasswordResetToken creates a password reset token
func (r *TokenRepository) CreatePasswordResetToken(ctx context.Context, token *models.PasswordResetToken) error {
	token.ID = primitive.NewObjectID()
	token.CreatedAt = time.Now()
	token.Used = false

	_, err := r.passwordResetCollection.InsertOne(ctx, token)
	return err
}

// GetPasswordResetToken retrieves a password reset token
func (r *TokenRepository) GetPasswordResetToken(ctx context.Context, token string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	err := r.passwordResetCollection.FindOne(ctx, bson.M{"token": token, "used": false}).Decode(&resetToken)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("invalid or expired token")
		}
		return nil, err
	}

	if time.Now().After(resetToken.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}

	return &resetToken, nil
}

// MarkPasswordResetTokenUsed marks a password reset token as used
func (r *TokenRepository) MarkPasswordResetTokenUsed(ctx context.Context, token string) error {
	_, err := r.passwordResetCollection.UpdateOne(
		ctx,
		bson.M{"token": token},
		bson.M{"$set": bson.M{"used": true}},
	)
	return err
}

// CreateEmailVerificationToken creates an email verification token
func (r *TokenRepository) CreateEmailVerificationToken(ctx context.Context, token *models.EmailVerificationToken) error {
	token.ID = primitive.NewObjectID()
	token.CreatedAt = time.Now()
	token.Used = false

	_, err := r.emailVerifCollection.InsertOne(ctx, token)
	return err
}

// GetEmailVerificationToken retrieves an email verification token
func (r *TokenRepository) GetEmailVerificationToken(ctx context.Context, token string) (*models.EmailVerificationToken, error) {
	var verifyToken models.EmailVerificationToken
	err := r.emailVerifCollection.FindOne(ctx, bson.M{"token": token, "used": false}).Decode(&verifyToken)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("invalid or expired token")
		}
		return nil, err
	}

	if time.Now().After(verifyToken.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}

	return &verifyToken, nil
}

// MarkEmailVerificationTokenUsed marks an email verification token as used
func (r *TokenRepository) MarkEmailVerificationTokenUsed(ctx context.Context, token string) error {
	_, err := r.emailVerifCollection.UpdateOne(
		ctx,
		bson.M{"token": token},
		bson.M{"$set": bson.M{"used": true}},
	)
	return err
}

// CreateOAuthState creates an OAuth state for CSRF protection
func (r *TokenRepository) CreateOAuthState(ctx context.Context, state *models.OAuthState) error {
	_, err := r.oauthStateCollection.InsertOne(ctx, state)
	return err
}

// GetOAuthState retrieves an OAuth state
func (r *TokenRepository) GetOAuthState(ctx context.Context, state string) (*models.OAuthState, error) {
	var oauthState models.OAuthState
	err := r.oauthStateCollection.FindOne(ctx, bson.M{"state": state}).Decode(&oauthState)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("invalid state")
		}
		return nil, err
	}

	if time.Now().After(oauthState.ExpiresAt) {
		return nil, fmt.Errorf("state expired")
	}

	return &oauthState, nil
}

