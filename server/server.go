package main

import (
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log"
	_ "obcsapi-go/dao" // init 检查数据模式 是 S3， CouchDb ..
	"obcsapi-go/tools"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

//go:embed static/index.html
var indeHtml string

//go:embed templates
var files embed.FS

var limiter = rate.NewLimiter(0.00001, 1) // 限制 token 发送到 email (0.01 ,1) 意思 100 秒，允许 1 次。用于 LimitMiddleware

func main() {
	ShowConfig() // 打印基础消息

	f, err := os.Create("gin.log")
	if err != nil {
		log.Println(err)
	}
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout) // 日志写入文件和控制台

	r := gin.Default()

	templ := template.Must(template.New("").ParseFS(files, "templates/*.html")) // 加载模板
	r.SetHTMLTemplate(templ)

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Token")
	r.Use(cors.New(config)) // cors 配置

	r.GET("/", IndexHandler)               // index.html
	r.GET("/404", BaseHandler)             // 404
	r.GET("/time", Greet)                  // 打招呼 测试使用 GET
	r.GET("/info", InfoHandler)            // Obcsapi info
	r.POST("/token", VerifyToken1Handler)  // 验证 Token1 有效性
	r.Any("/api/wechat", WeChatMpHandlers) // wecheet 机器人 用于公众测试号

	r.GET("/api/sendtoken2mail", LimitMiddleware(), SendTokenHandler) // 请求将 token发送到 email GET 请求

	obGroup1 := r.Group("/ob", Token1AuthMiddleware()) // 前端使用
	{
		obGroup1.Any("today", ObTodayHandler)             // Obsidian Token1 GET/POST 今日日记
		obGroup1.POST("today/all", ObPostTodayAllHandler) // Obsidian Token1 POST 整片修改今日日记
		obGroup1.GET("recent", ObGet3DaysHandler)         // Obsidian Token1 GET 近三天日记
	}
	obGroup2 := r.Group("/ob", Token2AuthMiddleware())
	{
		obGroup2.POST("fv", fvHandler)          // Obsidian Token2 POST 安卓 FV 悬浮球 快捷存储 文字，图片
		obGroup2.POST("sr/webhook", SRWebHook)  // Obsidian Token2 POST 简悦 Webhook 使用
		obGroup2.POST("general", GeneralHeader) // Obsidian Token2 POST 通用接口 今日日记
		obGroup2.POST("url", Url2MdHandler)     // Obsidian Token2 POST 页面转 md 存储 效果很一般 不如简悦
	}
	r.POST("/ob/moonreader", MoodReaderHandler)     // Obsidian Token2 POST 静读天下 api 此 API 使用 Authorization 头验证
	r.GET("/public/*fileName", ObsidianPublicFiles) // Obsidian Public Files GET

	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	r.POST("/api/upload", Token2AuthMiddleware(), ImagesHostUplaodHanler) //图床
	r.Static("/images", "./webdav/images")

	// Webdav
	r.Use(WebDavServe(
		"/webdav/",
		"./webdav",
		WebDavServeAuth,
	))

	RunCronJob() //  运行定时任务

	r.Run(fmt.Sprintf("%s:%s", tools.ConfigGetString("host"), tools.ConfigGetString("port"))) // 运行服务
}
