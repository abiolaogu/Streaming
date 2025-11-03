# StreamVerse iOS App

Native iOS application for StreamVerse platform with authentication, video playback, DRM support, and content browsing.

## Overview

This iOS app provides a native experience for the StreamVerse platform, supporting:
- User authentication with secure token storage
- Content browsing and search
- Video playback with AVPlayer
- FairPlay DRM protection
- SwiftUI-based UI

## Architecture

- **Language**: Swift 5.9+
- **UI Framework**: SwiftUI
- **Minimum iOS**: 16.0
- **Video Player**: AVPlayer / AVKit
- **DRM**: FairPlay
- **Networking**: URLSession
- **Architecture**: MVVM with Combine

## Project Structure

```
ios-app/
├── StreamVerse/
│   ├── Models/
│   │   ├── Content.swift
│   │   └── AuthModels.swift
│   ├── Services/
│   │   ├── APIService.swift
│   │   ├── AuthService.swift
│   │   ├── ContentService.swift
│   │   ├── TokenManager.swift
│   │   └── DRMHelper.swift
│   ├── ViewModels/
│   │   ├── AuthViewModel.swift
│   │   └── ContentViewModel.swift
│   ├── Views/
│   │   ├── StreamVerseApp.swift
│   │   ├── LoginView.swift
│   │   ├── ContentListView.swift
│   │   ├── ContentDetailView.swift
│   │   └── VideoPlayerView.swift
│   └── Info.plist
├── StreamVerseTests/
│   ├── AuthViewModelTests.swift
│   └── ContentViewModelTests.swift
└── README.md
```

## Features

### ✅ Implemented

- [x] **Authentication** - Login/logout with secure token storage in Keychain
- [x] **API Integration** - URLSession-based networking with auth headers
- [x] **Content Browsing** - Home screen with content rows
- [x] **Search** - Content search functionality
- [x] **Video Playback** - AVPlayer integration with full-screen playback
- [x] **DRM Support** - FairPlay configuration helper
- [x] **Token Management** - Secure storage using Keychain
- [x] **MVVM Architecture** - ViewModels with Combine for reactive updates

### 🔄 In Progress / TODO

- [ ] Complete FairPlay DRM implementation
- [ ] Settings screen
- [ ] Offline caching
- [ ] Push notifications
- [ ] Error handling UI
- [ ] Unit test coverage expansion

## Prerequisites

- Xcode 15.0 or later
- iOS 16.0+ SDK
- Swift 5.9+
- CocoaPods or Swift Package Manager (if needed)

## Setup

1. Open `StreamVerse.xcodeproj` in Xcode
2. Configure API endpoints in Build Settings:
   - `API_BASE_URL` = `https://api.streamverse.com/`
   - `DRM_LICENSE_SERVER` = `https://drm.streamverse.com/v1/fairplay/license`
3. Update Bundle Identifier in project settings
4. Build and run on simulator or device

## Configuration

### API Endpoints

Configure in Xcode Build Settings or Info.plist:
- `API_BASE_URL` - Base URL for StreamVerse API
- `DRM_LICENSE_SERVER` - FairPlay license server URL

### Authentication

The app expects the backend to provide:
- `POST /api/v1/auth/login` - Returns `AuthResponse` with token
- `POST /api/v1/auth/refresh` - Token refresh endpoint
- `POST /api/v1/auth/logout` - Logout endpoint

### DRM (FairPlay)

FairPlay license server should accept:
- `Authorization: Bearer <token>` header
- `X-Content-ID` header
- Standard FairPlay license request body

## Building

### Debug Build

```bash
xcodebuild -project StreamVerse.xcodeproj \
           -scheme StreamVerse \
           -configuration Debug \
           -sdk iphonesimulator \
           build
```

### Release Build

```bash
xcodebuild -project StreamVerse.xcodeproj \
           -scheme StreamVerse \
           -configuration Release \
           -sdk iphoneos \
           build
```

## Testing

### Unit Tests

```bash
xcodebuild test -project StreamVerse.xcodeproj \
                -scheme StreamVerse \
                -destination 'platform=iOS Simulator,name=iPhone 15'
```

Or run tests in Xcode: `Cmd + U`

## Key Components

### TokenManager

Manages authentication tokens securely using Keychain Services:
- Stores access and refresh tokens in Keychain
- Manages token expiration
- Provides user info from UserDefaults

### APIService

Base networking service:
- Configures URLSession with timeouts
- Automatically injects auth tokens
- Handles JSON encoding/decoding
- Provides error handling

### DRMHelper

FairPlay DRM configuration:
- Configures AVAssetResourceLoader
- Handles license requests
- Injects auth tokens for license server

### ViewModels

Reactive ViewModels using Combine:
- `AuthViewModel` - Authentication state
- `ContentViewModel` - Content operations

## API Integration

The app connects to StreamVerse backend API:
- **Content catalog** - `/api/v1/content/home`, `/api/v1/content/category/{category}`
- **Authentication** - `/api/v1/auth/login`, `/api/v1/auth/refresh`, `/api/v1/auth/logout`
- **Search** - `/api/v1/content/search`
- **Content details** - `/api/v1/content/{id}`

All authenticated requests include `Authorization: Bearer <token>` header automatically.

## Security

- **Keychain Storage** - Sensitive tokens stored in iOS Keychain
- **HTTPS Only** - All API calls use HTTPS
- **Token Expiration** - Automatic token expiry checking
- **Secure Coding** - Swift best practices for security

## Deployment

### App Store

1. Update version and build number
2. Configure signing certificates
3. Archive and upload via Xcode
4. Submit for App Store review

### TestFlight

1. Archive the app
2. Upload to App Store Connect
3. Add testers to TestFlight
4. Distribute beta builds

## Troubleshooting

### Build Errors

**Error**: `Cannot find type in scope`
- **Solution**: Ensure all files are added to the Xcode project target

**Error**: `Missing Info.plist key`
- **Solution**: Add required keys to Info.plist or use Build Settings

### Runtime Errors

**App crashes on launch**:
- Check Info.plist configuration
- Verify API endpoints are set
- Check console logs

**Authentication not working**:
- Verify `API_BASE_URL` is correct
- Check network connectivity
- Verify login endpoint exists

**Video playback fails**:
- Check DRM license server URL
- Verify stream URLs are accessible
- Check network connectivity
- Ensure FairPlay certificate is configured

## Related Documentation

- [Android TV App Documentation](../tv-apps/android-tv/README.md)
- [StreamVerse API Documentation](../../../docs/API.md)

## License

Copyright © 2025 StreamVerse. All rights reserved.

