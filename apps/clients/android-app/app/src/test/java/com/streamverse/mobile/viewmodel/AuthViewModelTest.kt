package com.streamverse.mobile.viewmodel

import com.streamverse.mobile.data.repository.AuthRepository
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.runTest
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test
import org.mockito.kotlin.*

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
    fun `checkAuthStatus should update logged in state`() {
        // Given
        whenever(repository.isLoggedIn()).thenReturn(true)

        // When
        viewModel.checkAuthStatus()

        // Then
        assertTrue(viewModel.isLoggedIn.value)
    }
}

