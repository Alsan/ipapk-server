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

### SSL 证书
该项目会自动生成自签名HTTPS证书，需要手动安装、信任，NGINX代理时需要配置`proxy`以及非自签HTTPS证书，

### 开启服务
`$ ./ipapk-server`

### 上传
path:
```
POST /upload
```
param:
```
package:安装包文件, reqiured
changelog:ChangeLog, 仅支持\n换行, optional
```
response:
```
{
    "uid": "820e8c1b-954a-4489-8d6c-7ea2c47d1ec1",
    "name": "xxxxxxx",
    "platform": "ios",
    "bundleId": "com.xxxxxx",
    "version": "1.0.0",
    "build": "100",
    "install_url": "itms-services://?action=download-manifest&url=https://127.0.0.1:8080/plist/820e8c1b-954a-4489-8d6c-7ea2c47d1ec1",
    "qrcode_url": "https://127.0.0.1:8080/qrcode/820e8c1b-954a-4489-8d6c-7ea2c47d1ec1",
    "icon_url": "https://127.0.0.1:8080/icon/820e8c1b-954a-4489-8d6c-7ea2c47d1ec1",
    "downloads": 0
}
```
示例:
`curl -X POST https://127.0.0.1:8080/upload -F "file=@test.ipa" -F "changelog=123" --insecure`

## 截图
![s1](s1.png)

