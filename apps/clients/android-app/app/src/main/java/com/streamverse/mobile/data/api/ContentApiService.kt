package com.streamverse.mobile.data.api

import com.streamverse.mobile.models.Content
import com.streamverse.mobile.models.ContentRow
import retrofit2.http.GET
import retrofit2.http.Path
import retrofit2.http.Query

/**
 * API service interface for content endpoints.
 */
interface ContentApiService {

    @GET("api/v1/content/home")
    suspend fun getHomeContent(): List<ContentRow>

    @GET("api/v1/content/category/{category}")
    suspend fun getContentByCategory(@Path("category") category: String): List<Content>

    @GET("api/v1/content/{id}")
    suspend fun getContentById(@Path("id") id: String): Content

    @GET("api/v1/content/search")
    suspend fun searchContent(@Query("q") query: String): List<Content>
}

