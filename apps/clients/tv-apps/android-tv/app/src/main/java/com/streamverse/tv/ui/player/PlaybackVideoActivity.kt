package com.streamverse.tv.ui.player

import android.content.Intent
import android.net.Uri
import android.os.Bundle
import android.view.View
import androidx.appcompat.app.AppCompatActivity
import androidx.media3.common.MediaItem
import androidx.media3.common.Player
import androidx.media3.exoplayer.ExoPlayer
import androidx.media3.exoplayer.drm.DefaultDrmSessionManagerProvider
import androidx.media3.exoplayer.drm.FrameworkMediaDrm
import androidx.media3.exoplayer.drm.HttpMediaDrmCallback
import androidx.media3.ui.PlayerView
import com.streamverse.tv.R
import com.streamverse.tv.data.model.Content
import com.streamverse.tv.ui.player.drm.DRMHelper
import okhttp3.OkHttpClient

/**
 * Video playback activity using ExoPlayer with HLS/DASH and DRM support.
 */
class PlaybackVideoActivity : AppCompatActivity() {

    private var player: ExoPlayer? = null
    private var playerView: PlayerView? = null
    private var content: Content? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_playback)

        content = intent.getParcelableExtra(EXTRA_CONTENT)
        playerView = findViewById(R.id.player_view)

        initializePlayer()
    }

    private fun initializePlayer() {
        val exoPlayer = ExoPlayer.Builder(this)
            .setDrmSessionManagerProvider(DefaultDrmSessionManagerProvider())
            .build()
            .also { player = it }

        content?.let { content ->
            // Build media item with DRM if needed
            val mediaItem = if (content.isDrmProtected) {
                buildDrmMediaItem(content)
            } else {
                MediaItem.fromUri(Uri.parse(content.streamUrl))
            }

            exoPlayer.setMediaItem(mediaItem)
            exoPlayer.prepare()
            exoPlayer.playWhenReady = true
        }

        playerView?.player = exoPlayer
    }

    private fun buildDrmMediaItem(content: Content): MediaItem {
        val drmConfig = DRMHelper.getWidevineConfig(content, this)
        
        return MediaItem.Builder()
            .setUri(Uri.parse(content.streamUrl))
            .setDrmConfiguration(
                androidx.media3.common.DrmConfiguration.Builder(
                    androidx.media3.common.C.WIDEVINE_UUID
                )
                    .setLicenseRequestHeaders(drmConfig.licenseRequestHeaders)
                    .setLicenseUri(Uri.parse(drmConfig.licenseServerUrl))
                    .build()
            )
            .build()
    }

    override fun onStart() {
        super.onStart()
        player?.let {
            if (it.playWhenReady) {
                it.play()
            }
        }
    }

    override fun onResume() {
        super.onResume()
        playerView?.onResume()
    }

    override fun onPause() {
        super.onPause()
        playerView?.onPause()
    }

    override fun onStop() {
        super.onStop()
        player?.pause()
    }

    override fun onDestroy() {
        super.onDestroy()
        releasePlayer()
    }

    private fun releasePlayer() {
        player?.release()
        player = null
        playerView?.player = null
    }

    companion object {
        const val EXTRA_CONTENT = "content"
        const val EXTRA_STREAM_URL = "stream_url"
    }
}

