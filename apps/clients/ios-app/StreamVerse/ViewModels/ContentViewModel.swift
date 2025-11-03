//
//  ContentViewModel.swift
//  StreamVerse
//
//  Created on 2025-01-XX.
//

import Foundation
import Combine

/// ViewModel for content operations.
@MainActor
class ContentViewModel: ObservableObject {
    private let contentService = ContentService.shared
    
    @Published var contentRows: [ContentRow] = []
    @Published var isLoading: Bool = false
    @Published var errorMessage: String?
    
    /// Load home content.
    func loadHomeContent() async {
        isLoading = true
        errorMessage = nil
        
        do {
            contentRows = try await contentService.getHomeContent()
        } catch {
            errorMessage = error.localizedDescription
            contentRows = []
        }
        
        isLoading = false
    }
    
    /// Search content.
    func search(query: String) async -> [Content] {
        guard query.count >= 2 else { return [] }
        
        do {
            return try await contentService.searchContent(query: query)
        } catch {
            errorMessage = error.localizedDescription
            return []
        }
    }
    
    /// Get content by ID.
    func getContent(byId id: String) async -> Content? {
        do {
            return try await contentService.getContentById(id)
        } catch {
            errorMessage = error.localizedDescription
            return nil
        }
    }
}

