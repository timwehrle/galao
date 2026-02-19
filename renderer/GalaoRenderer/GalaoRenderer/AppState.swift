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
    
    func send(event: OutgoingMessage) {
        guard let data = try? JSONEncoder().encode(event),
              let str = String(data: data, encoding: .utf8) else { return }
        print(str)
        fflush(stdout)
    }
    
    func startReading() {
            Thread.detachNewThread {
                while let line = readLine(strippingNewline: true) {
                    guard let data = line.data(using: .utf8),
                          let msg = try? JSONDecoder().decode(IncomingMessage.self, from: data)
                    else { continue }
                    
                    DispatchQueue.main.async {
                        if msg.type == "set_view" {
                            self.root = msg.tree
                        }
                    }
                }
            }
        }
}
