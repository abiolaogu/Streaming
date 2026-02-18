package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/streamverse/auth-service/models"
	"github.com/streamverse/common-go/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestContractIssueTokensContainExpectedClaims(t *testing.T) {
	authService := &AuthService{
		jwtSecret: "test_contract_secret_with_minimum_length_32",
	}

	user := &models.User{
		ID:    primitive.NewObjectID(),
		Email: "contract-user@streamverse.io",
		Roles: []string{"user", "subscriber"},
	}

	accessToken, refreshToken, err := authService.issueTokens(user)
	if err != nil {
		t.Fatalf("expected token issuance to succeed, got error: %v", err)
	}

	accessClaims, err := jwt.VerifyToken(accessToken, authService.jwtSecret)
	if err != nil {
		t.Fatalf("expected access token to verify, got error: %v", err)
	}
	if accessClaims.UserID != user.ID.Hex() {
		t.Fatalf("expected access token user id %q, got %q", user.ID.Hex(), accessClaims.UserID)
	}
	if accessClaims.Email != user.Email {
		t.Fatalf("expected access token email %q, got %q", user.Email, accessClaims.Email)
	}
	if len(accessClaims.Roles) != 2 {
		t.Fatalf("expected 2 access roles, got %d", len(accessClaims.Roles))
	}

	refreshClaims, err := jwt.VerifyToken(refreshToken, authService.jwtSecret)
	if err != nil {
		t.Fatalf("expected refresh token to verify, got error: %v", err)
	}
	if refreshClaims.UserID != user.ID.Hex() {
		t.Fatalf("expected refresh token user id %q, got %q", user.ID.Hex(), refreshClaims.UserID)
	}
}

func TestContractSyntheticOAuthEmailIsDeterministic(t *testing.T) {
	got := syntheticOAuthEmail(" Google ", "User@Example.Com")
	want := "google+user_example.com@oauth.streamverse.local"

	if got != want {
		t.Fatalf("expected synthetic oauth email %q, got %q", want, got)
	}
}

func TestContractParseEmailVerifiedSupportsCommonProviderFormats(t *testing.T) {
	if !parseEmailVerified(true) {
		t.Fatalf("expected bool true to parse as verified")
	}
	if !parseEmailVerified("TRUE") {
		t.Fatalf("expected string TRUE to parse as verified")
	}
	if parseEmailVerified("false") {
		t.Fatalf("expected string false to parse as unverified")
	}
}

func TestContractOAuthVerifierRejectsUnsupportedProvider(t *testing.T) {
	verifier := &OIDCVerifier{}
	_, err := verifier.Verify(context.Background(), "github", "token")
	if err == nil {
		t.Fatalf("expected unsupported provider error")
	}
	if err.Error() != "unsupported oauth provider: github" {
		t.Fatalf("expected unsupported provider mapping, got %v", err)
	}
}

func TestContractOAuthVerifierRejectsUnconfiguredProvider(t *testing.T) {
	verifier := &OIDCVerifier{}
	_, err := verifier.Verify(context.Background(), "google", "token")
	if err == nil {
		t.Fatalf("expected unconfigured provider error")
	}
	if err.Error() != "google oauth is not configured" {
		t.Fatalf("expected provider configuration mapping, got %v", err)
	}
}

func TestContractOAuthVerifierRejectsNilVerifier(t *testing.T) {
	var verifier *OIDCVerifier
	_, err := verifier.Verify(context.Background(), "google", "token")
	if err == nil {
		t.Fatalf("expected nil verifier error")
	}
	if err.Error() != "oauth verifier not configured" {
		t.Fatalf("expected nil verifier mapping, got %v", err)
	}
}

