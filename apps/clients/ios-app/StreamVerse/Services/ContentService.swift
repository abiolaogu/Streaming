//
//  ContentService.swift
//  StreamVerse
//
//  Created on 2025-01-XX.
//

import Foundation

/// Service for content-related API operations.
class ContentService {
    static let shared = ContentService()
    private let apiService = APIService.shared
    
    private init() {}
    
    /// Get home content rows.
    func getHomeContent() async throws -> [ContentRow] {
        let request = try apiService.createRequest(endpoint: "api/v1/content/home")
        return try await apiService.performRequest(request, responseType: [ContentRow].self)
    }
    
    /// Get content by category.
    func getContentByCategory(_ category: String) async throws -> [Content] {
        let request = try apiService.createRequest(endpoint: "api/v1/content/category/\(category)")
        return try await apiService.performRequest(request, responseType: [Content].self)
    }
    
    /// Get content by ID.
    func getContentById(_ id: String) async throws -> Content {
        let request = try apiService.createRequest(endpoint: "api/v1/content/\(id)")
        return try await apiService.performRequest(request, responseType: Content.self)
    }
    
    /// Search content.
    func searchContent(query: String) async throws -> [Content] {
        var components = URLComponents(string: apiService.baseURL + "api/v1/content/search")!
        components.queryItems = [URLQueryItem(name: "q", value: query)]
        
        guard let url = components.url else {
            throw APIError.invalidURL
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        
        if let token = TokenManager.shared.getAccessToken() {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }
        
        return try await apiService.performRequest(request, responseType: [Content].self)
    }
}

