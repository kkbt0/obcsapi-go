package main

import (
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"obcsapi-go/auth"
	_ "obcsapi-go/dao" // init 检查数据模式 是 S3， CouchDb ..
	"obcsapi-go/docs"
	"obcsapi-go/talk"
	"obcsapi-go/tools"
	"os"

	_ "obcsapi-go/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

//go:embed templates
var files embed.FS

var limiter = rate.NewLimiter(0.00001, 1)     // 限制 token 发送到 email (0.01 ,1) 意思 100 秒，允许 1 次。用于 LimitMiddleware
var limitPublicPage = rate.NewLimiter(0.1, 1) // 公开文档限制
var loginLimter = rate.NewLimiter(0.1, 5)     // 登录速率限制

// @title Obcsapi
// @version v4.2.5 版本
// @description 基于 Obsidian S3 存储， CouchDb ，本地存储和 WebDAV 的后端 API ,可借助 Obsidian 插件 Remotely-Save 插件，或者 Self-hosted LiveSync (ex:Obsidian-livesync) 插件 CouchDb 方式同步，保存消息到 Obsidian 库。该调试页面大部分仅提供对 Headers-Token 验证方式的支持，其他如 Query-token，Headers-Authorization 除了特殊的几个其他并不支持。可以使用 https://hoppscotch.io/ 或者 Postman 之类的工具，或者使用 VsCode 插件 REST Client ，使用 REST Client 可以在 https://gitee.com/kkbt/obcsapi-go/tree/master/http ，找到测试文件
// @contact.url https://gitee.com/kkbt/obcsapi-go
// @externalDocs.url https://kkbt.gitee.io/obcsapi-go/
// @securityDefinitions.apikey Token
// @in header
// @name Token
// @securityDefinitions.apikey AuthorizationToken
// @in header
// @name Authorization
// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
func main() {
	ShowConfig() // 打印基础消息

	f, err := os.Create("gin.log")
	if err != nil {
		log.Println(err)
	}
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout) // 日志写入文件和控制台

	r := gin.Default()

	r.NoMethod(NotFoundHandler)
	r.NoRoute(NotFoundHandler)

	templ := template.Must(template.New("").ParseFS(files, "templates/*.html")) // 加载模板
	r.SetHTMLTemplate(templ)

	tools.ReloadRunConfig() // 初始化运行配置

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Token", "Authorization")
	r.Use(cors.New(config)) // cors 配置

	r.GET("/", IndexHandler)      // index.html vue3 pwa
	r.Static("/web", "./website") // h5 静态文件

	r.GET("/info", InfoHandler)            // Obcsapi info
	r.Any("/api/wechat", WeChatMpHandlers) // wechat 机器人 用于公众测试号

	apiGroup := r.Group("/api", TokenAuthMiddleware("./token/token2.json")) // default token2
	{
		apiGroup.POST("wechatmpmsg", WeChatMpInfoHandler) // 公众测试号 模板消息通知
		apiGroup.POST("sendmail", SendMailHandler)        // 邮件消息通知
	}

	obGroup2 := r.Group("/ob", TokenAuthMiddleware("./token/token2.json")) // default token2
	{
		obGroup2.POST("fv", fvHandler)          // Obsidian Token2 POST 安卓 FV 悬浮球 快捷存储 文字，图片
		obGroup2.POST("sr/webhook", SRWebHook)  // Obsidian Token2 POST 简悦 Webhook 使用
		obGroup2.POST("general", GeneralHeader) // Obsidian Token2 POST 通用接口 今日日记
		obGroup2.POST("url", Url2MdHandler)     // Obsidian Token2 POST 页面转 md 存储 效果很一般 不如简悦

		obGroup2.POST("generalall", GeneralPostAllHandler) // Obsidian Token2 POST 通用接口 全部文件都可以
		obGroup2.GET("generalall", GeneralGetAllHandler)   // Obsidian Token2 POST 通用接口 全部文件都可以

	}
	r.POST("/ob/general/*paramtoken", SpecialTokenMiddleware("./token/token2.json"), GeneralHeader2) // Token2 flomo like api
	r.POST("/ob/moonreader", StandardTokenAuthMiddleware("./token/token2.json"), MoodReaderHandler)  // Obsidian POST 静读天下 api 此 API 使用 Authorization 头验证

	r.GET("/public/*fileName", ObsidianPublicFiles) // Obsidian Public Files GET

	r.POST("login", LimitLoginMiddleware(), auth.LoginHandler)

	authGroup := r.Group("/auth", LimitLoginMiddleware())
	{
		authGroup.GET("/oauth2-set", auth.JWTAuth(), auth.InitiateOAuthFlowAndSetState)
		authGroup.GET("/oauth2-login", auth.InitiateOAuthFlow)
		authGroup.POST("/oauth2-login-bycode", auth.LoginByCodeHandler)
		authGroup.GET("/oauth2-callback", auth.HandleCallback)
	}

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
	api1Group := r.Group("/api/v1", auth.JWTAuth())
	{
		api1Group.GET("/sayHello", JwtHello)
		api1Group.GET("/daily", ObV1GetDailyHandler)                // 使用一周前有缓存的 daily
		api1Group.GET("/daily/nocache", ObV1GetDailyNoCacheHandler) // 使用没有缓存的 daily （每次都请求服务器）
		api1Group.POST("/line", ObV1PostLineHandler)                // 行修改
		api1Group.POST("/cacheupdate", ObV1UpdateCacheHandler)      // 更新缓存 ?key=1.md
		api1Group.POST("/upload", ImagesHostUplaodHanler)           // jwt 图床
		// 运行时 可修改配置 GET/POST
		api1Group.GET("/config", tools.GetRunConfigHandler)
		api1Group.POST("/config", ConfigAllow("lock_config", true), tools.PostConfigHandler)
		api1Group.GET("/config/reset", ConfigAllow("lock_config", true), tools.ResetRunConfigHandler)
		// 更新 Viper config.yaml
		api1Group.GET("/updateconfig", ConfigAllow("lock_config", true), UpdateViperHandler)
		api1Group.GET("/updatebd", UpdateBdAccessTokenHandler) // 更新 BD Access Token

		api1Group.GET("/mailtest", MailTesterHandler) // 邮件测试
		api1Group.POST("/talk", talk.TalkHandler)     // 对话 API
		api1Group.GET("/mention", GetMentionHandler)  // 提示词
		api1Group.GET("/random", RandomMemosHandler)  // 随机 Memos
		// 实验性质
		// 前端编辑器 列出文件 GET/POST 文件
		api1Group.GET("/list", ConfigAllow("experimental_features", false), LisFileHandler)
		api1Group.GET("/text", ConfigAllow("experimental_features", false), TextGetHandler)
		api1Group.POST("/text", ConfigAllow("experimental_features", false), TextPostHandler)
		api1Group.POST("/search", ConfigAllow("experimental_features", false), KvSerchHandler)
	}

	r.GET("/ob/file", ObFileHanlder) // 需要带验证参数

	docs.SwaggerInfo.BasePath = tools.ConfigGetString("backend_base_path")
	r.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerfiles.Handler, "OBCSAPI_SWAGGER_DISABLE"))

	RunCronJob() //  运行定时任务

	if viper.GetString("server_cert_path") != "" {
		r.RunTLS(
			fmt.Sprintf("%s:%s", viper.GetString("host"), viper.GetString("port")),
			viper.GetString("server_cert_path"), viper.GetString("server_key_path"))
	} else {
		r.Run(fmt.Sprintf("%s:%s", viper.GetString("host"), viper.GetString("port"))) // 运行服务
	}
}
