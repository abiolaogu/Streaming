# Scheduler Service

Scheduler Service manages FAST channels and live TV scheduling, EPG (Electronic Program Guide) generation, and channel manifests.

## Features

- **Channel Management**: Create and manage FAST channels and live TV channels
- **Schedule Management**: Create, update, and delete schedule entries
- **EPG Generation**: Generate 7-day electronic program guide for channels
- **Manifest Generation**: Generate streaming manifest URLs for channels
- **Current Schedule**: Get currently playing content for a channel

## Endpoints

### Channel Management
- `GET /scheduler/channels` - List all channels (filter by status)
- `GET /scheduler/channels/{channel_id}/epg` - Get EPG for channel (next 7 days)
- `GET /scheduler/channels/{channel_id}/manifest` - Get streaming manifest URL
- `GET /scheduler/channels/{channel_id}/now` - Get currently playing schedule entry

### Schedule Management
- `POST /scheduler/schedule` - Create schedule entry
- `PUT /scheduler/schedule/{id}` - Update schedule entry
- `DELETE /scheduler/schedule/{id}` - Delete schedule entry

## Channel Types

### FAST Channels
- Pre-scheduled content from catalog
- 24/7 continuous playback
- EPG shows upcoming content for 7 days
- Manifest generated dynamically based on current schedule

### Live Channels
- Live streaming from ingest URL
- Direct manifest URL from OME
- Real-time schedule updates

## Usage Examples

### Get EPG for a Channel

```bash
curl -X GET "http://localhost:8080/scheduler/channels/pluto-drama/epg"
```

Response:
```json
{
  "channelId": "pluto-drama",
  "channelName": "Pluto Drama",
  "schedule": [
    {
      "title": "Movie Title",
      "startTime": "2025-11-01T08:00:00Z",
      "duration": 7200,
      "description": "...",
      "poster": "https://...",
      "contentId": "movie_123"
    }
  ],
  "generatedAt": "2025-11-01T12:00:00Z"
}
```

### Create Schedule Entry

```bash
curl -X POST "http://localhost:8080/scheduler/schedule" \
  -H "Content-Type: application/json" \
  -d '{
    "channel_id": "pluto-drama",
    "content_id": "movie_123",
    "start_time": "2025-11-01T08:00:00Z",
    "end_time": "2025-11-01T10:00:00Z",
    "title": "Movie Title",
    "description": "A great movie",
    "poster": "https://..."
  }'
```

## Environment Variables

```bash
DATABASE_URI=mongodb://localhost:27017
DATABASE_NAME=streamverse
CDN_BASE_URL=https://cdn.streamverse.io
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
```

## Development

```bash
go mod download
go run main.go
```

## Testing

```bash
go test ./...
```

## Docker

```bash
docker build -t scheduler-service .
docker run -p 8080:8080 scheduler-service
```

