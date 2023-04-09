#!/bin/bash
rm -rf output
go build -o server  -ldflags '-linkmode "external" -extldflags "-static"' .
# CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc CGO_LDFLAGS=“-static” go build -o server 
# zip server.zip server
ldd server
mkdir output
cp config.example.yaml output/config.yaml
cp server output/
cp -R static/ output/
cp -R templates/ output/
cd output 
echo "Hello" > tem.txt
mkdir webdav
cd webdav
mkdir images
cd ..
zip -r -m obcsapi.$1.zip * 