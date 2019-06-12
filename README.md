# Xiaomi
This plugin is a part of [Casa](https://github.com/getcasa), it's used to interact with xiaomi ecosystem.

## Downloads
Use the integrated store in casa or [github releases](https://github.com/getcasa/plugin-xiaomi/releases).

## Build
```
sudo env CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc go build -buildmode=plugin -o xiaomi.so *.go
```

## Install
1. Extract `xiaomi.zip`
2. Move `xiaomi` folder to casa `plugins` folder
3. Restart casa
