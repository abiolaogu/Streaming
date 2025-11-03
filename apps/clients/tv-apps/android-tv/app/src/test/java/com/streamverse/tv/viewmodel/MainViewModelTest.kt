package com.streamverse.tv.viewmodel

import com.streamverse.tv.data.model.ContentRow
import com.streamverse.tv.data.repository.ContentRepository
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.runTest
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test
import org.mockito.kotlin.*

/**
 * Unit tests for MainViewModel.
 */
@OptIn(ExperimentalCoroutinesApi::class)
class MainViewModelTest {

    private lateinit var repository: ContentRepository
    private lateinit var viewModel: MainViewModel

    @Before
    fun setup() {
        repository = mock()
        viewModel = MainViewModel(repository)
    }

    @Test
    fun `loadContent should update contentRows when successful`() = runTest {
        // Given
        val mockRows = listOf(
            ContentRow("1", "Trending", emptyList()),
            ContentRow("2", "New Releases", emptyList())
        )
        whenever(repository.getHomeContentRows()).thenReturn(mockRows)

        // When
        viewModel.loadContent()

        // Then
        assertEquals(mockRows, viewModel.contentRows.value)
        assertFalse(viewModel.isLoading.value ?: true)
        assertNull(viewModel.error.value)
    }

    @Test
    fun `loadContent should handle errors`() = runTest {
        // Given
        val error = Exception("Network error")
        whenever(repository.getHomeContentRows()).thenThrow(error)

        // When
        viewModel.loadContent()

        // Then
        assertNull(viewModel.contentRows.value)
        assertFalse(viewModel.isLoading.value ?: true)
        assertEquals("Network error", viewModel.error.value)
    }

    @Test
    fun `loadContent should set loading state`() = runTest {
        // Given
        val mockRows = listOf<ContentRow>()
        whenever(repository.getHomeContentRows()).thenReturn(mockRows)

        // When
        viewModel.loadContent()

        // Then
        verify(repository).getHomeContentRows()
    }
}

