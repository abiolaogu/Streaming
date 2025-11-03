# StreamVerse Web Player SDK

Cross-platform web video player SDK with ABR, DRM, and analytics support.

## Features

- **HLS/DASH Support**: Automatic detection and playback
- **DRM**: Widevine, FairPlay, PlayReady
- **Adaptive Bitrate (ABR)**: Automatic quality selection
- **Manual Quality Selection**: User-controlled quality switching
- **QoE Metrics**: Startup time, rebuffering, quality changes
- **Event System**: Play, pause, seek, ended, error, qualitychange
- **Captions/Subtitles**: Multi-language support (SRT, VTT)

## Installation

```bash
npm install @streamverse/player-web
```

## Usage

### Basic Usage

```javascript
import StreamVersePlayer from '@streamverse/player-web';

const player = new StreamVersePlayer({
  container: '#player',
  src: 'https://cdn.streamverse.io/hls/content_id/token.m3u8',
  logging: {
    enabled: true,
    analyticsUrl: 'https://api.streamverse.io/analytics/events',
    userId: 'user_123',
    contentId: 'content_456',
  },
});

player.play();
```

### With DRM

```javascript
const player = new StreamVersePlayer({
  container: '#player',
  src: 'https://cdn.streamverse.io/hls/content_id/token.m3u8',
  drm: {
    widevine: {
      licenseUrl: 'https://license.widevine.com',
      serverCertificate: base64Certificate,
    },
  },
  streaming: {
    bufferingGoal: 8, // seconds
    abrEnabled: true,
  },
});

player.play();
```

### Event Listeners

```javascript
player.on('play', () => {
  console.log('Playing');
});

player.on('qualitychange', (event) => {
  console.log('Quality changed:', event.quality);
});

player.on('error', (error) => {
  console.error('Error:', error);
});
```

### Quality Selection

```javascript
// Get available qualities
const qualities = player.getQualities();
console.log(qualities);

// Set quality
player.setQuality('level_2');
```

## API

### Methods

- `play()` - Start playback
- `pause()` - Pause playback
- `seek(time)` - Seek to time (seconds)
- `setQuality(qualityId)` - Set quality level
- `getQualities()` - Get available quality levels
- `on(event, callback)` - Add event listener
- `off(event, callback)` - Remove event listener
- `destroy()` - Destroy player instance

### Events

- `play` - Playback started
- `pause` - Playback paused
- `seek` - Seeking to new position
- `ended` - Playback ended
- `error` - Playback error
- `qualitychange` - Quality level changed
- `buffering` - Buffering state changed
- `ready` - Player ready

## Acceptance Criteria

- ✅ Player plays HLS/DASH streams
- ✅ ABR adapts correctly
- ✅ DRM license fetched and applied
- ✅ QoE metrics sent to Analytics Service
- ✅ Captions/subtitles work
- ✅ All player controls functional

