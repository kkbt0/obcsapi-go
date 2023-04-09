# Obsidian S3 存储的后端 API Golang 版本

python 老版本 https://gitee.com/kkbt/obsidian-csapi 

文档 Docs : [https://kkbt.gitee.io/obcsapi-go/#/](https://kkbt.gitee.io/obcsapi-go/#/)

基于 Obsidian S3 存储的后端 API ,保存到 S3 存储的 Obsidian 库。支持列表

可开启 Webdav 服务，进行本地存储和文件管理
一个简易前端（后有图）
微信测试号 微信到Obsidian  
支持简悦 SimpRead Webook  
支持 fv悬浮球文字图片分享保存  
静读天下 MoonReader 高亮标注 仿 ReadWise API  
通用 http api  
邮件发送登录链接


实现了一个邮件发送登录链接，从而实现前端登录。


## 部署

复制 config.examples.yaml 为 config.yaml 。部署时建议把项目文件夹内文件都复制过去。（至少包含 template , token 两个文件夹中，及其相关内容。 tem.txt 和 config.yaml 两个文件。


## 展示

后台发送的邮件

![](docs/images/Snipaste_2023-03-07_11-36-48.png)

点击进入的样子

![](docs/images/Snipaste_2023-03-07_11-37-38.png)

