package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/time/rate"
)

type IndexInfo struct {
	Title   string
	Content string
}

func main() {
	ShowConfig() // 打印基础消息
	local_ip, _ := LocalIPv4s()
	log.Printf("LocalIp http://%s:%s\n", local_ip[0], ConfigGetString("port"))

	http.HandleFunc("/", BaseHandler)
	http.HandleFunc("/token", VerifyToken1Handler)
	http.Handle("/api/sendtoken2mail", limit(http.HandlerFunc(SendTokenHandler))) // 对请求将 token发送到 email 速率进行限制
	http.ListenAndServe(fmt.Sprintf("%s:%s", ConfigGetString("host"), ConfigGetString("port")), nil)
}

var limiter = rate.NewLimiter(0.3, 1)

// 短时间多次请求限制
func limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
