# SSAI (Server-Side Ad Insertion) Setup

## Overview

Server-side ad insertion for FAST channels and VOD content. Integrates with Google DAI (Dynamic Ad Insertion) or similar ad decision service.

## Architecture

```
1. User requests manifest for FAST channel
2. Streaming Service calls ADS API (Google DAI)
3. ADS returns ad breaks (timestamp, duration, URL)
4. Streaming Service stitches ad segments into main manifest
5. Viewer sees seamless content + ads
6. Ad impressions tracked (for billing)
```

## Components

### 1. Ad Decision Service (ADS)

**Google DAI Integration**:

```javascript
const adsRequest = {
  assetKey: 'channel_id',
  contentSourceId: 'content_source_id',
  streamActivityMonitorId: 'monitor_id',
  userInfo: {
    timezoneOffset: '-0800',
    obfuscatedGaiaId: userId
  }
};

const response = await fetch('https://googleads.g.doubleclick.net/dai/api/v2', {
  method: 'POST',
  body: JSON.stringify(adsRequest),
  headers: {
    'Content-Type': 'application/json'
  }
});
```

### 2. Manifest Rewriting

**HLS Manifest with Ad Breaks**:

```m3u8
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:6

# Ad Break Start (Pre-roll)
#EXT-X-CUE-OUT:30.0
#EXTINF:6.0,
https://ads.example.com/ad_segment_1.ts
#EXTINF:6.0,
https://ads.example.com/ad_segment_2.ts
#EXTINF:6.0,
https://ads.example.com/ad_segment_3.ts
#EXTINF:6.0,
https://ads.example.com/ad_segment_4.ts
#EXTINF:6.0,
https://ads.example.com/ad_segment_5.ts
#EXT-X-CUE-IN
#EXTINF:6.0,
https://cdn.streamverse.io/content/segment_1.ts

# Content continues...
```

### 3. Ad Break Placement

- **Pre-roll**: Before content starts
- **Mid-roll**: Every 15-30 minutes during content
- **Post-roll**: After content ends

## Implementation

### Streaming Service Integration

Add to `services/streaming-service/service/streaming_service.go`:

```go
// GetAdBreaksForChannel gets ad breaks for a FAST channel
func (s *StreamingService) GetAdBreaksForChannel(ctx context.Context, channelID string) ([]AdBreak, error) {
    // Call Google DAI API
    adsResponse, err := s.callDAI(ctx, channelID)
    if err != nil {
        return nil, err
    }
    
    // Parse ad breaks from response
    adBreaks := parseAdBreaks(adsResponse)
    return adBreaks, nil
}

// GenerateHLSManifestWithAds generates HLS manifest with ad breaks inserted
func (s *StreamingService) GenerateHLSManifestWithAds(ctx context.Context, channelID string, baseManifest string) (string, error) {
    adBreaks, err := s.GetAdBreaksForChannel(ctx, channelID)
    if err != nil {
        return baseManifest, nil // Fallback to manifest without ads
    }
    
    // Insert ad segments into manifest
    manifest := insertAdBreaks(baseManifest, adBreaks)
    return manifest, nil
}
```

### Ad Service Integration

Ad Service already exists at `services/ad-service/`. Enhance with:

- Ad decision service client (Google DAI)
- Ad break stitching logic
- Impression tracking

## Configuration

### Environment Variables

```bash
# Google DAI
GOOGLE_DAI_API_URL=https://googleads.g.doubleclick.net/dai/api/v2
GOOGLE_DAI_CONTENT_SOURCE_ID=your_content_source_id
GOOGLE_DAI_ASSET_KEY_PREFIX=streamverse_

# Ad break settings
AD_PRE_ROLL_ENABLED=true
AD_MID_ROLL_INTERVAL=900  # 15 minutes
AD_POST_ROLL_ENABLED=true
```

## Ad Break Timing

### Pre-roll
- Duration: 15-30 seconds
- Before: First content segment

### Mid-roll
- Duration: 15-30 seconds
- Interval: Every 15-30 minutes
- Placement: At natural breaks (scene changes, chapter markers)

### Post-roll
- Duration: 15-30 seconds
- After: Last content segment

## Tracking

### Ad Impressions

Track ad impressions for billing:

```go
type AdImpression struct {
    AdBreakID    string
    ChannelID    string
    UserID       string
    AdURL        string
    Duration     int
    Watched      bool
    Timestamp    time.Time
}
```

Send to Analytics Service for aggregation.

## Testing

### Test FAST Channel

1. Create FAST channel with ad breaks enabled
2. Request manifest: `GET /scheduler/channels/{channel_id}/manifest`
3. Verify manifest includes ad segments
4. Playback should show ads seamlessly

### Test VOD Content

1. Enable ads for VOD content
2. Request manifest: `GET /streaming/manifest/{content_id}/{token}.m3u8`
3. Verify ad breaks inserted at correct timestamps

## Acceptance Criteria

- ✅ Manifest includes ad breaks
- ✅ Ad segments play seamlessly
- ✅ Impressions tracked correctly
- ✅ FAST channel playback includes ads
- ✅ VOD content can have ads (optional)

## Monitoring

Monitor SSAI:
- Ad break insertion rate
- Ad playback completion rate
- Ad impression count
- Revenue from ads

