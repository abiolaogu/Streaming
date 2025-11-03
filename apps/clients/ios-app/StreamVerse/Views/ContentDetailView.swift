//
//  ContentDetailView.swift
//  StreamVerse
//
//  Created on 2025-01-XX.
//

import SwiftUI

struct ContentDetailView: View {
    let content: Content
    @State private var isPlaying = false
    
    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 20) {
                // Backdrop Image
                AsyncImage(url: URL(string: content.backdropUrl)) { image in
                    image
                        .resizable()
                        .aspectRatio(contentMode: .fill)
                } placeholder: {
                    Rectangle()
                        .fill(Color.gray.opacity(0.3))
                }
                .frame(height: 300)
                .clipped()
                
                // Content Info
                VStack(alignment: .leading, spacing: 12) {
                    Text(content.title)
                        .font(.largeTitle)
                        .fontWeight(.bold)
                    
                    HStack {
                        Text("\(content.releaseYear)")
                        Text("•")
                        Text("\(String(format: "%.1f", content.rating)) ⭐")
                        Text("•")
                        Text(content.genre)
                    }
                    .font(.subheadline)
                    .foregroundColor(.secondary)
                    
                    Text(content.description)
                        .font(.body)
                    
                    Button(action: {
                        isPlaying = true
                    }) {
                        HStack {
                            Image(systemName: "play.fill")
                            Text("Play")
                        }
                        .frame(maxWidth: .infinity)
                        .padding()
                        .background(Color.blue)
                        .foregroundColor(.white)
                        .cornerRadius(10)
                    }
                    .padding(.top)
                }
                .padding(.horizontal)
            }
        }
        .navigationBarTitleDisplayMode(.inline)
        .fullScreenCover(isPresented: $isPlaying) {
            VideoPlayerView(content: content)
        }
    }
}

#Preview {
    NavigationStack {
        ContentDetailView(content: Content(
            id: "1",
            title: "Sample Movie",
            description: "A sample movie description",
            genre: "Action",
            category: "movie",
            posterUrl: "",
            backdropUrl: "",
            streamUrl: "",
            duration: 7200000,
            releaseYear: 2024,
            rating: 8.5,
            isDrmProtected: false,
            drmType: nil,
            thumbnailUrl: nil,
            cast: [],
            directors: [],
            tags: []
        ))
    }
}

