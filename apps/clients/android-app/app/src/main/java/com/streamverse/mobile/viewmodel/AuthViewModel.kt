package com.streamverse.mobile.viewmodel

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import com.streamverse.mobile.data.repository.AuthRepository
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

/**
 * ViewModel for authentication operations.
 */
class AuthViewModel(
    private val authRepository: AuthRepository
) : ViewModel() {

    private val _isLoggedIn = MutableStateFlow(false)
    val isLoggedIn: StateFlow<Boolean> = _isLoggedIn.asStateFlow()

    private val _isLoading = MutableStateFlow(false)
    val isLoading: StateFlow<Boolean> = _isLoading.asStateFlow()

    private val _errorMessage = MutableStateFlow<String?>(null)
    val errorMessage: StateFlow<String?> = _errorMessage.asStateFlow()

    private val _currentUser = MutableStateFlow(authRepository.getCurrentUser())
    val currentUser = _currentUser.asStateFlow()

    init {
        checkAuthStatus()
    }

    /**
     * Check current authentication status.
     */
    fun checkAuthStatus() {
        _isLoggedIn.value = authRepository.isLoggedIn()
        _currentUser.value = authRepository.getCurrentUser()
    }

    /**
     * Login with email and password.
     */
    fun login(email: String, password: String) {
        viewModelScope.launch {
            _isLoading.value = true
            _errorMessage.value = null

            authRepository.login(email, password)
                .onSuccess { response ->
                    _isLoggedIn.value = true
                    _currentUser.value = response.user
                }
                .onFailure { exception ->
                    _errorMessage.value = exception.message ?: "Login failed"
                    _isLoggedIn.value = false
                }
                .also {
                    _isLoading.value = false
                }
        }
    }

    /**
     * Logout current user.
     */
    fun logout() {
        viewModelScope.launch {
            authRepository.logout()
            _isLoggedIn.value = false
            _currentUser.value = null
        }
    }

    class Factory(
        private val repository: AuthRepository
    ) : ViewModelProvider.Factory {
        @Suppress("UNCHECKED_CAST")
        override fun <T : ViewModel> create(modelClass: Class<T>): T {
            return AuthViewModel(repository) as T
        }
    }
}

