# casa-xiaomi

## Build

```
sudo env CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc go build -buildmode=plugin -o xiaomi.so *.go
```