# 簡易圖床 | SimpleImageHosting

**注意**: 開發中 功能尚未實作完全

簡單的圖床伺服器

## 編譯

以 Ubuntu 20.04 為例:

### 後端

安裝 GO 套件(如果沒有)

```
sudo apt install golang
```

開始編譯:

```
cd backend
CGO_ENABLED=0 go build -o SIHBackned -tags=go_json main.go
```

### 前端

[安裝 quasar](https://v1.quasar.dev/quasar-cli/installation)

建置

```
cd frontend
quasar build
```
