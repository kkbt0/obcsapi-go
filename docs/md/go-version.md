## Go 语言版本

见 [https://gitee.com/kkbt/obcsapi-go](https://gitee.com/kkbt/obcsapi-go)

部分说明请见 python 版本说明更详细一些

若有更高性能需求，可自行修改 golang 源码，并编译打包。不过 golang 功能要少一些（没有收藏微信文章功能）。新特性：可以发送邮件发送登录链接，这样配合 token 的自动失效能稍微安全一些，不过我不太擅长前端，很简陋。


### Token 说明

后端定时更换 Token 。实现了一个邮件发送登录链接，从而实现前端登录。

两种 token 
1. token1 包括增删改查 （有效期内可用，配置中写明邮件，发送到邮箱从而获取有效 token）。可设置有效期
2. token2 用于其他服务调用，无限期可用
3. 运行时，会在终端显示

### 部署说明

复制 config.examples.yaml 为 config.yaml 。部署时建议把项目文件夹内文件都复制过去。（至少包含 template , token 两个文件夹中，及其相关内容。 tem.txt 和 config.yaml 两个文件。


现在这玩意算是能用，但是这个 token 系统还称不上好用。

部署方法

1. 拉取项目，构建部署
2. Docker 部署
3. 阿里云云函数 FC 运行构建的程序或 Docker

#### 配置文件说明

```yaml
name: obcsapi-go # 项目名称
version: v4.0.9 # 项目版本
description: by kkbt # 描述
host: 0.0.0.0 
port: 8900
token_path: ./token/ # token 文件位置。可用于云函数挂载硬盘位置
token1_live_time: 72h # 可设置有效期 至少 3 天 也可以设置 9999h

# 前端后地址 邮件中登录链接会用到
front_url: https://kkbt.gitee.io/obweb/#/Memos
backend_url: api.ftls.xyz # 这个值图床也会用到
wechat_return_str: "📩 已保存，<a href='https://note.ftls.xyz/#/ZK/202209050658'>点击查看今日笔记</a>" # 微信返回的字符串

# Obsidian 文件夹设置

ob_daily_dir: 日志/
ob_daily_attachment_dir: 日志/附件/
ob_daily_other_dir: 支持类文件/ # 用于第三方软件 如静读天下，简悦
# 图床文件 有三部分构成 文件夹及前缀，原名字，随机字符
# 图床文件夹及文件前缀 eg 2006-01-02 15:04:05 如 按月存放是 01/ ; 按 年存放 2006/ ; 文件前缀 200601 ; 文件夹和文件前缀 200601/200601_
images_hosted_fmt: 200601/kkbt_
images_hosted_use_raw_name: true # 图床文件是否使用原名字 true or false
images_hosted_random_name_length: 5 # 图床文件随机字符命名 随机字符长度
# 百度 OCR https://ai.baidu.com/ai-doc/OCR/zk3h7xz52 该项置空或删除此项则不进行 OCR 注意该项有效期 30 天
# https://ai.baidu.com/ai-doc/REFERENCE/Ck3dwjhhu
# bd_ocr_access_token: xxxxx.xxxxx.xxxxx.xxxxx.xxxxx-xxxxx

# S3 -> 1 ; CouchDb -> 2
data_source: 2

# CouchDb 配置
couchdb_url: http://admin:password@localhost:5984/ # admin 是账户名 ; password 是密码；
couchdb_db_name: note # 数据库

# S3 配置 几乎同 Remotely-Save S3 配置
access_key: xxxxxxxxxxxxxx
secret_key: xxxxxxxxxxxxxx
end_point: https://cos.ap-beijing.myqcloud.com
region: ap-beijing
bucket: obsidion-xxxxxxxxxxxxxx

# wechat 测试号/公众号
wechat_token: xxxxxxxxxxxxxx # 微信公众平台的Token
wechat_appid: xxxxxxxxxxxxxx # 微信公众平台的AppID
wechat_secret: xxxxxxxxxxxxxx # 微信公众平台的AppSecret
wechat_openid: xxxxxxxxxxxxxx # OpenID 自己关注测试号后，获取的

# 任务提醒
email_reminder_time: 0700 # 指每天 07:00 

# smtp 邮箱服务
smtp_mail:
  smtp_host: smtpdm.aliyun.com
  port: 80
  username: no-reply@mail.ftls.xyz
  password: xxxxxxxxxxxxxx
  # 以下一个是发送者 一个是接受者
  mail_sender_address: no-reply@mail.ftls.xyz
  mail_sender_name: "ObCSAPI"
  mail_send_to: yourmail@foxmail.com # 接受者邮箱
```

#### 构建 server

```bash
cd server/
go run .
go build -o server  -ldflags '-linkmode "external" -extldflags "-static"' .
```

#### Docker

```sh
# 构建镜像
docker build -t kkbt/obcsapi:v4.0.9 . 
# 运行 Docker
docker run -d -p 8900:8900 --name myObcsapi4.0.9 -v /home/kkbt/app/obcsapi-go/:/app/data/ kkbt/obcsapi:v4.0.9
# 或者通过 cp 方式修改好的 config.yaml
docker cp config.yaml myObcsapi4.0.9:/app/data/config.yaml
docker restart myObcsapi4.0.9
```
如果 -v 后文件出现没有权限访问的问题，可在宿主机执行 `sudo chmod 777 -R /home/kkbt/app/obcsapi-go/` 。

### 基本功能使用

接口一览

| 方式     |      验证      | 路径                | 说明                        |
| -------- | :------------: | ------------------- | --------------------------- |
| Any      | WeChat Mp 验证 | /api/wechat         | 微信测试号使用              |
| GET      |       /        | /api/sendtoken2mail | 邮件发送 token              |
| GET/POST |     token1     | /ob/today           | Web 使用                    |
| POST     |     token1     | /ob/today/all       | Web 使用                    |
| GET      |     token1     | /ob/recent          | Web 使用                    |
| POST     |     token2     | /ob/fv              | fv 悬浮球自定义任务图文分享 |
| POST     |     token2     | /ob/sr/webhook      | 简悦 Webhook                |
| POST     |     token2     | /ob/general         | General 通用接口            |
| POST     |     token2     | /ob/url             | 页面转 md 存储 效果很一般   |
| POST     |     token2     | /ob/moonreader      | 静读天下高亮和笔记          |
| POST     |     token2     | /api/upload         | 图床 PicGo 可用             |
| GET      |       /        | /public/*           | 公开文档功能                |

一些访问方式可见文件 http/server.http

#### 通用接口

```http
POST {{host}}/ob/general
Token: {{token2}}
Content-Type: application/json

{"content":"test token2"}
```

如 IOS 捷径可使用此接口

![](../images/IMG_5471_recompress.jpg)

![](../images/IMG_5472_recompress.jpg)

#### Public 公开文档功能

以配置项中 `ob_daily_other_dir: 支持类文件/` 为例，在 `支持类文件/Public/`下的文档可以公开访问。如 `支持类文件/Public/test.md` 的访问路径为 `https://127.0.0.1:8900/public/test.md` 。仅支持 Markdown 文件，图片文件公开请使用图床功能。

#### WeChat MP

微信测试号，配置完成后。微信测试号发送消息返回的链接，已保存的网址，进入网站。按下齿轮，勾选 Debug。拉到页面下面，两个输入框，第一个是后端api。第二个是 token。api填写类似 api.ftls.xyz/ob ，不需要协议头和尾部斜杠。按函数计算fc格式，第一个框如 someone.cn-hangzhou.fcapp.run/ob ，第二个框 `BxzbwiSIxeXCv71ktQe6cAIcLVy6qd96` 为 Token1 的值。

填写完成后，点击 updateConfig 按钮并刷新页面。

Enjoy it !


##### 微信公众号说明

微信测试号发送消息，保存 Obsidian 中。前端在 https://gitee.com/kkbt/obweb 中，实例 https://kkbt.gitee.io/obweb/#/ 点击右侧图片。使用 localStorge 存储 api 地址和 token 。你可以将 https://gitee.com/kkbt/obweb frok 下来，用 Memos.md 文件内容覆盖 HOME.md ，然后借助 github pages 静态部署整个项目，这个项目整个文件夹就是一个 Obsidian 库，也是一个 Docsify 文档文件夹。如果使用对象存储的静态网站功能，可以轻松部署成网页。


- 支持图片和文字
- 图片下载到存储本地，而非链接(微信发送的图片，会给咱们的服务器返回图片URL)
- 对用户的判断，仅限特定用户存储笔记。(根据 OpenID 判断)
- 检索文字中含有 "todo" ，则生成勾选框。如 `- [ ] 13:11 somethingtodo`
- 正常生成 `- 13:11 something`
- 内容能在 Obsidian 插件 Memos 中正常显示
- 提供三天查询 一天修改的数据
- 支持消息类型: 文字，图片，链接(收藏中的)，地图位置，语音消息(直接调用微信转文字存储)
 
**不建议做的事情：**

- 不推荐批量传图片，推荐显示已保存后依次上传。
- 不推荐一秒内上传多个文件，图片命名精确到1S。1S内多图片会覆盖。
- 不要使用微信自带的表情符号，请使用输入法表情。
- 如果微信输入框换行或分段，只会在这一条消息最开始有 `- 13:11 `。也就是说，第二行、第二段不会在 Memos 插件中显示。

#### 邮件发送 

访问 /api/sendtoken2mail 路径，这个路径大约 3 天可以访问一次。会更新 token1 并发送邮件到指定邮箱。

#### 任务提醒


##### 微信分钟级别提醒

[登录微信测试号](https://mp.weixin.qq.com/debug/cgi-bin/sandboxinfo?action=showinfo&t=sandbox/index),模板消息接口新增测试模板，标题随意。内容处包含 `{{content.DATA}}` 即可。如

```
待办任务： {{content.DATA}}
```

创建完成后，模板ID 复制到配置文件，作为 `wechat_template_id`  的值。

库的根目录下创建 `提醒任务` .( Obsidian 会创建 提醒任务.md 文件，所以不必加后缀 .md)。内容示例如下

```md
这是一个任务提醒功能
20230407 1322 测试提醒1
20230407 1353 测试提醒2<br>第二行
```

示例中,`20230407 1322 测试提醒1` 会在 2023年04月07日 13:22 发送微信模板消息。发送后的那一行消息，按配置会被移动到 `支持类文件/WeChatSended/202304.md` 中。

该功能每分钟查询一次，所以此 `提醒任务.md` 文件不要过大。否则会引起流量过大等问题。**程序不会自动删除过期任务，或者是不符合任务提醒格式的行，需要手动确认删除无用内容**

快捷创建方式: 微信语音或文章发送包含 `提醒我` 的句子，即可添加到 `提醒任务.md`。如发送 `提醒我后天早上10:35的会议，需要及时参与` ，`提醒任务.md` 会新增 `20230407 1322 测试提醒1` 。同时后天的日记中会出现 `- [ ] 提醒我后天早上10:35的会议，需要及时参与`。

![](../images/20230407175219.jpg)

![](../images/20230407175618.jpg)

使用库: https://github.com/DanPlayer/timefinder

<details>
<summary>更多示例</summary>
  <p>6月9日有一场show要去观看<br>
2021-06-09 00:00:00<br>
后天早上10:35的会议，需要及时参与<br>
2021-06-04 10:35:00<br>
明天下午三点的飞机，提醒我坐车<br>
2021-06-03 15:00:00<br>
一个小时后提醒我喝水<br>
2021-06-02 18:40:07<br>
明天早上8:00喊我起床<br>
2021-06-03 08:00:00<br>
明天早上8点喊我起床<br>
2021-06-03 08:00:00<br>
明早十点喊我喝水<br>
2021-06-03 10:00:00<br>
明天早上十点喊我喝水<br>
2021-06-03 10:00:00<br>
明天下午三点提醒我喝水<br>
2021-06-03 15:00:00<br>
一天后提醒我喝水<br>
2021-06-03 17:40:07<br>
一年后提醒我喝水<br>
2022-06-02 17:40:07<br>
一个月后提醒我喝水<br>
2021-07-02 17:40:07<br>
一月后提醒我喝水<br>
2021-07-02 17:40:07<br>
我要住到大后天<br>
2021-06-05 00:00:00<br>
我要住到明天<br>
2021-06-03 00:00:00<br>
下个月到上个月再到这个月<br>
2021-07-02 00:00:00<br>
我要住到明天下午三点十分<br>
2021-06-03 15:10:00<br>
帮我预定明天凌晨3点的飞机<br>
2021-06-03 03:00:00<br>
今天13:00的飞机<br>
2021-06-02 13:00:00<br>
3月15号的飞机<br>
2021-03-15 00:00:00<br>
昨天凌晨2点<br>
2021-06-01 02:00:00<br>
十分钟后提醒我喝水<br>
2021-06-02 17:50:07<br>
</p>
</details>


##### 邮件每日提醒

邮件配置正确情况下，并且配置文件中 `email_reminder_time` 处于可取的值。 Obcsapi 会读取根目录 `每日提醒.md` 和三天日记（今天，昨天，前天）中的 `- [ ]` 开头的行发送到指定邮箱中。

如果想关闭此项功能可将 `email_reminder_time` 设置为 `9999` 等数字。


#### fv 悬浮球

设置自定义任务，设置分享菜单为文字，图片 POST /ob/fv ，Header Token 为 token2 的值。

[图文教程 中的 fv 部分](https://www.ftls.xyz/posts/obcsapi-fc-simple/#%E4%B8%8B%E8%BD%BD%E5%8E%8B%E7%BC%A9%E5%8C%85)

演示和教程 见 https://www.bilibili.com/video/BV1Ad4y1s7EP/

#### 简悦 Webhook 配置

简悦-服务 Webhook 填写样例

```json
{"name": "WH2COS","url": "http://127.0.0.1:9000/webhook","type": "POST","headers": {"Content-Type": "application/json","Token": "your_simp_read_token"},"body": {"url": "{{url}}","title": "{{title}}","desc": "{{desc}}","content": "{{content}}","tags": "{{tags}}","note": "{{note}}"}}
```


### 图床功能 PicGo 配置

![](../images/imagesHost1.png)

![](../images/imagesHost2.png)

![](../images/imagesHost3.png)

POST {{host}}/api/upload
Token: {{token2}}
Content-Type: form-data

字段名: file

```json
{"data":{"url":"http://example.com/images/test.jpg","url2":"http://example.com/images/test.jpg"}}
```
~~第三行是选 url or url2 。url 是 http ；url2 是 https~~ 现在 url url2 都一样的，http/https 由配置文件决定

### 图床 PicGo-Core 类似 cli 工具

概述：用于上传到 Obcsapi 图床的专用命令行工具。可配合 Obsidian 插件 Image Auto upload Plugin ，实现 Obsidian 图片上传。

下载地址： [https://gitee.com/kkbt/obcsapi-go/tree/master/images-host-cli/out](https://gitee.com/kkbt/obcsapi-go/tree/master/images-host-cli/out)

windows，linux，mac 下

```cmd
obcsapi-picgo.exe :: 没有配置则生成配置
obcsapi-picgo.exe u .\1.png .\2.jpg :: 上传本地图片 1.png 2.jpg
obcsapi-picgo.exe u  :: 上传剪切板图片，如果有的话
obcsapi-picgo.exe u https://cdn.ftls.xyz/images/2022/11/zhishi2.jpg :: 下载网络图片，并上传
```

配置说明，第一次运行会自动生产配置。有三行
```
http://127.0.0.1:8900/api/upload
fQbzONJAAw
url
```

第一行是上传链接
第二行是 token2 的值，程序会在请求 Header Token 设置为该值
~~第三行是选 url or url2 。url 是 http ；url2 是 https~~ 现在 url url2 都一样的，http/https 由配置文件决定

除了命令行运行之外，还可配合 Obsidian 插件 Image Auto upload Plugin .选择 Picgo-Core 模式，路径选择可执行文件在的位置 如 `C:\CLI\obcsapi-picgo.exe`。也可以改名，为`picgo.exe`。

### 图床 OCR（试验性质）

图床可选配置，使用百度 OCR 进行文字识别，结果保存到 Obsidian 便于后续查询。使用 https://ai.baidu.com/tech/ocr

使用百度 OCR 通用文字识别（标准版）

配置文件中将 bd_ocr_access_token 取消注释，替换成自己的 access_token 即可，注意 access_token 有效期 30 天。  
access_token 申请 [https://ai.baidu.com/ai-doc/REFERENCE/Ck3dwjhhu](https://ai.baidu.com/ai-doc/REFERENCE/Ck3dwjhhu)

```yaml
# 百度 OCR https://ai.baidu.com/ai-doc/OCR/zk3h7xz52 该项置空或删除此项则不进行 OCR 注意该项有效期 30 天
# https://ai.baidu.com/ai-doc/REFERENCE/Ck3dwjhhu
# bd_ocr_access_token: xxxxx.xxxxx.xxxxx.xxxxx.xxxxx-xxxxx
```

数据按配置文件将保存在: `支持类文件/OcrData/bdocr.md` 。此文件可以随意删除，移动，更换名字。如果上传图片时没有此文件，程序会自动创建追加写入 OCR 结果。

### 展示

后台发送的邮件

![](../images/Snipaste_2023-03-07_11-36-48.png)

点击进入的样子

![](../images/Snipaste_2023-03-07_11-37-38.png)

### 开发说明

Go 语言开发

文件说明 server

- dao 数据操作代码
- dao/couchdb.go couchdb 处理操作相关的代码
- dao/s3.go s3 处理操作相关的代码
- dao/type.go 一些 struct 定义
- static/* 前端文件
- template/* 前端文件
- token/token1 存储 token1 前端 token ，有效期取决于配置文件
- token/token2 存储只发送 token2 用于一些只发送的其他第三方程序 api、
- tools/token.go Token 相关代码 。Token 的生成，修改，验证代码
- tools/tools.go 一些辅助代码
- assest.go 一些辅助代码
- build.sh 将 go 打包成无链接的静态文件
- config.example.yaml 示例配置文件
- docker-entrypoint.sh Docker entry 脚本
- Dockerfile Docker 构建代码
- handles.go 一些路径处理
- mail.go 发邮件相关代码
- middleware.go 中间件代码
- obsidian.go 发送到 Obsidian 的相关代码
- server.go 程序主要入口
- wechat.go 微信公众号（测试号）处理代码
- http/server.http vscode 插件 postcode 使用的，用于测试
- image-hosting-cli 图床 CLI 工具

### 更新记录 版本说明

4.0.2 每次部署 Docker 更新token，会在日志中显示一次
4.0.3 可以自定义日记等文件夹
4.0.4 可自定义微信返回内容 可 -v 映射文件夹 /app/data/
4.0.5 新增一个简易图床
4.0.6 图床增加一些自定义功能 ，增加可配合 Obsidian 插件 Image auto upload Plugin ，使用 PicGo-Core 即可上传
4.0.7 增加 Public 目录公开访问文档功能;开启日志功能;修复了一些已知错误
4.0.8 增加百度图片OCR功能进行测试
4.0.9 任务提醒功能