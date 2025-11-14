# StreamVerse for Apple tvOS

StreamVerse streaming app for Apple TV built with Swift and UIKit/SwiftUI.

## Features

- 🎬 4K/HDR10/Dolby Vision playback
- 🎮 Siri Remote support with touch gestures
- 🔐 Sign in with Apple
- 🎯 Top Shelf integration
- 📥 Watchlist sync via iCloud
- 🔍 Siri search integration
- 🌐 Universal Search
- ⚡ Adaptive bitrate streaming (HLS)
- 🛡️ FairPlay DRM
- 📺 AirPlay sender support
- 🎮 Game controller support

## Prerequisites

- macOS 13.0+
- Xcode 15.0+
- tvOS 16.0+ target
- Apple TV 4K (recommended) or Apple TV HD
- Apple Developer Program membership

## Installation

### 1. Clone Repository

```bash
cd apps/clients/tv-apps/apple-tvos
```

### 2. Install Dependencies

```bash
# Using CocoaPods
pod install

# Using Swift Package Manager (SPM)
# Dependencies are managed in Xcode
```

### 3. Open Project

```bash
open StreamVerse.xcworkspace
```

### 4. Configure Signing

1. Select project in Xcode
2. Select target "StreamVerse"
3. Signing & Capabilities tab
4. Select your team
5. Enable required capabilities

### 5. Build and Run

```bash
# Command line
xcodebuild -workspace StreamVerse.xcworkspace \
           -scheme StreamVerse \
           -destination 'platform=tvOS Simulator,name=Apple TV' \
           build

# Or press Cmd+R in Xcode
```

## Project Structure

```
StreamVerse/
├── StreamVerse.xcodeproj       # Xcode project
├── StreamVerse/
│   ├── App/
│   │   ├── StreamVerseApp.swift # App entry point
│   │   └── SceneDelegate.swift  # Scene lifecycle
│   ├── Views/                   # SwiftUI/UIKit views
│   │   ├── HomeView.swift
│   │   ├── PlayerView.swift
│   │   ├── SearchView.swift
│   │   ├── ProfileView.swift
│   │   └── Components/
│   │       ├── ContentCard.swift
│   │       ├── ContentRow.swift
│   │       └── VideoPlayerView.swift
│   ├── ViewModels/              # MVVM view models
│   │   ├── HomeViewModel.swift
│   │   ├── PlayerViewModel.swift
│   │   └── SearchViewModel.swift
│   ├── Models/                  # Data models
│   │   ├── Content.swift
│   │   ├── User.swift
│   │   └── Video.swift
│   ├── Services/                # Business logic
│   │   ├── APIService.swift
│   │   ├── AuthService.swift
│   │   ├── ContentService.swift
│   │   ├── StreamingService.swift
│   │   └── DRMService.swift
│   ├── Resources/               # Assets
│   │   ├── Assets.xcassets
│   │   ├── Info.plist
│   │   └── TopShelf/
│   └── Utilities/
│       ├── Constants.swift
│       ├── Extensions/
│       └── Helpers/
└── StreamVerseTopShelf/         # Top Shelf extension
    └── ContentProvider.swift
```

## Key Technologies

### Swift & SwiftUI
- Swift 5.9+
- SwiftUI for declarative UI
- Combine for reactive programming
- Async/await for concurrency

### AVFoundation
- AVPlayer for video playback
- AVPlayerViewController
- FairPlay Streaming (FPS)

### Frameworks
- UIKit (legacy support)
- TVUIKit (tvOS-specific components)
- CloudKit (iCloud sync)
- AuthenticationServices (Sign in with Apple)

## Architecture

### MVVM Pattern
```swift
// Model
struct Content: Codable, Identifiable {
    let id: String
    let title: String
    let description: String
    let thumbnailURL: URL
    let streamURL: URL
}

// ViewModel
@MainActor
class HomeViewModel: ObservableObject {
    @Published var featuredContent: [Content] = []
    @Published var trendingContent: [Content] = []
    @Published var isLoading = false

    private let contentService: ContentService

    func loadContent() async {
        isLoading = true
        defer { isLoading = false }

        do {
            featuredContent = try await contentService.getFeaturedContent()
            trendingContent = try await contentService.getTrendingContent()
        } catch {
            print("Error loading content: \(error)")
        }
    }
}

// View
struct HomeView: View {
    @StateObject private var viewModel = HomeViewModel()

    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 40) {
                FeaturedContentView(content: viewModel.featuredContent)
                ContentRowView(title: "Trending", items: viewModel.trendingContent)
            }
        }
        .task {
            await viewModel.loadContent()
        }
    }
}
```

