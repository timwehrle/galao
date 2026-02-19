//
//  GalaoRendererApp.swift
//  GalaoRenderer
//
//  Created by Tim on 19.02.26.
//

import SwiftUI

@main
struct GalaoRendererApp: App {
    @StateObject var state = AppState()
    
    var body: some Scene {
        WindowGroup {
            ContentView().environmentObject(state)
                .onAppear { state.startReading()
                    NSApplication.shared.activate(ignoringOtherApps: true)
                }
        }
    }
}
