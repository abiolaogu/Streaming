package com.streamverse.mobile.ui

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.viewModels
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.ui.Modifier
import com.streamverse.mobile.data.repository.AuthRepository
import com.streamverse.mobile.ui.theme.StreamVerseTheme
import com.streamverse.mobile.viewmodel.AuthViewModel

class MainActivity : ComponentActivity() {
    
    private val authViewModel: AuthViewModel by viewModels {
        AuthViewModel.Factory(AuthRepository.create(this))
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            StreamVerseTheme {
                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = MaterialTheme.colorScheme.background
                ) {
                    StreamVerseApp(authViewModel = authViewModel)
                }
            }
        }
    }
}

