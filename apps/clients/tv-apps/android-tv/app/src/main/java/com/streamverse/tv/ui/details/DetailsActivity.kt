package com.streamverse.tv.ui.details

import android.app.Activity
import android.content.Intent
import android.os.Bundle
import androidx.fragment.app.FragmentActivity
import com.streamverse.tv.R
import com.streamverse.tv.data.model.Content
import com.streamverse.tv.ui.player.PlaybackVideoActivity

/**
 * Activity displaying content details (poster, description, cast, etc.).
 * Launches video playback when user selects play.
 */
class DetailsActivity : FragmentActivity() {

    private var content: Content? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_details)

        content = intent.getParcelableExtra(EXTRA_CONTENT)

        if (savedInstanceState == null) {
            content?.let {
                val fragment = DetailsFragment.newInstance(it)
                supportFragmentManager.beginTransaction()
                    .replace(R.id.details_fragment, fragment)
                    .commit()
            }
        }
    }

    companion object {
        const val EXTRA_CONTENT = "content"
        const val EXTRA_CONTENT_ID = "content_id"

        fun launch(activity: Activity, content: Content) {
            val intent = Intent(activity, DetailsActivity::class.java).apply {
                putExtra(EXTRA_CONTENT, content)
            }
            activity.startActivity(intent)
        }
    }
}

