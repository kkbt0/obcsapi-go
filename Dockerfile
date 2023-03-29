FROM golang:1.18 as builder
MAINTAINER 恐咖兵糖<0@ftls.xyz>

ENV GO111MODULE=on 
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /home/workspace
COPY . .
RUN go build -o server  -ldflags '-linkmode "external" -extldflags "-static"' .
RUN mkdir app && cp server app/ && cp docker-entrypoint.sh app/

FROM alpine:latest
MAINTAINER 恐咖兵糖<0@ftls.xyz>

ENV VERSION 4.0.5
ENV GIN_MODE release

WORKDIR /app
COPY --from=builder /home/workspace/app/ .
RUN chmod +x docker-entrypoint.sh
RUN pwd && ls
EXPOSE 8900

ENTRYPOINT ["/app/docker-entrypoint.sh"]

