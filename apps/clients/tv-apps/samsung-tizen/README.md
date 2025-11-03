# Samsung Tizen TV App - StreamVerse

Native Samsung Tizen TV application for StreamVerse platform.

## Overview

This Tizen TV app provides a native TV experience for the StreamVerse platform, supporting:
- User authentication
- Content browsing and search
- Video playback with Tizen AVPlay
- PlayReady DRM protection
- Remote control navigation

## Architecture

- **Platform**: Samsung Tizen TV (Tizen 6.0+)
- **Technology**: Web technologies (HTML5, CSS3, JavaScript)
- **Video Player**: Tizen AVPlay API
- **DRM**: PlayReady
- **Navigation**: Spatial navigation for remote control

## Project Structure

```
samsung-tizen/
├── config.xml              # Tizen app configuration
├── index.html              # Main HTML file
├── css/
│   └── style.css           # Styles
├── js/
│   ├── api.js              # API service
│   ├── auth.js             # Authentication
│   ├── content.js          # Content management
│   ├── player.js           # Tizen AVPlay player
│   ├── navigation.js       # TV navigation
│   └── main.js             # Main application logic
├── icon.png                # App icon
└── README.md
```

## Features

### ✅ Implemented

- [x] **Tizen Web App Setup** - Project structure and configuration
- [x] **TV-Optimized UI** - Spatial navigation support
- [x] **Home Screen** - Content rows with horizontal scrolling
- [x] **Content Details** - Full content information
- [x] **Video Player** - Tizen AVPlay integration
- [x] **DRM Support** - PlayReady configuration
- [x] **Remote Control Navigation** - D-pad and button support
- [x] **Search** - Content search functionality
- [x] **Settings** - User preferences
- [x] **Authentication** - Login/logout with token storage

## Prerequisites

- **Tizen Studio** - Samsung's development IDE
- **Tizen SDK** - Tizen 6.0 or later
- **Samsung TV** - For testing (or Tizen emulator)

## Setup

1. **Install Tizen Studio**
   - Download from [Samsung Developers](https://developer.samsung.com/tizen)
   - Install Tizen Studio with TV extension

2. **Open Project**
   - Launch Tizen Studio
   - File → Import → Tizen → Tizen Web Project
   - Select this directory

3. **Configure API**
   - Edit `js/api.js`
   - Update `API_BASE_URL` if needed

## Building

### Build WGT Package

1. In Tizen Studio:
   - Right-click project → Run As → Tizen Web Application
   - Or use CLI: `tizen package -t wgt -s <profile>`

2. Command Line:
   ```bash
   tizen package -t wgt
   ```

### Certificate and Signing

1. Create certificate in Tizen Studio:
   - Tools → Certificate Manager
   - Create new certificate
   - Author certificate

2. Sign package:
   ```bash
   tizen sign -p <package> -c <certificate>
   ```

## Testing

### Emulator

1. Launch Tizen TV Emulator
2. Install app: `tizen install -n <package.wgt> -t <device-id>`
3. Test navigation and playback

### Real Device

1. Enable Developer Mode on Samsung TV
2. Connect TV to same network as development machine
3. Install via Tizen Studio or CLI

## Deployment

### Samsung App Store

1. Build signed WGT package
2. Create app in Samsung Seller Office
3. Upload WGT package
4. Submit for certification

### Sideloading

1. Build WGT package
2. Transfer to TV via USB or network
3. Install via Smart Hub

## Configuration

### API Endpoints

Edit `js/api.js`:
```javascript
const API_BASE_URL = 'https://api.streamverse.com/';
```

### DRM License Server

Edit `js/player.js`:
```javascript
const licenseServer = 'https://drm.streamverse.com/v1/playready/license';
```

## Key Components

### TVNavigation

Handles remote control navigation:
- D-pad movement (Up, Down, Left, Right)
- Enter/Return key
- Back button
- Focus management

### TizenPlayer

Tizen AVPlay wrapper:
- Stream playback
- DRM configuration
- Playback controls
- Error handling

### ApiService

Network layer:
- Authentication
- Content fetching
- Token management
- Error handling

## Tizen Certification Requirements

- [ ] App must not crash
- [ ] All features functional
- [ ] DRM must work correctly
- [ ] Navigation must be smooth
- [ ] Performance requirements met
- [ ] Privacy policy included
- [ ] Terms of service included

## Troubleshooting

### Player Not Initializing

- Check Tizen version (requires 6.0+)
- Verify AVPlay API availability
- Check console for errors

### DRM Not Working

- Verify PlayReady certificate
- Check license server URL
- Verify token is being sent

### Navigation Issues

- Ensure elements have tabindex
- Check focus styles
- Verify key event handlers

## Resources

- [Tizen Developer Portal](https://developer.tizen.org/)
- [Tizen TV Documentation](https://developer.samsung.com/tv)
- [AVPlay API Reference](https://developer.samsung.com/tv/develop/api-references/tizen-web-device-api/avplay)
- [Samsung Seller Office](https://seller.samsungapps.com/)

## License

Copyright © 2025 StreamVerse. All rights reserved.

