# 简介
IPA、APK内测应用分发服务，纯 golang 开发

## 用法
### 配置文件
```
{
  "host": "127.0.0.1",
  "port": "8080",
  "proxy": "http://127.0.0.1:8080" 
}
```
NGINX代理时需要配置`proxy`

### 开启服务
`$ ./ipapk-server`

### 上传
path:
```
POST /ipapk/upload
```
param:
```
package:安装包文件, reqiured
changelog:ChangeLog, 仅支持\n换行, optional
```
response:
```
{
    "uuid": "820e8c1b-954a-4489-8d6c-7ea2c47d1ec1",
    "name": "xxxxxxx",
    "platform": "ios",
    "bundleId": "com.xxxxxx",
    "version": "1.0.0",
    "build": "100",
    "install_url": "itms-services://?action=download-manifest&url=https://127.0.0.1:8080/ipapk/plist/820e8c1b-954a-4489-8d6c-7ea2c47d1ec1",
    "qrcode_url": "https://127.0.0.1:8080/ipapk/qrcode/820e8c1b-954a-4489-8d6c-7ea2c47d1ec1",
    "icon_url": "https://127.0.0.1:8080/ipapk/icon/820e8c1b-954a-4489-8d6c-7ea2c47d1ec1.png",
    "downloads": 0
}
```
示例:
`curl -X POST https://127.0.0.1:8080/ipapk/upload -F "file=@test.ipa -F "changelog=123" --insecure`


### SSL 证书
该项目需要HTTPS证书支持