## Video Playback

### AVPlayer Implementation
```swift
import AVFoundation
import AVKit

class VideoPlayerViewController: AVPlayerViewController {
    private var player: AVPlayer?

    func playContent(_ content: Content) {
        // Create player item
        let asset = AVURLAsset(url: content.streamURL)
        let playerItem = AVPlayerItem(asset: asset)

        // Configure FairPlay DRM if needed
        if content.requiresDRM {
            configureFairPlayDRM(for: asset, with: content.drmLicenseURL)
        }

        // Create and configure player
        player = AVPlayer(playerItem: playerItem)
        self.player = player

        // Add observers
        addObservers()

        // Start playback
        player?.play()
    }

    private func addObservers() {
        // Playback status
        player?.currentItem?.addObserver(self,
            forKeyPath: "status",
            options: [.new, .initial],
            context: nil)

        // Time observation
        player?.addPeriodicTimeObserver(
            forInterval: CMTime(seconds: 1, preferredTimescale: 1),
            queue: .main
        ) { [weak self] time in
            self?.updatePlaybackProgress(time)
        }
    }
}
```

### FairPlay DRM
```swift
import AVFoundation

class FairPlayManager: NSObject, AVAssetResourceLoaderDelegate {
    func configureFairPlay(for asset: AVURLAsset, licenseURL: URL) {
        asset.resourceLoader.setDelegate(self, queue: DispatchQueue.main)
    }

    func resourceLoader(_ resourceLoader: AVAssetResourceLoader,
                       shouldWaitForLoadingOfRequestedResource loadingRequest: AVAssetResourceLoadingRequest) -> Bool {
        guard let url = loadingRequest.request.url else { return false }

        // Get certificate
        guard let certificate = getCertificate() else {
            loadingRequest.finishLoading(with: NSError(domain: "FairPlay", code: -1))
            return false
        }

        // Get content identifier
        guard let contentId = url.host?.data(using: .utf8) else {
            return false
        }

        // Create SPC (Server Playback Context)
        do {
            let spcData = try loadingRequest.streamingContentKeyRequestData(
                forApp: certificate,
                contentIdentifier: contentId,
                options: nil
            )

            // Request CKC (Content Key Context) from license server
            requestCKC(spcData: spcData) { ckcData in
                loadingRequest.dataRequest?.respond(with: ckcData)
                loadingRequest.finishLoading()
            }

            return true
        } catch {
            loadingRequest.finishLoading(with: error)
            return false
        }
    }

    private func getCertificate() -> Data? {
        // Fetch FairPlay certificate from server
        // Implementation depends on your DRM provider
        return nil
    }

    private func requestCKC(spcData: Data, completion: @escaping (Data) -> Void) {
        // Request CKC from license server
        // Implementation depends on your DRM provider
    }
}
```

## Focus Engine

### Focus Management
```swift
import UIKit

class ContentCardView: UIView {
    override var canBecomeFocused: Bool {
        return true
    }

    override func didUpdateFocus(in context: UIFocusUpdateContext,
                                 with coordinator: UIFocusAnimationCoordinator) {
        super.didUpdateFocus(in: context, with: coordinator)

        coordinator.addCoordinatedAnimations {
            if self.isFocused {
                // Scale up when focused
                self.transform = CGAffineTransform(scaleX: 1.1, y: 1.1)
                self.layer.shadowOpacity = 0.5
            } else {
                // Scale down when unfocused
                self.transform = .identity
                self.layer.shadowOpacity = 0
            }
        }
    }
}
```

## Siri Integration

### Siri Intents
```swift
import Intents

class IntentHandler: INExtension, INPlayMediaIntentHandling {
    func handle(intent: INPlayMediaIntent,
                completion: @escaping (INPlayMediaIntentResponse) -> Void) {
        // Handle "Hey Siri, play [content] on StreamVerse"
        guard let mediaItem = intent.mediaItems?.first else {
            completion(INPlayMediaIntentResponse(code: .failure, userActivity: nil))
            return
        }

        // Find and play content
        let userActivity = NSUserActivity(activityType: "com.streamverse.play")
        userActivity.userInfo = ["contentId": mediaItem.identifier]

        completion(INPlayMediaIntentResponse(code: .handleInApp, userActivity: userActivity))
    }
}
```

