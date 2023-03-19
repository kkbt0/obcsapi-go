package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var guestInputText = []byte("")

//go:embed daily.md
var dailymd string

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Token")
	r.Use(cors.New(config)) // cors 配置

	r.GET("/ob/time", GreetHandler)                                    // 验证服务和服务器时间
	r.Any("/ob/today", TokenAuthMiddleware(), TodayHandler)            // Obsidian Token1 GET/POST 今日日记
	r.Any("/ob/today/all", TokenAuthMiddleware(), PostTodayAllHandler) // Obsidian Token1 POST 整片修改今日日记
	r.Any("/ob/recent", TokenAuthMiddleware(), Get3DaysHandler)

	r.Run(":3015") // 运行服务
}

func GreetHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 2000,
		"msg":  time.Now(),
	})
}

func TodayHandler(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		// 注意是数组
		c.JSON(200, []Daily{{
			Data:       string(guestInputText),
			MdShowData: string(guestInputText),
			Date:       timeFmt("2006-01-02"),
			ServerTime: timeFmt("2006-01-02-15-04"),
		}})
	case "POST":
		decoder := json.NewDecoder(c.Request.Body)
		var postJson PostJson
		err := decoder.Decode(&postJson)
		if err != nil {
			log.Println(err)
			c.Status(500)
			return
		}
		// EscapeString 预防 xss 个人使用不需要
		memos := fmt.Sprintf("\n- %s %s", timeFmt("15:04"), html.EscapeString(postJson.Content))
		guestInputText = append(guestInputText, []byte(memos)...)
		c.Status(200)
	case "OPTIONS":
		c.Status(200)
	default:
		c.Status(404)
	}
}

func PostTodayAllHandler(c *gin.Context) {
	if len(guestInputText) > 5000 {
		guestInputText = []byte{}
	}
	switch c.Request.Method {
	case "POST":
		decoder := json.NewDecoder(c.Request.Body)
		var postJson PostJson
		err := decoder.Decode(&postJson)
		if err != nil {
			log.Println(err)
			c.Status(500)
			return
		}
		// EscapeString 预防 xss 个人使用不需要
		guestInputText = []byte(html.EscapeString(postJson.Content))
		c.Status(200)
	case "OPTIONS":
		c.Status(200)
	default:
		c.Status(404)
	}
}

func Get3DaysHandler(c *gin.Context) {
	serverTime := timeFmt("2006-01-02-15-04")
	switch c.Request.Method {
	case "GET":
		c.JSON(200, []Daily{
			{Data: "No such file: 日志/1999-12-30.md",
				MdShowData: "No such file: 日志/1999-12-30.md",
				Date:       "1999-12-31",
				ServerTime: serverTime},
			{Data: dailymd,
				MdShowData: dailymd,
				Date:       "2000-01-01",
				ServerTime: serverTime},
			{Data: string(guestInputText),
				MdShowData: string(guestInputText),
				Date:       timeFmt("2006-01-02"),
				ServerTime: serverTime}})
	case "OPTIONS":
		c.Status(200)
	default:
		c.Status(404)
	}

}

// TokenAuthMiddleware 认证中间件
func TokenAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Token") != "yourtoken" {
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

type PostJson struct {
	Content string `json:"content"`
}
type Daily struct {
	Data       string `json:"data"`
	MdShowData string `json:"md_show_data"`
	Date       string `json:"date"`
	ServerTime string `json:"serverTime"`
}

// Time fmt eg 2006-01-02 15:04:05
func timeFmt(fmt string) string {
	// fmt.Println(time.Now().In(cstZone).Format("2006-01-02 15:04:05"))
	var cstZone = time.FixedZone("CST", 8*3600)
	return time.Now().In(cstZone).Format(fmt)
}