func TestContractOAuthLoginPropagatesVerifierFailuresWithoutSideEffects(t *testing.T) {
	cases := []struct {
		name      string
		provider  string
		verifyErr error
	}{
		{
			name:      "provider_not_configured",
			provider:  "google",
			verifyErr: errors.New("google oauth is not configured"),
		},
		{
			name:      "unsupported_provider",
			provider:  "github",
			verifyErr: errors.New("unsupported oauth provider: github"),
		},
		{
			name:      "invalid_token",
			provider:  "apple",
			verifyErr: errors.New("invalid apple id_token: token is malformed"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := newTestAuthUserRepo()
			deviceRepo := &testAuthDeviceRepo{}
			authService := &AuthService{
				userRepo:   userRepo,
				deviceRepo: deviceRepo,
				jwtSecret:  "test_contract_secret_with_minimum_length_32",
				oauth: &testOAuthVerifier{
					err: tc.verifyErr,
				},
			}

			accessToken, refreshToken, user, err := authService.OAuthLogin(
				context.Background(),
				tc.provider,
				"id-token",
				"Mozilla",
				"127.0.0.1",
			)

			if err == nil {
				t.Fatalf("expected verifier error for %s", tc.name)
			}
			if err.Error() != tc.verifyErr.Error() {
				t.Fatalf("expected propagated verifier error %q, got %q", tc.verifyErr.Error(), err.Error())
			}
			if accessToken != "" || refreshToken != "" || user != nil {
				t.Fatalf("expected empty oauth login response values on verifier failure")
			}
			if len(userRepo.created) != 0 || len(userRepo.updated) != 0 {
				t.Fatalf("expected no user persistence writes on verifier failure")
			}
			if len(deviceRepo.created) != 0 {
				t.Fatalf("expected no device persistence writes on verifier failure")
			}
		})
	}
}

func TestContractOAuthLoginLinksExistingEmailUser(t *testing.T) {
	existing := &models.User{
		ID:            primitive.NewObjectID(),
		Email:         "existing@streamverse.io",
		Roles:         []string{"user"},
		EmailVerified: false,
	}

	userRepo := newTestAuthUserRepo()
	userRepo.byEmail[existing.Email] = existing
	deviceRepo := &testAuthDeviceRepo{}
	verifier := &testOAuthVerifier{
		identity: &OAuthIdentity{
			Subject:       "google-sub-1",
			Email:         existing.Email,
			Name:          "Existing User",
			EmailVerified: true,
		},
	}

	authService := &AuthService{
		userRepo:   userRepo,
		deviceRepo: deviceRepo,
		jwtSecret:  "test_contract_secret_with_minimum_length_32",
		oauth:      verifier,
	}

	accessToken, refreshToken, user, err := authService.OAuthLogin(
		context.Background(),
		"google",
		"id-token",
		"Mozilla",
		"127.0.0.1",
	)
	if err != nil {
		t.Fatalf("expected oauth login success, got error: %v", err)
	}

	if accessToken == "" || refreshToken == "" {
		t.Fatalf("expected access and refresh tokens to be returned")
	}
	if user == nil {
		t.Fatalf("expected user in oauth login response")
	}
	if len(userRepo.updated) != 1 {
		t.Fatalf("expected exactly one user update for provider linking, got %d", len(userRepo.updated))
	}
	if user.OAuthProviders["google"] != "google-sub-1" {
		t.Fatalf("expected linked google subject to be persisted")
	}
	if !user.EmailVerified || user.EmailVerifiedAt == nil {
		t.Fatalf("expected verified email to be persisted from oauth identity")
	}
	if len(deviceRepo.created) != 1 {
		t.Fatalf("expected one device record created, got %d", len(deviceRepo.created))
	}
}

