# StreamVerse API Documentation

Complete REST API documentation for the StreamVerse platform.

## Base URL

```
Production: https://api.streamverse.io/v1
Staging: https://api-staging.streamverse.io/v1
Development: http://localhost:8080/v1
```

## Authentication

All API requests (except `/auth/login` and `/auth/register`) require an authentication token in the `Authorization` header.

```http
Authorization: Bearer YOUR_ACCESS_TOKEN
```

---

## Table of Contents

1. [Authentication API](#authentication-api)
2. [User API](#user-api)
3. [Content API](#content-api)
4. [Streaming API](#streaming-api)
5. [Search API](#search-api)
6. [Recommendation API](#recommendation-api)
7. [Payment API](#payment-api)
8. [Analytics API](#analytics-api)
9. [Notification API](#notification-api)
10. [WebSocket API](#websocket-api)
11. [Error Codes](#error-codes)

---

## Authentication API

### POST /auth/register

Register a new user account.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "firstName": "John",
  "lastName": "Doe",
  "dateOfBirth": "1990-01-01"
}
```

**Response:** `201 Created`
```json
{
  "userId": "usr_1234567890",
  "email": "user@example.com",
  "accessToken": "eyJhbGciOiJIUzI1NiIs...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
  "expiresIn": 3600
}
```

### POST /auth/login

Authenticate and receive access tokens.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```

**Response:** `200 OK`
```json
{
  "userId": "usr_1234567890",
  "accessToken": "eyJhbGciOiJIUzI1NiIs...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
  "expiresIn": 3600,
  "user": {
    "id": "usr_1234567890",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "subscription": "premium"
  }
}
```

### POST /auth/refresh

Refresh access token using refresh token.

**Request:**
```json
{
  "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response:** `200 OK`
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIs...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
  "expiresIn": 3600
}
```

### POST /auth/logout

Invalidate current session.

**Response:** `204 No Content`

### POST /auth/forgot-password

Request password reset email.

**Request:**
```json
{
  "email": "user@example.com"
}
```

**Response:** `200 OK`
```json
{
  "message": "Password reset email sent"
}
```

### POST /auth/reset-password

Reset password with token from email.

**Request:**
```json
{
  "token": "reset_token_from_email",
  "newPassword": "NewSecurePassword123!"
}
```

**Response:** `200 OK`
```json
{
  "message": "Password successfully reset"
}
```

---

## User API

### GET /users/me

Get current user profile.

**Response:** `200 OK`
```json
{
  "id": "usr_1234567890",
  "email": "user@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "dateOfBirth": "1990-01-01",
  "subscription": {
    "tier": "premium",
    "status": "active",
    "expiresAt": "2025-12-31T23:59:59Z"
  },
  "profiles": [
    {
      "id": "prof_111",
      "name": "John",
      "avatar": "https://cdn.streamverse.io/avatars/1.png",
      "isKids": false
    },
    {
      "id": "prof_222",
      "name": "Kids",
      "avatar": "https://cdn.streamverse.io/avatars/kids.png",
      "isKids": true
    }
  ],
  "preferences": {
    "language": "en",
    "autoPlayNextEpisode": true,
    "enableNotifications": true
  }
}
```

### PUT /users/me

Update current user profile.

**Request:**
```json
{
  "firstName": "John",
  "lastName": "Smith",
  "preferences": {
    "language": "en",
    "autoPlayNextEpisode": false
  }
}
```

**Response:** `200 OK`

### GET /users/me/profiles

Get all user profiles.

**Response:** `200 OK`
```json
{
  "profiles": [
    {
      "id": "prof_111",
      "name": "John",
      "avatar": "https://cdn.streamverse.io/avatars/1.png",
      "isKids": false,
      "createdAt": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### POST /users/me/profiles

Create a new user profile.

**Request:**
```json
{
  "name": "Guest",
  "avatar": "https://cdn.streamverse.io/avatars/5.png",
  "isKids": false
}
```

**Response:** `201 Created`

### GET /users/me/watchlist

Get user's watchlist.

**Response:** `200 OK`
```json
{
  "items": [
    {
      "id": "cnt_123",
      "title": "Movie Title",
      "type": "movie",
      "thumbnail": "https://cdn.streamverse.io/thumbnails/123.jpg",
      "addedAt": "2024-01-15T12:00:00Z"
    }
  ],
  "total": 25
}
```

### POST /users/me/watchlist/{contentId}

Add content to watchlist.

**Response:** `201 Created`

### DELETE /users/me/watchlist/{contentId}

Remove content from watchlist.

**Response:** `204 No Content`

### GET /users/me/watch-history

Get viewing history.

**Query Parameters:**
- `limit` (optional): Number of items (default: 20)
- `offset` (optional): Pagination offset (default: 0)

**Response:** `200 OK`
```json
{
  "items": [
    {
      "content": {
        "id": "cnt_456",
        "title": "Series Title",
        "type": "series"
      },
      "episode": {
        "id": "ep_789",
        "season": 1,
        "episode": 5,
        "title": "Episode Title"
      },
      "progress": 1245,
      "duration": 2400,
      "watchedAt": "2024-01-20T20:30:00Z"
    }
  ],
  "total": 150
}
```

---

## Content API

### GET /content/home

Get personalized home page content.

**Query Parameters:**
- `profileId` (required): Profile ID

**Response:** `200 OK`
```json
{
  "featured": {
    "id": "cnt_001",
    "title": "Featured Movie",
    "description": "An epic adventure...",
    "type": "movie",
    "thumbnail": "https://cdn.streamverse.io/featured/001.jpg",
    "poster": "https://cdn.streamverse.io/posters/001.jpg",
    "rating": 4.5,
    "year": 2024,
    "duration": 7200
  },
  "rows": [
    {
      "id": "row_trending",
      "title": "Trending Now",
      "items": [
        {
          "id": "cnt_002",
          "title": "Popular Series",
          "type": "series",
          "thumbnail": "https://cdn.streamverse.io/thumbnails/002.jpg",
          "rating": 4.8
        }
      ]
    },
    {
      "id": "row_continue",
      "title": "Continue Watching",
      "items": []
    }
  ]
}
```

### GET /content/{contentId}

Get detailed content information.

**Response:** `200 OK`
```json
{
  "id": "cnt_123",
  "title": "Amazing Movie",
  "description": "A thrilling story about...",
  "type": "movie",
  "genres": ["Action", "Adventure", "Sci-Fi"],
  "rating": 4.7,
  "year": 2024,
  "duration": 7800,
  "ageRating": "PG-13",
  "director": "John Director",
  "cast": [
    "Actor One",
    "Actor Two",
    "Actor Three"
  ],
  "thumbnail": "https://cdn.streamverse.io/thumbnails/123.jpg",
  "poster": "https://cdn.streamverse.io/posters/123.jpg",
  "trailer": {
    "url": "https://cdn.streamverse.io/trailers/123.m3u8",
    "duration": 120
  },
  "availableQualities": ["720p", "1080p", "4K"],
  "subtitles": [
    {"language": "en", "label": "English"},
    {"language": "es", "label": "Spanish"}
  ],
  "releaseDate": "2024-06-15",
  "tags": ["Blockbuster", "Award Winner"]
}
```

### GET /content/{seriesId}/seasons

Get all seasons of a series.

**Response:** `200 OK`
```json
{
  "series": {
    "id": "cnt_456",
    "title": "Great Series"
  },
  "seasons": [
    {
      "seasonNumber": 1,
      "episodeCount": 10,
      "year": 2023,
      "thumbnail": "https://cdn.streamverse.io/seasons/456-s1.jpg"
    },
    {
      "seasonNumber": 2,
      "episodeCount": 12,
      "year": 2024,
      "thumbnail": "https://cdn.streamverse.io/seasons/456-s2.jpg"
    }
  ]
}
```

### GET /content/{seriesId}/seasons/{seasonNumber}/episodes

Get all episodes in a season.

**Response:** `200 OK`
```json
{
  "season": {
    "seriesId": "cnt_456",
    "seasonNumber": 1
  },
  "episodes": [
    {
      "id": "ep_001",
      "episodeNumber": 1,
      "title": "Pilot",
      "description": "The beginning of an epic journey...",
      "duration": 2400,
      "thumbnail": "https://cdn.streamverse.io/episodes/ep001.jpg",
      "airDate": "2023-09-01"
    }
  ]
}
```

### GET /content/categories

Get all content categories.

**Response:** `200 OK`
```json
{
  "categories": [
    {
      "id": "cat_action",
      "name": "Action",
      "slug": "action",
      "thumbnail": "https://cdn.streamverse.io/categories/action.jpg"
    },
    {
      "id": "cat_comedy",
      "name": "Comedy",
      "slug": "comedy",
      "thumbnail": "https://cdn.streamverse.io/categories/comedy.jpg"
    }
  ]
}
```

### GET /content/categories/{categoryId}

Get content in a category.

**Query Parameters:**
- `limit` (optional): Number of items (default: 20)
- `offset` (optional): Pagination offset (default: 0)
- `sortBy` (optional): `popular`, `newest`, `rating` (default: `popular`)

**Response:** `200 OK`
```json
{
  "category": {
    "id": "cat_action",
    "name": "Action"
  },
  "items": [
    {
      "id": "cnt_789",
      "title": "Action Movie",
      "type": "movie",
      "thumbnail": "https://cdn.streamverse.io/thumbnails/789.jpg",
      "rating": 4.5
    }
  ],
  "total": 500,
  "hasMore": true
}
```

---

## Streaming API

### GET /streaming/{contentId}/play

Get streaming URL and DRM license for content.

**Query Parameters:**
- `quality` (optional): `720p`, `1080p`, `4K` (auto-selected if not specified)
- `profileId` (required): User profile ID

**Response:** `200 OK`
```json
{
  "contentId": "cnt_123",
  "streamUrl": "https://cdn.streamverse.io/streams/123/master.m3u8",
  "protocol": "hls",
  "qualities": [
    {
      "quality": "720p",
      "url": "https://cdn.streamverse.io/streams/123/720p.m3u8",
      "bitrate": 3000000
    },
    {
      "quality": "1080p",
      "url": "https://cdn.streamverse.io/streams/123/1080p.m3u8",
      "bitrate": 6000000
    },
    {
      "quality": "4K",
      "url": "https://cdn.streamverse.io/streams/123/4k.m3u8",
      "bitrate": 25000000
    }
  ],
  "drm": {
    "type": "widevine",
    "licenseUrl": "https://drm.streamverse.io/widevine/license",
    "certificateUrl": "https://drm.streamverse.io/widevine/cert"
  },
  "subtitles": [
    {
      "language": "en",
      "label": "English",
      "url": "https://cdn.streamverse.io/subtitles/123-en.vtt"
    }
  ],
  "sessionId": "sess_abc123",
  "expiresAt": "2024-01-20T23:59:59Z"
}
```

### POST /streaming/sessions/{sessionId}/heartbeat

Send playback heartbeat to track viewing progress.

**Request:**
```json
{
  "currentTime": 1245,
  "duration": 7200,
  "quality": "1080p",
  "bufferHealth": 15
}
```

**Response:** `200 OK`
```json
{
  "acknowledged": true
}
```

### POST /streaming/sessions/{sessionId}/complete

Mark playback session as complete.

**Request:**
```json
{
  "watchedDuration": 7200,
  "completionPercentage": 100
}
```

**Response:** `200 OK`

---

## Search API

### GET /search

Search for content.

**Query Parameters:**
- `q` (required): Search query
- `type` (optional): `movie`, `series`, `all` (default: `all`)
- `limit` (optional): Results per page (default: 20)
- `offset` (optional): Pagination offset (default: 0)

**Response:** `200 OK`
```json
{
  "query": "action adventure",
  "results": [
    {
      "id": "cnt_999",
      "title": "Action Adventure Movie",
      "type": "movie",
      "year": 2024,
      "rating": 4.6,
      "thumbnail": "https://cdn.streamverse.io/thumbnails/999.jpg",
      "relevanceScore": 0.95
    }
  ],
  "total": 45,
  "facets": {
    "types": {
      "movie": 30,
      "series": 15
    },
    "genres": {
      "Action": 25,
      "Adventure": 20
    }
  }
}
```

### GET /search/autocomplete

Get search suggestions.

**Query Parameters:**
- `q` (required): Partial search query

**Response:** `200 OK`
```json
{
  "suggestions": [
    "Action movies 2024",
    "Action adventure series",
    "Action comedy"
  ]
}
```

---

## Recommendation API

### GET /recommendations/for-you

Get personalized content recommendations.

**Query Parameters:**
- `profileId` (required): User profile ID
- `limit` (optional): Number of recommendations (default: 20)

**Response:** `200 OK`
```json
{
  "recommendations": [
    {
      "id": "cnt_555",
      "title": "Recommended Movie",
      "type": "movie",
      "thumbnail": "https://cdn.streamverse.io/thumbnails/555.jpg",
      "rating": 4.7,
      "reason": "Because you watched 'Similar Movie'"
    }
  ]
}
```

### GET /recommendations/similar/{contentId}

Get similar content recommendations.

**Query Parameters:**
- `limit` (optional): Number of recommendations (default: 10)

**Response:** `200 OK`
```json
{
  "similar": [
    {
      "id": "cnt_666",
      "title": "Similar Content",
      "type": "movie",
      "thumbnail": "https://cdn.streamverse.io/thumbnails/666.jpg",
      "rating": 4.5,
      "similarity": 0.88
    }
  ]
}
```

---

## Payment API

### GET /payments/subscriptions

Get available subscription plans.

**Response:** `200 OK`
```json
{
  "plans": [
    {
      "id": "plan_basic",
      "name": "Basic",
      "price": 9.99,
      "currency": "USD",
      "interval": "month",
      "features": [
        "SD quality",
        "1 screen",
        "Ads"
      ]
    },
    {
      "id": "plan_premium",
      "name": "Premium",
      "price": 19.99,
      "currency": "USD",
      "interval": "month",
      "features": [
        "4K quality",
        "4 screens",
        "No ads",
        "Downloads"
      ]
    }
  ]
}
```

### POST /payments/subscribe

Subscribe to a plan.

**Request:**
```json
{
  "planId": "plan_premium",
  "paymentMethodId": "pm_card_visa"
}
```

**Response:** `201 Created`
```json
{
  "subscriptionId": "sub_123456",
  "status": "active",
  "currentPeriodEnd": "2025-02-20T00:00:00Z"
}
```

### PUT /payments/subscriptions/{subscriptionId}

Update subscription (change plan, cancel, etc.).

**Request:**
```json
{
  "action": "change_plan",
  "newPlanId": "plan_basic"
}
```

**Response:** `200 OK`

### GET /payments/invoices

Get payment history.

**Response:** `200 OK`
```json
{
  "invoices": [
    {
      "id": "inv_001",
      "amount": 19.99,
      "currency": "USD",
      "status": "paid",
      "date": "2024-01-20T00:00:00Z",
      "pdfUrl": "https://api.streamverse.io/v1/payments/invoices/inv_001/pdf"
    }
  ]
}
```

---

## Analytics API

### POST /analytics/events

Track custom events.

**Request:**
```json
{
  "event": "content_view",
  "properties": {
    "contentId": "cnt_123",
    "duration": 120,
    "source": "home_page"
  }
}
```

**Response:** `202 Accepted`

---

## Notification API

### GET /notifications

Get user notifications.

**Query Parameters:**
- `limit` (optional): Number of notifications (default: 20)
- `unreadOnly` (optional): Filter unread (default: false)

**Response:** `200 OK`
```json
{
  "notifications": [
    {
      "id": "notif_001",
      "type": "new_content",
      "title": "New episodes available",
      "body": "3 new episodes of 'Your Favorite Show' are now available",
      "contentId": "cnt_789",
      "read": false,
      "createdAt": "2024-01-20T10:00:00Z"
    }
  ],
  "unreadCount": 5
}
```

### PUT /notifications/{notificationId}/read

Mark notification as read.

**Response:** `200 OK`

### POST /notifications/register-device

Register device for push notifications.

**Request:**
```json
{
  "deviceToken": "fcm_or_apns_token",
  "platform": "ios",
  "deviceId": "device_unique_id"
}
```

**Response:** `201 Created`

---

## WebSocket API

### Connection

```
wss://api.streamverse.io/v1/ws
```

**Query Parameters:**
- `token`: Authentication token

### Events

#### Client → Server

**Subscribe to content updates:**
```json
{
  "type": "subscribe",
  "channel": "content_updates"
}
```

**Join watch party:**
```json
{
  "type": "join_watch_party",
  "partyId": "party_123"
}
```

#### Server → Client

**New content notification:**
```json
{
  "type": "new_content",
  "content": {
    "id": "cnt_new",
    "title": "Brand New Series",
    "type": "series"
  }
}
```

**Watch party sync:**
```json
{
  "type": "watch_party_sync",
  "partyId": "party_123",
  "currentTime": 1245,
  "state": "playing"
}
```

---

## Error Codes

| Code | Message | Description |
|------|---------|-------------|
| 400 | Bad Request | Invalid request parameters |
| 401 | Unauthorized | Missing or invalid authentication |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Resource conflict (e.g., duplicate email) |
| 422 | Unprocessable Entity | Validation errors |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Server error |
| 503 | Service Unavailable | Service temporarily unavailable |

**Error Response Format:**
```json
{
  "error": {
    "code": "invalid_credentials",
    "message": "Invalid email or password",
    "details": {
      "field": "password"
    }
  }
}
```

---

## Rate Limiting

- **Default**: 100 requests per minute per IP
- **Authenticated**: 1000 requests per minute per user
- **Streaming**: 10 concurrent streams per account

**Rate Limit Headers:**
```http
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640000000
```

---

## Versioning

API version is specified in the URL path:
- Current: `/v1`
- Deprecated versions are supported for 6 months after new version release

---

## Support

- **API Status**: [status.streamverse.io/api](https://status.streamverse.io/api)
- **Support Email**: api-support@streamverse.io
- **Developer Portal**: [developers.streamverse.io](https://developers.streamverse.io)

---

**API Documentation Version**: 2.0
**Last Updated**: 2025
