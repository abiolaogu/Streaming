# DRM License Server Integration

## Overview

DRM license server configuration and integration guide for Widevine, FairPlay, and PlayReady.

## DRM Providers

### Widevine (Google)
- **License Server**: `https://license.widevine.com`
- **Levels**:
  - Level 1: 4K content (hardware security)
  - Level 2: 1080p content
  - Level 3: SD content (software security, free tier)

### FairPlay (Apple)
- **License Server**: `https://license.fairplay.com`
- **Certificate**: `priv_key.pem` (stored in Vault)
- **Platforms**: iOS, macOS, tvOS

### PlayReady (Microsoft)
- **License Server**: `https://license.playready.com`
- **Platforms**: Windows, Smart TVs (Tizen, webOS)

## License Token Generation

License tokens are JWT tokens signed by the Streaming Service:

```json
{
  "content_id": "content_123",
  "user_id": "user_456",
  "drm_level": "1",
  "exp": 3600
}
```

## Configuration

### Environment Variables

```bash
# Widevine
WIDEVINE_LICENSE_URL=https://license.widevine.com
WIDEVINE_CONTENT_ID_KEY=...

# FairPlay
FAIRPLAY_LICENSE_URL=https://license.fairplay.com
FAIRPLAY_CERT_PATH=/etc/certs/fairplay_cert.pem
FAIRPLAY_PRIVATE_KEY_PATH=/etc/certs/fairplay_key.pem

# PlayReady
PLAYREADY_LICENSE_URL=https://license.playready.com
PLAYREADY_CUSTOM_DATA=...
```

## Client-Side Integration

### Web (Widevine)

```javascript
const video = document.getElementById('video');
const drmConfig = {
  widevine: {
    licenseUrl: 'https://license.widevine.com',
    serverCertificate: base64Certificate
  }
};

// Shaka Player
const player = new shaka.Player(video);
await player.configure({
  drm: drmConfig
});
await player.load(manifestUrl);
```

### iOS (FairPlay)

```swift
let asset = AVURLAsset(url: manifestURL)
asset.resourceLoader.setDelegate(self, queue: .main)

let playerItem = AVPlayerItem(asset: asset)
let player = AVPlayer(playerItem: playerItem)

// FairPlay certificate request
func resourceLoader(_ resourceLoader: AVAssetResourceLoader,
                    shouldWaitForLoadingOfRequestedResource loadingRequest: AVAssetResourceLoadingRequest) -> Bool {
    // Handle FairPlay license request
    return true
}
```

### Android (Widevine)

```kotlin
val drmConfig = DrmConfig(
    uuid = C.WIDEVINE_UUID,
    licenseServerUrl = "https://license.widevine.com",
    authToken = getAuthToken()
)

player.setDrmConfig(drmConfig)
```

## Server-Side Token Generation

Token generation implemented in Streaming Service (`services/streaming-service/service/streaming_service.go`):

```go
func GenerateDRMToken(contentID, userID, drmLevel string) (string, error) {
    claims := jwt.MapClaims{
        "content_id": contentID,
        "user_id":    userID,
        "drm_level":  drmLevel,
        "exp":        time.Now().Add(time.Hour).Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(jwtSecret))
}
```

## Testing

### Test Content

Create test content with DRM protection:
- 4K content → Widevine Level 1
- 1080p content → Widevine Level 2 / PlayReady
- SD content → Widevine Level 3
- iOS content → FairPlay

### Test Clients

1. **Widevine**: Chrome browser, Android app
2. **FairPlay**: Safari, iOS app, tvOS app
3. **PlayReady**: Edge browser, Windows app, Smart TVs

## Acceptance Criteria

- ✅ License requests from test clients work
- ✅ Different content tiers (4K, HD, SD) enforce correctly
- ✅ Streaming with DRM protection verified
- ✅ Fallback to non-DRM playback for test content

## Monitoring

Monitor DRM license requests:
- Success rate
- Latency (P95, P99)
- Error rate by DRM provider
- License rejections (unauthorized content)

## Security

- DRM certificates stored in Vault
- License tokens expire after 1 hour
- User authentication required
- Content ID validation

