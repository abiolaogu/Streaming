//
//  AuthViewModel.swift
//  StreamVerse
//
//  Created on 2025-01-XX.
//

import Foundation
import Combine

/// ViewModel for authentication operations.
@MainActor
class AuthViewModel: ObservableObject {
    private let authService = AuthService.shared
    
    @Published var isLoggedIn: Bool = false
    @Published var isLoading: Bool = false
    @Published var errorMessage: String?
    @Published var currentUser: UserInfo?
    
    init() {
        checkAuthStatus()
    }
    
    /// Check current authentication status.
    func checkAuthStatus() {
        isLoggedIn = TokenManager.shared.isLoggedIn()
        currentUser = TokenManager.shared.getCurrentUser()
    }
    
    /// Login with email and password.
    func login(email: String, password: String) async {
        isLoading = true
        errorMessage = nil
        
        do {
            let response = try await authService.login(email: email, password: password)
            isLoggedIn = true
            currentUser = response.user
        } catch {
            errorMessage = error.localizedDescription
            isLoggedIn = false
        }
        
        isLoading = false
    }
    
    /// Logout current user.
    func logout() async {
        isLoading = true
        
        do {
            try await authService.logout()
        } catch {
            // Continue with logout even if API call fails
            print("Logout error: \(error.localizedDescription)")
        }
        
        isLoggedIn = false
        currentUser = nil
        isLoading = false
    }
}

