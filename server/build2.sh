#!/bin/bash
export GOOS=windows
export GOARCH=amd64

go build -o server.exe -ldflags "-s -w" -tags netgo -installsuffix netgo .