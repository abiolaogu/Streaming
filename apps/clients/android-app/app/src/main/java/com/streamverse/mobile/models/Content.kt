package com.streamverse.mobile.models

import android.os.Parcelable
import kotlinx.parcelize.Parcelize

/**
 * Content model representing videos (movies, TV shows, live events).
 */
@Parcelize
data class Content(
    val id: String,
    val title: String,
    val description: String,
    val genre: String,
    val category: String, // "movie", "show", "live"
    val posterUrl: String,
    val backdropUrl: String,
    val streamUrl: String, // HLS/DASH manifest URL
    val duration: Long, // in milliseconds
    val releaseYear: Int,
    val rating: Float,
    val isDrmProtected: Boolean = false,
    val drmType: String? = null, // "widevine", "fairplay", etc.
    val thumbnailUrl: String? = null,
    val cast: List<String> = emptyList(),
    val directors: List<String> = emptyList(),
    val tags: List<String> = emptyList()
) : Parcelable

/**
 * Content row for home screen categories.
 */
data class ContentRow(
    val id: String,
    val title: String,
    val items: List<Content>
)

