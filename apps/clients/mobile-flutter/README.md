# StreamVerse Mobile App (Flutter)

A comprehensive, Netflix-level streaming platform mobile application built with Flutter for iOS and Android.

## Features

### 🎬 Content & Streaming
- Adaptive Bitrate Streaming (HLS/DASH)
- Multi-DRM Support (Widevine, FairPlay)
- 4K/HDR Video Playback
- Offline Downloads with background sync
- Picture-in-Picture mode
- Chromecast & AirPlay support
- Continue watching across devices

### 🔐 Authentication
- Email/Password login
- Google Sign-In
- Apple Sign-In
- Biometric authentication (Face ID/Touch ID)
- Multi-factor authentication

### 🎯 Discovery
- AI-powered recommendations
- Smart search with autocomplete
- Voice search
- Visual search
- Trending & Popular content
- Personalized categories

### 👤 User Experience
- Multiple user profiles
- Parental controls
- Watchlist & Favorites
- Custom playlists
- Watch history
- Content ratings & reviews

### 📱 Platform Features
- Push notifications
- In-app purchases
- Subscription management
- Social sharing
- Watch parties
- Dark/Light themes
- Offline mode
- Multi-language support (50+ languages)

## Prerequisites

- Flutter SDK 3.2.0 or higher
- Dart 3.0 or higher
- Xcode 14+ (for iOS)
- Android Studio with SDK 21+ (for Android)
- CocoaPods (for iOS dependencies)

## Installation

### 1. Clone the repository
```bash
git clone https://github.com/yourusername/streamverse.git
cd streamverse/apps/clients/mobile-flutter
```

### 2. Install dependencies
```bash
flutter pub get
```

### 3. Set up environment variables
Create a `.env` file in the root directory:
```env
API_BASE_URL=https://api.streamverse.io
GEMINI_API_KEY=your_gemini_api_key
SENTRY_DSN=your_sentry_dsn
ENVIRONMENT=development
```

### 4. Configure Firebase
- Download `google-services.json` for Android and place in `android/app/`
- Download `GoogleService-Info.plist` for iOS and place in `ios/Runner/`

### 5. Run the app
```bash
# For Android
flutter run -d android

# For iOS
flutter run -d ios

# For both with environment
flutter run --dart-define=API_BASE_URL=https://api.streamverse.io
```

## Project Structure

```
lib/
├── config/               # App configuration & theme
│   ├── app_config.dart
│   └── theme_config.dart
├── models/               # Data models
│   ├── user.dart
│   ├── content.dart
│   ├── video.dart
│   └── subscription.dart
├── screens/              # UI Screens
│   ├── splash_screen.dart
│   ├── login_screen.dart
│   ├── home_screen.dart
│   ├── player_screen.dart
│   ├── search_screen.dart
│   ├── profile_screen.dart
│   └── downloads_screen.dart
├── widgets/              # Reusable widgets
│   ├── content_card.dart
│   ├── content_row.dart
│   ├── video_player.dart
│   └── search_bar.dart
├── services/             # API & Business logic
│   ├── api_service.dart
│   ├── auth_service.dart
│   ├── content_service.dart
│   ├── streaming_service.dart
│   ├── download_service.dart
│   └── notification_service.dart
├── providers/            # State management
│   ├── auth_provider.dart
│   ├── content_provider.dart
│   └── user_provider.dart
├── utils/                # Utilities
│   ├── constants.dart
│   ├── helpers.dart
│   └── validators.dart
└── main.dart             # App entry point
```

## Key Screens

### 1. Splash Screen
- App initialization
- Authentication check
- Preload critical data

### 2. Login Screen
- Email/Password login
- Social sign-in (Google, Apple)
- Biometric authentication
- Registration flow

### 3. Home Screen
- Featured content carousel
- Personalized content rows
- Continue watching
- Trending & Popular
- Category browsing

### 4. Player Screen
- Adaptive video player
- Quality selection
- Subtitle/Audio track selection
- Chromecast integration
- PiP support
- Skip intro/credits

### 5. Search Screen
- Text search with autocomplete
- Voice search
- Visual search
- Filter & sort options
- Search history

### 6. Profile Screen
- User profile management
- Account settings
- Subscription management
- Download management
- App settings

### 7. Downloads Screen
- Downloaded content list
- Download progress
- Storage management
- Offline playback

## Services

### Auth Service
Handles authentication and user management:
- Login/Logout
- Registration
- Password reset
- Token management
- Biometric auth

### Content Service
Manages content discovery and metadata:
- Fetch content lists
- Content details
- Categories & genres
- Recommendations
- Search

### Streaming Service
Handles video playback:
- Video URL generation
- DRM license acquisition
- Quality adaptation
- Playback tracking

### Download Service
Manages offline content:
- Download queue
- Background downloads
- Storage management
- Expiration handling

### Notification Service
Handles push notifications:
- FCM integration
- Local notifications
- In-app notifications
- Notification preferences

## State Management

This app uses **Riverpod** for state management:

```dart
// Example provider
final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  return AuthNotifier();
});

// Usage in widgets
class HomeScreen extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authState = ref.watch(authProvider);
    // Build UI based on state
  }
}
```

