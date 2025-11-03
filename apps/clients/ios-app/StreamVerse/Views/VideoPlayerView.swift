//
//  VideoPlayerView.swift
//  StreamVerse
//
//  Created on 2025-01-XX.
//

import SwiftUI
import AVKit
import AVFoundation

struct VideoPlayerView: View {
    let content: Content
    @Environment(\.dismiss) var dismiss
    @StateObject private var playerViewModel = VideoPlayerViewModel()
    
    var body: some View {
        ZStack {
            Color.black.ignoresSafeArea()
            
            if let player = playerViewModel.player {
                VideoPlayer(player: player)
                    .ignoresSafeArea()
                    .onAppear {
                        playerViewModel.loadContent(content)
                    }
            } else {
                ProgressView("Loading video...")
            }
        }
        .toolbar {
            ToolbarItem(placement: .navigationBarTrailing) {
                Button("Done") {
                    playerViewModel.stop()
                    dismiss()
                }
            }
        }
    }
}

/// ViewModel for video playback.
@MainActor
class VideoPlayerViewModel: ObservableObject {
    @Published var player: AVPlayer?
    
    private var drmHelper = DRMHelper.shared
    
    func loadContent(_ content: Content) {
        guard let url = URL(string: content.streamUrl) else {
            return
        }
        
        let asset = AVURLAsset(url: url)
        
        // Configure DRM if needed
        if content.isDrmProtected, let drmType = content.drmType {
            drmHelper.configureDRM(for: asset, content: content, drmType: drmType)
        }
        
        let playerItem = AVPlayerItem(asset: asset)
        player = AVPlayer(playerItem: playerItem)
        player?.play()
    }
    
    func stop() {
        player?.pause()
        player = nil
    }
}

#Preview {
    VideoPlayerView(content: Content(
        id: "1",
        title: "Sample",
        description: "Sample",
        genre: "Action",
        category: "movie",
        posterUrl: "",
        backdropUrl: "",
        streamUrl: "https://example.com/video.m3u8",
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

