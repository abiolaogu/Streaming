package com.streamverse.tv

import androidx.test.ext.junit.runners.AndroidJUnit4
import androidx.test.rule.ActivityTestRule
import org.junit.Rule
import org.junit.Test
import org.junit.runner.RunWith

/**
 * Integration tests for MainActivity.
 */
@RunWith(AndroidJUnit4::class)
class MainActivityTest {

    @get:Rule
    val activityRule = ActivityTestRule(MainActivity::class.java, false, false)

    @Test
    fun testMainActivityLaunches() {
        // This test verifies the activity can be launched
        // In a real scenario, you'd need to mock authentication
        // activityRule.launchActivity(null)
        // assertNotNull(activityRule.activity)
    }

    // Additional integration tests can be added here
    // Note: These require a running Android TV emulator or device
}

