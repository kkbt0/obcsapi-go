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

func setupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Token")
}

func main() {
	ShowConfig() // 打印基础消息
	http.Handle("/", http.FileServer(http.Dir("./template")))
	http.HandleFunc("/404", Logging(BaseHandler))                                 // 404
	http.HandleFunc("/token", Logging(VerifyToken1Handler))                       // Token 验证 测试使用
	http.HandleFunc("/api/wechat", Logging(wechatmpfunc))                         // wecheet 机器人 用于公众测试号
	http.HandleFunc("/ob/today", Logging(ob_today))                               // Obsidian Token1 GET/POST 今日日记
	http.HandleFunc("/ob/today/all", Logging(ob_today_all))                       // Obsidian Token1 POST 整片修改今日日记
	http.HandleFunc("/ob/recent", Logging(get_3_day))                             // Obsidian Token1 GET 近三天日记
	http.HandleFunc("/ob/moonreader", Logging(moodreaderHandler))                 // Obsidian Token2 POST 静读天下 api
	http.HandleFunc("/ob/fv", Logging(fvHandler))                                 // Obsidian Token2 POST 安卓 FV 悬浮球 快捷存储 文字，图片
	http.HandleFunc("/ob/sr/webhook", Logging(SRWebHook))                         // Obsidian Token2 POST 简悦 Webhook 使用
	http.HandleFunc("/ob/general", Logging(GeneralHeader))                        // Obsidian Token2 POST 通用接口 今日日记
	http.HandleFunc("/ob/url", Logging(Url2MdHandler))                            // Obsidian Token2 POST 页面转 md 存储
	http.HandleFunc("/time", Greet)                                               // 打招呼 测试使用 GET
	http.Handle("/api/sendtoken2mail", limit(http.HandlerFunc(SendTokenHandler))) // 请求将 token发送到 email GET 请求
	http.ListenAndServe(fmt.Sprintf("%s:%s", ConfigGetString("host"), ConfigGetString("port")), nil)
}

var limiter = rate.NewLimiter(0.00001, 1) // 限制 token 发送到 email (0.01 ,1) 意思 100 秒，允许 1 次。

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

func Logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		f(w, r)
	}
}
