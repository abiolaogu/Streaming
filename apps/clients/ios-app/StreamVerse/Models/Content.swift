//
//  Content.swift
//  StreamVerse
//
//  Created on 2025-01-XX.
//

import Foundation

/// Content model representing videos (movies, TV shows, live events).
struct Content: Codable, Identifiable, Hashable {
    let id: String
    let title: String
    let description: String
    let genre: String
    let category: String // "movie", "show", "live"
    let posterUrl: String
    let backdropUrl: String
    let streamUrl: String // HLS/DASH manifest URL
    let duration: Int64 // in milliseconds
    let releaseYear: Int
    let rating: Float
    let isDrmProtected: Bool
    let drmType: String? // "fairplay", "widevine", etc.
    let thumbnailUrl: String?
    let cast: [String]
    let directors: [String]
    let tags: [String]
    
    enum CodingKeys: String, CodingKey {
        case id
        case title
        case description
        case genre
        case category
        case posterUrl
        case backdropUrl
        case streamUrl
        case duration
        case releaseYear
        case rating
        case isDrmProtected
        case drmType
        case thumbnailUrl
        case cast
        case directors
        case tags
    }
}

/// Content row for home screen categories.
struct ContentRow: Codable, Identifiable, Hashable {
    let id: String
    let title: String
    let items: [Content]
}

