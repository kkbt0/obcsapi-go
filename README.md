# Obsidian 云存储的后端 API Golang 版本


基于 Obsidian 云存储， CouchDb ，本地存储和 WebDAV 的后端 API ,可借助 Obsidian 插件 Remotely-Save 插件，或者 Self-hosted LiveSync (ex:Obsidian-livesync) 插件 CouchDb 方式，保存消息到 Obsidian 库。

文档 Docs : [https://kkbt.gitee.io/obcsapi-go/#/](https://kkbt.gitee.io/obcsapi-go/#/)
如果你不使用 Obsidian ，也可以借助坚果云，或者 WebDav 进行文件同步，配合其他文本编辑器使用。

![](docs/images/default_canvas.svg)

绘图 PowerBy [Handraw](https://handraw.top/)

![](docs/images/canvas_2_show.svg)

基于 Obsidian 云存储的后端 API ,保存到 S3 存储的 Obsidian 库。支持列表

可开启 Webdav 服务，进行本地存储和文件管理
一个简易前端（后有图）
微信测试号 微信到Obsidian  
支持简悦 SimpRead Webook  
支持 fv悬浮球文字图片分享保存  
静读天下 MoonReader 高亮标注 仿 ReadWise API  
通用 http api  
基于 Obsidian S3 存储， CouchDb ，本地存储和WebDAV 的后端 API ,可借助 Obsidian 插件 Remotely-Save 插件，或者 Self-hosted LiveSync (ex:Obsidian-livesync) 插件 CouchDb 方式，保存消息到 Obsidian 库。特点

- 前端添加 Memos / 简答编辑 ， 支持指令模式，有黑暗主题 ，是 PWA 应用
- 微信测试号 微信到 Obsidian
- 支持简悦 SimpRead Webook 裁剪网页文章
- 支持 fv悬浮球文字图片分享保存
- 静读天下 MoonReader 高亮标注 仿 ReadWise API
- 通用 http api ， 有 Swagger 配合调试
- 可提供 WebDAV 服务
- 一个简易图床，附带命令行上传工具
- 自定义指令机器人，可执行脚本和命令。可接入其他服务
- 云函数 或者 Dokcer 部署
- 开放 API ， 可自行配合其他软件使用。如 Quicker 


更多功能说明见文档: [https://kkbt.gitee.io/obcsapi-go/#/](https://kkbt.gitee.io/obcsapi-go/#/)


## 展示

### PWA Web 应用

![](docs/images/Snipaste_2023-05-09_21-21-34.png)

![](docs/images/Snipaste_2023-05-09_21-22-36.png)

![](docs/images/Snipaste_2023-05-09_21-26-04.png)

![](docs/images/Snipaste_2023-05-09_21-26-13.png)

### Tauri 桌面端应用 

[Tauri 可以构建跨平台的快速、安全、前端隔离应用](https://tauri.app/zh-cn/)。图片为 桌面端应用  - Windows - 虽然大概直接用 Obsidian 更方便一些。

![](docs/images/Snipaste_2023-08-01_12-57-50-tauri-windows.png)

可从 [https://gitee.com/kkbt/obcsapi-go/releases](https://gitee.com/kkbt/obcsapi-go/releases) 下载

## 其他

python 老版本 https://gitee.com/kkbt/obsidian-csapi 
发展历程: [obsidian-使用指南](https://www.ftls.xyz/series/obsidian-%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97/)