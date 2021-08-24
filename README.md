# 簡易圖床 | SimpleImageHosting

**注意**: 開發中 功能尚未實作完全

簡單的圖床伺服器

## 截圖

首頁
![image](https://raw.githubusercontent.com/orangesobeautiful/SimpleImageHosting/master/demo/homepage.jpg)
使用者
![image](https://raw.githubusercontent.com/orangesobeautiful/SimpleImageHosting/master/demo/user-images.jpg)
圖片資訊
![image](https://raw.githubusercontent.com/orangesobeautiful/SimpleImageHosting/master/demo/image-info.jpg)

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
npm install
quasar build
```
