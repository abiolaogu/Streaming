# StreamVerse for LG webOS

StreamVerse streaming app for LG webOS Smart TVs built with Enact framework (React-based).

## Features

- 🎬 4K/HDR video playback
- 🎮 Magic Remote support with cursor and voice control
- 📺 Full D-pad navigation
- 🔐 User authentication
- 🎯 Personalized recommendations
- 📥 Watchlist management
- 🔍 Advanced search with voice
- 🌐 Multi-language support
- ⚡ Adaptive bitrate streaming (HLS/DASH)
- 🛡️ DRM support (Widevine)

## Prerequisites

- Node.js 16+
- webOS TV SDK (CLI tools)
- LG Developer Account
- Physical LG TV or webOS TV emulator

## Installation

### 1. Install webOS TV SDK

```bash
# macOS
brew install webos-tv-cli

# Linux/Windows - Download from:
# http://webostv.developer.lge.com/sdk/download/download-sdk/
```

### 2. Install Dependencies

```bash
npm install
```

### 3. Configure webOS TV CLI

```bash
# Add your TV
ares-setup-device

# Verify connection
ares-device-info -d tv
```

### 4. Build the App

```bash
# Development build
npm run build

# Production build
npm run build
```

## Development

### Run Development Server

```bash
npm run dev
```

### Test on TV/Emulator

```bash
# Package the app
npm run package

# Install on TV
npm run install-device

# Launch app
npm run launch
```

## Project Structure

```
lg-webos/
├── appinfo.json              # App configuration
├── package.json              # npm dependencies
├── webpack.config.js         # Webpack configuration
├── index.html                # Entry HTML
├── src/
│   ├── App.jsx               # Main App component
│   ├── index.js              # Entry point
│   ├── components/           # Reusable components
│   │   ├── ContentCard.jsx
│   │   ├── ContentRow.jsx
│   │   ├── VideoPlayer.jsx
│   │   ├── Navigation.jsx
│   │   └── SearchBar.jsx
│   ├── screens/              # App screens
│   │   ├── HomeScreen.jsx
│   │   ├── PlayerScreen.jsx
│   │   ├── SearchScreen.jsx
│   │   ├── ProfileScreen.jsx
│   │   └── LoginScreen.jsx
│   ├── services/             # API services
│   │   ├── ApiService.js
│   │   ├── AuthService.js
│   │   ├── ContentService.js
│   │   └── StreamingService.js
│   └── utils/                # Utilities
│       ├── navigation.js
│       └── constants.js
├── assets/
│   ├── images/               # App images
│   └── css/                  # Stylesheets
└── config/
    └── config.js             # App configuration
```

## Key Technologies

### Enact Framework
Enact is LG's React-based framework optimized for webOS TVs:
- **Spotlight**: Focus management system
- **Moonstone**: UI component library
- **i18n**: Internationalization
- **Virtualization**: Efficient rendering of large lists

### webOS Services
Access to native webOS features:
```javascript
// Example: Launch webOS app
webOS.service.request("luna://com.webos.service.applicationmanager", {
    method: "launch",
    parameters: {
        id: "com.streamverse.webos"
    }
});
```

## Navigation

### Magic Remote Support
The app supports LG Magic Remote with:
- **Cursor mode**: Point and click
- **Wheel**: Scroll through content
- **Voice control**: Search via voice
- **D-pad**: Traditional navigation

### Keyboard Support
Full D-pad navigation:
- ↑↓←→: Navigate
- Enter: Select
- Back: Go back
- Home: Return to home

## Video Playback

### Supported Formats
- MP4, WebM
- HLS (HTTP Live Streaming)
- DASH (Dynamic Adaptive Streaming)
- HDR10, Dolby Vision
- Up to 4K resolution

### DRM Support
- Google Widevine Level 1

### Implementation
```javascript
import videojs from 'video.js';

const player = videojs('video-player', {
    controls: true,
    autoplay: false,
    preload: 'auto',
    fluid: true,
    playbackRates: [0.5, 1, 1.5, 2],
    plugins: {
        dash: {
            fastSwitching: true
        }
    }
});

player.src({
    src: 'https://cdn.streamverse.io/content.mpd',
    type: 'application/dash+xml',
    keySystems: {
        'com.widevine.alpha': {
            url: 'https://drm.streamverse.io/widevine'
        }
    }
});
```

## API Integration

### Configuration
```javascript
// config/config.js
export const API_CONFIG = {
    baseURL: 'https://api.streamverse.io/v1',
    timeout: 30000,
    headers: {
        'Content-Type': 'application/json'
    }
};
```

### Service Example
```javascript
// services/ContentService.js
import axios from 'axios';
import { API_CONFIG } from '../config/config';

const api = axios.create(API_CONFIG);

export const getHomeContent = async () => {
    const response = await api.get('/content/home');
    return response.data;
};

export const getContentDetails = async (contentId) => {
    const response = await api.get(`/content/${contentId}`);
    return response.data;
};
```

## webOS Specific Features

### Deep Linking
Handle deep links to specific content:
```javascript
webOS.fetchAppInfo((info) => {
    if (info.params && info.params.contentId) {
        navigateToContent(info.params.contentId);
    }
});
```

### Memory Management
```javascript
// Monitor memory usage
const memoryInfo = webOS.deviceInfo.getMemoryInfo();
console.log('Available memory:', memoryInfo.availableMemory);
```

### Platform Detection
```javascript
const platformInfo = {
    model: webOS.platform.tv.model,
    version: webOS.platform.tv.version,
    resolution: webOS.platform.tv.screenResolution
};
```

## Deployment

### App Packaging

```bash
# Package for distribution
npm run package

# This creates: com.streamverse.webos_1.0.0_all.ipk
```

### LG Content Store Submission

1. **Prepare Assets**:
   - App icon (80x80, 130x130)
   - Screenshots (1920x1080)
   - Promotional images

2. **Test on Multiple TVs**:
   - webOS 3.0+
   - Different screen sizes
   - Various regions

3. **Submit to LG Seller Lounge**:
   - https://seller.lgappstv.com

## Testing

### Unit Tests
```bash
npm test
```

### Manual Testing Checklist
- [ ] Navigation with Magic Remote
- [ ] Navigation with D-pad
- [ ] Video playback
- [ ] DRM content
- [ ] Login/Logout
- [ ] Search functionality
- [ ] Memory leaks
- [ ] Network error handling

## Performance Optimization

### Image Optimization
- Use WebP format when possible
- Lazy load images
- Implement image caching

### Memory Management
- Dispose video players when not in use
- Clean up event listeners
- Use virtualized lists for large content

### Network Optimization
- Implement request caching
- Use CDN for static assets
- Compress API responses

## Troubleshooting

### Common Issues

**1. App won't install on TV**
```bash
# Check TV connection
ares-device-info -d tv

# Verify app package
ares-package -c dist
```

**2. Video playback issues**
- Check video codec compatibility
- Verify DRM license server
- Test on actual TV hardware

**3. Navigation issues**
- Ensure Spotlight is properly configured
- Check focus management
- Test with physical remote

## Resources

- [webOS TV Developer Site](http://webostv.developer.lge.com/)
- [Enact Framework Docs](https://enactjs.com/docs/)
- [LG Seller Lounge](https://seller.lgappstv.com)
- [webOS TV CLI Tools](http://webostv.developer.lge.com/sdk/tools/cli-installation/)

## License

MIT License - see [LICENSE](../../../../LICENSE) for details.

## Support

- Email: tv-support@streamverse.io
- Docs: docs.streamverse.io/webos