func TestContractOAuthLoginCreatesNewUserWithSyntheticEmail(t *testing.T) {
	userRepo := newTestAuthUserRepo()
	deviceRepo := &testAuthDeviceRepo{}
	verifier := &testOAuthVerifier{
		identity: &OAuthIdentity{
			Subject:       "apple-sub-1",
			Email:         "",
			Name:          "New Apple User",
			EmailVerified: false,
		},
	}

	authService := &AuthService{
		userRepo:   userRepo,
		deviceRepo: deviceRepo,
		jwtSecret:  "test_contract_secret_with_minimum_length_32",
		oauth:      verifier,
	}

	_, _, user, err := authService.OAuthLogin(
		context.Background(),
		"apple",
		"id-token",
		"Safari",
		"127.0.0.1",
	)
	if err != nil {
		t.Fatalf("expected oauth signup success, got error: %v", err)
	}

	if len(userRepo.created) != 1 {
		t.Fatalf("expected one created user, got %d", len(userRepo.created))
	}
	if user == nil {
		t.Fatalf("expected user in oauth signup response")
	}

	expectedEmail := "apple+apple-sub-1@oauth.streamverse.local"
	if user.Email != expectedEmail {
		t.Fatalf("expected synthetic oauth email %q, got %q", expectedEmail, user.Email)
	}
	if user.OAuthProviders["apple"] != "apple-sub-1" {
		t.Fatalf("expected oauth provider mapping for apple subject")
	}
	if len(deviceRepo.created) != 1 {
		t.Fatalf("expected one device for oauth signup, got %d", len(deviceRepo.created))
	}
}

func TestContractOAuthLoginUsesProviderLinkedUserWithoutRelinking(t *testing.T) {
	existing := &models.User{
		ID:            primitive.NewObjectID(),
		Email:         "linked@streamverse.io",
		Name:          "Linked User",
		Roles:         []string{"user"},
		EmailVerified: true,
		OAuthProviders: map[string]string{
			"google": "google-linked-sub",
		},
	}

	userRepo := newTestAuthUserRepo()
	userRepo.byProvider["google|google-linked-sub"] = existing
	deviceRepo := &testAuthDeviceRepo{}
	verifier := &testOAuthVerifier{
		identity: &OAuthIdentity{
			Subject:       "google-linked-sub",
			Email:         existing.Email,
			Name:          "Linked User",
			EmailVerified: true,
		},
	}

	authService := &AuthService{
		userRepo:   userRepo,
		deviceRepo: deviceRepo,
		jwtSecret:  "test_contract_secret_with_minimum_length_32",
		oauth:      verifier,
	}

	_, _, user, err := authService.OAuthLogin(
		context.Background(),
		"google",
		"id-token",
		"Mozilla",
		"127.0.0.1",
	)
	if err != nil {
		t.Fatalf("expected oauth login success, got error: %v", err)
	}

	if user == nil || user.ID != existing.ID {
		t.Fatalf("expected existing linked user to be reused")
	}
	if len(userRepo.updated) != 0 {
		t.Fatalf("expected no user update when provider is already linked, got %d", len(userRepo.updated))
	}
}

func TestContractOAuthLoginRelinksProviderWhenSubjectChanges(t *testing.T) {
	existing := &models.User{
		ID:            primitive.NewObjectID(),
		Email:         "relink@streamverse.io",
		Name:          "Relink User",
		Roles:         []string{"user"},
		EmailVerified: true,
		OAuthProviders: map[string]string{
			"google": "old-google-sub",
		},
	}

	userRepo := newTestAuthUserRepo()
	userRepo.byEmail[existing.Email] = existing
	userRepo.byProvider["google|old-google-sub"] = existing

	deviceRepo := &testAuthDeviceRepo{}
	verifier := &testOAuthVerifier{
		identity: &OAuthIdentity{
			Subject:       "new-google-sub",
			Email:         existing.Email,
			Name:          existing.Name,
			EmailVerified: true,
		},
	}

	authService := &AuthService{
		userRepo:   userRepo,
		deviceRepo: deviceRepo,
		jwtSecret:  "test_contract_secret_with_minimum_length_32",
		oauth:      verifier,
	}

	_, _, user, err := authService.OAuthLogin(
		context.Background(),
		"GOOGLE",
		"id-token",
		"Mozilla",
		"127.0.0.1",
	)
	if err != nil {
		t.Fatalf("expected oauth login success, got error: %v", err)
	}
	if user == nil || user.ID != existing.ID {
		t.Fatalf("expected existing email user to be reused during relink")
	}
	if len(userRepo.updated) != 1 {
		t.Fatalf("expected one user update for provider relink, got %d", len(userRepo.updated))
	}
	if got := user.OAuthProviders["google"]; got != "new-google-sub" {
		t.Fatalf("expected google provider subject to be relinked, got %q", got)
	}
}

