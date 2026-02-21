//
//  Framing.swift
//  GalaoRenderer
//
//  Created by Tim on 19.02.26.
//

import Foundation

final class FramedIO {
    private let input = FileHandle.standardInput
    private let output = FileHandle.standardOutput
    
    func writeJSON<T: Encodable>(_ value: T) {
        guard let payload = try? JSONEncoder().encode(value) else { return }
        var len = UInt32(payload.count).bigEndian
        let lenData = Data(bytes: &len, count: 4)
        output.write(lenData)
        output.write(payload)
    }
    
    func readFrame() -> Data? {
           let lenData = input.readData(ofLength: 4)
           if lenData.count != 4 { return nil }

           let n = lenData.withUnsafeBytes { ptr -> UInt32 in
               ptr.load(as: UInt32.self).bigEndian
           }
           if n == 0 { return nil }

           let payload = input.readData(ofLength: Int(n))
           if payload.count != Int(n) { return nil }
           return payload
       }
}
