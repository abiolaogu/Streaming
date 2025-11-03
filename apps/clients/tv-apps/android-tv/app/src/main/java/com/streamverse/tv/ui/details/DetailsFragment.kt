package com.streamverse.tv.ui.details

import android.content.Intent
import android.os.Bundle
import android.view.View
import androidx.leanback.app.DetailsSupportFragment
import androidx.leanback.widget.*
import com.bumptech.glide.Glide
import com.bumptech.glide.request.target.SimpleTarget
import com.bumptech.glide.request.transition.Transition
import com.streamverse.tv.R
import com.streamverse.tv.data.model.Content
import com.streamverse.tv.ui.presenter.DetailsDescriptionPresenter
import com.streamverse.tv.ui.player.PlaybackVideoActivity

/**
 * Fragment showing content details with poster, description, actions (Play, Watchlist).
 */
class DetailsFragment : DetailsSupportFragment() {

    private var content: Content? = null

    override fun onActivityCreated(savedInstanceState: Bundle?) {
        super.onActivityCreated(savedInstanceState)

        content = arguments?.getParcelable(ARG_CONTENT)

        setupDetailsOverviewRow()
        setupActionsRow()
    }

    private fun setupDetailsOverviewRow() {
        content?.let { content ->
            val overviewRow = DetailsOverviewRow(content).apply {
                imageDrawable = null
                // Load poster image
                Glide.with(requireContext())
                    .asDrawable()
                    .load(content.posterUrl)
                    .into(object : SimpleTarget<android.graphics.drawable.Drawable>() {
                        override fun onResourceReady(
                            resource: android.graphics.drawable.Drawable,
                            transition: Transition<in android.graphics.drawable.Drawable>?
                        ) {
                            imageDrawable = resource
                        }
                    })
            }

            val presenter = FullWidthDetailsOverviewRowPresenter(DetailsDescriptionPresenter())
            val adapter = ArrayObjectAdapter(presenter).apply {
                add(overviewRow)
            }
            this.adapter = adapter
        }
    }

    private fun setupActionsRow() {
        val actionAdapter = ArrayObjectAdapter().apply {
            add(Action(1, getString(R.string.play_button), null))
            add(Action(2, getString(R.string.add_to_watchlist), null))
            add(Action(3, getString(R.string.action_share), null))
        }

        onActionClickedListener = OnActionClickedListener { action ->
            when (action.id) {
                1L -> {
                    // Play
                    content?.let {
                        val intent = Intent(requireContext(), PlaybackVideoActivity::class.java).apply {
                            putExtra(PlaybackVideoActivity.EXTRA_CONTENT, it)
                        }
                        startActivity(intent)
                    }
                }
                2L -> {
                    // Add to watchlist
                    // TODO: Implement watchlist functionality
                }
                3L -> {
                    // Share
                    // TODO: Implement share functionality
                }
            }
        }

            val actionsRowAdapter = ArrayObjectAdapter(ListRowPresenter()).apply {
            add(ListRow(HeaderItem(0, "Actions"), actionAdapter))
        }
        
        // Merge with existing adapter
        val currentAdapter = adapter as? ArrayObjectAdapter
        if (currentAdapter != null) {
            currentAdapter.addAll(currentAdapter.size(), actionsRowAdapter)
        } else {
            adapter = actionsRowAdapter
        }
    }

    companion object {
        private const val ARG_CONTENT = "content"

        fun newInstance(content: Content): DetailsFragment {
            return DetailsFragment().apply {
                arguments = Bundle().apply {
                    putParcelable(ARG_CONTENT, content)
                }
            }
        }
    }
}

