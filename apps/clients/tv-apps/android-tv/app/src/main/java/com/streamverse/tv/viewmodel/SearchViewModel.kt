package com.streamverse.tv.viewmodel

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import com.streamverse.tv.data.model.Content
import com.streamverse.tv.data.repository.ContentRepository
import kotlinx.coroutines.launch

/**
 * ViewModel for search functionality.
 */
class SearchViewModel(
    private val contentRepository: ContentRepository
) : ViewModel() {

    private val _searchResults = MutableLiveData<List<Content>>()
    val searchResults: LiveData<List<Content>> = _searchResults

    private val _isLoading = MutableLiveData<Boolean>()
    val isLoading: LiveData<Boolean> = _isLoading

    private val _error = MutableLiveData<String?>()
    val error: LiveData<String?> = _error

    /**
     * Search for content.
     */
    fun search(query: String) {
        if (query.length < 2) {
            _searchResults.value = emptyList()
            return
        }

        viewModelScope.launch {
            _isLoading.value = true
            _error.value = null

            try {
                val results = contentRepository.searchContent(query)
                _searchResults.value = results
            } catch (e: Exception) {
                _error.value = e.message ?: "Search failed"
                _searchResults.value = emptyList()
            } finally {
                _isLoading.value = false
            }
        }
    }

    /**
     * Clear search results.
     */
    fun clearSearch() {
        _searchResults.value = emptyList()
        _error.value = null
    }

    class Factory(
        private val repository: ContentRepository
    ) : ViewModelProvider.Factory {
        @Suppress("UNCHECKED_CAST")
        override fun <T : ViewModel> create(modelClass: Class<T>): T {
            return SearchViewModel(repository) as T
        }
    }
}

