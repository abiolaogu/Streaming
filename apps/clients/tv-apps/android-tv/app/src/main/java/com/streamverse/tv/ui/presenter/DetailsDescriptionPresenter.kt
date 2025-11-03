package com.streamverse.tv.ui.presenter

import android.view.ViewGroup
import androidx.leanback.widget.AbstractDetailsDescriptionPresenter
import com.streamverse.tv.data.model.Content

/**
 * Presenter for content details description.
 */
class DetailsDescriptionPresenter : AbstractDetailsDescriptionPresenter() {

    override fun onBindDescription(
        viewHolder: AbstractDetailsDescriptionPresenter.ViewHolder,
        item: Any
    ) {
        val content = item as Content
        viewHolder.title.text = content.title
        viewHolder.subtitle.text = "${content.genre} • ${content.releaseYear}"
        viewHolder.body.text = content.description
    }
}

