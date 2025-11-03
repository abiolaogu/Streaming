# User Service

Microservice for user profile management, preferences, and account settings.

## Features

- ✅ User profile management (display name, bio, avatar, date of birth, location)
- ✅ User preferences (language, subtitles, content ratings, notifications, playback)
- ✅ Multi-profile support (family accounts with PINs and kids mode)
- ✅ Watch history tracking and continue watching
- ✅ Favorites/watchlist
- ✅ Account deletion (GDPR compliance)
- ✅ Data export ready (GDPR compliance)

## API Endpoints

All endpoints require authentication.

- `GET /api/v1/users/me/profile` - Get user profile
- `PUT /api/v1/users/me/profile` - Update user profile
- `GET /api/v1/users/me/preferences` - Get user preferences
- `PUT /api/v1/users/me/preferences` - Update user preferences
- `GET /api/v1/users/me/profiles` - Get all sub-profiles
- `POST /api/v1/users/me/profiles` - Create sub-profile
- `GET /api/v1/users/me/watch-history` - Get watch history
- `PUT /api/v1/users/me/watch-history` - Update watch history
- `DELETE /api/v1/users/me/watch-history` - Clear watch history
- `GET /api/v1/users/me/continue-watching` - Get continue watching items
- `GET /api/v1/users/me/watchlist` - Get watchlist
- `POST /api/v1/users/me/watchlist` - Add to watchlist
- `DELETE /api/v1/users/me/watchlist/:contentId` - Remove from watchlist
- `DELETE /api/v1/users/me` - Delete user data (GDPR)

## Environment Variables

- `SERVER_PORT` - Server port (default: 8080)
- `SERVER_HOST` - Server host (default: 0.0.0.0)
- `DATABASE_URI` - MongoDB connection URI
- `DATABASE_NAME` - Database name (default: streamverse)
- `JWT_SECRET_KEY` - JWT secret key (required)
- `LOG_LEVEL` - Log level (default: info)

## Running

```bash
go run main.go
```

## Docker

```bash
docker build -t user-service .
docker run -p 8080:8080 user-service
```

