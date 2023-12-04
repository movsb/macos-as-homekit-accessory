// Swift: How to observe if screen is locked in macOS
// https://stackoverflow.com/a/54356794/3628322

// build with: `swiftc lock.swift`

import Foundation

NSLog("lock started")

let dnc = DistributedNotificationCenter.default()

let lockObserver = dnc.addObserver(forName: .init("com.apple.screenIsLocked"), object: nil, queue: .main) { _ in
	NSLog("Screen Locked")
}

let unlockObserver = dnc.addObserver(forName: .init("com.apple.screenIsUnlocked"), object: nil, queue: .main) { _ in
	NSLog("Screen Unlocked")
}

RunLoop.current.run()
