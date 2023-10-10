# Obsidian 云存储的后端 API Golang 版本

Back-end APIs based on Obsidian S3 storage , CouchDb Local and WebDAV can save messages to the Obsidian library with the help of the Obsidian plugin Remotely-Save plugin, or Self-hosted LiveSync (ex:Obsidian-livesync) plugin CouchDb mode.Or a text editor that supports local folders. peculiarity

- Add Memos / Short Answer Editor on the front end, support instruction mode, have dark theme, and be a PWA app
- WeChat MP to Obsidian
- Support for SimpRead Webook to crop web articles
- Support FooView hoverball text picture sharing and saving
- MoonReader highlights
- Universal HTTP API
- Extend functionality with Lua & Bash . Users can process any request
- WebDAV Server
- A simple graph bed with a command line upload tool. <sup>1</sup>
- SCF or Dokcer deployment

基于 Obsidian S3 存储， CouchDb ，本地存储和WebDAV 的后端 API ,可借助 Obsidian 插件 Remotely-Save 插件，或者 Self-hosted LiveSync (ex:Obsidian-livesync) 插件 CouchDb 方式，保存消息到 Obsidian 库。或者支持本地文件夹的文本编辑器。特点

- 前端添加 Memos / 简答编辑 ， 支持指令模式，有黑暗主题 ，是 PWA 应用
- 微信测试号 微信到 Obsidian
- 支持简悦 SimpRead Webook 裁剪网页文章
- 支持 fv悬浮球文字图片分享保存
- 静读天下 MoonReader 高亮标注 仿 ReadWise API
- 通用 http api
- 使用 Lua & Bash 拓展功能。用户可以处理任何请求
- WebDAV 服务
- 一个简易图床，附带命令行上传工具。<sup>1</sup>
- 云函数 或者 Dokcer 部署


更多功能说明见文档: [https://kkbt.gitee.io/obcsapi-go/#/](https://kkbt.gitee.io/obcsapi-go/#/)

---

文档 Docs : [https://kkbt.gitee.io/obcsapi-go/#/](https://kkbt.gitee.io/obcsapi-go/#/)
如果你不使用 Obsidian ，也可以借助坚果云，或者 WebDav 进行文件同步，配合其他文本编辑器使用。

![](docs/images/default_canvas.svg)

绘图 PowerBy [Handraw](https://handraw.top/)

![](docs/images/canvas_2_show.svg)



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