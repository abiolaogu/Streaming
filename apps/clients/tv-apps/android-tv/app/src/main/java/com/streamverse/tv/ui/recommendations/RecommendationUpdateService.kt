package com.streamverse.tv.ui.recommendations

import android.app.Notification
import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.content.Context
import android.content.Intent
import android.graphics.BitmapFactory
import android.os.Build
import androidx.core.app.JobIntentService
import androidx.core.app.NotificationCompat
import com.streamverse.tv.MainActivity
import com.streamverse.tv.R
import com.streamverse.tv.data.model.Content
import com.streamverse.tv.data.repository.ContentRepository
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch

/**
 * Service that updates Android TV home screen recommendations.
 */
class RecommendationUpdateService : JobIntentService() {
    
    private val serviceScope = CoroutineScope(Dispatchers.IO)

    companion object {
        private const val JOB_ID = 1001
        private const val CHANNEL_ID = "recommendations"
        private const val MAX_RECOMMENDATIONS = 50

        fun startService(context: Context) {
            val intent = Intent(context, RecommendationUpdateService::class.java)
            enqueueWork(context, RecommendationUpdateService::class.java, JOB_ID, intent)
        }
    }

    override fun onHandleWork(intent: Intent) {
        updateRecommendations()
    }

    private fun updateRecommendations() {
        serviceScope.launch {
            try {
                // 1. Fetch recommended content from repository
                val repository = ContentRepository.create(applicationContext)
                val rows = repository.getHomeContentRows()
                
                // Get recommended content (first row is typically "Recommended for You")
                val recommended = rows.firstOrNull()?.items?.take(MAX_RECOMMENDATIONS) ?: emptyList()
                
                // 2. Build notifications using TvContract
                val notificationManager = getSystemService(Context.NOTIFICATION_SERVICE) as NotificationManager
                createNotificationChannel(notificationManager)
                
                recommended.forEachIndexed { index, content ->
                    val notification = buildRecommendationNotification(content, index)
                    notificationManager.notify(index, notification)
                }
                
                // 3. Publish to Android TV home screen via TvContract
                publishRecommendations(recommended)
            } catch (e: Exception) {
                // Handle error silently in background service
            }
        }
    }
    
    private fun createNotificationChannel(notificationManager: NotificationManager) {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel = NotificationChannel(
                CHANNEL_ID,
                "Recommendations",
                NotificationManager.IMPORTANCE_DEFAULT
            ).apply {
                description = "Content recommendations for Android TV"
            }
            notificationManager.createNotificationChannel(channel)
        }
    }
    
    private fun buildRecommendationNotification(
        content: Content,
        id: Int
    ): Notification {
        val intent = Intent(this, MainActivity::class.java).apply {
            putExtra("content_id", content.id)
        }
        val pendingIntent = PendingIntent.getActivity(
            this,
            id,
            intent,
            PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_IMMUTABLE
        )
        
        return NotificationCompat.Builder(this, CHANNEL_ID)
            .setContentTitle(content.title)
            .setContentText(content.description)
            .setContentIntent(pendingIntent)
            .setLargeIcon(BitmapFactory.decodeResource(resources, R.drawable.ic_launcher_banner))
            .setAutoCancel(true)
            .setCategory(NotificationCompat.CATEGORY_RECOMMENDATION)
            .build()
    }
    
    private fun publishRecommendations(contents: List<Content>) {
        // Use TvContract to publish recommendations to Android TV home screen
        // This requires additional setup and proper TV Provider implementation
        // For now, this is a placeholder - full implementation requires:
        // 1. ContentProvider extending TvProviderContract
        // 2. Proper content URIs and MIME types
        // 3. Watch next program setup
        // See: https://developer.android.com/training/tv/discovery/recommendations-content
        
        // TODO: Implement full TV Provider integration
        // This would involve:
        // - Creating a ContentProvider
        // - Implementing TvContract methods
        // - Publishing programs via ContentResolver
    }
}
