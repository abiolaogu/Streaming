//
//  StreamVerseApp.swift
//  StreamVerse
//
//  Created on 2025-01-XX.
//

import SwiftUI

@main
struct StreamVerseApp: App {
    @StateObject private var authViewModel = AuthViewModel()
    
    var body: some Scene {
        WindowGroup {
            if authViewModel.isLoggedIn {
                ContentListView()
                    .environmentObject(authViewModel)
            } else {
                LoginView()
                    .environmentObject(authViewModel)
            }
        }
    }
}

