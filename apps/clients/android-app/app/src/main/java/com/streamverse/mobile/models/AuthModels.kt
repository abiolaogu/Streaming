package com.streamverse.mobile.models

import com.google.gson.annotations.SerializedName

/**
 * Login request model.
 */
data class LoginRequest(
    val email: String,
    val password: String
)

/**
 * Authentication response model.
 */
data class AuthResponse(
    @SerializedName("token") val accessToken: String,
    @SerializedName("refresh_token") val refreshToken: String? = null,
    @SerializedName("user") val user: UserInfo? = null,
    @SerializedName("expires_at") val expiresAt: String? = null,
    @SerializedName("expires_in") val expiresIn: Long? = null
)

/**
 * User information model.
 */
data class UserInfo(
    val id: String,
    val email: String,
    val name: String? = null,
    val avatar: String? = null,
    val roles: List<String> = emptyList()
)

