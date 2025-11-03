package com.streamverse.mobile.ui.player.drm

import android.content.Context
import android.content.SharedPreferences
import androidx.media3.common.util.UnstableApi
import androidx.media3.exoplayer.drm.DefaultDrmSessionManagerProvider
import androidx.media3.exoplayer.drm.FrameworkMediaDrm
import androidx.media3.exoplayer.drm.HttpMediaDrmCallback
import androidx.media3.exoplayer.ExoPlayer
import androidx.media3.common.C
import com.streamverse.mobile.BuildConfig
import com.streamverse.mobile.models.Content
import com.streamverse.mobile.data.repository.AuthRepository
import okhttp3.OkHttpClient
import java.util.*

/**
 * Helper for DRM configuration (Widevine).
 */
object DRMHelper {
    private const val WIDEVINE_UUID = "edef8ba9-79d6-4ace-a3c8-27dcd51d21ed"
    private const val PREFS_NAME = "streamverse_prefs"
    private const val KEY_AUTH_TOKEN = "auth_token"

    @UnstableApi
    fun configureDRM(player: ExoPlayer, content: Content) {
        val licenseServerUrl = content.drmType?.let {
            when (it.lowercase()) {
                "widevine" -> BuildConfig.DRM_LICENSE_SERVER
                else -> BuildConfig.DRM_LICENSE_SERVER
            }
        } ?: BuildConfig.DRM_LICENSE_SERVER

        val drmCallback = HttpMediaDrmCallback(
            licenseServerUrl,
            OkHttpClient()
        )

        // Add auth token to license requests
        val context = player.context
        val authToken = getAuthToken(context)
        if (authToken.isNotEmpty()) {
            drmCallback.setKeyRequestProperty("Authorization", "Bearer $authToken")
        }
        drmCallback.setKeyRequestProperty("Content-Type", "application/octet-stream")
        drmCallback.setKeyRequestProperty("X-Content-ID", content.id)

        val drmSessionManagerProvider = DefaultDrmSessionManagerProvider()
        drmSessionManagerProvider.setDrmHttpCallbackFactory { drmCallback }
        drmSessionManagerProvider.setDrmUuid(UUID.fromString(WIDEVINE_UUID))

        player.drmSessionManagerProvider = drmSessionManagerProvider
    }

    private fun getAuthToken(context: Context): String {
        val prefs: SharedPreferences = context.getSharedPreferences(PREFS_NAME, Context.MODE_PRIVATE)
        return prefs.getString(KEY_AUTH_TOKEN, "") ?: ""
    }
}

