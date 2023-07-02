#!/bin/bash
cd server
go mod tidy 
go mod vendor
version="4.2.6"
docker build -t kkbt/obcsapi:v$version .
docker build -t kkbt/obcsapi:latest .
docker save -o ob4.2.tar kkbt/obcsapi:v$version && gzip ob4.2.tar
bash build.sh $version
