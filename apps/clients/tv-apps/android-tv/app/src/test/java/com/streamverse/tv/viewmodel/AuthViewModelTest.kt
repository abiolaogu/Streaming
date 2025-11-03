package com.streamverse.tv.viewmodel

import com.streamverse.tv.data.model.AuthResponse
import com.streamverse.tv.data.model.UserInfo
import com.streamverse.tv.data.repository.AuthRepository
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.runTest
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test
import org.mockito.kotlin.*

/**
 * Unit tests for AuthViewModel.
 */
@OptIn(ExperimentalCoroutinesApi::class)
class AuthViewModelTest {

    private lateinit var repository: AuthRepository
    private lateinit var viewModel: AuthViewModel

    @Before
    fun setup() {
        repository = mock()
        viewModel = AuthViewModel(repository)
    }

    @Test
    fun `login should succeed with valid credentials`() = runTest {
        // Given
        val email = "test@example.com"
        val password = "password123"
        val authResponse = AuthResponse(
            accessToken = "token123",
            refreshToken = "refresh123",
            user = UserInfo("user1", email, "Test User")
        )
        whenever(repository.login(email, password)).thenReturn(Result.success(authResponse))

        // When
        viewModel.login(email, password)

        // Then
        assertTrue(viewModel.loginSuccess.value ?: false)
        assertNull(viewModel.error.value)
    }

    @Test
    fun `login should fail with invalid credentials`() = runTest {
        // Given
        val email = "test@example.com"
        val password = "wrong"
        val error = Exception("Invalid credentials")
        whenever(repository.login(email, password)).thenReturn(Result.failure(error))

        // When
        viewModel.login(email, password)

        // Then
        assertFalse(viewModel.loginSuccess.value ?: true)
        assertEquals("Invalid credentials", viewModel.error.value)
    }

    @Test
    fun `isLoggedIn should return repository result`() {
        // Given
        whenever(repository.isLoggedIn()).thenReturn(true)

        // When
        val result = viewModel.isLoggedIn()

        // Then
        assertTrue(result)
    }
}

