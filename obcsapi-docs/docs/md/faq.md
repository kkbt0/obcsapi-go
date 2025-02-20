---
title: FAQ
---

# Q&A 常见问题

### Q: 使用云函数，有时响应很慢
A: 云函数会销毁长期不用的实例。这个时间可能只有几分钟。过了一定时间，再次访问会重启启动实例运行。这就是冷启动，想要减少这个时间，参考[函数计算冷启动优化最佳实践](https://help.aliyun.com/document_detail/140338.html)。最简单低廉的方法，大约是每 5 分钟左右就请求一次。这个请求可以在 阿里云-云函数-定时任务 中设置。

### Q: 微信测试号没有响应
A: 接口配置信息 URL 为 公网访问地址+ /api/wechat 。如 https://someone.cn-hangzhou.fcapp.run/api/wechat 。token 为 config 中 WeChat 的 Token。如果配置正确，接口配置信息能成功提交。

### Q: 提示  “SecondLevelDomainForbidden” 错误 如[I656IE: 不知道哪里出问题](https://gitee.com/kkbt/obsidian-csapi/issues/I656IE)
A: 类似
![](https://foruda.gitee.com/images/1670427617474643168/f95aeb1f_5177166.png)。解决方法参考 [https://help.aliyun.com/document_detail/375241.html](https://help.aliyun.com/document_detail/375241.html)。一般错误原因，是配置中使用 https://examplebucket.oss-cn-hangzhou.aliyuncs.com 类似格式 。

推荐使用Virtual-Hosted Style的访问方式。因为这个可以提高访问性能，少一跳。
配置示例：

|配置项目名	|OSS	|COS|
|---|---|---|
|s3Endpoint	|oss-cn-beijing.aliyuncs.com	|cos.ap-beijing.myqcloud.com|
|s3Region	|oss-cn-beijing	|ap-beijing|
|s3AccessKeyID	|AccessKey 获取	|AccessKey 获取|
|s3SecretAccessKey|	AccessKey 获取	|AccessKey 获取|
|s3BucketName|	obsidian	|obsidian-123|

有使用者反馈， OSS 如果错误可以这么配置

```
Endpoint = obbaocun.oss-cn-zhangjiakou.aliyuncs.com
Region = oss-cn-zhangjiakou
SecretId = xxx
SecretKey = xxx
bucket = root
```

如果是 Minio ，参考配置如下，配置 path_style: true 即可。 

```
# S3 配置
access_key: xxxxx
secret_key: xxxxx
end_point: http://127.0.0.1:9000
region: cn-e-0
bucket: obsidian
path_style: true
s3_wiki_link_use_presign: false
```

### Q: go 版本 https 问题

内网可申请自签名证书，放到容器里，配置中指定证书和私钥容器目录位置，开启 https 。
公网且有域名，可申请免费证书，然后配置反向代理等等。

如果没有 https ，部分版本较新的浏览器，https 和 http 不能混用。因此可能出现使用 https 的前端无法获取后端数据问题。
此外，没有 https ，前端 PWA 无法生效，不能自动弹出安装通知。

### Q: go 版本前端登录 404

找到齿轮，清除缓存，注销重新登录。可能出现于首次使用。

### Q: 前端错误
A: F12 打开开发者工具，可以看到元素，控制台等。选择应用程序，在弹出的页面中，左侧有应用程序，存储等。选择存储，本地存储，下面会有一个网址，右键，然后清除即可。

### Q: 微信模板消息为空
A: **！！！注意：微信模板消息施行掐头去尾，很有可能不好使！！！** 参考 [关于规范公众号模板消息的再次公告 2023 03 30](https://developers.weixin.qq.com/community/develop/doc/000a2ae286cdc0f41a8face4c51801?blockType=1&page=14#comment-list) 此外 5 月 4 日后中间的主内容中，单个字段内容不超过20个字，且不支持换行。