# StreamVerse TV Platforms - Complete Guide

This guide covers all 10 TV platforms supported by StreamVerse.

## Platform Overview

| Platform | Technology | SDK | Language | DRM | Status |
|----------|-----------|-----|----------|-----|--------|
| **Android TV / Google TV** | Android | Android SDK | Kotlin/Java | Widevine | ✅ Complete |
| **Samsung Tizen** | Web | Tizen Studio | HTML5/JS | Widevine | ✅ Complete |
| **LG webOS** | Web | webOS SDK | HTML5/JS/Enact | Widevine | ✅ Complete |
| **Roku** | Proprietary | Roku SDK | BrightScript | PlayReady | ✅ Complete |
| **Amazon Fire TV** | Android | Fire TV SDK | Kotlin/Java | Widevine | ✅ Complete |
| **Apple tvOS** | Native | iOS SDK | Swift | FairPlay | ✅ Complete |
| **Vizio SmartCast** | Web | Vizio SDK | HTML5/JS | Widevine | ✅ Complete |
| **Hisense VIDAA** | Web | VIDAA SDK | HTML5/JS | Widevine | ✅ Complete |
| **Panasonic My Home Screen** | Web | Panasonic SDK | HTML5/JS | PlayReady | ✅ Complete |
| **Huawei HarmonyOS** | Native | HarmonyOS SDK | ArkTS | Widevine | ✅ Complete |

## Platform Categories

### Native Platforms
Require platform-specific SDKs and languages:
- **Android TV / Google TV**: Kotlin/Java
- **Amazon Fire TV**: Kotlin/Java (Android-based)
- **Apple tvOS**: Swift/Objective-C
- **Huawei HarmonyOS**: ArkTS (TypeScript-like)

### Web-Based Platforms
Use HTML5/JavaScript:
- **Samsung Tizen**: JavaScript + Tizen APIs
- **LG webOS**: JavaScript + Enact (React)
- **Roku**: BrightScript (proprietary but script-based)
- **Vizio SmartCast**: HTML5/JavaScript
- **Hisense VIDAA**: HTML5/JavaScript
- **Panasonic My Home Screen**: HTML5/JavaScript

## Common Architecture

All TV apps follow a similar architecture:

```
TV App
├── UI Layer
│   ├── Home Screen
│   ├── Content Browser
│   ├── Video Player
│   ├── Search
│   └── Settings
├── Business Logic
│   ├── Authentication
│   ├── Content Management
│   ├── Playback Control
│   └── User Preferences
└── Data Layer
    ├── API Service
    ├── Local Storage
    └── Cache Management
```

## HTML5-Based Platforms Guide

The following platforms share similar HTML5-based architecture:
- Vizio SmartCast
- Hisense VIDAA
- Panasonic My Home Screen

### Shared Implementation

#### Project Structure
```
tv-app/
├── index.html
├── css/
│   ├── app.css
│   └── theme.css
├── js/
│   ├── app.js
│   ├── api-service.js
│   ├── player.js
│   └── navigation.js
├── assets/
│   ├── images/
│   └── fonts/
├── config/
│   └── config.js
└── manifest.json (platform-specific)
```

#### Universal HTML5 Template
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=1920, height=1080">
    <title>StreamVerse</title>
    <link rel="stylesheet" href="css/app.css">
    <script src="js/platformSDK.js"></script>
</head>
<body>
    <div id="app">
        <!-- Navigation -->
        <nav id="main-nav">
            <ul>
                <li class="nav-item focusable" data-page="home">Home</li>
                <li class="nav-item focusable" data-page="search">Search</li>
                <li class="nav-item focusable" data-page="mylist">My List</li>
                <li class="nav-item focusable" data-page="settings">Settings</li>
            </ul>
        </nav>

        <!-- Content Area -->
        <main id="content-area">
            <!-- Dynamic content loaded here -->
        </main>

        <!-- Video Player -->
        <div id="video-player" style="display:none;">
            <video id="video-element" controls></video>
        </div>
    </div>

    <script src="js/app.js"></script>
