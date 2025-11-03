//
//  ContentView.swift
//  StreamVerseTV
//
//  Created on 2025-01-XX.
//

import SwiftUI

struct ContentView: View {
    @StateObject private var viewModel = ContentViewModel()
    
    var body: some View {
        NavigationStack {
            ZStack {
                if viewModel.isLoading {
                    ProgressView("Loading content...")
                } else if let error = viewModel.errorMessage {
                    VStack {
                        Text("Error: \(error)")
                            .foregroundColor(.red)
                        Button("Retry") {
                            Task {
                                await viewModel.loadHomeContent()
                            }
                        }
                    }
                } else {
                    ScrollView {
                        VStack(alignment: .leading, spacing: 40) {
                            ForEach(viewModel.contentRows) { row in
                                ContentRowView(row: row)
                            }
                        }
                        .padding()
                    }
                }
            }
            .navigationTitle("StreamVerse")
            .task {
                await viewModel.loadHomeContent()
            }
        }
    }
}

struct ContentRowView: View {
    let row: ContentRow
    
    var body: some View {
        VStack(alignment: .leading, spacing: 20) {
            Text(row.title)
                .font(.title2)
                .foregroundColor(.white)
            
            ScrollView(.horizontal, showsIndicators: false) {
                HStack(spacing: 20) {
                    ForEach(row.items) { content in
                        NavigationLink(destination: ContentDetailView(content: content)) {
                            ContentCardView(content: content)
                        }
                        .buttonStyle(PlainButtonStyle())
                    }
                }
            }
        }
    }
}

struct ContentCardView: View {
    let content: Content
    
    var body: some View {
        VStack(alignment: .leading) {
            AsyncImage(url: URL(string: content.posterUrl)) { image in
                image
                    .resizable()
                    .aspectRatio(contentMode: .fill)
            } placeholder: {
                Rectangle()
                    .fill(Color.gray.opacity(0.3))
            }
            .frame(width: 300, height: 450)
            .cornerRadius(12)
            
            Text(content.title)
                .font(.headline)
                .foregroundColor(.white)
                .lineLimit(2)
                .frame(width: 300)
        }
    }
}

#Preview {
    ContentView()
        .background(Color.black)
}

