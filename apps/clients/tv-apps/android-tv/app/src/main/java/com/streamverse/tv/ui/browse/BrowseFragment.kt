package com.streamverse.tv.ui.browse

import android.os.Bundle
import androidx.leanback.app.BrowseSupportFragment
import androidx.leanback.widget.*
import com.streamverse.tv.R
import com.streamverse.tv.data.model.Content
import com.streamverse.tv.ui.presenter.ContentCardPresenter
import com.streamverse.tv.viewmodel.BrowseViewModel
import org.koin.androidx.viewmodel.ext.android.viewModel

/**
 * Browse fragment for specific categories (Movies, TV Shows, Live).
 * Can be used as alternative to MainFragment or as category-specific view.
 */
class BrowseFragment : BrowseSupportFragment() {

    private val viewModel: BrowseViewModel by viewModel()
    private lateinit var rowsAdapter: ArrayObjectAdapter

    companion object {
        const val ARG_CATEGORY = "category"
    }

    override fun onActivityCreated(savedInstanceState: Bundle?) {
        super.onActivityCreated(savedInstanceState)

        val category = arguments?.getString(ARG_CATEGORY) ?: "all"
        setupUI(category)
        loadContent(category)
    }

    private fun setupUI(category: String) {
        title = when (category) {
            "movies" -> getString(R.string.category_movies)
            "shows" -> getString(R.string.category_shows)
            "live" -> getString(R.string.category_live)
            else -> getString(R.string.browse_title)
        }

        headersState = HEADERS_ENABLED
        isHeadersTransitionOnBackEnabled = true
    }

    private fun loadContent(category: String) {
        rowsAdapter = ArrayObjectAdapter(ListRowPresenter())
        adapter = rowsAdapter

        viewModel.getContentByCategory(category).observe(viewLifecycleOwner) { contentList ->
            // Group by genre or create single row
            val genreMap = contentList.groupBy { it.genre }
            
            genreMap.forEach { (genre, items) ->
                val listRowAdapter = ArrayObjectAdapter(ContentCardPresenter())
                items.forEach { content ->
                    listRowAdapter.add(content)
                }
                rowsAdapter.add(
                    ListRow(HeaderItem(genre.hashCode().toLong(), genre), listRowAdapter)
                )
            }
            
            rowsAdapter.notifyArrayItemRangeChanged(0, rowsAdapter.size())
        }
    }
}