## Video Player

Uses **better_player** for advanced video playback features:
- HLS/DASH support
- DRM support (Widevine, FairPlay)
- Adaptive bitrate
- Subtitle support
- Picture-in-Picture
- Background playback

## Offline Downloads

Implements background downloads using **flutter_downloader**:
- Queue management
- Parallel downloads
- Resume capability
- Download notifications
- Storage limits

## DRM Integration

### Widevine (Android)
```dart
BetterPlayerDataSource dataSource = BetterPlayerDataSource(
  BetterPlayerDataSourceType.network,
  videoUrl,
  drmConfiguration: BetterPlayerDrmConfiguration(
    drmType: BetterPlayerDrmType.widevine,
    licenseUrl: AppConfig.widevineLicenseUrl,
  ),
);
```

### FairPlay (iOS)
```dart
BetterPlayerDataSource dataSource = BetterPlayerDataSource(
  BetterPlayerDataSourceType.network,
  videoUrl,
  drmConfiguration: BetterPlayerDrmConfiguration(
    drmType: BetterPlayerDrmType.fairplay,
    licenseUrl: AppConfig.fairPlayLicenseUrl,
    certificateUrl: AppConfig.fairPlayCertificateUrl,
  ),
);
```

## Chromecast Integration

```dart
import 'package:flutter_cast_video/flutter_cast_video.dart';

// Initialize Chromecast
ChromeCastController chromecastController = ChromeCastController();

// Cast video
chromecastController.loadMedia(
  videoUrl,
  title: content.title,
  image: content.thumbnailUrl,
);
```

## Push Notifications

### Setup FCM
1. Configure Firebase project
2. Add Firebase dependencies
3. Request permissions
4. Handle notification callbacks

```dart
FirebaseMessaging.onMessage.listen((RemoteMessage message) {
  // Handle foreground notification
  showNotification(message);
});
```

## Testing

### Unit Tests
```bash
flutter test
```

### Integration Tests
```bash
flutter test integration_test/
```

### Widget Tests
```bash
flutter test test/widgets/
```

## Building for Production

### Android
```bash
# Build APK
flutter build apk --release

# Build App Bundle
flutter build appbundle --release
```

### iOS
```bash
# Build IPA
flutter build ios --release

# Archive with Xcode
open ios/Runner.xcworkspace
```

## Performance Optimization

- Image caching with `cached_network_image`
- Lazy loading for content lists
- Code splitting with deferred loading
- Tree shaking for unused code
- Ahead-of-time (AOT) compilation

## Security

- API token encryption with `flutter_secure_storage`
- Certificate pinning for API calls
- DRM for content protection
- Biometric authentication
- Secure local storage with Hive encryption

## Analytics

Integrated with Firebase Analytics:
- User engagement tracking
- Content view tracking
- Custom events
- User properties
- Crash reporting with Sentry

## Accessibility

- Screen reader support
- High contrast mode
- Font scaling
- Keyboard navigation
- Audio descriptions

## Localization

Supports 50+ languages using Flutter's built-in i18n:
```dart
MaterialApp(
  localizationsDelegates: [
    GlobalMaterialLocalizations.delegate,
    GlobalWidgetsLocalizations.delegate,
    AppLocalizations.delegate,
  ],
  supportedLocales: [
    const Locale('en', 'US'),
    const Locale('es', 'ES'),
    const Locale('fr', 'FR'),
    // ... 47 more locales
  ],
);
```

## Platform-Specific Features

### iOS
- Face ID / Touch ID
- AirPlay
- Picture-in-Picture
- Universal Links
- Siri Shortcuts

### Android
- Fingerprint / Face Unlock
- Chromecast
- Picture-in-Picture
- App Shortcuts
- Android TV Launcher

## CI/CD

GitHub Actions workflow for automated builds and tests:
```yaml
# .github/workflows/flutter.yml
name: Flutter CI
on: [push, pull_request]
jobs:
  build:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v2
      - uses: subosito/flutter-action@v2
      - run: flutter pub get
      - run: flutter test
      - run: flutter build apk
      - run: flutter build ios --no-codesign
```

## Troubleshooting

### Common Issues

**1. Build fails on iOS**
- Run `pod install` in `ios/` directory
- Clean build folder: `flutter clean`
- Update CocoaPods: `pod repo update`

**2. Video playback issues**
- Check DRM license configuration
- Verify video URL accessibility
- Check network permissions

**3. Download failures**
- Verify storage permissions
- Check available storage space
- Review download service logs

## Contributing

Please see [CONTRIBUTING.md](../../../CONTRIBUTING.md) for contribution guidelines.

## License

MIT License - see [LICENSE](../../../LICENSE) for details.

## Support

- Documentation: [docs.streamverse.io](https://docs.streamverse.io)
- Issues: [GitHub Issues](https://github.com/yourusername/streamverse/issues)
- Email: mobile-support@streamverse.io

---

**Built with Flutter 3.2+ • Supports iOS 13+ and Android 7.0+ (API 24+)**
