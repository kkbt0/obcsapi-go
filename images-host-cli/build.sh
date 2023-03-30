set GOOS windows
go build -ldflags "-s -w" -o obcsapi-picgo.exe obcsapi-picgo.go
set GOOS linux
go build -ldflags "-s -w" -o obcsapi-picgo.linux obcsapi-picgo.go
set GOOS darwin
go build -ldflags "-s -w" -o obcsapi-picgo.mac obcsapi-picgo.go

upx -9 obcsapi-picgo.exe -o obcsapi-picgo.small.exe
upx -9 obcsapi-picgo.linux -o obcsapi-picgo.samll.linux
upx -9 obcsapi-picgo.mac -o obcsapi-picgo.samll.mac
