FROM golang:1.18 as builder
LABEL maintainer="恐咖兵糖<0@ftls.xyz>"

ENV GO111MODULE=on 
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /home/workspace
COPY . .
RUN go build -o server  -ldflags '-linkmode "external" -extldflags "-static"' . && strip -s server
RUN mkdir app && mkdir app/static/ && \
    cp server app/ && cp docker-entrypoint.sh app/ && \
    cp -R static/ app/ && \
    cp -R website/ app/ && \
    cp -R sh/ app/ 
    # cp zoneinfo.zip app/
    # docker-entrypoint.sh 运行后还会复制 static/ 到相应目录
FROM alpine:latest
LABEL maintainer="恐咖兵糖<0@ftls.xyz>"

ENV VERSION 4.2.1
ENV GIN_MODE release
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /home/workspace/app/ .
RUN chmod +x docker-entrypoint.sh && sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \ 
    apk add --no-cache tzdata bash curl jq
# bash curl about +5MB jq: a json cli
# RUN pwd && ls
VOLUME /app/data
EXPOSE 8900

ENTRYPOINT ["/app/docker-entrypoint.sh"]

