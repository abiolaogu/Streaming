//
//  AuthModels.swift
//  StreamVerse
//
//  Created on 2025-01-XX.
//

import Foundation

/// Login request model.
struct LoginRequest: Codable {
    let email: String
    let password: String
}

/// Authentication response model.
struct AuthResponse: Codable {
    let token: String
    let refreshToken: String?
    let user: UserInfo?
    let expiresAt: String?
    let expiresIn: Int64?
    
    enum CodingKeys: String, CodingKey {
        case token
        case refreshToken = "refresh_token"
        case user
        case expiresAt = "expires_at"
        case expiresIn = "expires_in"
    }
}

/// User information model.
struct UserInfo: Codable, Identifiable {
    let id: String
    let email: String
    let name: String?
    let avatar: String?
    let roles: [String]
}