### Universal Search
```swift
import CoreSpotlight
import MobileCoreServices

func indexContent(_ content: Content) {
    // Create searchable item
    let attributeSet = CSSearchableItemAttributeSet(itemContentType: kUTTypeVideo as String)
    attributeSet.title = content.title
    attributeSet.contentDescription = content.description
    attributeSet.thumbnailURL = content.thumbnailURL
    attributeSet.contentRating = NSNumber(value: content.rating)
    attributeSet.duration = NSNumber(value: content.duration)

    let item = CSSearchableItem(
        uniqueIdentifier: content.id,
        domainIdentifier: "com.streamverse.content",
        attributeSet: attributeSet
    )

    // Index item
    CSSearchableIndex.default().indexSearchableItems([item]) { error in
        if let error = error {
            print("Indexing error: \(error)")
        }
    }
}
```

## Top Shelf

### Top Shelf Extension
```swift
import TVServices

class ContentProvider: TVTopShelfContentProvider {
    override func loadTopShelfContent(completionHandler: @escaping (TVTopShelfContent?) -> Void) {
        // Fetch featured content
        fetchFeaturedContent { items in
            let identifier = TVTopShelfSectionedIdentifier(
                sectionIdentifiers: ["featured"],
                itemIdentifiers: items.map { $0.id }
            )

            let items = items.map { content -> TVTopShelfSectionedItem in
                let item = TVTopShelfSectionedItem(identifier: content.id)
                item.setImageURL(content.thumbnailURL, for: .screenScale1x)
                item.title = content.title
                item.displayAction = TVTopShelfAction(url: URL(string: "streamverse://play/\(content.id)")!)
                return item
            }

            let section = TVTopShelfItemCollection(items: items)
            section.title = "Featured"

            let content = TVTopShelfSectionedContent(sections: [identifier], items: [section])
            completionHandler(content)
        }
    }
}
```

## Authentication

### Sign in with Apple
```swift
import AuthenticationServices

class AuthService: NSObject, ASAuthorizationControllerDelegate {
    func signInWithApple() {
        let provider = ASAuthorizationAppleIDProvider()
        let request = provider.createRequest()
        request.requestedScopes = [.fullName, .email]

        let controller = ASAuthorizationController(authorizationRequests: [request])
        controller.delegate = self
        controller.performRequests()
    }

    func authorizationController(controller: ASAuthorizationController,
                                 didCompleteWithAuthorization authorization: ASAuthorization) {
        guard let credential = authorization.credential as? ASAuthorizationAppleIDCredential else {
            return
        }

        // Send to backend
        authenticateWithBackend(
            userId: credential.user,
            identityToken: credential.identityToken,
            authorizationCode: credential.authorizationCode
        )
    }
}
```

## Deployment

### App Store Connect

1. **Create App Record**
   - Go to App Store Connect
   - Create new tvOS app
   - Fill in metadata

2. **Prepare Assets**
   - App icon (1280x768)
   - Screenshots (1920x1080)
   - App preview video

3. **Archive and Upload**
   ```bash
   xcodebuild -workspace StreamVerse.xcworkspace \
              -scheme StreamVerse \
              -archivePath StreamVerse.xcarchive \
              archive

   xcodebuild -exportArchive \
              -archivePath StreamVerse.xcarchive \
              -exportPath StreamVerse.ipa \
              -exportOptionsPlist ExportOptions.plist
   ```

4. **Submit for Review**

## Testing

### Unit Tests
```swift
import XCTest
@testable import StreamVerse

class ContentServiceTests: XCTestCase {
    var sut: ContentService!

    override func setUp() {
        super.setUp()
        sut = ContentService()
    }

    func testFetchFeaturedContent() async throws {
        let content = try await sut.getFeaturedContent()
        XCTAssertFalse(content.isEmpty)
    }
}
```

### UI Tests
```swift
import XCTest

class StreamVerseUITests: XCTestCase {
    func testNavigation() {
        let app = XCUIApplication()
        app.launch()

        // Test focus movement
        let remote = XCUIRemote.shared
        remote.press(.down)
        remote.press(.select)

        XCTAssertTrue(app.otherElements["PlayerView"].exists)
    }
}
```

## Resources

- [tvOS Documentation](https://developer.apple.com/tvos/)
- [Human Interface Guidelines](https://developer.apple.com/design/human-interface-guidelines/tvos)
- [FairPlay Streaming](https://developer.apple.com/streaming/fps/)
- [App Store Review Guidelines](https://developer.apple.com/app-store/review/guidelines/)

## License

MIT License

## Support

- Email: tvos-support@streamverse.io
- Docs: docs.streamverse.io/tvos
