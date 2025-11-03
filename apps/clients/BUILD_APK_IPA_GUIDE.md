# Build APK and IPA Guide

Complete guide for building Android APK and iOS IPA files for StreamVerse mobile apps.

## Android APK

### Prerequisites

- Android Studio Hedgehog (2023.1.1) or later
- JDK 17 or later
- Android SDK 34

### Build Debug APK

**Method 1: Android Studio**

1. Open project:
   ```bash
   cd apps/clients/android-app
   # Open in Android Studio: File → Open → Select android-app directory
   ```

2. Build APK:
   - Build → Build Bundle(s) / APK(s) → Build APK(s)
   - Or: `Build → Make Project` (Cmd+B / Ctrl+B)

3. Locate APK:
   ```
   app/build/outputs/apk/debug/app-debug.apk
   ```

**Method 2: Command Line**

```bash
cd apps/clients/android-app

# Build debug APK
./gradlew assembleDebug

# APK location:
# app/build/outputs/apk/debug/app-debug.apk
```

### Build Release APK

**Step 1: Create Keystore**

```bash
keytool -genkey -v -keystore streamverse-release.keystore -alias streamverse -keyalg RSA -keysize 2048 -validity 10000
```

**Step 2: Configure Signing**

Edit `app/build.gradle.kts`:

```kotlin
android {
    signingConfigs {
        create("release") {
            storeFile = file("../streamverse-release.keystore")
            storePassword = "YOUR_STORE_PASSWORD"
            keyAlias = "streamverse"
            keyPassword = "YOUR_KEY_PASSWORD"
        }
    }

    buildTypes {
        release {
            isMinifyEnabled = true
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro"
            )
            signingConfig = signingConfigs.getByName("release")
        }
    }
}
```

**Step 3: Build Release APK**

```bash
./gradlew assembleRelease

# APK location:
# app/build/outputs/apk/release/app-release.apk
```

### Build App Bundle (AAB) for Play Store

```bash
./gradlew bundleRelease

# AAB location:
# app/build/outputs/bundle/release/app-release.aab
```

### Install APK on Device

**Via ADB:**
```bash
adb install app/build/outputs/apk/debug/app-debug.apk
```

**Via USB:**
1. Enable Developer Options on Android device
2. Enable USB Debugging
3. Transfer APK to device
4. Open APK file on device to install

## iOS IPA

### Prerequisites

- **macOS** (required for iOS development)
- Xcode 15.0 or later
- Apple Developer Account (free account works for development)
- CocoaPods (if using pods)

### Build for Development/Testing

**Method 1: Xcode (Recommended)**

1. Open project:
   ```bash
   cd apps/clients/ios-app
   open StreamVerse/StreamVerse.xcodeproj
   ```

2. Select target device:
   - Choose simulator or connected device
   - Or select "Any iOS Device" for archive

3. Build:
   - Product → Build (Cmd+B)
   - For device: Product → Run (Cmd+R)

**Method 2: Command Line**

```bash
cd apps/clients/ios-app/StreamVerse

# Build for simulator
xcodebuild -project StreamVerse.xcodeproj \
           -scheme StreamVerse \
           -sdk iphonesimulator \
           -configuration Debug \
           build

# Build for device
xcodebuild -project StreamVerse.xcodeproj \
           -scheme StreamVerse \
           -sdk iphoneos \
           -configuration Debug \
           build
```

### Build IPA for Distribution

**Method 1: Xcode Archive (Recommended)**

1. **Configure Signing:**
   - Select project in Xcode
   - Select target "StreamVerse"
   - Signing & Capabilities → Enable "Automatically manage signing"
   - Select your Team

2. **Archive:**
   - Product → Archive (or Cmd+Shift+B)
   - Wait for build to complete
   - Organizer window opens

3. **Distribute:**
   - Click "Distribute App"
   - Choose distribution method:
     - **App Store Connect** - For App Store submission
     - **Ad Hoc** - For specific devices (requires UDIDs)
     - **Enterprise** - For enterprise distribution
     - **Development** - For testing

4. **Export IPA:**
   - Follow distribution wizard
   - Choose options (include bitcode, etc.)
   - Export to desired location
   - IPA file will be created

**Method 2: Command Line (Fastlane)**

1. Install Fastlane:
   ```bash
   sudo gem install fastlane
   ```

