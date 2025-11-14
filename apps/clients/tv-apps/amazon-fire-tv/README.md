# StreamVerse for Amazon Fire TV

StreamVerse streaming app for Amazon Fire TV devices built with Kotlin/Java (Android-based).

## Overview

Fire TV is based on Android, so development is similar to Android TV with some Fire-specific optimizations and features. This app extends the Android TV implementation with:

- Alexa voice integration
- Fire TV UI guidelines
- Amazon IAP (In-App Purchasing)
- Fire TV-specific recommendations
- Fire App Builder compatibility

## Prerequisites

- Android Studio
- Fire TV device or emulator
- Amazon Developer Account
- Java 11+ or Kotlin 1.9+

## Key Differences from Android TV

### 1. Alexa Integration
```kotlin
import com.amazon.device.messaging.ADM

class AlexaIntegrationManager {
    fun handleAlexaIntent(intent: Intent) {
        when (intent.action) {
            "com.amazon.alexa.action.PLAY" -> {
                val contentId = intent.getStringExtra("contentId")
                playContent(contentId)
            }
            "com.amazon.alexa.action.SEARCH" -> {
                val query = intent.getStringExtra("query")
                searchContent(query)
            }
        }
    }
}
```

### 2. Amazon In-App Purchasing
```kotlin
import com.amazon.device.iap.PurchasingService
import com.amazon.device.iap.PurchasingListener
import com.amazon.device.iap.model.*

class IAPManager : PurchasingListener {
    fun init(context: Context) {
        PurchasingService.registerListener(context, this)
        PurchasingService.getUserData()
    }

    override fun onProductDataResponse(response: ProductDataResponse) {
        when (response.requestStatus) {
            ProductDataResponse.RequestStatus.SUCCESSFUL -> {
                // Handle available products
                val products = response.productData
            }
            else -> {
                // Handle error
            }
        }
    }

    override fun onPurchaseResponse(response: PurchaseResponse) {
        when (response.requestStatus) {
            PurchaseResponse.RequestStatus.SUCCESSFUL -> {
                // Fulfill purchase
                val receipt = response.receipt
                verifyPurchase(receipt)
            }
            else -> {
                // Handle failed purchase
            }
        }
    }
}
```

### 3. Amazon Device Messaging (ADM)
```kotlin
import com.amazon.device.messaging.ADMMessageHandlerBase

class StreamVerseADMHandler : ADMMessageHandlerBase("StreamVerseADMHandler") {
    override fun onMessage(intent: Intent) {
        val message = intent.extras?.getString("message")
        // Handle push notification
        showNotification(message)
    }

    override fun onRegistrationError(errorId: String) {
        Log.e(TAG, "ADM registration error: $errorId")
    }

    override fun onRegistered(registrationId: String) {
        // Send registration ID to server
        sendTokenToServer(registrationId)
    }

    override fun onUnregistered(registrationId: String) {
        // Handle unregistration
    }
}
```

## Fire TV Specific Features

### Content Recommendations
```kotlin
import android.app.Notification
import android.app.NotificationManager
import android.media.tv.TvContract

fun addRecommendation(content: Content) {
    val notification = Notification.Builder(context)
        .setContentTitle(content.title)
        .setContentText(content.description)
        .setLargeIcon(content.poster)
        .setSmallIcon(R.drawable.ic_notification)
        .setContentIntent(createContentIntent(content))
        .setExtras(Bundle().apply {
            putString(TvContract.EXTRA_WATCH_NEXT_TYPE,
                TvContract.WatchNextPrograms.WATCH_NEXT_TYPE_CONTINUE)
        })
        .build()

    val notificationManager = context.getSystemService(Context.NOTIFICATION_SERVICE) as NotificationManager
    notificationManager.notify(content.id.hashCode(), notification)
}
```

### Fire TV Home Screen Integration
```kotlin
// Update home screen content
class HomeScreenContentProvider : ContentProvider() {
    override fun query(uri: Uri, projection: Array<String>?,
                      selection: String?, selectionArgs: Array<String>?,
                      sortOrder: String?): Cursor? {
        // Return recommended content for Fire TV home screen
        return createRecommendationsCursor()
    }
}
```

## Setup Instructions

See the main [Android TV README](../android-tv/README.md) for base implementation, then add Fire-specific features:

1. **Add Amazon dependencies**
```gradle
dependencies {
    // Amazon IAP
    implementation files('libs/in-app-purchasing-2.0.76.jar')

    // Amazon Device Messaging
    implementation files('libs/amazon-device-messaging-1.2.0.jar')

    // Alexa Mobile Voice SDK (optional)
    implementation 'com.amazon.alexa:alexa-mobile-voice-sdk:1.0.0'
}
```

2. **Configure AndroidManifest.xml**
```xml
<!-- Amazon permissions -->
<permission android:name="com.streamverse.permission.RECEIVE_ADM_MESSAGE"
    android:protectionLevel="signature" />
<uses-permission android:name="com.streamverse.permission.RECEIVE_ADM_MESSAGE" />
<uses-permission android:name="com.amazon.device.messaging.permission.RECEIVE" />

<!-- Alexa integration -->
<intent-filter>
    <action android:name="com.amazon.alexa.action.PLAY" />
    <category android:name="android.intent.category.DEFAULT" />
</intent-filter>
```

## Testing on Fire TV

```bash
# Connect via ADB
adb connect <FIRE_TV_IP>:5555

# Install app
adb install app-release.apk

# Launch app
adb shell am start -n com.streamverse.firetv/.MainActivity

# View logs
adb logcat -s StreamVerse
```

## Submission to Amazon Appstore

1. Create app on Amazon Developer Console
2. Fill in metadata and screenshots
3. Upload APK
4. Submit for review

## Resources

- [Fire TV Development](https://developer.amazon.com/fire-tv)
- [Fire App Builder](https://developer.amazon.com/docs/fire-app-builder/overview.html)
- [Amazon IAP](https://developer.amazon.com/docs/in-app-purchasing/iap-overview.html)

## License

MIT License

## Support

- Email: firetv-support@streamverse.io
