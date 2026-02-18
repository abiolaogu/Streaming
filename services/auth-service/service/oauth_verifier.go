package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
)

const (
	googleIssuer = "https://accounts.google.com"
	appleIssuer  = "https://appleid.apple.com"
)

// OAuthIdentity contains verified identity claims from an OAuth ID token.
type OAuthIdentity struct {
	Subject       string
	Email         string
	Name          string
	EmailVerified bool
}

// OIDCVerifier verifies Google and Apple ID tokens using OIDC discovery + JWKS.
type OIDCVerifier struct {
	googleVerifier *oidc.IDTokenVerifier
	appleVerifier  *oidc.IDTokenVerifier
}

// NewOIDCVerifier creates an OIDC verifier for configured providers.
func NewOIDCVerifier(ctx context.Context, googleClientID, appleClientID string) (*OIDCVerifier, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	verifier := &OIDCVerifier{}

	if googleClientID != "" {
		googleProvider, err := oidc.NewProvider(ctx, googleIssuer)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Google OIDC provider: %w", err)
		}
		verifier.googleVerifier = googleProvider.Verifier(&oidc.Config{ClientID: googleClientID})
	}

	if appleClientID != "" {
		appleProvider, err := oidc.NewProvider(ctx, appleIssuer)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Apple OIDC provider: %w", err)
		}
		verifier.appleVerifier = appleProvider.Verifier(&oidc.Config{ClientID: appleClientID})
	}

	return verifier, nil
}

// Verify validates a provider token and extracts identity claims.
func (v *OIDCVerifier) Verify(ctx context.Context, provider, token string) (*OAuthIdentity, error) {
	if v == nil {
		return nil, fmt.Errorf("oauth verifier not configured")
	}

	var verifier *oidc.IDTokenVerifier
	switch strings.ToLower(provider) {
	case "google":
		verifier = v.googleVerifier
	case "apple":
		verifier = v.appleVerifier
	default:
		return nil, fmt.Errorf("unsupported oauth provider: %s", provider)
	}

	if verifier == nil {
		return nil, fmt.Errorf("%s oauth is not configured", provider)
	}

	idToken, err := verifier.Verify(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("invalid %s id_token: %w", provider, err)
	}

	var claims struct {
		Sub           string      `json:"sub"`
		Email         string      `json:"email"`
		Name          string      `json:"name"`
		GivenName     string      `json:"given_name"`
		EmailVerified interface{} `json:"email_verified"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse oauth claims: %w", err)
	}

	name := strings.TrimSpace(claims.Name)
	if name == "" {
		name = strings.TrimSpace(claims.GivenName)
	}

	return &OAuthIdentity{
		Subject:       claims.Sub,
		Email:         strings.TrimSpace(claims.Email),
		Name:          name,
		EmailVerified: parseEmailVerified(claims.EmailVerified),
	}, nil
}

func parseEmailVerified(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		return strings.EqualFold(v, "true")
	default:
		return false
	}
}
