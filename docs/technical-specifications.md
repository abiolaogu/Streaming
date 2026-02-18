# API Reference — StreamVerse Streaming Platform

## 1. Base URLs

| Environment | URL |
|-------------|-----|
| Production | `https://api.streamverse.io/v1` |
| Staging | `https://api-staging.streamverse.io/v1` |
| Development | `http://localhost:8080/v1` |

## 2. Authentication

All authenticated endpoints require a Bearer token:
```
Authorization: Bearer {accessToken}
```

Rate limits:
- Anonymous: 100 requests/minute per IP
- Authenticated: 1000 requests/minute per user
- Streaming: 10 concurrent streams per account

---

## 3. Authentication API

### POST /auth/register
Create a new user account.

**Request**: `{ email, password, firstName, lastName, dateOfBirth }`
**Response 201**: `{ userId, email, accessToken, refreshToken, expiresIn }`

### POST /auth/login
Authenticate with email and password.

**Request**: `{ email, password }`
**Response 200**: `{ userId, accessToken, refreshToken, expiresIn, user: { id, email, firstName, lastName, subscription } }`

### POST /auth/refresh
Refresh an expired access token.

**Request**: `{ refreshToken }`
**Response 200**: `{ accessToken, refreshToken, expiresIn }`

### POST /auth/logout
Invalidate current session. **Response 204**.

### POST /auth/forgot-password
Request password reset email.

### POST /auth/reset-password
Reset password with email token.

---

## 4. User API

### GET /users/me
Get current user profile including subscription details and profiles array.

### PUT /users/me
Update user profile fields and preferences.

### GET /users/me/profiles
List all profiles for the current user.

### POST /users/me/profiles
Create a new profile. Maximum 5 per account.

### GET /users/me/watchlist
Get user's saved content watchlist with pagination.

### POST /users/me/watchlist/{contentId}
Add content to watchlist. **Response 201**.

### DELETE /users/me/watchlist/{contentId}
Remove content from watchlist. **Response 204**.

### GET /users/me/watch-history
Get viewing history with pagination. Returns content metadata, episode info, progress position, and duration.

**Query Parameters**: `limit` (default 20), `offset` (default 0)

---

## 5. Content API

### GET /content/home
Personalized home page content for a profile.

**Query Parameters**: `profileId` (required)
**Response**: Featured content, content rows (trending, continue watching, recommendations, genre-based rows).

### GET /content/{contentId}
Detailed content information including metadata, cast, trailer, available qualities, subtitles.

### GET /content/{seriesId}/seasons
All seasons for a series with episode counts.

### GET /content/{seriesId}/seasons/{seasonNumber}/episodes
All episodes in a season with descriptions, durations, air dates.

### GET /content/categories
All content categories/genres.

### GET /content/categories/{categoryId}
Browse content within a category with sorting and pagination.

**Query Parameters**: `limit`, `offset`, `sortBy` (popular|newest|rating)

---

## 6. Streaming API

### GET /streaming/{contentId}/play
Get streaming URLs, DRM configuration, and subtitles for playback.

**Query Parameters**: `quality` (720p|1080p|4K|auto), `profileId` (required)
**Response**:
```json
{
  "streamUrl": "master.m3u8",
  "protocol": "hls",
  "qualities": [{ "quality": "1080p", "url": "...", "bitrate": 6000000 }],
  "drm": { "type": "widevine", "licenseUrl": "...", "certificateUrl": "..." },
  "subtitles": [{ "language": "en", "url": "...en.vtt" }],
  "sessionId": "sess_abc123",
  "expiresAt": "2026-01-20T23:59:59Z"
}
```

### POST /streaming/sessions/{sessionId}/heartbeat
Send playback progress heartbeat every 30 seconds.

**Request**: `{ currentTime, duration, quality, bufferHealth }`

### POST /streaming/sessions/{sessionId}/complete
Mark playback session as complete.

**Request**: `{ watchedDuration, completionPercentage }`

---

## 7. Search API

### GET /search
Full-text content search with faceted results.

**Query Parameters**: `q` (required), `type` (movie|series|all), `limit`, `offset`
**Response**: Results with relevance scores, facets by type and genre.

### GET /search/autocomplete
Search suggestions as user types.

**Query Parameters**: `q` (required)
**Response**: `{ suggestions: ["Action movies 2024", ...] }`

---

## 8. Recommendation API

### GET /recommendations/for-you
Personalized recommendations using hybrid ML strategy.

**Query Parameters**: `profileId` (required), `limit` (default 20)
**Response**: Recommendations with reason strings (e.g., "Because you watched...").

### GET /recommendations/similar/{contentId}
Similar content based on genre, cast, and viewer overlap.

---

## 9. Payment API

### GET /payments/subscriptions
Available subscription plans with pricing and feature lists.

### POST /payments/subscribe
Subscribe to a plan. **Request**: `{ planId, paymentMethodId }`

### PUT /payments/subscriptions/{subscriptionId}
Change plan or cancel subscription.

### GET /payments/invoices
Payment history with invoice PDF links.

---

## 10. Analytics API

### POST /analytics/events
Track custom analytics events. **Response 202 Accepted** (async processing).

---

## 11. Notification API

### GET /notifications
Get user notifications with optional unread filter.

### PUT /notifications/{notificationId}/read
Mark notification as read.

### POST /notifications/register-device
Register device for push notifications (FCM/APNs).

---

## 12. WebSocket API

**Connection**: `wss://api.streamverse.io/v1/ws?token={jwt}`

**Client Events**:
- `subscribe`: Subscribe to content_updates channel
- `join_watch_party`: Join synchronized viewing session

**Server Events**:
- `new_content`: New content published notification
- `watch_party_sync`: Playback synchronization (currentTime, state)

---

## 13. gRPC Services

Protobuf definitions in `packages/proto/`:
- `auth/v1/auth.proto`: Register, Login, ValidateToken
- `content/v1/content.proto`: GetContent, ListContent, CreateContent
- `streaming/streaming.proto`: GetPlaybackURL, CreateSession
- `payment/payment.proto`: CheckEntitlement, ProcessPayment

---

## 14. Error Response Format

```json
{
  "error": {
    "code": "invalid_credentials",
    "message": "Invalid email or password",
    "details": { "field": "password" }
  }
}
```

| HTTP Code | Description |
|-----------|-------------|
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 409 | Conflict |
| 422 | Validation Error |
| 429 | Rate Limit Exceeded |
| 500 | Internal Server Error |
| 503 | Service Unavailable |

---

## 15. Rate Limit Headers

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640000000
```

---

**Document Version**: 2.0
**Last Updated**: 2026-02-17
