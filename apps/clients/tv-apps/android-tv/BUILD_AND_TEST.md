# Build and Test Guide

## Prerequisites

- Android Studio Hedgehog (2023.1.1) or later
- Android SDK 34
- JDK 17
- Android TV emulator or physical device

## Building the App

### Debug Build

```bash
cd apps/clients/tv-apps/android-tv
./gradlew assembleDebug
```

Output: `app/build/outputs/apk/debug/app-debug.apk`

### Release Build

```bash
./gradlew assembleRelease
```

**Note**: Requires signing configuration in `build.gradle.kts`

### Install on Device

```bash
# Install debug APK
./gradlew installDebug

# Or use adb
adb install app/build/outputs/apk/debug/app-debug.apk
```

## Configuration

### API Endpoints

Configure in `app/build.gradle.kts`:
```kotlin
buildConfigField("String", "API_BASE_URL", "\"https://api.streamverse.com/\"")
```

For different environments, use build variants or `buildConfigField` per build type.

### DRM License Server

Configure in `app/build.gradle.kts`:
```kotlin
buildConfigField("String", "DRM_LICENSE_SERVER", "\"https://drm.streamverse.com/v1/widevine/license\"")
```

## Testing on Android TV Emulator

1. **Create Android TV Emulator**:
   - Open Android Studio
   - AVD Manager → Create Virtual Device
   - Select TV device (e.g., "TV 1080p")
   - System Image: Android 14 (API 34) TV

2. **Launch Emulator**:
   ```bash
   emulator -avd <your_tv_avd_name>
   ```

3. **Install and Run**:
   ```bash
   ./gradlew installDebug
   ```

4. **Navigate**:
   - Use arrow keys or D-pad on remote
   - Enter key selects
   - Back button returns

## Testing Checklist

### Functional Testing

- [ ] App launches and shows login screen
- [ ] Login with valid credentials
- [ ] Home screen displays content rows
- [ ] Browse by category works
- [ ] Search functionality works
- [ ] Voice search works
- [ ] Content details screen displays
- [ ] Video playback starts
- [ ] DRM-protected content plays
- [ ] Settings screen accessible
- [ ] Recommendations appear on TV home screen

### Performance Testing

- [ ] Smooth navigation
- [ ] Fast content loading
- [ ] No memory leaks
- [ ] Efficient image loading

### TV-Specific Testing

- [ ] D-pad navigation works
- [ ] Remote control works
- [ ] Focus indicators visible
- [ ] UI readable from 10 feet away
- [ ] Landscape orientation maintained

## Troubleshooting

### Build Errors

**Error**: `Unresolved reference: BuildConfig`
- **Solution**: Ensure `buildConfig = true` in `build.gradle.kts`

**Error**: `Cannot find symbol: R`
- **Solution**: Sync Gradle files (`./gradlew --refresh-dependencies`)

### Runtime Errors

**App crashes on launch**:
- Check Logcat for errors
- Verify all dependencies are included
- Check AndroidManifest.xml permissions

**Authentication not working**:
- Verify API_BASE_URL is correct
- Check network connectivity
- Verify login endpoint exists

**Video playback fails**:
- Check DRM license server URL
- Verify stream URLs are accessible
- Check network connectivity

## Debugging

### View Logs

```bash
# View all logs
adb logcat

# Filter by package
adb logcat | grep com.streamverse.tv

# Clear logs
adb logcat -c
```

### Enable Debug Mode

Add to `AndroidManifest.xml`:
```xml
<application
    android:debuggable="true"
    ... >
```

## Release Preparation

1. **Signing Configuration**:
   - Create keystore
   - Add signing config to `build.gradle.kts`

2. **ProGuard Rules**:
   - Configure in `proguard-rules.pro`
   - Test release build thoroughly

3. **Version Management**:
   - Update `versionCode` and `versionName`
   - In `build.gradle.kts`

4. **Testing**:
   - Full functional test on real device
   - Performance testing
   - Security audit