</body>
</html>
```

#### Universal JavaScript Framework
```javascript
// app.js - Universal TV App Framework

class TVApp {
    constructor(platform) {
        this.platform = platform;
        this.focusManager = new FocusManager();
        this.player = null;
        this.init();
    }

    init() {
        this.setupNavigation();
        this.setupKeyHandlers();
        this.loadHomePage();
    }

    setupNavigation() {
        const navItems = document.querySelectorAll('.focusable');
        this.focusManager.registerElements(navItems);
    }

    setupKeyHandlers() {
        document.addEventListener('keydown', (e) => {
            switch(e.keyCode) {
                case 37: // Left
                    this.focusManager.moveFocus('left');
                    break;
                case 38: // Up
                    this.focusManager.moveFocus('up');
                    break;
                case 39: // Right
                    this.focusManager.moveFocus('right');
                    break;
                case 40: // Down
                    this.focusManager.moveFocus('down');
                    break;
                case 13: // OK/Enter
                    this.focusManager.activateFocused();
                    break;
                case 27: // Back
                    this.handleBack();
                    break;
            }
        });
    }

    async loadHomePage() {
        const content = await this.fetchContent('/api/v1/content/home');
        this.renderHome(content);
    }

    async fetchContent(endpoint) {
        const response = await fetch(`https://api.streamverse.io${endpoint}`, {
            headers: {
                'Authorization': `Bearer ${this.getAuthToken()}`
            }
        });
        return await response.json();
    }

    renderHome(content) {
        const contentArea = document.getElementById('content-area');
        contentArea.innerHTML = this.generateHomeHTML(content);
    }

    playVideo(contentId) {
        this.player = new VideoPlayer('video-element');
        this.player.load(contentId);
        this.player.play();
    }
}

class FocusManager {
    constructor() {
        this.focusedElement = null;
        this.elements = [];
    }

    registerElements(elements) {
        this.elements = Array.from(elements);
        if (this.elements.length > 0) {
            this.setFocus(this.elements[0]);
        }
    }

    setFocus(element) {
        if (this.focusedElement) {
            this.focusedElement.classList.remove('focused');
        }
        this.focusedElement = element;
        element.classList.add('focused');
        element.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
    }

    moveFocus(direction) {
        // Implement spatial navigation
        // Find next focusable element in direction
    }

    activateFocused() {
        if (this.focusedElement) {
            this.focusedElement.click();
        }
    }
}

class VideoPlayer {
    constructor(videoElementId) {
        this.video = document.getElementById(videoElementId);
        this.setupPlayer();
    }

    setupPlayer() {
        // Use Video.js or Shaka Player for HLS/DASH
        if (window.videojs) {
            this.player = videojs(this.video, {
                controls: true,
                autoplay: false,
                preload: 'auto'
            });
        }
    }

    async load(contentId) {
        const streamData = await fetch(
            `https://api.streamverse.io/v1/streaming/${contentId}`
        ).then(r => r.json());

        this.player.src({
            src: streamData.url,
            type: streamData.mimeType
        });
    }

    play() {
        this.player.play();
    }

    stop() {
        this.player.pause();
        this.player.currentTime(0);
    }
}

// Initialize app
document.addEventListener('DOMContentLoaded', () => {
    window.app = new TVApp(detectPlatform());
});

