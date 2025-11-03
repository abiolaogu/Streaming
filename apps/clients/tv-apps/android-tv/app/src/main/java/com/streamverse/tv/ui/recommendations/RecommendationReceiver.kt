package com.streamverse.tv.ui.recommendations

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent

/**
 * Broadcast receiver for updating Android TV home screen recommendations.
 * Triggered on boot and periodically.
 */
class RecommendationReceiver : BroadcastReceiver() {

    override fun onReceive(context: Context, intent: Intent) {
        if (Intent.ACTION_BOOT_COMPLETED == intent.action) {
            RecommendationUpdateService.startService(context)
        }
    }
}

