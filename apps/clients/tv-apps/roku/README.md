# StreamVerse for Roku

StreamVerse streaming app for Roku devices built with BrightScript and SceneGraph.

## Features

- 🎬 4K/HDR video playback
- 🎮 Remote control navigation
- 🔐 User authentication
- 🎯 Personalized recommendations
- 📥 Watchlist management
- 🔍 Search functionality
- 🌐 Multi-language support
- ⚡ Adaptive bitrate streaming (HLS)
- 🛡️ DRM support (PlayReady)
- 📺 Deep linking support

## Prerequisites

- Roku device (Roku OS 9.0+) or Roku emulator
- Roku Developer Account
- Eclipse/VS Code with BrightScript extension (optional)

## Installation

### 1. Enable Developer Mode on Roku

1. Press **Home** button 3 times
2. Press **Up** 2 times
3. Press **Right, Left, Right, Left, Right**
4. Set a development password

### 2. Access Developer Portal

Navigate to `http://YOUR_ROKU_IP` in a browser

### 3. Package and Install

#### Using Web Interface:
1. Go to `http://YOUR_ROKU_IP`
2. Upload the zipped project folder
3. Click "Install"

#### Using Command Line:
```bash
# Install roku-deploy tool
npm install -g roku-deploy

# Create deployment config
echo '{
  "host": "YOUR_ROKU_IP",
  "password": "YOUR_DEV_PASSWORD",
  "rootDir": ".",
  "files": ["**/*"],
  "retainStagingFolder": false
}' > roku-deploy.json

# Deploy
roku-deploy
```

## Project Structure

```
roku/
├── manifest                    # App configuration
├── source/                     # BrightScript source files
│   └── main.brs                # Entry point
├── components/                 # SceneGraph components
│   ├── MainScene.xml           # Main scene XML
│   ├── MainScene.brs           # Main scene logic
│   ├── ContentRow.xml          # Content row component
│   ├── VideoPlayer.xml         # Custom video player
│   └── SearchScreen.xml        # Search interface
├── images/                     # App images
│   ├── channel-icon_HD.png     # Channel icon (290x218)
│   ├── channel-icon_FHD.png    # Channel icon (540x405)
│   ├── splash-screen_HD.jpg    # Splash screen (1280x720)
│   └── splash-screen_FHD.jpg   # Splash screen (1920x1080)
└── fonts/                      # Custom fonts (optional)
```

## Key Technologies

### BrightScript
Roku's scripting language based on BASIC:
- Dynamically typed
- Event-driven
- Optimized for Roku hardware

### SceneGraph (RSG)
Roku's declarative UI framework:
- XML-based component definitions
- Render thread optimization
- Focus management
- Animation support

## Navigation

### Remote Control Support
- **D-pad (↑↓←→)**: Navigate through content
- **OK/Select**: Select content
- **Back**: Go back one screen
- **Home**: Exit app
- *****: Options menu
- **Play/Pause**: Control playback
- **Rev/Fwd**: Seek backward/forward

### Focus Management
```xml
<Button id="playButton"
    focusable="true"
    focusBitmapUri="pkg:/images/button_focus.png" />
```

## Video Playback

### Supported Formats
- HLS (HTTP Live Streaming) - Primary
- Smooth Streaming
- MP4 (Progressive download)
- HDR10, Dolby Vision (on supported devices)
- Up to 4K resolution

### DRM Support
- PlayReady DRM

### Video Player Implementation
```brightscript
' Create video content node
videoContent = CreateObject("roSGNode", "ContentNode")
videoContent.url = "https://cdn.streamverse.io/content.m3u8"
videoContent.title = "Content Title"
videoContent.streamFormat = "hls"

' DRM configuration
videoContent.encodingType = "PlayReadyLicenseAcquisitionUrl"
videoContent.encodingKey = "https://drm.streamverse.io/playready"

' Set content and play
m.videoPlayer.content = videoContent
m.videoPlayer.control = "play"
```

### Adaptive Bitrate
HLS automatically adapts quality based on:
- Network bandwidth
- Device capabilities
- Buffer health

## API Integration

### HTTP Requests
```brightscript
port = CreateObject("roMessagePort")
urlTransfer = CreateObject("roUrlTransfer")
urlTransfer.SetPort(port)
urlTransfer.SetUrl("https://api.streamverse.io/v1/content/home")
urlTransfer.AddHeader("Authorization", "Bearer " + authToken)
urlTransfer.SetMessagePort(port)

if urlTransfer.AsyncGetToString() then
    msg = wait(5000, port)
    if type(msg) = "roUrlEvent" then
        if msg.GetResponseCode() = 200 then
            response = ParseJson(msg.GetString())
            ' Process response
        end if
    end if
end if
```

