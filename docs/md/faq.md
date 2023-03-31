
## Q&A 常见问题

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