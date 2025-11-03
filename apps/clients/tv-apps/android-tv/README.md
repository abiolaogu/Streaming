# Android TV / Google TV App

Android TV application optimized for 10-foot UI experience with Leanback library, ExoPlayer, and DRM support.

## Overview

This Android TV app provides a native TV-optimized experience for the StreamVerse platform, supporting:
- HLS/DASH streaming with adaptive bitrate
- DRM (Widevine) protection
- Voice search
- Android TV home screen recommendations
- Leanback UI components

## Architecture

- **Language**: Kotlin
- **Min SDK**: 26 (Android 8.0)
- **Target SDK**: 34 (Android 14)
- **TV UI Framework**: AndroidX Leanback Library
- **Video Player**: ExoPlayer
- **DRM**: Widevine
- **Architecture**: MVVM with Jetpack Compose (optional) or traditional Views

## Project Structure

```
android-tv/
├── app/
│   ├── build.gradle
│   ├── src/
│   │   ├── main/
│   │   │   ├── AndroidManifest.xml
│   │   │   ├── java/com/streamverse/tv/
│   │   │   │   ├── MainActivity.kt
│   │   │   │   ├── MainFragment.kt
│   │   │   │   ├── BrowseFragment.kt
│   │   │   │   ├── DetailsFragment.kt
│   │   │   │   ├── PlaybackVideoActivity.kt
│   │   │   │   ├── SearchFragment.kt
│   │   │   │   ├── SettingsFragment.kt
│   │   │   │   ├── data/
│   │   │   │   │   ├── api/
│   │   │   │   │   ├── local/
│   │   │   │   │   └── model/
│   │   │   │   ├── ui/
│   │   │   │   │   ├── adapter/
│   │   │   │   │   ├── player/
│   │   │   │   │   └── presenter/
│   │   │   │   └── viewmodel/
│   │   │   └── res/
│   │   │       ├── layout/
│   │   │       ├── values/
│   │   │       └── drawable/
│   │   └── androidTest/
│   ├── proguard-rules.pro
│   └── build.gradle.kts
├── build.gradle.kts
├── settings.gradle.kts
├── gradle.properties
└── README.md
```

## Features

### ✅ Implemented (Issue #32)

- [x] **Project structure setup** - Complete with Gradle configuration
- [x] **Leanback library integration** - All fragments use Leanback components
- [x] **TV-optimized navigation** - D-pad and remote control support throughout
- [x] **Home screen with content rows** - MainFragment with dynamic content rows
- [x] **Browse fragments** - Movies, shows, live TV, FAST channels
- [x] **Details screen** - Full details with poster, description, actions
- [x] **Video player** - ExoPlayer with HLS/DASH and adaptive bitrate
- [x] **DRM support** - Widevine integration with license server
- [x] **Search with voice input** - Voice recognition enabled
- [x] **Settings** - Preference screen with playback and subtitle options
- [x] **Android TV recommendations** - Foundation complete (service + receiver)

### ✅ Completed (Latest Updates)

- [x] **Authentication flow** - Complete login system with token management
- [x] **API integration** - Retrofit with auth interceptors and error handling
- [x] **Search API integration** - Full SearchViewModel with API connection
- [x] **BuildConfig** - API URLs and DRM configuration via BuildConfig
- [x] **Drawable resources** - Icons, banners, button backgrounds
- [x] **Unit tests** - ViewModel and Repository test suites
- [x] **Testing documentation** - Complete testing guide and build instructions

### 🔄 Remaining TODOs

- [ ] Full TvContract ContentProvider implementation for recommendations
- [ ] Error UI screens
- [ ] Integration tests on real device/emulator
- [ ] Production signing configuration
- [ ] ProGuard rules for release builds

### Recent Updates (2025-01-XX)

- ✅ **Authentication System** - Complete login flow with AuthViewModel, AuthRepository, and LoginActivity/LoginFragment
- ✅ **API Integration** - Retrofit interceptors for automatic token injection, BuildConfig for endpoints
- ✅ **Search Enhancement** - Full SearchViewModel integration with API, voice search support
- ✅ **Testing Framework** - Unit tests for ViewModels and Repositories with Mockito-Kotlin
- ✅ **Build Configuration** - BuildConfig fields for API_BASE_URL and DRM_LICENSE_SERVER
- ✅ **Drawable Resources** - Launcher banner, button backgrounds, edit text backgrounds
- ✅ **Documentation** - Testing guide, build instructions, and updated README

## Prerequisites

- Android Studio Hedgehog (2023.1.1) or later
- Android SDK 34
- JDK 17 or later
- Kotlin 1.9+
- Android TV emulator or physical device for testing

See [BUILD_AND_TEST.md](BUILD_AND_TEST.md) for detailed setup and testing instructions.

## Setup

1. Open the project in Android Studio
2. Sync Gradle files
3. Configure your API endpoint in `BuildConfig` or `strings.xml`
4. Add your Widevine license server URL
5. Build and run on Android TV emulator or device

## Building

```bash
./gradlew assembleDebug
./gradlew assembleRelease
```

## Testing

### Unit Tests

```bash
# Run all unit tests
./gradlew test

# Run specific test class
./gradlew test --tests "com.streamverse.tv.viewmodel.MainViewModelTest"

# Generate test coverage report
./gradlew testDebugUnitTest jacocoTestReport
```

### Integration Tests

```bash
# Run on connected Android TV device/emulator
./gradlew connectedAndroidTest
```

See [TESTING_GUIDE.md](TESTING_GUIDE.md) for detailed testing instructions.

## Deployment

For sideloading or Play Store submission, see [DEPLOYMENT.md](DEPLOYMENT.md)

## API Integration

The app connects to the StreamVerse backend API:
- **Content catalog endpoint** - `/api/v1/content/home`, `/api/v1/content/category/{category}`
- **Authentication endpoint** - `/api/v1/auth/login`, `/api/v1/auth/refresh`, `/api/v1/auth/logout`
- **Search endpoint** - `/api/v1/content/search`
- **Content details** - `/api/v1/content/{id}`
- **DRM license server** - Configured via BuildConfig

### Configuration

Endpoints are configured in `app/build.gradle.kts`:

```kotlin
buildConfigField("String", "API_BASE_URL", "\"https://api.streamverse.com/\"")
buildConfigField("String", "DRM_LICENSE_SERVER", "\"https://drm.streamverse.com/v1/widevine/license\"")
```

Authentication tokens are automatically injected via OkHttp interceptors in `ContentRepository` and `AuthRepository`.

## DRM Configuration

Widevine DRM requires:
- License server URL
- API key for authentication
- Content ID mapping

See `app/src/main/java/com/streamverse/tv/ui/player/DRMHelper.kt`

## Android TV Recommendations

The app provides recommendations to Android TV home screen:
- Continue watching
- Trending content
- Personalized recommendations
- New releases

Implemented via `TvContract` API.