2. Initialize Fastlane:
   ```bash
   cd apps/clients/ios-app
   fastlane init
   ```

3. Create Fastfile:
   ```ruby
   # fastlane/Fastfile
   default_platform(:ios)

   platform :ios do
     desc "Build IPA"
     lane :build_ipa do
       build_app(
         workspace: "StreamVerse.xcworkspace",
         scheme: "StreamVerse",
         export_method: "ad-hoc"
       )
     end
   end
   ```

4. Run:
   ```bash
   fastlane build_ipa
   ```

**Method 3: xcodebuild Command**

```bash
cd apps/clients/ios-app/StreamVerse

# Archive
xcodebuild archive \
  -project StreamVerse.xcodeproj \
  -scheme StreamVerse \
  -configuration Release \
  -archivePath ./build/StreamVerse.xcarchive

# Export IPA
xcodebuild -exportArchive \
  -archivePath ./build/StreamVerse.xcarchive \
  -exportPath ./build \
  -exportOptionsPlist ExportOptions.plist
```

**Create ExportOptions.plist:**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>method</key>
    <string>ad-hoc</string>
    <key>teamID</key>
    <string>YOUR_TEAM_ID</string>
</dict>
</plist>
```

### Install IPA on Device

**Via Xcode:**
1. Connect device via USB
2. In Xcode: Window → Devices and Simulators
3. Select device
4. Drag IPA to "Installed Apps"

**Via TestFlight:**
1. Upload IPA to App Store Connect
2. Add testers in TestFlight section
3. Testers receive email invitation
4. Install via TestFlight app

**Via Ad Hoc Distribution:**
1. Register device UDIDs in Apple Developer Portal
2. Build IPA with Ad Hoc distribution
3. Transfer IPA to device (email, AirDrop, etc.)
4. Install via Settings → General → VPN & Device Management

### Build for App Store Submission

1. **Archive in Xcode:**
   - Product → Archive
   - Wait for completion

2. **Upload to App Store Connect:**
   - In Organizer, click "Distribute App"
   - Choose "App Store Connect"
   - Follow wizard
   - Click "Upload"

3. **Submit for Review:**
   - Go to App Store Connect
   - My Apps → Select app
   - Create new version
   - Uploaded build appears
   - Fill in metadata and submit

## Troubleshooting

### Android

**Build fails:**
- Sync Gradle: `./gradlew --refresh-dependencies`
- Clean build: `./gradlew clean`
- Invalidate caches in Android Studio

**APK won't install:**
- Uninstall existing app first
- Check device architecture matches APK
- Enable "Install from unknown sources"

### iOS

**Code signing errors:**
- Check Team selected in project settings
- Verify certificates in Keychain Access
- Regenerate provisioning profiles

**Archive fails:**
- Clean build folder: Product → Clean Build Folder (Cmd+Shift+K)
- Check build settings
- Verify all required files included in target

**IPA won't install:**
- Verify device UDID registered (for Ad Hoc)
- Check certificate matches device
- Verify provisioning profile includes device

## Quick Reference

### Android Commands

```bash
# Debug APK
./gradlew assembleDebug

# Release APK
./gradlew assembleRelease

# Release AAB (Play Store)
./gradlew bundleRelease

# Install on connected device
./gradlew installDebug
```

### iOS Commands

```bash
# Build for simulator
xcodebuild -project StreamVerse.xcodeproj -scheme StreamVerse -sdk iphonesimulator build

# Archive
xcodebuild archive -project StreamVerse.xcodeproj -scheme StreamVerse

# Clean
xcodebuild clean -project StreamVerse.xcodeproj -scheme StreamVerse
```

## Output Locations

### Android
- Debug APK: `app/build/outputs/apk/debug/app-debug.apk`
- Release APK: `app/build/outputs/apk/release/app-release.apk`
- Release AAB: `app/build/outputs/bundle/release/app-release.aab`

### iOS
- Archive: `~/Library/Developer/Xcode/Archives/`
- IPA: Location specified during export (usually `~/Desktop` or build directory)

## Next Steps

1. **Test APK/IPA** on real devices
2. **Beta Testing**:
   - Android: Upload to Play Console → Internal/Closed testing
   - iOS: Upload to TestFlight
3. **Production Release**:
   - Android: Submit to Play Store
   - iOS: Submit to App Store

