#!/bin/bash
cd server
go mod tidy 
go mod vendor
version="4.2.9"
# docker buildx build -t kkbt/obcsapi:v$version --platform=linux/arm,linux/arm64,linux/amd64 . --push
# docker buildx build -t kkbt/obcsapi:latest --platform=linux/arm,linux/arm64,linux/amd64 . --push
docker build -t kkbt/obcsapi:v$version .
docker build -t kkbt/obcsapi:latest . 
# docker buildx build -t kkbt/obcsapi:latest --platform=linux/arm,linux/arm64,linux/amd64 . --push
docker save -o ob4.2.tar kkbt/obcsapi:v$version && gzip ob4.2.tar
bash build.sh $version
