//
//  Models.swift
//  GalaoRenderer
//
//  Created by Tim on 19.02.26.
//

import Foundation

struct ViewNode: Codable {
    let kind: String
    let id: String?
    let value: String?
    let label: String?
    let children: [ViewNode]?
}

struct IncomingMessage: Codable {
    let type: String
    let tree: ViewNode?
}

struct OutgoingMessage: Codable {
    let type: String
    let id: String
    let event: String
}
