package main

import (
	_ "embed"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IndexInfo struct {
	Title   string
	Content string
}

//go:embed template/index.html
var indeHtml string

func main() {
	ShowConfig() // 打印基础消息

	// f, _ := os.Create("gin.log")
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout) // 日志写入文件和控制台
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Token")
	r.Use(cors.New(config)) // cors 配置

	r.GET("/", IndexHandler)                                // index.html
	r.GET("/404", BaseHandler)                              // 404
	r.GET("/time", Greet)                                   // 打招呼 测试使用 GET
	r.POST("/token", VerifyToken1Handler)                   // 验证 Token1 有效性
	r.Any("/api/wechat", WeChatMpHandlers)                  // wecheet 机器人 用于公众测试号
	r.GET("/api/sendtoken2mail", SendTokenHandler, limit()) // 请求将 token发送到 email GET 请求

	obGroup := r.Group("/ob")
	{
		obGroup.Any("today", Token1AuthMiddleware(), ObTodayHandler)             // Obsidian Token1 GET/POST 今日日记
		obGroup.POST("today/all", Token1AuthMiddleware(), ObPostTodayAllHandler) // Obsidian Token1 POST 整片修改今日日记
		obGroup.GET("recent", Token1AuthMiddleware(), ObGet3DaysHandler)         // Obsidian Token1 GET 近三天日记

		obGroup.POST("moonreader", Token2AuthMiddleware(), MoodReaderHandler) // Obsidian Token2 POST 静读天下 api
		obGroup.POST("fv", Token2AuthMiddleware(), fvHandler)                 // Obsidian Token2 POST 安卓 FV 悬浮球 快捷存储 文字，图片
		obGroup.POST("sr/webhook", Token2AuthMiddleware(), SRWebHook)         // Obsidian Token2 POST 简悦 Webhook 使用
		obGroup.POST("general", Token2AuthMiddleware(), GeneralHeader)        // Obsidian Token2 POST 通用接口 今日日记
		obGroup.POST("url", Token2AuthMiddleware(), Url2MdHandler)            // Obsidian Token2 POST 页面转 md 存储 效果很一般 不如简悦
	}
	r.Run(fmt.Sprintf("%s:%s", ConfigGetString("host"), ConfigGetString("port"))) // 运行服务
}

var limiter = rate.NewLimiter(0.00001, 1) // 限制 token 发送到 email (0.01 ,1) 意思 100 秒，允许 1 次。

// 短时间多次请求限制
func limit() func(c *gin.Context) {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			http.Error(c.Writer, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			c.Abort()
			return
		}
		c.Next()
	}
}

func Token1AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		if !VerifyToken1(c.Request.Header.Get("Token")) {
			c.JSON(401, gin.H{
				"code": 401,
				"msg":  "验证错误",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func Token2AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		if !VerifyToken2(c.Request.Header.Get("Token")) {
			c.JSON(401, gin.H{
				"code": 401,
				"msg":  "验证错误",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
