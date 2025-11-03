package com.streamverse.tv.ui.settings

import android.os.Bundle
import androidx.leanback.preference.LeanbackPreferenceFragmentCompat
import androidx.preference.Preference
import androidx.preference.PreferenceScreen
import com.streamverse.tv.R

/**
 * Settings fragment for app preferences.
 */
class SettingsFragment : LeanbackPreferenceFragmentCompat() {

    override fun onCreatePreferences(savedInstanceState: Bundle?, rootKey: String?) {
        setPreferencesFromResource(R.xml.preferences, rootKey)
    }

    override fun onPreferenceTreeClick(preference: Preference): Boolean {
        when (preference.key) {
            "account" -> {
                // Navigate to account settings
                return true
            }
            "quality" -> {
                // Show quality options
                return true
            }
            "about" -> {
                // Show about screen
                return true
            }
        }
        return super.onPreferenceTreeClick(preference)
    }
}

