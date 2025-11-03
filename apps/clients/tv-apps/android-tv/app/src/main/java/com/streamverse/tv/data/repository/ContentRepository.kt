package com.streamverse.tv.data.repository

import android.content.Context
import com.streamverse.tv.data.api.ContentApiService
import com.streamverse.tv.data.model.Content
import com.streamverse.tv.data.model.ContentRow
import okhttp3.OkHttpClient
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory

/**
 * Repository for content data.
 * Fetches from API and provides to ViewModels.
 */
class ContentRepository(
    private val apiService: ContentApiService
) {
    suspend fun getHomeContentRows(): List<ContentRow> {
        return apiService.getHomeContent()
    }

    suspend fun getContentByCategory(category: String): List<Content> {
        return apiService.getContentByCategory(category)
    }

    suspend fun getContentById(id: String): Content {
        return apiService.getContentById(id)
    }

    suspend fun searchContent(query: String): List<Content> {
        return apiService.searchContent(query)
    }

    companion object {
        /**
         * Factory method to create ContentRepository instance.
         * In production, use dependency injection (Koin/Dagger/Hilt).
         */
        fun create(context: Context): ContentRepository {
            val baseUrl = com.streamverse.tv.BuildConfig.API_BASE_URL
            
            val okHttpClient = OkHttpClient.Builder()
                .connectTimeout(30, java.util.concurrent.TimeUnit.SECONDS)
                .readTimeout(30, java.util.concurrent.TimeUnit.SECONDS)
                .writeTimeout(30, java.util.concurrent.TimeUnit.SECONDS)
                .addInterceptor { chain ->
                    val request = chain.request().newBuilder().apply {
                        // Add auth token if available
                        val authRepo = AuthRepository.create(context)
                        authRepo.getAccessToken()?.let { token ->
                            addHeader("Authorization", "Bearer $token")
                        }
                        addHeader("Accept", "application/json")
                        addHeader("Content-Type", "application/json")
                    }.build()
                    chain.proceed(request)
                }
                .build()

            val retrofit = Retrofit.Builder()
                .baseUrl(baseUrl)
                .client(okHttpClient)
                .addConverterFactory(GsonConverterFactory.create())
                .build()

            val apiService = retrofit.create(ContentApiService::class.java)
            return ContentRepository(apiService)
        }
    }
}

