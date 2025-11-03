//
//  AuthService.swift
//  StreamVerse
//
//  Created on 2025-01-XX.
//

import Foundation

/// Service for authentication operations.
class AuthService {
    static let shared = AuthService()
    private let apiService = APIService.shared
    
    private init() {}
    
    /// Login with email and password.
    func login(email: String, password: String) async throws -> AuthResponse {
        let request = try apiService.createRequest(
            endpoint: "api/v1/auth/login",
            method: "POST",
            body: try JSONEncoder().encode(LoginRequest(email: email, password: password))
        )
        
        let response: AuthResponse = try await apiService.performRequest(request, responseType: AuthResponse.self)
        
        // Save tokens
        TokenManager.shared.saveAuthData(response)
        
        return response
    }
    
    /// Refresh access token.
    func refreshToken() async throws -> AuthResponse {
        guard let refreshToken = TokenManager.shared.getRefreshToken() else {
            throw APIError.unauthorized
        }
        
        let request = try apiService.createRequest(
            endpoint: "api/v1/auth/refresh",
            method: "POST",
            body: refreshToken.data(using: .utf8)
        )
        
        let response: AuthResponse = try await apiService.performRequest(request, responseType: AuthResponse.self)
        
        // Update tokens
        TokenManager.shared.saveAuthData(response)
        
        return response
    }
    
    /// Logout current user.
    func logout() async throws {
        let request = try apiService.createRequest(
            endpoint: "api/v1/auth/logout",
            method: "POST"
        )
        
        do {
            // Send logout request (response doesn't matter)
            _ = try await apiService.session.data(for: request)
        } catch {
            // Continue with logout even if API call fails
            print("Logout API error: \(error.localizedDescription)")
        }
        
        // Clear tokens regardless of API response
        TokenManager.shared.clearAuthData()
    }
}

