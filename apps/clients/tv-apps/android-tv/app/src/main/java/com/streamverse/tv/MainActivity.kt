package com.streamverse.tv

import android.content.Intent
import android.os.Bundle
import androidx.fragment.app.FragmentActivity
import com.streamverse.tv.data.repository.AuthRepository
import com.streamverse.tv.ui.auth.LoginActivity

/**
 * Main activity for Android TV app.
 * Uses Leanback BrowseFragment for TV-optimized UI.
 * Checks authentication and redirects to login if needed.
 */
class MainActivity : FragmentActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        
        // Check if user is logged in
        val authRepository = AuthRepository.create(this)
        if (!authRepository.isLoggedIn()) {
            // Redirect to login
            startActivity(Intent(this, LoginActivity::class.java))
            finish()
            return
        }

        setContentView(R.layout.activity_main)

        if (savedInstanceState == null) {
            supportFragmentManager.beginTransaction()
                .replace(R.id.main_browse_fragment, MainFragment())
                .commitNow()
        }
    }
}

