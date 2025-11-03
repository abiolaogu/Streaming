package com.streamverse.tv.viewmodel

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import com.streamverse.tv.data.repository.AuthRepository
import kotlinx.coroutines.launch

/**
 * ViewModel for authentication operations.
 */
class AuthViewModel(
    private val authRepository: AuthRepository
) : ViewModel() {

    private val _isLoading = MutableLiveData<Boolean>()
    val isLoading: LiveData<Boolean> = _isLoading

    private val _error = MutableLiveData<String?>()
    val error: LiveData<String?> = _error

    private val _loginSuccess = MutableLiveData<Boolean>()
    val loginSuccess: LiveData<Boolean> = _loginSuccess

    /**
     * Login with email and password.
     */
    fun login(email: String, password: String) {
        viewModelScope.launch {
            _isLoading.value = true
            _error.value = null

            authRepository.login(email, password)
                .onSuccess {
                    _loginSuccess.value = true
                }
                .onFailure { exception ->
                    _error.value = exception.message ?: "Login failed"
                    _loginSuccess.value = false
                }
                .also {
                    _isLoading.value = false
                }
        }
    }

    /**
     * Check if user is logged in.
     */
    fun isLoggedIn(): Boolean {
        return authRepository.isLoggedIn()
    }

    /**
     * Get current user info.
     */
    fun getCurrentUser() = authRepository.getCurrentUser()

    /**
     * Logout current user.
     */
    fun logout() {
        viewModelScope.launch {
            authRepository.logout()
            _loginSuccess.value = false
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

