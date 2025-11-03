package com.streamverse.mobile.data.repository

import android.content.Context
import android.content.SharedPreferences
import com.streamverse.mobile.BuildConfig
import com.streamverse.mobile.data.api.AuthApiService
import com.streamverse.mobile.models.AuthResponse
import com.streamverse.mobile.models.LoginRequest
import com.streamverse.mobile.models.UserInfo
import okhttp3.OkHttpClient
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import java.util.concurrent.TimeUnit

/**
 * Repository for authentication data.
 * Handles login, logout, token management, and user session.
 */
class AuthRepository(
    private val authApiService: AuthApiService,
    private val context: Context
) {
    
    private val prefs: SharedPreferences = context.getSharedPreferences(
        PREFS_NAME,
        Context.MODE_PRIVATE
    )

    companion object {
        private const val PREFS_NAME = "streamverse_auth_prefs"
        private const val KEY_ACCESS_TOKEN = "access_token"
        private const val KEY_REFRESH_TOKEN = "refresh_token"
        private const val KEY_USER_ID = "user_id"
        private const val KEY_USER_EMAIL = "user_email"
        private const val KEY_USER_NAME = "user_name"
        private const val KEY_TOKEN_EXPIRES_AT = "token_expires_at"

        /**
         * Factory method to create AuthRepository instance.
         */
        fun create(context: Context): AuthRepository {
            val baseUrl = BuildConfig.API_BASE_URL
            
            val okHttpClient = OkHttpClient.Builder()
                .connectTimeout(30, TimeUnit.SECONDS)
                .readTimeout(30, TimeUnit.SECONDS)
                .writeTimeout(30, TimeUnit.SECONDS)
                .build()

            val retrofit = Retrofit.Builder()
                .baseUrl(baseUrl)
                .client(okHttpClient)
                .addConverterFactory(GsonConverterFactory.create())
                .build()

            val authApiService = retrofit.create(AuthApiService::class.java)
            return AuthRepository(authApiService, context)
        }
    }

    /**
     * Login with email and password.
     */
    suspend fun login(email: String, password: String): Result<AuthResponse> {
        return try {
            val request = LoginRequest(email, password)
            val response = authApiService.login(request)
            saveAuthData(response)
            Result.success(response)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    /**
     * Logout current user.
     */
    suspend fun logout() {
        try {
            authApiService.logout()
        } catch (e: Exception) {
            // Log error but continue with local logout
        } finally {
            clearAuthData()
        }
    }

    /**
     * Refresh access token.
     */
    suspend fun refreshToken(): Result<AuthResponse> {
        return try {
            val refreshToken = getRefreshToken()
            if (refreshToken == null) {
                return Result.failure(Exception("No refresh token available"))
            }
            val response = authApiService.refreshToken(refreshToken)
            saveAuthData(response)
            Result.success(response)
        } catch (e: Exception) {
            clearAuthData()
            Result.failure(e)
        }
    }

    /**
     * Get current access token.
     */
    fun getAccessToken(): String? {
        return prefs.getString(KEY_ACCESS_TOKEN, null)
    }

    /**
     * Get refresh token.
     */
    fun getRefreshToken(): String? {
        return prefs.getString(KEY_REFRESH_TOKEN, null)
    }

    /**
     * Check if user is logged in.
     */
    fun isLoggedIn(): Boolean {
        val token = getAccessToken()
        if (token == null) return false
        
        // Check if token is expired
        val expiresAt = prefs.getLong(KEY_TOKEN_EXPIRES_AT, 0)
        return expiresAt > System.currentTimeMillis()
    }

    /**
     * Get current user info.
     */
    fun getCurrentUser(): UserInfo? {
        val id = prefs.getString(KEY_USER_ID, null) ?: return null
        val email = prefs.getString(KEY_USER_EMAIL, null) ?: return null
        val name = prefs.getString(KEY_USER_NAME, null)
        
        return UserInfo(
            id = id,
            email = email,
            name = name
        )
    }

    /**
     * Save authentication data to SharedPreferences.
     */
    private fun saveAuthData(response: AuthResponse) {
        prefs.edit().apply {
            putString(KEY_ACCESS_TOKEN, response.accessToken)
            response.refreshToken?.let { putString(KEY_REFRESH_TOKEN, it) }
            response.user?.let { user ->
                putString(KEY_USER_ID, user.id)
                putString(KEY_USER_EMAIL, user.email)
                user.name?.let { putString(KEY_USER_NAME, it) }
            }
            response.expiresIn?.let { expiresIn ->
                putLong(KEY_TOKEN_EXPIRES_AT, System.currentTimeMillis() + expiresIn * 1000)
            }
            apply()
        }
    }

    /**
     * Clear authentication data.
     */
    private fun clearAuthData() {
        prefs.edit().clear().apply()
    }
}

