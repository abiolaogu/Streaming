package com.streamverse.mobile.ui

import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import com.streamverse.mobile.ui.login.LoginScreen
import com.streamverse.mobile.ui.content.ContentListScreen
import com.streamverse.mobile.viewmodel.AuthViewModel

@Composable
fun StreamVerseApp(authViewModel: AuthViewModel) {
    val navController = rememberNavController()
    val isLoggedIn by authViewModel.isLoggedIn.collectAsState()

    NavHost(
        navController = navController,
        startDestination = if (isLoggedIn) "content_list" else "login"
    ) {
        composable("login") {
            LoginScreen(
                authViewModel = authViewModel,
                onLoginSuccess = {
                    navController.navigate("content_list") {
                        popUpTo("login") { inclusive = true }
                    }
                }
            )
        }

        composable("content_list") {
            ContentListScreen(
                authViewModel = authViewModel,
                onLogout = {
                    navController.navigate("login") {
                        popUpTo("content_list") { inclusive = true }
                    }
                }
            )
        }
    }
}

