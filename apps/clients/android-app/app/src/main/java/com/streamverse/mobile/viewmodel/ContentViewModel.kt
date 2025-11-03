package com.streamverse.mobile.viewmodel

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import com.streamverse.mobile.data.repository.ContentRepository
import com.streamverse.mobile.models.Content
import com.streamverse.mobile.models.ContentRow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

/**
 * ViewModel for content operations.
 */
class ContentViewModel(
    private val contentRepository: ContentRepository
) : ViewModel() {

    private val _contentRows = MutableStateFlow<List<ContentRow>>(emptyList())
    val contentRows: StateFlow<List<ContentRow>> = _contentRows.asStateFlow()

    private val _isLoading = MutableStateFlow(false)
    val isLoading: StateFlow<Boolean> = _isLoading.asStateFlow()

    private val _errorMessage = MutableStateFlow<String?>(null)
    val errorMessage: StateFlow<String?> = _errorMessage.asStateFlow()

    /**
     * Load home content.
     */
    fun loadHomeContent() {
        viewModelScope.launch {
            _isLoading.value = true
            _errorMessage.value = null

            try {
                _contentRows.value = contentRepository.getHomeContentRows()
            } catch (e: Exception) {
                _errorMessage.value = e.message ?: "Failed to load content"
                _contentRows.value = emptyList()
            } finally {
                _isLoading.value = false
            }
        }
    }

    /**
     * Search content.
     */
    fun search(query: String, onResult: (List<Content>) -> Unit) {
        if (query.length < 2) {
            onResult(emptyList())
            return
        }

        viewModelScope.launch {
            try {
                val results = contentRepository.searchContent(query)
                onResult(results)
            } catch (e: Exception) {
                _errorMessage.value = e.message ?: "Search failed"
                onResult(emptyList())
            }
        }
    }

    /**
     * Get content by ID.
     */
    fun getContentById(id: String, onResult: (Content?) -> Unit) {
        viewModelScope.launch {
            try {
                val content = contentRepository.getContentById(id)
                onResult(content)
            } catch (e: Exception) {
                _errorMessage.value = e.message ?: "Failed to load content"
                onResult(null)
            }
        }
    }

    class Factory(
        private val repository: ContentRepository
    ) : ViewModelProvider.Factory {
        @Suppress("UNCHECKED_CAST")
        override fun <T : ViewModel> create(modelClass: Class<T>): T {
            return ContentViewModel(repository) as T
        }
    }
}

