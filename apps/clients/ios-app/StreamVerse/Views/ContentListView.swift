//
//  ContentListView.swift
//  StreamVerse
//
//  Created on 2025-01-XX.
//

import SwiftUI

struct ContentListView: View {
    @StateObject private var viewModel = ContentViewModel()
    @EnvironmentObject var authViewModel: AuthViewModel
    @State private var searchText: String = ""
    @State private var isSearching: Bool = false
    @State private var searchResults: [Content] = []
    
    var body: some View {
        NavigationStack {
            ZStack {
                if viewModel.isLoading && viewModel.contentRows.isEmpty {
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
                        VStack(alignment: .leading, spacing: 24) {
                            ForEach(viewModel.contentRows) { row in
                                ContentRowView(row: row)
                            }
                        }
                        .padding(.vertical)
                    }
                }
            }
            .navigationTitle("StreamVerse")
            .navigationBarTitleDisplayMode(.large)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Menu {
                        NavigationLink("Settings") {
                            SettingsView()
                                .environmentObject(authViewModel)
                        }
                        Button("Logout", role: .destructive) {
                            Task {
                                await authViewModel.logout()
                            }
                        }
                    } label: {
                        Image(systemName: "person.circle")
                    }
                }
            }
            .searchable(text: $searchText, isPresented: $isSearching) {
                if !searchResults.isEmpty {
                    ForEach(searchResults) { content in
                        NavigationLink(destination: ContentDetailView(content: content)) {
                            ContentCardView(content: content)
                        }
                    }
                }
            }
            .onSubmit(of: .search) {
                Task {
                    searchResults = await viewModel.search(query: searchText)
                }
            }
            .task {
                await viewModel.loadHomeContent()
            }
        }
    }
}

struct ContentRowView: View {
    let row: ContentRow
    
    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text(row.title)
                .font(.headline)
                .padding(.horizontal)
            
            ScrollView(.horizontal, showsIndicators: false) {
                HStack(spacing: 16) {
                    ForEach(row.items) { content in
                        NavigationLink(destination: ContentDetailView(content: content)) {
                            ContentCardView(content: content)
                        }
                    }
                }
                .padding(.horizontal)
            }
        }
        .padding(.vertical, 8)
    }
}

struct ContentCardView: View {
    let content: Content
    
    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            AsyncImage(url: URL(string: content.posterUrl)) { image in
                image
                    .resizable()
                    .aspectRatio(contentMode: .fill)
            } placeholder: {
                Rectangle()
                    .fill(Color.gray.opacity(0.3))
            }
            .frame(width: 150, height: 225)
            .cornerRadius(12)
            
            Text(content.title)
                .font(.caption)
                .lineLimit(2)
                .frame(width: 150, alignment: .leading)
        }
    }
}

#Preview {
    ContentListView()
        .environmentObject(AuthViewModel())
}

