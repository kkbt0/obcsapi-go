package main

import (
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log"
	_ "obcsapi-go/dao" // init 检查数据模式 是 S3， CouchDb ..
	"obcsapi-go/jwt"
	"obcsapi-go/talk"
	"obcsapi-go/tools"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

//go:embed templates
var files embed.FS

var limiter = rate.NewLimiter(0.00001, 1)     // 限制 token 发送到 email (0.01 ,1) 意思 100 秒，允许 1 次。用于 LimitMiddleware
var limitPublicPage = rate.NewLimiter(0.1, 1) // 公开文档限制
var loginLimter = rate.NewLimiter(0.1, 3)     // 登录速率限制

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

	tools.ReloadRunConfig() // 初始化运行配置

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Token", "Authorization")
	r.Use(cors.New(config)) // cors 配置

	r.GET("/", IndexHandler)      // index.html vue3 pwa
	r.Static("/web", "./website") // h5 静态文件

	r.GET("/404", BaseHandler)             // 404
	r.GET("/time", Greet)                  // 打招呼 测试使用 GET
	r.GET("/info", InfoHandler)            // Obcsapi info
	r.Any("/api/wechat", WeChatMpHandlers) // wechat 机器人 用于公众测试号

	apiGroup := r.Group("/api", TokenAuthMiddleware("./token/token2.json")) // default token2
	{
		apiGroup.GET("testtoken", Greet)                  // test token
		apiGroup.POST("wechatmpmsg", WeChatMpInfoHandler) // 公众测试号 模板消息通知
		apiGroup.POST("sendmail", SendMailHandler)        // 邮件消息通知

	}

	obGroup2 := r.Group("/ob", TokenAuthMiddleware("./token/token2.json")) // default token2
	{
		obGroup2.POST("fv", fvHandler)                 // Obsidian Token2 POST 安卓 FV 悬浮球 快捷存储 文字，图片
		obGroup2.POST("sr/webhook", SRWebHook)         // Obsidian Token2 POST 简悦 Webhook 使用
		obGroup2.POST("general", GeneralHeader)        // Obsidian Token2 POST 通用接口 今日日记
		obGroup2.POST("url", Url2MdHandler)            // Obsidian Token2 POST 页面转 md 存储 效果很一般 不如简悦
		obGroup2.POST("generalall", GeneralAllHandler) // Obsidian Token2 POST 通用接口 全部文件都可以
	}
	r.POST("/ob/general/*paramtoken", SpecialTokenMiddleware("./token/token2.json"), GeneralHeader2) // Token2 flomo like api
	r.POST("/ob/moonreader", StandardTokenAuthMiddleware("./token/token2.json"), MoodReaderHandler)  // Obsidian POST 静读天下 api 此 API 使用 Authorization 头验证

	r.GET("/public/*fileName", ObsidianPublicFiles) // Obsidian Public Files GET

	r.POST("login", LimitLoginMiddleware(), jwt.LoginHandler)

	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	r.POST("/api/upload", TokenAuthMiddleware("./token/token2.json"), ImagesHostUplaodHanler) //图床
	r.Static("/images", "./webdav/images")

	// Webdav
	r.Use(WebDavServe(
		"/webdav/",
		"./webdav",
		WebDavServeAuth,
	))

	r.Use(AllowOPTIONS())
	api1Group := r.Group("/api/v1", jwt.JWTAuth())
	{
		api1Group.GET("/sayHello", JwtHello)
		api1Group.GET("/daily", ObV1GetDailyHandler)                // 使用一周前有缓存的 daily
		api1Group.GET("/daily/nocache", ObV1GetDailyNoCacheHandler) // 使用没有缓存的 daily （每次都请求服务器）
		api1Group.POST("/line", ObV1PostLineHandler)                // 行修改
		api1Group.POST("/cacheupdate", ObV1UpdateCacheHandler)      // 更新缓存 ?key=1.md
		api1Group.POST("/upload", ImagesHostUplaodHanler)           // jwt 图床
		api1Group.GET("/config", tools.GetRunConfigHandler)         // 运行时 可修改配置
		api1Group.POST("/config", tools.PostConfigHandler)          //运行时 可修改配置
		api1Group.GET("/mailtest", MailTesterHandler)               // 邮件测试
		api1Group.POST("/talk", talk.TalkHandler)                   // 对话 API
	}

	r.GET("/ob/file", ObFileHanlder) // 需要带验证参数

	RunCronJob() //  运行定时任务

	r.Run(fmt.Sprintf("%s:%s", tools.ConfigGetString("host"), tools.ConfigGetString("port"))) // 运行服务
}
