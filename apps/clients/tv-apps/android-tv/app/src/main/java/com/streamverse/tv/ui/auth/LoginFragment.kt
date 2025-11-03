package com.streamverse.tv.ui.auth

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Button
import android.widget.EditText
import android.widget.TextView
import androidx.fragment.app.Fragment
import androidx.lifecycle.ViewModelProvider
import com.streamverse.tv.R
import com.streamverse.tv.data.repository.AuthRepository
import com.streamverse.tv.viewmodel.AuthViewModel

/**
 * Login fragment for Android TV.
 * Optimized for remote control navigation.
 */
class LoginFragment : Fragment() {

    private lateinit var viewModel: AuthViewModel
    private lateinit var emailEditText: EditText
    private lateinit var passwordEditText: EditText
    private lateinit var loginButton: Button
    private lateinit var errorTextView: TextView

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        val view = inflater.inflate(R.layout.fragment_login, container, false)
        
        emailEditText = view.findViewById(R.id.email_edit_text)
        passwordEditText = view.findViewById(R.id.password_edit_text)
        loginButton = view.findViewById(R.id.login_button)
        errorTextView = view.findViewById(R.id.error_text_view)

        return view
    }

    override fun onActivityCreated(savedInstanceState: Bundle?) {
        super.onActivityCreated(savedInstanceState)

        val repository = AuthRepository.create(requireContext())
        viewModel = ViewModelProvider(
            this,
            AuthViewModel.Factory(repository)
        )[AuthViewModel::class.java]

        setupObservers()
        setupListeners()
    }

    private fun setupObservers() {
        viewModel.isLoading.observe(viewLifecycleOwner) { isLoading ->
            loginButton.isEnabled = !isLoading
            loginButton.text = if (isLoading) {
                getString(R.string.logging_in)
            } else {
                getString(R.string.login_button)
            }
        }

        viewModel.error.observe(viewLifecycleOwner) { error ->
            errorTextView.text = error ?: ""
            errorTextView.visibility = if (error != null) View.VISIBLE else View.GONE
        }

        viewModel.loginSuccess.observe(viewLifecycleOwner) { success ->
            if (success) {
                // Navigate to main screen
                activity?.finish()
            }
        }
    }

    private fun setupListeners() {
        loginButton.setOnClickListener {
            val email = emailEditText.text.toString().trim()
            val password = passwordEditText.text.toString()

            if (validateInput(email, password)) {
                viewModel.login(email, password)
            }
        }

        // TV-optimized: Enter key submits
        passwordEditText.setOnEditorActionListener { _, _, _ ->
            loginButton.performClick()
            true
        }
    }

    private fun validateInput(email: String, password: String): Boolean {
        if (email.isEmpty()) {
            showError(getString(R.string.error_email_required))
            return false
        }
        if (!android.util.Patterns.EMAIL_ADDRESS.matcher(email).matches()) {
            showError(getString(R.string.error_invalid_email))
            return false
        }
        if (password.isEmpty()) {
            showError(getString(R.string.error_password_required))
            return false
        }
        return true
    }

    private fun showError(message: String) {
        errorTextView.text = message
        errorTextView.visibility = View.VISIBLE
    }
}

