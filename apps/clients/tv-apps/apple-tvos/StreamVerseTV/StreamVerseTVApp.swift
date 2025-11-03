//
//  StreamVerseTVApp.swift
//  StreamVerseTV
//
//  Created on 2025-01-XX.
//

import SwiftUI

@main
struct StreamVerseTVApp: App {
    @StateObject private var authViewModel = AuthViewModel()
    
    var body: some Scene {
        WindowGroup {
            if authViewModel.isLoggedIn {
                ContentView()
                    .environmentObject(authViewModel)
            } else {
                LoginView()
                    .environmentObject(authViewModel)
            }
        }
    }
}

