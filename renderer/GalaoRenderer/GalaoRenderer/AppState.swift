//
//  AppState.swift
//  Galao
//
//  Created by Tim on 19.02.26.
//

import SwiftUI
import Combine

class AppState: ObservableObject {
    @Published var root: ViewNode? = nil
    private let io = FramedIO()
    
    func send(event: OutgoingMessage) {
        io.writeJSON(event)
    }
    
    func startReading() {
        io.writeJSON(["type": "ready"])

        Thread.detachNewThread {
            while let data = self.io.readFrame() {
                DispatchQueue.main.async {
                    guard let msg = try? JSONDecoder().decode(IncomingMessage.self, from: data) else { return }
                    if msg.type == "set_view" {
                        self.root = msg.tree
                    }
                }
            }
        }
    }
}

