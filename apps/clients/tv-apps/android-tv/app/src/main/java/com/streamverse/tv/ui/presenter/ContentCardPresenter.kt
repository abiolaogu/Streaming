package com.streamverse.tv.ui.presenter

import android.view.ViewGroup
import androidx.leanback.widget.ImageCardView
import androidx.leanback.widget.Presenter
import com.bumptech.glide.Glide
import com.streamverse.tv.R
import com.streamverse.tv.data.model.Content

/**
 * Presenter for content cards in rows.
 * Displays poster image with title overlay.
 */
class ContentCardPresenter : Presenter() {

    private val cardWidth = 400
    private val cardHeight = 600

    override fun onCreateViewHolder(parent: ViewGroup): ViewHolder {
        val cardView = ImageCardView(parent.context).apply {
            isFocusable = true
            isFocusableInTouchMode = true
            setMainImageDimensions(cardWidth, cardHeight)
        }
        return ViewHolder(cardView)
    }

    override fun onBindViewHolder(viewHolder: ViewHolder, item: Any) {
        val content = item as Content
        val cardView = viewHolder.view as ImageCardView

        cardView.titleText = content.title
        cardView.contentText = content.releaseYear.toString()

        // Load poster image
        Glide.with(viewHolder.view.context)
            .load(content.posterUrl)
            .centerCrop()
            .into(cardView.mainImageView)
    }

    override fun onUnbindViewHolder(viewHolder: ViewHolder) {
        val cardView = viewHolder.view as ImageCardView
        cardView.badgeImage = null
        cardView.mainImage = null
    }
}

