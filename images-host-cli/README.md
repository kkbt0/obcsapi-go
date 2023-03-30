http://127.0.0.1:8900/api/upload
fQbzONJAAw
url

第一行是上传链接
第二行是 token2 的值
第三行可选 url or url2 。url 是 http ；url2 是 https

构建 
```
CGO_ENABLED=0 GOOS=windows GOARCH=amd64  go build  -o obcsapi-picgo.exe  obcsapi-picgo.go
go build -ldflags "-s -w" -o obcsapi-picgo.exe obcsapi-picgo.go
```