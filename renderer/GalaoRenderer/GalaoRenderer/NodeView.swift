//
//  NodeView.swift
//  GalaoRenderer
//
//  Created by Tim on 19.02.26.
//

import SwiftUI

struct NodeView: View {
    let node: ViewNode
    let state: AppState
    
    var body: some View {
        switch node.kind {
        case "vstack":
            VStack {
                ForEach(node.children ?? [], id: \.id) { child in
                    NodeView(node: child, state: state)
                }
            }
        case "hstack":
            HStack {
                ForEach(node.children ?? [], id: \.id) { child in
                    NodeView(node: child, state: state)
                }
            }
        case "text":
            Text(node.value ?? "")
        case "button":
            Button(node.label ?? "") {
                state.send(event: OutgoingMessage(
                    type: "event",
                    id: node.id ?? "",
                    event: "tap"
                ))
            }
        default:
            EmptyView()
        }
    }
}
