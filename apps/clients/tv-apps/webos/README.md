# LG webOS TV App - StreamVerse

Native LG webOS TV application for StreamVerse platform.

## Overview

This webOS app provides a native TV experience for the StreamVerse platform, supporting:
- User authentication
- Content browsing and search
- Video playback with webOS MediaPlayer
- PlayReady/Widevine DRM protection
- Magic Remote navigation

## Architecture

- **Platform**: LG webOS TV (webOS 7.0+)
- **Technology**: Web technologies (HTML5, CSS3, JavaScript)
- **Video Player**: HTML5 Video with webOS MediaPlayer
- **DRM**: PlayReady, Widevine
- **Navigation**: Magic Remote (pointer + D-pad)

## Setup

1. **Install webOS CLI**
   - Download from [LG Developer](https://webostv.developer.lge.com/)
   - Install webOS CLI tools

2. **Open Project**
   - Use webOS IDE or CLI
   - Import this directory

3. **Configure API**
   - Edit `js/api.js`
   - Update `API_BASE_URL` if needed

## Building

### Build IPK Package

```bash
# Using webOS CLI
ares-package .

# Or using webOS IDE
# Right-click project → Package as webOS app
```

### Install on TV

```bash
# Connect TV to same network
ares-install com.streamverse.webos_1.0.0_all.ipk -d <device-ip>
```

## Features

- TV-optimized UI with Magic Remote support
- Home screen with content rows
- Content details and playback
- Search functionality
- Settings screen
- DRM support (PlayReady, Widevine)

## Resources

- [webOS Developer Portal](https://webostv.developer.lge.com/)
- [webOS TV API Reference](https://webostv.developer.lge.com/develop/specifications/web-api-reference)

