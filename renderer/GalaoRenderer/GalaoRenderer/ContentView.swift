//
//  ContentView.swift
//  GalaoRenderer
//
//  Created by Tim on 19.02.26.
//

import SwiftUI


struct ContentView: View {
    @EnvironmentObject var state: AppState
    
    var body: some View {
        Group {
            if let root = state.root {
                NodeView(node: root, state: state)
                    .padding()
            } else {
                Text("Waiting...")
            }
        }
        .frame(minWidth: 300, minHeight: 200)
    }
}