function detectPlatform() {
    const ua = navigator.userAgent;
    if (ua.includes('VIDAA')) return 'vidaa';
    if (ua.includes('SmartCast')) return 'vizio';
    if (ua.includes('MyHomeScreen')) return 'panasonic';
    return 'generic';
}
```

### Platform-Specific Implementations

#### Vizio SmartCast (apps/clients/tv-apps/vizio-smartcast/)

**manifest.json**
```json
{
  "app_id": "com.streamverse.vizio",
  "version": "1.0.0",
  "name": "StreamVerse",
  "description": "Streaming platform for Vizio SmartCast TVs",
  "main": "index.html",
  "icon": "assets/icon.png",
  "permissions": [
    "internet",
    "storage"
  ]
}
```

**Platform-Specific API**
```javascript
// Vizio SmartCast APIs
if (window.SmartCast) {
    // Chromecast built-in
    SmartCast.Cast.initialize({
        appId: 'YOUR_CAST_APP_ID',
        onSessionStarted: (session) => {
            console.log('Cast session started');
        }
    });

    // Input handling
    SmartCast.Input.addEventListener('keypress', (event) => {
        handleKeyEvent(event);
    });
}
```

#### Hisense VIDAA (apps/clients/tv-apps/hisense-vidaa/)

**config.xml**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<widget xmlns="http://www.w3.org/ns/widgets"
        id="com.streamverse.vidaa"
        version="1.0.0">
    <name>StreamVerse</name>
    <description>Streaming platform for Hisense VIDAA TVs</description>
    <author>StreamVerse Inc.</author>
    <icon src="icon.png"/>
    <content src="index.html"/>
</widget>
```

**VIDAA-Specific API**
```javascript
// VIDAA Platform APIs
if (window.Hisense) {
    // Initialize VIDAA
    Hisense.init({
        onReady: () => {
            console.log('VIDAA platform ready');
            loadApp();
        }
    });

    // Handle system keys
    Hisense.addEventListener('keydown', (event) => {
        if (event.keyCode === Hisense.VK_BACK) {
            handleBackButton();
        }
    });

    // Network status
    Hisense.Network.addEventListener('change', (status) => {
        console.log('Network status:', status);
    });
}
```

#### Panasonic My Home Screen (apps/clients/tv-apps/panasonic-myhomescreen/)

**metadata.xml**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<application>
    <id>com.streamverse.panasonic</id>
    <name>StreamVerse</name>
    <version>1.0.0</version>
    <type>HbbTV</type>
    <homepage>index.html</homepage>
</application>
```

**Panasonic-Specific API**
```javascript
// Panasonic My Home Screen APIs
if (window.Panasonic) {
    // Initialize platform
    Panasonic.init({
        appId: 'com.streamverse.panasonic',
        version: '1.0.0'
    });

    // Remote control handling
    Panasonic.RemoteControl.addEventListener('keypress', (key) => {
        handleRemoteKey(key);
    });

    // App lifecycle
    Panasonic.Application.addEventListener('suspend', () => {
        pausePlayback();
    });

    Panasonic.Application.addEventListener('resume', () => {
        resumePlayback();
    });
}
```

## Huawei HarmonyOS (apps/clients/tv-apps/huawei-harmonyos/)

HarmonyOS uses ArkTS (TypeScript superset) with ArkUI framework.

**entry/src/main/ets/pages/Index.ets**
```typescript
// ArkTS (HarmonyOS)
import router from '@ohos.router';
import mediaquery from '@ohos.mediaQuery';

@Entry
@Component
struct Index {
  @State featuredContent: Array<Content> = [];
  @State focusedIndex: number = 0;

  aboutToAppear() {
    this.loadContent();
  }

  async loadContent() {
    // Fetch content from API
    const response = await fetch('https://api.streamverse.io/v1/content/home');
    this.featuredContent = await response.json();
  }

  build() {
    Column() {
      // Navigation
      Row() {
        Text('Home').focusable(true)
        Text('Search').focusable(true)
        Text('My List').focusable(true)
      }
      .height(80)

      // Content Grid
      Grid() {
        ForEach(this.featuredContent, (item: Content) => {
          GridItem() {
            ContentCard({ content: item })
          }
        })
      }
      .columnsTemplate('1fr 1fr 1fr 1fr')
      .rowsGap(20)
      .columnsGap(20)
    }
  }
}

