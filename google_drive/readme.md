Google Drive API
===
都需要用 OAuth2.0 來認證，所以會跳網頁出來請你登入

## v3 的版本
直接照著 https://developers.google.com/drive/api/v3/quickstart/go 
就一切順利了

首先到管理頁註冊 OAuth2 的用戶端憑證
https://console.cloud.google.com/apis/api/drive.googleapis.com/credentials

下載他的 `credentials.json` 放進來  
直接 `go run .` 就出去了，第一次需要認證他會再長出一個 token.json  
這樣下一次就不需要認證了

## v2 的版本
可以看這個範例
https://github.com/googleapis/google-api-go-client/tree/master/examples

首先到管理頁註冊 OAuth2 的用戶端憑證
https://console.cloud.google.com/apis/api/drive.googleapis.com/credentials

建立好憑證後，選進去編輯找到你的 client id 還有 client secret 就了事了  
跑的話因為 github 上把各種專案都放一起所以會比較麻煩
```
cd v2
export CLIENT_ID=xxxxxxxxxx.apps.googleusercontent.com
export CLIENT_SECRET=mLbZRCoFTCMgP3G4hP0hfD5k
go run . -clientid=$CLIENT_ID -secret=$CLIENT_SECRET drive main.go
```

然後你的 main.go 就上傳上你的 drive 囉 XDD