go build -o demo-server  -ldflags '-linkmode "external" -extldflags "-static"' .


# nohup ./server  > log.log 2>&1 &
# nohup ./demo-server  > demo-server.log 2>&1 &