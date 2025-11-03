package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// GenerateSignedURL generates a signed CDN URL with token-based authentication
func GenerateSignedURL(baseURL, contentPath, secretKey string, ttl time.Duration) string {
	expires := time.Now().Add(ttl).Unix()
	
	// Create token
	data := fmt.Sprintf("%s:%d", contentPath, expires)
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(data))
	token := hex.EncodeToString(mac.Sum(nil))

	// Build URL
	u, _ := url.Parse(baseURL)
	u.Path = contentPath
	q := u.Query()
	q.Set("token", token)
	q.Set("expires", strconv.FormatInt(expires, 10))
	u.RawQuery = q.Encode()

	return u.String()
}

// GenerateCDNURL generates a CDN URL for content
func GenerateCDNURL(cdnBaseURL, contentID, quality, segment string) string {
	return fmt.Sprintf("%s/%s/%s/%s", cdnBaseURL, contentID, quality, segment)
}

// ValidateSignedURL validates a signed URL token
func ValidateSignedURL(token, contentPath string, expires int64, secretKey string) bool {
	// Check expiration
	if time.Now().Unix() > expires {
		return false
	}

	// Verify token
	data := fmt.Sprintf("%s:%d", contentPath, expires)
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(data))
	expectedToken := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(token), []byte(expectedToken))
}

// GeneratePlaybackToken generates a playback session token
func GeneratePlaybackToken(userID, contentID, secretKey string, ttl time.Duration) string {
	expires := time.Now().Add(ttl).Unix()
	data := fmt.Sprintf("%s:%s:%d", userID, contentID, expires)
	
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(data))
	token := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	
	return fmt.Sprintf("%s:%d", token, expires)
}

