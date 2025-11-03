//
//  SettingsView.swift
//  StreamVerse
//
//  Created on 2025-01-XX.
//

import SwiftUI

struct SettingsView: View {
    @EnvironmentObject var authViewModel: AuthViewModel
    @AppStorage("videoQuality") private var videoQuality = "auto"
    @AppStorage("subtitlesEnabled") private var subtitlesEnabled = true
    @AppStorage("autoplayEnabled") private var autoplayEnabled = false
    
    var body: some View {
        NavigationStack {
            Form {
                Section("Account") {
                    if let user = authViewModel.currentUser {
                        HStack {
                            Text("Email")
                            Spacer()
                            Text(user.email)
                                .foregroundColor(.secondary)
                        }
                        
                        if let name = user.name {
                            HStack {
                                Text("Name")
                                Spacer()
                                Text(name)
                                    .foregroundColor(.secondary)
                            }
                        }
                        
                        Button("Logout", role: .destructive) {
                            Task {
                                await authViewModel.logout()
                            }
                        }
                    }
                }
                
                Section("Playback") {
                    Picker("Video Quality", selection: $videoQuality) {
                        Text("Auto").tag("auto")
                        Text("1080p").tag("1080p")
                        Text("720p").tag("720p")
                        Text("480p").tag("480p")
                    }
                    
                    Toggle("Subtitles", isOn: $subtitlesEnabled)
                    Toggle("Autoplay", isOn: $autoplayEnabled)
                }
                
                Section("About") {
                    HStack {
                        Text("Version")
                        Spacer()
                        Text("1.0.0")
                            .foregroundColor(.secondary)
                    }
                }
            }
            .navigationTitle("Settings")
            .navigationBarTitleDisplayMode(.inline)
        }
    }
}

#Preview {
    SettingsView()
        .environmentObject(AuthViewModel())
}

