# MacOS as HomeKit Accessory

This app enables you controlling your MacOS over your HomeKit devices, like iPhone and Apple Watch.

## Features

- Volume control: Increase, decrease, and mute toggle
- Lock Screen: Lock, Current Lock Status (unlock is impossible)

## How to use

Precompiled binaries can be downloaded from the [Release pages](https://github.com/movsb/macos-as-homekit-accessory/releases).

1. Run `make` to get `lock` & `maha`
2. Run `./maha`
3. Open your iPhone's Home app, click on *Add Accessory* (be sure that your MacOS and your iPhone are in the same WiFi network)
4. Your MacOS should be listed there, which is named *MacOS ...*
5. Manually pair the accessory by entering the PIN code `00102003`
6. Finish the pairing process
7. Try to play with it 😀

## A Lightbulb? 💡

Yes, a lightbulb. Apple's HomeKit doesn't support Speaker/Mute characteristics anymore.
So instead, a lightbulb is shown there. Before you use Siri to control it, make sure it has a better name to not let Siri misunderstand it. 🥵
