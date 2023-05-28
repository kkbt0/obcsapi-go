# Obsidian S3 存储的后端 API Golang 版本

python 老版本 https://gitee.com/kkbt/obsidian-csapi 
基于 Obsidian S3 存储， CouchDb ，本地存储和 WebDAV 的后端 API ,可借助 Obsidian 插件 Remotely-Save 插件，或者 Self-hosted LiveSync (ex:Obsidian-livesync) 插件 CouchDb 方式，保存消息到 Obsidian 库。

文档 Docs : [https://kkbt.gitee.io/obcsapi-go/#/](https://kkbt.gitee.io/obcsapi-go/#/)
如果你不使用 Obsidian ，也可以借助坚果云，或者 WebDav 进行文件同步，配合其他文本编辑器使用。

基于 Obsidian S3 存储的后端 API ,保存到 S3 存储的 Obsidian 库。支持列表

可开启 Webdav 服务，进行本地存储和文件管理
一个简易前端（后有图）
微信测试号 微信到Obsidian  
支持简悦 SimpRead Webook  
支持 fv悬浮球文字图片分享保存  
静读天下 MoonReader 高亮标注 仿 ReadWise API  
通用 http api  
基于 Obsidian S3 存储， CouchDb ，本地存储和WebDAV 的后端 API ,可借助 Obsidian 插件 Remotely-Save 插件，或者 Self-hosted LiveSync (ex:Obsidian-livesync) 插件 CouchDb 方式，保存消息到 Obsidian 库。特点

- 微信测试号 微信到 Obsidian
- 支持简悦 SimpRead Webook 裁剪网页文章
- 支持 fv悬浮球文字图片分享保存
- 静读天下 MoonReader 高亮标注 仿 ReadWise API
- 通用 http api
- 一个简易图床，附带命令行上传工具
- 云函数 或者 Dokcer 部署


更多功能说明见文档: [https://kkbt.gitee.io/obcsapi-go/#/](https://kkbt.gitee.io/obcsapi-go/#/)


## 展示

![](docs/images/Snipaste_2023-05-09_21-21-34.png)

![](docs/images/Snipaste_2023-05-09_21-22-36.png)

![](docs/images/Snipaste_2023-05-09_21-26-04.png)

![](docs/images/Snipaste_2023-05-09_21-26-13.png)


测试网站
```js
fetch('https://jsonplaceholder.typicode.com/todos/1')
      .then(response => response.json())
      .then(json => console.log(json))
```