func TestContractOAuthLoginRejectsMissingSubject(t *testing.T) {
	userRepo := newTestAuthUserRepo()
	deviceRepo := &testAuthDeviceRepo{}
	verifier := &testOAuthVerifier{
		identity: &OAuthIdentity{
			Subject:       "",
			Email:         "nosub@streamverse.io",
			Name:          "No Subject",
			EmailVerified: true,
		},
	}

	authService := &AuthService{
		userRepo:   userRepo,
		deviceRepo: deviceRepo,
		jwtSecret:  "test_contract_secret_with_minimum_length_32",
		oauth:      verifier,
	}

	_, _, user, err := authService.OAuthLogin(
		context.Background(),
		"google",
		"id-token",
		"Mozilla",
		"127.0.0.1",
	)
	if err == nil {
		t.Fatalf("expected oauth login to fail when subject is missing")
	}
	if err.Error() != "oauth subject is missing" {
		t.Fatalf("expected missing subject error, got %v", err)
	}
	if user != nil {
		t.Fatalf("expected nil user on oauth login failure")
	}
	if len(userRepo.created) != 0 || len(userRepo.updated) != 0 {
		t.Fatalf("expected no persistence side effects on oauth login failure")
	}
	if len(deviceRepo.created) != 0 {
		t.Fatalf("expected no device records created on oauth login failure")
	}
}

type testOAuthVerifier struct {
	identity *OAuthIdentity
	err      error
}

func (v *testOAuthVerifier) Verify(ctx context.Context, provider, token string) (*OAuthIdentity, error) {
	if v.err != nil {
		return nil, v.err
	}
	return v.identity, nil
}

type testAuthUserRepo struct {
	byEmail    map[string]*models.User
	byProvider map[string]*models.User
	created    []*models.User
	updated    []*models.User
}

func newTestAuthUserRepo() *testAuthUserRepo {
	return &testAuthUserRepo{
		byEmail:    make(map[string]*models.User),
		byProvider: make(map[string]*models.User),
	}
}

func (r *testAuthUserRepo) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}
	clone := *user
	r.created = append(r.created, &clone)

	r.byEmail[user.Email] = user
	for provider, subject := range user.OAuthProviders {
		r.byProvider[provider+"|"+subject] = user
	}
	return user, nil
}

func (r *testAuthUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, ok := r.byEmail[email]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *testAuthUserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	for _, user := range r.byEmail {
		if user.ID.Hex() == id {
			return user, nil
		}
	}
	for _, user := range r.byProvider {
		if user.ID.Hex() == id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *testAuthUserRepo) GetByOAuthProvider(ctx context.Context, provider, subject string) (*models.User, error) {
	user, ok := r.byProvider[provider+"|"+subject]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *testAuthUserRepo) Update(ctx context.Context, id string, user *models.User) error {
	clone := *user
	r.updated = append(r.updated, &clone)

	r.byEmail[user.Email] = user
	for provider, subject := range user.OAuthProviders {
		r.byProvider[provider+"|"+subject] = user
	}
	return nil
}

func (r *testAuthUserRepo) IncrementFailedLoginAttempts(ctx context.Context, id string) error {
	return nil
}

func (r *testAuthUserRepo) ResetFailedLoginAttempts(ctx context.Context, id string) error {
	return nil
}

func (r *testAuthUserRepo) LockAccount(ctx context.Context, id string, until time.Time) error {
	return nil
}

type testAuthDeviceRepo struct {
	created []*models.Device
}

func (r *testAuthDeviceRepo) CreateDevice(ctx context.Context, device *models.Device) error {
	clone := *device
	r.created = append(r.created, &clone)
	return nil
}

func (r *testAuthDeviceRepo) GetDevicesByUserID(ctx context.Context, userID string) ([]models.Device, error) {
	return nil, nil
}

func (r *testAuthDeviceRepo) DeleteDevice(ctx context.Context, deviceID, userID string) error {
	return nil
}
