package com.streamverse.tv.data.api

import com.streamverse.tv.data.model.AuthResponse
import com.streamverse.tv.data.model.LoginRequest
import retrofit2.http.Body
import retrofit2.http.POST

/**
 * API service interface for authentication endpoints.
 */
interface AuthApiService {

    @POST("api/v1/auth/login")
    suspend fun login(@Body request: LoginRequest): AuthResponse

    @POST("api/v1/auth/refresh")
    suspend fun refreshToken(@Body refreshToken: String): AuthResponse

    @POST("api/v1/auth/logout")
    suspend fun logout()
}

