package main

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

type IndexInfo struct {
	Title   string
	Content string
}

func setupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Token")
}

func main() {
	ShowConfig() // 打印基础消息

	http.Handle("/", http.FileServer(http.Dir("./template")))
	http.HandleFunc("/404", BaseHandler)                                          // 404
	http.HandleFunc("/token", VerifyToken1Handler)                                // Token 验证
	http.HandleFunc("/api/wechat", wechatmpfunc)                                  // wecheet 机器人 用于公众测试号
	http.HandleFunc("/ob/today", ob_today)                                        // Obsidian Token1
	http.HandleFunc("/ob/today/all", ob_today_all)                                // Obsidian Token1
	http.HandleFunc("/ob/recent", get_3_day)                                      // Obsidian Token1
	http.HandleFunc("/ob/moonreader", moodreader)                                 // Obsidian Token2
	http.HandleFunc("/ob/fv", fvHandler)                                          // Obsidian Token2
	http.HandleFunc("/time", greet)                                               // 打招呼 测试使用
	http.Handle("/api/sendtoken2mail", limit(http.HandlerFunc(SendTokenHandler))) // 请求将 token发送到 email
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
