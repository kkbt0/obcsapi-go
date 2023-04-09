rm -rf server && rm -rf server.zip
go build -o server  -ldflags '-linkmode "external" -extldflags "-static"' .
# CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc CGO_LDFLAGS=“-static” go build -o server 
zip server.zip server
ldd server 