package com.streamverse.tv.data.repository

import com.streamverse.tv.data.api.ContentApiService
import com.streamverse.tv.data.model.Content
import com.streamverse.tv.data.model.ContentRow
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.runTest
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test
import org.mockito.kotlin.*

/**
 * Unit tests for ContentRepository.
 */
@OptIn(ExperimentalCoroutinesApi::class)
class ContentRepositoryTest {

    private lateinit var apiService: ContentApiService
    private lateinit var repository: ContentRepository

    @Before
    fun setup() {
        apiService = mock()
        repository = ContentRepository(apiService)
    }

    @Test
    fun `getHomeContentRows should return content rows from API`() = runTest {
        // Given
        val expectedRows = listOf(
            ContentRow("1", "Trending", emptyList()),
            ContentRow("2", "New Releases", emptyList())
        )
        whenever(apiService.getHomeContent()).thenReturn(expectedRows)

        // When
        val result = repository.getHomeContentRows()

        // Then
        assertEquals(expectedRows, result)
        verify(apiService).getHomeContent()
    }

    @Test
    fun `searchContent should return search results`() = runTest {
        // Given
        val query = "action"
        val expectedResults = listOf(
            Content(
                id = "1",
                title = "Action Movie",
                description = "An action movie",
                genre = "Action",
                category = "movie",
                posterUrl = "",
                backdropUrl = "",
                streamUrl = "",
                duration = 60000,
                releaseYear = 2024,
                rating = 8.5f
            )
        )
        whenever(apiService.searchContent(query)).thenReturn(expectedResults)

        // When
        val result = repository.searchContent(query)

        // Then
        assertEquals(expectedResults, result)
        verify(apiService).searchContent(query)
    }

    @Test
    fun `getContentById should return content`() = runTest {
        // Given
        val contentId = "123"
        val expectedContent = Content(
            id = contentId,
            title = "Test Movie",
            description = "A test movie",
            genre = "Drama",
            category = "movie",
            posterUrl = "",
            backdropUrl = "",
            streamUrl = "",
            duration = 60000,
            releaseYear = 2024,
            rating = 7.5f
        )
        whenever(apiService.getContentById(contentId)).thenReturn(expectedContent)

        // When
        val result = repository.getContentById(contentId)

        // Then
        assertEquals(expectedContent, result)
        verify(apiService).getContentById(contentId)
    }
}

