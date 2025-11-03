//
//  PlayerView.swift
//  StreamVerseTV
//
//  Created on 2025-01-XX.
//

import SwiftUI
import AVKit

struct PlayerView: View {
    let content: Content
    @Environment(\.dismiss) var dismiss
    @StateObject private var playerViewModel = PlayerViewModel()
    
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

@MainActor
class PlayerViewModel: ObservableObject {
    @Published var player: AVPlayer?
    
    func loadContent(_ content: Content) {
        guard let url = URL(string: content.streamUrl) else {
            return
        }
        
        let asset = AVURLAsset(url: url)
        
        // Configure FairPlay DRM if needed
        if content.isDrmProtected {
            configureDRM(asset: asset, content: content)
        }
        
        let playerItem = AVPlayerItem(asset: asset)
        player = AVPlayer(playerItem: playerItem)
        player?.play()
    }
    
    func configureDRM(asset: AVURLAsset, content: Content) {
        // FairPlay DRM configuration
        // Requires certificate and license server setup
    }
    
    func stop() {
        player?.pause()
        player = nil
    }
}

