package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// Default secret key (should be set via config)
	DefaultSecretKey = "your-secret-key-change-in-production"
	// Default access token expiration
	DefaultAccessTokenExpiration = 15 * time.Minute
	// Default refresh token expiration
	DefaultRefreshTokenExpiration = 7 * 24 * time.Hour
)

// Claims represents JWT claims
type Claims struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates a new access token
func GenerateAccessToken(userID, email string, roles []string, secretKey string, expiration time.Duration) (string, error) {
	if secretKey == "" {
		secretKey = DefaultSecretKey
	}
	if expiration == 0 {
		expiration = DefaultAccessTokenExpiration
	}

	claims := Claims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// GenerateRefreshToken generates a new refresh token
func GenerateRefreshToken(userID string, secretKey string, expiration time.Duration) (string, error) {
	if secretKey == "" {
		secretKey = DefaultSecretKey
	}
	if expiration == 0 {
		expiration = DefaultRefreshTokenExpiration
	}

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// VerifyToken verifies and parses a JWT token
func VerifyToken(tokenString, secretKey string) (*Claims, error) {
	if secretKey == "" {
		secretKey = DefaultSecretKey
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

