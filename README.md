

前端后端交互

前端完全不依赖后端，后端仅仅类似 Nginx 代理静态文件。Host 默认使用，但是不强制。如浏览器初始化默认存储当前域名和子域名。但是可能不会正常显示。

后端定时更换 Token 
打算实现一个邮件发送登录链接，或者是发送 5 位的验证码。从而实现前端登录。

两种 token 
1. 全权限 token 包括增删改查 （有效期内可用，配置中写明邮件，发送到邮箱从而获取有效 token）
2. 只能发送信息的 token (只要不改一直有效)


## tests commands

```bash
go test -v tools.go main_test.go 
```

http://127.0.0.1:8080/index.html?backend_address=http://127.0.0.1:8080&token=u80ZNcx1JzvIffo3CTWtFR11yhpgDyC9


            getQueryString(name) {
                var reg = new RegExp('(^|&)' + name + '=([^&]*)(&|$)', 'i');
                var r = window.location.search.substr(1).match(reg);
                if (r != null) {
                    return unescape(r[2]);
                }
                return null;
            }