### JSON Parsing
```brightscript
' Parse JSON string
jsonString = '{"title":"Movie","id":123}'
jsonObject = ParseJson(jsonString)
print jsonObject.title  ' Outputs: Movie

' Convert to JSON
brightscriptObject = {title: "Movie", id: 123}
jsonString = FormatJson(brightscriptObject)
```

## Deep Linking

### Handle Deep Links
```brightscript
sub Main(args as Dynamic)
    if args.contentID <> invalid then
        ' Launch directly to content
        launchContent(args.contentID, args.mediaType)
    else
        ' Normal launch
        showHomeScreen()
    end if
end sub
```

### Create Deep Link
```
http://my.roku.com/launch/YOUR_CHANNEL_ID?contentID=12345&mediaType=movie
```

## State Management

### Registry (Persistent Storage)
```brightscript
' Write to registry
sec = CreateObject("roRegistrySection", "StreamVerseData")
sec.Write("authToken", token)
sec.Write("userId", userId)
sec.Flush()

' Read from registry
authToken = sec.Read("authToken")

' Delete from registry
sec.Delete("authToken")
sec.Flush()
```

## Performance Optimization

### Memory Management
```brightscript
' Check available memory
memoryInfo = CreateObject("roDeviceInfo").GetGeneralMemoryLevel()
print "Memory level: "; memoryInfo

' Clear unused objects
m.largeArray = invalid
m.videoPlayer.content = invalid
```

### Content Caching
```brightscript
' Use roUrlTransfer with caching
urlTransfer.EnableCookies()
urlTransfer.SetCertificatesFile("common:/certs/ca-bundle.crt")
urlTransfer.InitClientCertificates()
```

### List Virtualization
```xml
<RowList id="contentList"
    numRows="4"
    rowItemSize="[[280, 420]]"
    itemComponentName="ContentCard" />
```

## Testing

### Manual Testing
1. Side-load app on device
2. Test all navigation paths
3. Test video playback
4. Test different resolutions
5. Test on multiple Roku models

### Debug Console
Enable debugging:
```brightscript
' Print to console
print "Debug message: "; variable

' Telnet to Roku for logs
telnet YOUR_ROKU_IP 8085
```

### Profiling
```brightscript
' Measure execution time
startTime = CreateObject("roTimespan")
' ... code to profile ...
print "Execution time: "; startTime.TotalMilliseconds(); "ms"
```

## Certification Requirements

Before submitting to Roku Channel Store:

### Technical Requirements
- [ ] App must work on all Roku devices
- [ ] Support 720p and 1080p
- [ ] No crashes or freezes
- [ ] Memory usage < 100MB (recommend < 80MB)
- [ ] Fast channel launch (< 5 seconds)
- [ ] Proper error handling

### Content Requirements
- [ ] All content must be rated
- [ ] Parental controls implementation
- [ ] Terms of service
- [ ] Privacy policy

### UI Requirements
- [ ] Consistent navigation
- [ ] Proper focus indicators
- [ ] Loading indicators
- [ ] Error messages

## Deployment

### Create Package

1. **Via Web Interface**:
   - Go to `http://YOUR_ROKU_IP`
   - Installer > Select app > Package
   - Sign with production keys

2. **Required Assets**:
   - Channel poster (540x405)
   - Screenshots (1920x1080 or 1280x720)
   - HD channel poster for older devices

### Submit to Roku

1. **Developer Dashboard**:
   - https://developer.roku.com/

2. **Upload Package**:
   - Upload signed .pkg file
   - Provide metadata
   - Submit screenshots

3. **Review Process**:
   - Typically 2-5 business days
   - Address any feedback

## Troubleshooting

### Common Issues

**1. App won't side-load**
```bash
# Verify Roku IP
ping YOUR_ROKU_IP

# Check developer mode is enabled
# Re-enable if needed
```

**2. Video won't play**
- Check video URL is accessible
- Verify HLS manifest format
- Test on Roku hardware (emulator has limitations)
- Check DRM license server

**3. Memory issues**
```brightscript
' Monitor memory
deviceInfo = CreateObject("roDeviceInfo")
print "Memory level: "; deviceInfo.GetGeneralMemoryLevel()

' Clear unused resources
m.unusedNode = invalid
```

**4. Navigation issues**
- Ensure focusable="true" on interactive elements
- Check focus chain with setFocus()
- Verify key event handling

## Resources

- [Roku Developer Documentation](https://developer.roku.com/docs/developer-program/getting-started/roku-dev-prog.md)
- [BrightScript Reference](https://developer.roku.com/docs/references/brightscript/language/brightscript-language-reference.md)
- [SceneGraph API Reference](https://developer.roku.com/docs/references/scenegraph/xml-elements/xml-elements-overview.md)
- [Roku Forums](https://community.roku.com/t5/Developers/ct-p/channel-developers)

## License

MIT License - see [LICENSE](../../../../LICENSE) for details.

## Support

- Email: roku-support@streamverse.io
- Docs: docs.streamverse.io/roku
