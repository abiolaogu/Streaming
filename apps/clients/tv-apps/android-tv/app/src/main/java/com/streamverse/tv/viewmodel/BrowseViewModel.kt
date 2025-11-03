package com.streamverse.tv.viewmodel

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.streamverse.tv.data.model.Content
import com.streamverse.tv.data.repository.ContentRepository
import kotlinx.coroutines.launch

class BrowseViewModel(
    private val contentRepository: ContentRepository
) : ViewModel() {

    private val _content = MutableLiveData<List<Content>>()
    val content: LiveData<List<Content>> = _content

    fun getContentByCategory(category: String): LiveData<List<Content>> {
        viewModelScope.launch {
            try {
                val items = contentRepository.getContentByCategory(category)
                _content.value = items
            } catch (e: Exception) {
                _content.value = emptyList()
            }
        }
        return _content
    }
}

