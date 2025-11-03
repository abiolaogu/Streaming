package com.streamverse.tv.ui.auth

import android.os.Bundle
import androidx.fragment.app.FragmentActivity
import com.streamverse.tv.R

/**
 * Login activity for authentication.
 * Shows login screen if user is not authenticated.
 */
class LoginActivity : FragmentActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_login)

        if (savedInstanceState == null) {
            supportFragmentManager.beginTransaction()
                .replace(R.id.login_fragment, LoginFragment())
                .commitNow()
        }
    }
}

