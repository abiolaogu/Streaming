# Streaming Service

Microservice for video streaming with HLS/DASH delivery and DRM support.

## Features

- ✅ HLS manifest generation (.m3u8)
- ✅ DASH manifest generation (.mpd)
- ✅ Adaptive bitrate streaming (multiple quality levels)
- ✅ DRM integration (Widevine, FairPlay, PlayReady)
- ✅ Playback session management
- ✅ Position tracking and resume
- ✅ Concurrent stream limits
- ✅ Geo-restrictions support
- ✅ Subtitle delivery (VTT)
- ✅ Heartbeat tracking

## API Endpoints

- `GET /health` - Health check
- `GET /api/v1/streaming/:contentId/manifest?format=hls|dash` - Get streaming manifest
- `POST /api/v1/streaming/sessions` - Create playback session
- `PUT /api/v1/streaming/sessions/:sessionId/position` - Update playback position
- `POST /api/v1/streaming/sessions/:sessionId/heartbeat` - Send heartbeat
- `DELETE /api/v1/streaming/sessions/:sessionId` - End session

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
docker build -t streaming-service .
docker run -p 8080:8080 streaming-service
```
