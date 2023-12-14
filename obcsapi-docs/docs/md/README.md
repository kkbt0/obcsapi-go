---
title: Obcsapi 文档
index: false
icon: laptop-code
category:
  - 使用指南
---

Obsidian Cloud Storage API 文档

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

基于 Obsidian S3 存储， CouchDb ，本地存储和 WebDAV 的后端 API ,可借助 Obsidian 插件 Remotely-Save 插件，或者 Self-hosted LiveSync (ex:Obsidian-livesync) 插件 CouchDb 方式，保存消息到 Obsidian 库。或者支持本地文件夹的文本编辑器。特点

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

**注意如果文档由于浏览器缓存，或者更新不及时无法查看最新内容。请到项目 https://gitee.com/kkbt/obcsapi-go/tree/master/obcsapi-docs/docs 中查看**


两个版本:  
[Obsidian 云存储后端 API Go 版本](https://gitee.com/kkbt/obcsapi-go)
[Obsidian S3 存储的后端 API python 版本](https://gitee.com/kkbt/obsidian-csapi)  


|                                          | python           | go             |
| ---------------------------------------- | ---------------- | -------------- |
| 体积                                     | 未压缩 100Mb+    | 未压缩50Mb-    |
| 修改难度                                 | 简单<sup>2</sup> | 比python复杂   |
| 微信公众号(测试号)<sup>3</sup>           | √                | √              |
| 微信文章裁剪                             | √                | √              |
| 简悦 SimpRead Webook                     | √                | √              |
| fv悬浮球文字图片分享保存                 | √                | √              |
| 静读天下 MoonReader 高亮标注<sup>4</sup> | √                | √              |
| 通用 http api                            | √                | √              |
| S3 对象存储                              | √                | √              |
| CouchDb                                  | ×                | √              |
| Local (Webdav Server)<sup>5</sup>        | ×                | √              |
| Webdav                                   | ×                | √              |
| Web 网页支持                             | √                | √ <sup>6</sup> |
| 图床  和 CLI上传工具                     | ×                | √              |
| 公开文档功能                             | ×                | √              |
| 邮件/微信任务提醒                        | ×                | √              |
| 前端                                     | √                 | √              |
| 对话/指令模式                            | ×                | √              |
| 自定义运行脚本                           | ×                | √              |
| JSONSchema 表单                          | ×                | √              |
| Docker                                   | <=3              | >=4            |

Docker 3.0 之前的是 python 版本，之后的是 Go 版本。

Python 版本容易在阿里云-云函数服务处修改和部署，使用者完全不需要安装任何 Python 环境就可以简单修改。Go 版本性能高一些，实现了邮箱登录链接的功能，但是缺少微信文章剪裁的功能，而且修改需要 Golang 语言环境。

---

相关链接:

Docker [https://hub.docker.com/r/kkbt/obcsapi](https://hub.docker.com/r/kkbt/obcsapi)

博客教程及效果（小白图文版）: [https://www.ftls.xyz/posts/obcsapi-fc-simple/](https://www.ftls.xyz/posts/obcsapi-fc-simple/)
视频效果演示和 Python 版部署教程：[Obsidian 从本地到云端-哔哩哔哩](https://b23.tv/uJFvw3A)

前端 Demo [https://kkbt.gitee.io/obweb/#/Memos](https://kkbt.gitee.io/obweb/#/Memos) 第一次加载不正常显示，后端加载慢属于正常现象。


---

[1] CLI 支持剪贴板图片上传。可配合 Obsidian 插件 Image Auto upload Plugin ，实现图片上传。类似 PicGO-Core
[2] python 版本在阿里云函数计算 FC 上，可以方便的修改文件。python 无需编译，并且生态丰富
[3] 支持文字，图像，语音（转文字存储）等
[4] 类似 ReadWise API 增加接口
[5] 本地服务 LocalStorage 本地存储，开启 Webdav 服务为 Remotely Save 提供同步。同时 WebDav 服务可连接 RAIDrive (Windows) ， Mix (安卓) 等进行文件管理。
[6] Go 版本支持 Web 网页，S3 支持 Obsidian 库内的图片链接（需要基于库的路径或相对路径）。CouchDb,Local 也支持。

---

文档技术支持: [Docsify](https://docsify.js.org/#/)