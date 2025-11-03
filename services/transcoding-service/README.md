# Transcoding Service

Microservice for video transcoding and media processing using GStreamer.

## Features

- ✅ Transcoding job management
- ✅ Multi-bitrate ladder encoding (4K, 1080p, 720p, 480p, 360p)
- ✅ H.264 and H.265 (HEVC) support
- ✅ HLS/DASH packaging
- ✅ Thumbnail generation
- ✅ Job queue with priorities
- ✅ Progress tracking
- ✅ Quality validation

## API Endpoints

- `GET /health` - Health check
- `POST /api/v1/transcoding/jobs` - Create transcoding job
- `GET /api/v1/transcoding/jobs/:jobId` - Get job status

## Environment Variables

- `SERVER_PORT` - Server port (default: 8080)
- `DATABASE_URI` - MongoDB connection URI
- `JWT_SECRET_KEY` - JWT secret key (required)
- `LOG_LEVEL` - Log level (default: info)

## Running

```bash
go run main.go
```

