---
home: true
icon: home
title: 项目主页
heroImage: /logo.svg
bgImage: https://theme-hope-assets.vuejs.press/bg/6-light.svg
bgImageDark: https://theme-hope-assets.vuejs.press/bg/6-dark.svg
bgImageStyle:
  background-attachment: fixed
heroText: Obcsapi
tagline: Obsidian 云存储 API
actions:
  - text: 文档
    icon: lightbulb
    link: ./md/
    type: primary

  - text: 源码
    link: https://gitee.com/kkbt/obcsapi-go

highlights:
  - header: 功能
    image: /assets/image/box.svg
    bgImage: https://theme-hope-assets.vuejs.press/bg/3-light.svg
    bgImageDark: https://theme-hope-assets.vuejs.press/bg/3-dark.svg
    highlights:
      - title: 基于 Obsidian S3 存储， CouchDb ，本地存储和 WebDAV 的后端 API ,可借助 Obsidian 插件 Remotely-Save 插件，或者 Self-hosted LiveSync (ex:Obsidian-livesync) 插件 CouchDb 方式，保存消息到 Obsidian 库。或者支持本地文件夹的文本编辑器。

  - header: 后端 API
    description: 给云存储中的文本文件添加 API ，保存消息到云存储
    bgImage: https://theme-hope-assets.vuejs.press/bg/2-light.svg
    bgImageDark: https://theme-hope-assets.vuejs.press/bg/2-dark.svg
    bgImageStyle:
      background-repeat: repeat
      background-size: initial
    features:
      - title: 网页
        details: PWA 网页应用添加
      - title: 微信
        details: 微信测试号
      - title: 简悦
        details: 简悦 SimpRead Webook 裁剪网页文章
      - title: FV 悬浮球
        details: FV 悬浮球调用 API 添加图文
      - title: 静读天下
        details: 高亮标注添加
      - title: 其他
        details: 通用 API / Lua 脚本自定义处理逻辑添加

  - header: 前端应用
    description: 一个 PWA 应用。
    image: /assets/image/layout.svg
    bgImage: https://theme-hope-assets.vuejs.press/bg/5-light.svg
    bgImageDark: https://theme-hope-assets.vuejs.press/bg/5-dark.svg
    highlights:
      - title: 添加修改
        details: 添加修改基础功能
      - title: 上传图片
        details: 可以上传图片
      - title: 自定义表单
        details: 通过 JSONSchema 定义表单，交给后端 Lua 脚本处理
      - title: 指令模式
        details: 简易对话匹配 或 运行 Bash / Lua 脚本
      - title: 修改配置
        details: 修改前后端配置
      - title: 深色模式
        details: 可以自由切换浅色模式与深色模式

  - header: 高级
    description: 高级自定义功能
    image: /assets/image/advanced.svg
    bgImage: https://theme-hope-assets.vuejs.press/bg/4-light.svg
    bgImageDark: https://theme-hope-assets.vuejs.press/bg/4-dark.svg
    highlights:
      - title: 对话模式 Bash / Lua 脚本
        details: 运行自定义命令
      - title: 自定义请求处理
        details: 用户编写 Lua 脚本处理请求
      - title: 自定义定时任务
        details: 用户编写 Lua 脚本处理定时任务
      - title: 自定义表单
        details: 使用 JSONSchema 自定义表单，用户编写 Lua 脚本处理自定义表单     

copyright: false
footer: 使用 <a href="https://theme-hope.vuejs.press/zh/" target="_blank">VuePress Theme Hope</a> 主题 | MIT 协议, 版权所有 © 2023 恐咖兵糖
---
