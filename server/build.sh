rm -rf server && rm -rf server.zip
go build -o server  -ldflags '-linkmode "external" -extldflags "-static"' .
zip server.zip server
ldd server 