# MacOS as HomeKit Accessory

This app enables you controlling your MacOS over your HomeKit devices, like iPhone and Apple Watch.

## Features

- Volume control: Increase, decrease, and mute toggle
- Lock Screen: Lock, Current Lock Status (Unlock is impossible)

## How to compile

### Precompiled binaries

Precompiled binaries can be downloaded from the [Release pages](https://github.com/movsb/macos-as-homekit-accessory/releases).

Allowing to run in Terminal:

```bash
xattr -d com.apple.quarantine ./maha-darwin-arm64
xattr -d com.apple.quarantine ./maha-darwin-amd64
xattr -d com.apple.quarantine ./lock
```

### Manually compile

Just run `make` to get `lock` & `maha`

## Run

1. Run `./maha` (a directory named `db` will be crated to store the pairing information)
2. Open your iPhone's Home app, click on *Add Accessory* (be sure that your MacOS and your iPhone are in the same WiFi network)
3. Your MacOS should be listed there as a Bridge accessory, which is named *MacOS HomeKit*
4. Manually pair the accessory by entering the PIN code `00102003`
5. Finish the pairing process
6. Try to play with it 😀

## A Lightbulb? 💡

Yes, a lightbulb. Apple's HomeKit doesn't support Speaker/Mute characteristics anymore.
So instead, a lightbulb is shown there. Before you use Siri to control it, make sure it has a better name to not let Siri misunderstand it. 🥵
