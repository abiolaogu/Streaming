package com.streamverse.tv.ui.search

import android.app.SearchManager
import android.content.Context
import android.os.Bundle
import android.speech.SpeechRecognizer
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.lifecycleScope
import androidx.leanback.app.SearchSupportFragment
import androidx.leanback.widget.*
import com.streamverse.tv.R
import com.streamverse.tv.data.model.Content
import com.streamverse.tv.data.repository.ContentRepository
import com.streamverse.tv.ui.details.DetailsActivity
import com.streamverse.tv.ui.presenter.ContentCardPresenter
import com.streamverse.tv.viewmodel.SearchViewModel
import kotlinx.coroutines.launch

/**
 * Search fragment with voice input support.
 */
class SearchFragment : SearchSupportFragment(), SearchSupportFragment.SearchResultProvider,
    OnItemViewClickedListener {

    private val rowsAdapter = ArrayObjectAdapter(ListRowPresenter())
    private val resultsAdapter = ArrayObjectAdapter(ContentCardPresenter())
    private lateinit var viewModel: SearchViewModel

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setSearchResultProvider(this)
        setOnItemViewClickedListener(this)
        
        // Initialize ViewModel
        val repository = ContentRepository.create(requireContext())
        viewModel = ViewModelProvider(
            this,
            SearchViewModel.Factory(repository)
        )[SearchViewModel::class.java]
        
        setupObservers()
        
        // Enable voice search
        if (SpeechRecognizer.isRecognitionAvailable(context)) {
            setSpeechRecognitionCallback(object : SearchSupportFragment.SpeechRecognitionCallback {
                override fun recognizeSpeech() {
                    // Voice search is enabled
                }
            })
        }
    }
    
    private fun setupObservers() {
        viewModel.searchResults.observe(viewLifecycleOwner) { results ->
            resultsAdapter.clear()
            results.forEach { content ->
                resultsAdapter.add(content)
            }
            rowsAdapter.clear()
            if (results.isNotEmpty()) {
                rowsAdapter.add(ListRow(HeaderItem(0, getString(R.string.search_results)), resultsAdapter))
            }
        }
        
        viewModel.isLoading.observe(viewLifecycleOwner) { isLoading ->
            // Could show loading indicator
        }
        
        viewModel.error.observe(viewLifecycleOwner) { error ->
            // Could show error message
            if (error != null) {
                resultsAdapter.clear()
                rowsAdapter.clear()
            }
        }
    }

    override fun getResultsAdapter(): ObjectAdapter {
        return rowsAdapter
    }

    fun performSearch(query: String) {
        if (query.isEmpty()) {
            resultsAdapter.clear()
            rowsAdapter.clear()
            viewModel.clearSearch()
            return
        }
        
        viewModel.search(query)
    }

    override fun onQueryTextChange(newQuery: String): Boolean {
        performSearch(newQuery)
        return true
    }

    override fun onQueryTextSubmit(query: String): Boolean {
        performSearch(query)
        return true
    }

    override fun onItemClicked(
        itemViewHolder: Presenter.ViewHolder?,
        item: Any?,
        rowViewHolder: RowPresenter.ViewHolder?,
        row: Row?
    ) {
        if (item is Content) {
            DetailsActivity.launch(requireActivity(), item)
        }
    }
}

