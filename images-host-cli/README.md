http://127.0.0.1:8900/api/upload
fQbzONJAAw
url

第一行是上传链接
第二行是 token2 的值
第三行可选 url or url2 。url 是 http ；url2 是 https

构建 

CGO_ENABLED=0 剪切板功能失效

```bash
# linux
SET CGO_ENABLED=0  // 禁用CGO
SET GOOS=linux  // 目标平台是linux
SET GOARCH=amd64  // 目标处理器架构是amd64

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build # linux
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build # mac
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build # windows
```

```bash
go build -o obcsapi-picgo.linux  obcsapi-picgo.go
GOOS=windows GOARCH=amd64  go build  -o obcsapi-picgo.exe  obcsapi-picgo.go
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o obcsapi-picgo.linux obcsapi-picgo.go
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o obcsapi-picgo.mac obcsapi-picgo.go


upx obcsapi-picgo.exe -o obcsapi-picgo.small.exe 
upx -9 obcsapi-picgo.exe -o obcsapi-picgo.small.exe 
upx --ultra-brute obcsapi-picgo.exe -o obcsapi-picgo.small.exe 
```