package com.streamverse.tv.ui.player.drm

import android.content.Context
import android.content.SharedPreferences
import com.streamverse.tv.data.model.Content

/**
 * Helper class for DRM configuration, specifically Widevine.
 * Handles Widevine license server communication and token management.
 */
object DRMHelper {

    data class DRMConfig(
        val licenseServerUrl: String,
        val licenseRequestHeaders: Map<String, String>
    )

    private const val PREFS_NAME = "streamverse_prefs"
    private const val KEY_AUTH_TOKEN = "auth_token"
    
    private fun getDefaultLicenseServer(): String {
        return try {
            com.streamverse.tv.BuildConfig.DRM_LICENSE_SERVER
        } catch (e: Exception) {
            "https://drm.streamverse.com/v1/widevine/license"
        }
    }

    /**
     * Get Widevine DRM configuration for content.
     */
    fun getWidevineConfig(content: Content, context: Context? = null): DRMConfig {
        val token = context?.let { getAuthToken(it) } ?: ""
        
        return DRMConfig(
            licenseServerUrl = content.drmType?.let { 
                // Could support different DRM types
                when (it) {
                    "widevine" -> getDefaultLicenseServer()
                    else -> getDefaultLicenseServer()
                }
            } ?: getDefaultLicenseServer(),
            licenseRequestHeaders = buildMap {
                put("Content-Type", "application/json")
                if (token.isNotEmpty()) {
                    put("Authorization", "Bearer $token")
                }
                // Add content-specific headers if needed
                put("X-Content-ID", content.id)
            }
        )
    }

    /**
     * Get authentication token from SharedPreferences or auth service.
     */
    private fun getAuthToken(context: Context): String {
        val prefs: SharedPreferences = context.getSharedPreferences(
            PREFS_NAME,
            Context.MODE_PRIVATE
        )
        return prefs.getString(KEY_AUTH_TOKEN, "") ?: ""
    }

    /**
     * Save authentication token (called after login).
     */
    fun saveAuthToken(context: Context, token: String) {
        val prefs: SharedPreferences = context.getSharedPreferences(
            PREFS_NAME,
            Context.MODE_PRIVATE
        )
        prefs.edit().putString(KEY_AUTH_TOKEN, token).apply()
    }

    /**
     * Clear authentication token (called on logout).
     */
    fun clearAuthToken(context: Context) {
        val prefs: SharedPreferences = context.getSharedPreferences(
            PREFS_NAME,
            Context.MODE_PRIVATE
        )
        prefs.edit().remove(KEY_AUTH_TOKEN).apply()
    }
}

