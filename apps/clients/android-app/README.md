# StreamVerse Android Mobile App

Native Android mobile application for StreamVerse platform with authentication, video playback, DRM support, and content browsing.

## Overview

This Android app provides a native mobile experience for the StreamVerse platform, supporting:
- User authentication with secure token storage
- Content browsing and search
- Video playback with ExoPlayer
- Widevine DRM protection
- Modern Jetpack Compose UI

## Architecture

- **Language**: Kotlin
- **UI Framework**: Jetpack Compose
- **Min SDK**: 24 (Android 7.0)
- **Target SDK**: 34 (Android 14)
- **Video Player**: ExoPlayer (Media3)
- **DRM**: Widevine
- **Networking**: Retrofit + OkHttp
- **Architecture**: MVVM with StateFlow

## Project Structure

```
android-app/
├── app/
│   ├── src/
│   │   ├── main/
│   │   │   ├── java/com/streamverse/mobile/
│   │   │   │   ├── models/
│   │   │   │   │   ├── Content.kt
│   │   │   │   │   └── AuthModels.kt
│   │   │   │   ├── data/
│   │   │   │   │   ├── api/
│   │   │   │   │   │   ├── AuthApiService.kt
│   │   │   │   │   │   └── ContentApiService.kt
│   │   │   │   │   └── repository/
│   │   │   │   │       ├── AuthRepository.kt
│   │   │   │   │       └── ContentRepository.kt
│   │   │   │   ├── viewmodel/
│   │   │   │   │   ├── AuthViewModel.kt
│   │   │   │   │   └── ContentViewModel.kt
│   │   │   │   ├── ui/
│   │   │   │   │   ├── MainActivity.kt
│   │   │   │   │   ├── StreamVerseApp.kt
│   │   │   │   │   ├── login/
│   │   │   │   │   │   └── LoginScreen.kt
│   │   │   │   │   ├── content/
│   │   │   │   │   │   ├── ContentListScreen.kt
│   │   │   │   │   │   └── detail/
│   │   │   │   │   │       └── ContentDetailScreen.kt
│   │   │   │   │   ├── player/
│   │   │   │   │   │   ├── VideoPlayerScreen.kt
│   │   │   │   │   │   └── drm/
│   │   │   │   │   │       └── DRMHelper.kt
│   │   │   │   │   └── theme/
│   │   │   │   │       ├── Color.kt
│   │   │   │   │       ├── Theme.kt
│   │   │   │   │       └── Type.kt
│   │   │   │   └── res/
│   │   │   └── test/
│   │   └── build.gradle.kts
│   └── build.gradle.kts
└── README.md
```

## Features

### ✅ Implemented

- [x] **Authentication** - Login/logout with SharedPreferences token storage
- [x] **API Integration** - Retrofit-based networking with auth headers
- [x] **Content Browsing** - Home screen with content rows
- [x] **Search** - Content search functionality
- [x] **Video Playback** - ExoPlayer integration with full-screen playback
- [x] **DRM Support** - Widevine configuration helper
- [x] **Token Management** - Secure storage using SharedPreferences
- [x] **MVVM Architecture** - ViewModels with StateFlow for reactive updates
- [x] **Jetpack Compose** - Modern declarative UI

### 🔄 In Progress / TODO

- [ ] Settings screen
- [ ] Offline caching
- [ ] Push notifications
- [ ] Enhanced error handling UI
- [ ] Unit test coverage expansion

## Prerequisites

- Android Studio Hedgehog (2023.1.1) or later
- Android SDK 34
- JDK 17 or later
- Kotlin 1.9+

## Setup

1. Open the project in Android Studio
2. Sync Gradle files
3. Configure API endpoints in `build.gradle.kts`:
   - `API_BASE_URL` = `"https://api.streamverse.com/"`
   - `DRM_LICENSE_SERVER` = `"https://drm.streamverse.com/v1/widevine/license"`
4. Update `applicationId` in `build.gradle.kts`
5. Build and run on emulator or device

## Building

```bash
./gradlew assembleDebug
./gradlew assembleRelease
```

## Testing

```bash
./gradlew test
./gradlew connectedAndroidTest
```

## Configuration

### API Endpoints

Configure in `app/build.gradle.kts`:
```kotlin
buildConfigField("String", "API_BASE_URL", "\"https://api.streamverse.com/\"")
buildConfigField("String", "DRM_LICENSE_SERVER", "\"https://drm.streamverse.com/v1/widevine/license\"")
```

### Authentication

The app expects the backend to provide:
- `POST /api/v1/auth/login` - Returns `AuthResponse` with token
- `POST /api/v1/auth/refresh` - Token refresh endpoint
- `POST /api/v1/auth/logout` - Logout endpoint

### DRM (Widevine)

Widevine license server should accept:
- `Authorization: Bearer <token>` header
- `X-Content-ID` header
- Standard Widevine license request body

## Key Components

### AuthRepository

Manages authentication tokens:
- Stores access and refresh tokens in SharedPreferences
- Manages token expiration
- Provides user info

### ContentRepository

Content data repository:
- Fetches from API with automatic auth token injection
- Handles content operations (search, details, categories)

### ViewModels

Reactive ViewModels using StateFlow:
- `AuthViewModel` - Authentication state
- `ContentViewModel` - Content operations state

## API Integration

The app connects to StreamVerse backend API:
- **Content catalog** - `/api/v1/content/home`, `/api/v1/content/category/{category}`
- **Authentication** - `/api/v1/auth/login`, `/api/v1/auth/refresh`, `/api/v1/auth/logout`
- **Search** - `/api/v1/content/search`
- **Content details** - `/api/v1/content/{id}`

All authenticated requests include `Authorization: Bearer <token>` header automatically.

## Related Documentation

- [iOS App Documentation](../ios-app/README.md)
- [Android TV App Documentation](../tv-apps/android-tv/README.md)

