package com.streamverse.tv

import android.content.Intent
import android.os.Bundle
import android.view.View
import androidx.leanback.app.BrowseSupportFragment
import androidx.leanback.widget.*
import androidx.lifecycle.ViewModelProvider
import com.streamverse.tv.data.model.Content
import com.streamverse.tv.data.repository.ContentRepository
import com.streamverse.tv.ui.details.DetailsActivity
import com.streamverse.tv.ui.search.SearchActivity
import com.streamverse.tv.ui.presenter.ContentCardPresenter
import com.streamverse.tv.ui.recommendations.RecommendationUpdateService
import com.streamverse.tv.viewmodel.MainViewModel

/**
 * Main browse fragment showing content rows (trending, new releases, genres).
 * TV-optimized with Leanback components.
 */
class MainFragment : BrowseSupportFragment(), OnItemViewSelectedListener,
    OnItemViewClickedListener {

    private lateinit var viewModel: MainViewModel
    private lateinit var rowsAdapter: ArrayObjectAdapter

    override fun onActivityCreated(savedInstanceState: Bundle?) {
        super.onActivityCreated(savedInstanceState)

        // Initialize ViewModel
        val repository = ContentRepository.create(requireContext())
        viewModel = ViewModelProvider(
            this,
            MainViewModel.Factory(repository)
        )[MainViewModel::class.java]

        setupUIElements()
        setupEventListeners()
        loadContentRows()

        // Start recommendation service
        activity?.let {
            RecommendationUpdateService.startService(it)
        }
    }

    private fun setupUIElements() {
        title = getString(R.string.app_name)
        headersState = HEADERS_ENABLED
        isHeadersTransitionOnBackEnabled = true

        brandColor = resources.getColor(R.color.brand_color, null)
        searchAffordanceColor = resources.getColor(R.color.search_color, null)
        
        // Enable search
        setOnSearchClickedListener { startActivity(Intent(activity, SearchActivity::class.java)) }
    }

    private fun setupEventListeners() {
        onItemViewSelectedListener = this
        onItemViewClickedListener = this
    }

    private fun loadContentRows() {
        rowsAdapter = ArrayObjectAdapter(ListRowPresenter())
        adapter = rowsAdapter

        // Observe content categories from ViewModel
        viewModel.contentRows.observe(viewLifecycleOwner) { rows ->
            rowsAdapter.clear()
            rows.forEach { row ->
                val listRowAdapter = ArrayObjectAdapter(ContentCardPresenter())
                row.items.forEach { content ->
                    listRowAdapter.add(content)
                }
                rowsAdapter.add(ListRow(HeaderItem(row.id.toLong(), row.title), listRowAdapter))
            }
        }

        // Observe loading state
        viewModel.isLoading.observe(viewLifecycleOwner) { isLoading ->
            if (isLoading) {
                progressBarManager.show()
            } else {
                progressBarManager.hide()
            }
        }

        // Observe errors
        viewModel.error.observe(viewLifecycleOwner) { error ->
            error?.let {
                // Show error message
                // Could use a Toast or error fragment
            }
        }

        // Load data
        viewModel.loadContent()
    }

    override fun onItemSelected(
        itemViewHolder: Presenter.ViewHolder?,
        item: Any?,
        rowViewHolder: RowPresenter.ViewHolder?,
        row: Row?
    ) {
        // Handle item selection (focus) - can show preview, update description, etc.
    }

    override fun onItemClicked(
        itemViewHolder: Presenter.ViewHolder?,
        item: Any?,
        rowViewHolder: RowPresenter.ViewHolder?,
        row: Row?
    ) {
        if (item is Content) {
            val intent = Intent(requireContext(), DetailsActivity::class.java).apply {
                putExtra(DetailsActivity.EXTRA_CONTENT, item)
            }
            startActivity(intent)
        }
    }
}

