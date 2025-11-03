# Build and Test Guide - iOS App

## Prerequisites

- Xcode 15.0 or later
- macOS 13.0 or later
- iOS 16.0+ SDK
- Apple Developer Account (for device testing)
- CocoaPods (optional, for additional dependencies)

## Building the App

### Debug Build

1. **Open in Xcode**:
   ```bash
   open StreamVerse/StreamVerse.xcodeproj
   ```

2. **Select Scheme**:
   - Choose `StreamVerse` scheme
   - Select target device (Simulator or Physical Device)

3. **Build**:
   - Press `Cmd + B` or
   - Product → Build

### Release Build

1. **Archive**:
   - Product → Archive
   - Wait for build to complete

2. **Distribute**:
   - Click "Distribute App"
   - Choose distribution method (App Store, Ad Hoc, Enterprise)

## Configuration

### API Endpoints

Configure in Xcode Build Settings:

1. Select project in navigator
2. Select target "StreamVerse"
3. Go to Build Settings
4. Search for "User-Defined Settings"
5. Add:
   - `API_BASE_URL` = `https://api.streamverse.com/`
   - `DRM_LICENSE_SERVER` = `https://drm.streamverse.com/v1/fairplay/license`

Alternatively, edit `Info.plist`:
```xml
<key>API_BASE_URL</key>
<string>https://api.streamverse.com/</string>
<key>DRM_LICENSE_SERVER</key>
<string>https://drm.streamverse.com/v1/fairplay/license</string>
```

### Bundle Identifier

Update Bundle Identifier in:
- Project Settings → General → Bundle Identifier
- Example: `com.streamverse.app`

## Testing

### Unit Tests

**In Xcode**:
1. Press `Cmd + U` or
2. Product → Test

**Command Line**:
```bash
xcodebuild test \
  -project StreamVerse.xcodeproj \
  -scheme StreamVerse \
  -destination 'platform=iOS Simulator,name=iPhone 15'
```

### UI Tests

1. Create UI test target (if needed)
2. Write UI tests in `StreamVerseUITests`
3. Run tests: `Cmd + U`

### Manual Testing Checklist

- [ ] App launches successfully
- [ ] Login screen displays
- [ ] Login with valid credentials works
- [ ] Home screen displays content rows
- [ ] Content cards are clickable
- [ ] Content detail screen displays
- [ ] Video playback starts
- [ ] DRM-protected content plays
- [ ] Search functionality works
- [ ] Logout works correctly

## Running on Device

### Requirements

- Apple Developer Account
- Provisioning Profile
- Connected iOS device

### Steps

1. **Connect Device**:
   - Connect iPhone/iPad via USB
   - Trust computer if prompted

2. **Configure Signing**:
   - Select target device in Xcode
   - Xcode → Preferences → Accounts
   - Add Apple ID
   - Select team in project settings

3. **Run**:
   - Press `Cmd + R` or
   - Product → Run

## Debugging

### View Logs

**In Xcode**:
- View → Debug Area → Show Debug Area (`Cmd + Shift + Y`)
- Console shows runtime logs

**Command Line**:
```bash
# View device logs
xcrun simctl spawn booted log stream --level=debug

# View specific device
instruments -t "System Trace" -D trace.trace your.app
```

### Breakpoints

1. Set breakpoint: Click line number in editor
2. Run app in debug mode
3. App pauses at breakpoint
4. Inspect variables in debug area

### Network Debugging

Use tools like:
- **Charles Proxy**
- **Proxyman**
- **Network Link Conditioner** (built-in)

## Performance Testing

### Instruments

1. Product → Profile (`Cmd + I`)
2. Choose profiling template:
   - **Time Profiler** - CPU usage
   - **Allocations** - Memory usage
   - **Leaks** - Memory leaks
   - **Network** - Network activity

### Memory Management

Check for:
- Retain cycles
- Memory leaks
- Excessive allocations
- Image caching issues

## Troubleshooting

### Build Errors

**Error**: `No such module 'Foundation'`
- **Solution**: Clean build folder (`Cmd + Shift + K`) and rebuild

**Error**: `Signing for "StreamVerse" requires a development team`
- **Solution**: Select team in project settings → Signing & Capabilities

**Error**: `Command PhaseScriptExecution failed`
- **Solution**: Check script phases in Build Phases

### Runtime Errors

**App crashes on launch**:
- Check Info.plist configuration
- Verify all required files are included in target
- Check console logs for crash details

**Authentication not working**:
- Verify `API_BASE_URL` is correct
- Check network connectivity
- Verify backend API is accessible
- Check token storage in Keychain

**Video playback fails**:
- Check DRM license server URL
- Verify stream URLs are accessible
- Check network connectivity
- Ensure FairPlay certificate is configured
- Check AVPlayer error logs

**Keychain access denied**:
- Verify Keychain Sharing entitlement
- Check app's Keychain access group
- Ensure proper code signing

## Code Signing

### Development

1. Automatic Signing:
   - Project Settings → Signing & Capabilities
   - Enable "Automatically manage signing"
   - Select development team

### Distribution

1. Manual Signing:
   - Disable automatic signing
   - Select distribution certificate
   - Select provisioning profile

### Certificates

Manage in:
- Xcode → Preferences → Accounts
- [Apple Developer Portal](https://developer.apple.com/account)

## Release Preparation

1. **Version Management**:
   - Update version in Info.plist or Build Settings
   - Increment build number

2. **Code Signing**:
   - Ensure proper distribution certificate
   - Configure provisioning profile

3. **Archive**:
   - Product → Archive
   - Validate archive
   - Distribute to App Store

4. **Testing**:
   - TestFlight beta testing
   - Full functional testing
   - Performance testing

## CI/CD Integration

### Fastlane

Example `Fastfile`:
```ruby
lane :beta do
  build_app(scheme: "StreamVerse")
  upload_to_testflight
end
```

### GitHub Actions

Example workflow:
```yaml
name: Build and Test
on: [push]
jobs:
  test:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v2
      - name: Build
        run: xcodebuild build -project StreamVerse.xcodeproj
      - name: Test
        run: xcodebuild test -project StreamVerse.xcodeproj
```

## Additional Resources

- [Apple Developer Documentation](https://developer.apple.com/documentation/)
- [SwiftUI Tutorials](https://developer.apple.com/tutorials/swiftui)
- [AVFoundation Guide](https://developer.apple.com/av-foundation/)
- [FairPlay Streaming](https://developer.apple.com/streaming/fps/)

