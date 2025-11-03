# Ad Service

Microservice for AVOD (Ad-supported VOD) monetization.

## Features

- ✅ Pre-roll, mid-roll, post-roll ads
- ✅ Ad targeting (demographics, content, geography)
- ✅ VAST/VMAP support
- ✅ Ad tracking (impressions, clicks, completion)
- ✅ SSAI support
- ✅ Ad-free tier checks

## API Endpoints

- `POST /api/v1/ads/request` - Get ads for content
- `POST /api/v1/ads/track` - Track ad events

## Running

```bash
go run main.go
```

