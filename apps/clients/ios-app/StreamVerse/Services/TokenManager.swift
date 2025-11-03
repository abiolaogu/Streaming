//
//  TokenManager.swift
//  StreamVerse
//
//  Created on 2025-01-XX.
//

import Foundation
import Security

/// Manager for authentication tokens using Keychain for secure storage.
class TokenManager {
    static let shared = TokenManager()
    
    private let accessTokenKey = "com.streamverse.accessToken"
    private let refreshTokenKey = "com.streamverse.refreshToken"
    private let userIdKey = "com.streamverse.userId"
    private let userEmailKey = "com.streamverse.userEmail"
    private let userNameKey = "com.streamverse.userName"
    private let tokenExpiresAtKey = "com.streamverse.tokenExpiresAt"
    
    private init() {}
    
    /// Save authentication data.
    func saveAuthData(_ response: AuthResponse) {
        saveToKeychain(key: accessTokenKey, value: response.token)
        
        if let refreshToken = response.refreshToken {
            saveToKeychain(key: refreshTokenKey, value: refreshToken)
        }
        
        if let user = response.user {
            UserDefaults.standard.set(user.id, forKey: userIdKey)
            UserDefaults.standard.set(user.email, forKey: userEmailKey)
            if let name = user.name {
                UserDefaults.standard.set(name, forKey: userNameKey)
            }
        }
        
        if let expiresIn = response.expiresIn {
            let expiresAt = Date().addingTimeInterval(TimeInterval(expiresIn))
            UserDefaults.standard.set(expiresAt, forKey: tokenExpiresAtKey)
        }
    }
    
    /// Get access token.
    func getAccessToken() -> String? {
        return getFromKeychain(key: accessTokenKey)
    }
    
    /// Get refresh token.
    func getRefreshToken() -> String? {
        return getFromKeychain(key: refreshTokenKey)
    }
    
    /// Check if user is logged in.
    func isLoggedIn() -> Bool {
        guard let token = getAccessToken() else { return false }
        
        // Check if token is expired
        if let expiresAt = UserDefaults.standard.object(forKey: tokenExpiresAtKey) as? Date {
            return expiresAt > Date()
        }
        
        return !token.isEmpty
    }
    
    /// Get current user info.
    func getCurrentUser() -> UserInfo? {
        guard let id = UserDefaults.standard.string(forKey: userIdKey),
              let email = UserDefaults.standard.string(forKey: userEmailKey) else {
            return nil
        }
        
        let name = UserDefaults.standard.string(forKey: userNameKey)
        return UserInfo(id: id, email: email, name: name, avatar: nil, roles: [])
    }
    
    /// Clear all authentication data.
    func clearAuthData() {
        deleteFromKeychain(key: accessTokenKey)
        deleteFromKeychain(key: refreshTokenKey)
        UserDefaults.standard.removeObject(forKey: userIdKey)
        UserDefaults.standard.removeObject(forKey: userEmailKey)
        UserDefaults.standard.removeObject(forKey: userNameKey)
        UserDefaults.standard.removeObject(forKey: tokenExpiresAtKey)
    }
    
    // MARK: - Keychain Helpers
    
    private func saveToKeychain(key: String, value: String) {
        let data = value.data(using: .utf8)!
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: key,
            kSecValueData as String: data
        ]
        
        // Delete existing item first
        SecItemDelete(query as CFDictionary)
        
        // Add new item
        SecItemAdd(query as CFDictionary, nil)
    }
    
    private func getFromKeychain(key: String) -> String? {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: key,
            kSecReturnData as String: true
        ]
        
        var result: AnyObject?
        let status = SecItemCopyMatching(query as CFDictionary, &result)
        
        guard status == errSecSuccess,
              let data = result as? Data,
              let value = String(data: data, encoding: .utf8) else {
            return nil
        }
        
        return value
    }
    
    private func deleteFromKeychain(key: String) {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: key
        ]
        
        SecItemDelete(query as CFDictionary)
    }
}