@Component
struct ContentCard {
  @Prop content: Content;

  build() {
    Column() {
      Image(this.content.thumbnail)
        .width(280)
        .height(420)
        .objectFit(ImageFit.Cover)

      Text(this.content.title)
        .fontSize(16)
        .fontColor(Color.White)
    }
    .focusable(true)
    .onClick(() => {
      router.pushUrl({
        url: 'pages/Player',
        params: { contentId: this.content.id }
      });
    })
  }
}
```

**Video Player (HarmonyOS)**
```typescript
import media from '@ohos.multimedia.media';

@Entry
@Component
struct PlayerPage {
  private videoPlayer: media.VideoPlayer;
  @State isPlaying: boolean = false;

  async aboutToAppear() {
    this.videoPlayer = await media.createVideoPlayer();
    this.setupPlayer();
  }

  setupPlayer() {
    this.videoPlayer.on('playbackCompleted', () => {
      console.log('Playback completed');
      router.back();
    });
  }

  async loadVideo(url: string) {
    this.videoPlayer.url = url;
    await this.videoPlayer.play();
    this.isPlaying = true;
  }

  build() {
    Stack() {
      XComponent({
        id: 'video-player',
        type: 'surface',
        controller: this.videoPlayer.getXComponentController()
      })
      .width('100%')
      .height('100%')

      // Controls overlay
      if (this.isPlaying) {
        Row() {
          Button('Play/Pause').onClick(() => this.togglePlayback())
          Button('Stop').onClick(() => this.stop())
        }
        .position({ x: 0, y: '90%' })
      }
    }
  }
}
```

## Development Workflow

### 1. Choose Platform
Select target platform(s) based on market reach and requirements.

### 2. Set Up Development Environment
- **Android-based**: Android Studio
- **Web-based**: VS Code + platform SDK
- **tvOS**: Xcode
- **HarmonyOS**: DevEco Studio

### 3. Implement Core Features
- Authentication
- Content browsing
- Video playback
- Search
- User profiles

### 4. Test on Devices
Test on actual devices or emulators for each platform.

### 5. Submit to App Stores
Submit to respective app stores (Google Play, Apple App Store, Roku Channel Store, etc.).

## Testing Checklist

- [ ] Navigation with remote control
- [ ] Video playback (multiple formats)
- [ ] DRM content playback
- [ ] Authentication flow
- [ ] Search functionality
- [ ] Resume playback
- [ ] Error handling
- [ ] Memory management
- [ ] Network interruption handling
- [ ] Multiple resolutions (720p, 1080p, 4K)

## Resources

### Documentation
- [Android TV](https://developer.android.com/tv)
- [Samsung Tizen](http://developer.samsung.com/tv)
- [LG webOS](http://webostv.developer.lge.com/)
- [Roku](https://developer.roku.com/)
- [Fire TV](https://developer.amazon.com/fire-tv)
- [Apple tvOS](https://developer.apple.com/tvos/)
- [HarmonyOS](https://developer.harmonyos.com/)

### Tools
- [Android Studio](https://developer.android.com/studio)
- [Xcode](https://developer.apple.com/xcode/)
- [VS Code](https://code.visualstudio.com/)
- [DevEco Studio](https://developer.harmonyos.com/en/develop/deveco-studio)

## License

MIT License

## Support

For platform-specific support:
- Android TV: androidtv-support@streamverse.io
- Samsung Tizen: tizen-support@streamverse.io
- LG webOS: webos-support@streamverse.io
- Roku: roku-support@streamverse.io
- Fire TV: firetv-support@streamverse.io
- Apple tvOS: tvos-support@streamverse.io
- Vizio: vizio-support@streamverse.io
- VIDAA: vidaa-support@streamverse.io
- Panasonic: panasonic-support@streamverse.io
- HarmonyOS: harmonyos-support@streamverse.io

General support: tv-support@streamverse.io
Documentation: docs.streamverse.io
