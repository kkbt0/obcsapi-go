package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"obcsapi-go/dao"
	"obcsapi-go/jwt"
	"obcsapi-go/skv"
	"obcsapi-go/tools"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"
)

func IndexHandler(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/web/")
}

func Greet(c *gin.Context) {
	c.String(200, "Hello World! %s", time.Now())
}

func BaseHandler(c *gin.Context) {
	c.String(404, "404")
}

// 一个简易图床
func ImagesHostUplaodHanler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	config := tools.NowRunConfig.ImageHosting
	log.Println("ImagesHost Upload:", file.Filename, "=>", tools.ReplaceUnAllowedChars(file.Filename))
	// filePath: /webdav/images/202303/test.jpg
	typeName := path.Ext(file.Filename)
	filePath := []string{tools.TimeFmt(config.Prefix)}
	if config.UseRawName {
		filePath = append(filePath, strings.TrimSuffix(tools.ReplaceUnAllowedChars(file.Filename), typeName)) // 200601/test
	}
	if config.RandomCharLength != 0 {
		filePath = append(filePath, "_"+tools.RandomString(config.RandomCharLength)) // 200601/test_e5md1
	}
	filePath = append(filePath, typeName) // 200601/test_e5md1.jpg
	c.SaveUploadedFile(file, "./webdav/images/"+strings.Join(filePath, ""))
	// Bd ocr
	if config.BdOcrAccessToken != "" {
		switch typeName {
		case ".jpg", ".jpeg", ".png", ".bmp":
			ans, err := tools.BdGeneralBasicOcr("./webdav/images/" + strings.Join(filePath, ""))
			if err != nil {
				c.Error(err)
			}
			var textList []string
			for _, v := range ans {
				textList = append(textList, v.Words)
			}
			inMdText := fmt.Sprintf("%s%s", config.BaseURL, strings.Join(filePath, ""))
			inMdText += fmt.Sprintf("\n%s\n\n---\n", strings.Join(textList, "\n"))
			dao.TextAppend(tools.NowRunConfig.OtherDataDir()+"OcrData/bdocr-"+tools.TimeFmt("200601")+".md", inMdText)
			fmt.Println(ans)
		default:
			log.Println("UnSupported file type: ", typeName)
		}
	}
	// END ocr
	c.JSON(200, gin.H{
		"data": gin.H{
			"url":  fmt.Sprintf("%s%s", config.BaseURL, strings.Join(filePath, "")),
			"url2": fmt.Sprintf("%s%s", config.BaseURL, strings.Join(filePath, "")),
		},
	})
}

type VersionResponseJosn struct {
	Code          int       `json:"code"`
	ServerVersion string    `json:"server_version"`
	ConfigVersion string    `json:"config_version"`
	ServerTime    time.Time `json:"server_time"`
	Msg           string    `json:"msg"`
}

func InfoHandler(c *gin.Context) {
	c.JSON(200, VersionResponseJosn{
		Code:          200,
		ServerVersion: version,
		ConfigVersion: tools.ConfigGetString("version"),
		ServerTime:    time.Now(),
		Msg:           "道可道，非常道；名可名，非常名。无，名天地之始；有，名万物之母。故常无，欲以观其妙；常有，欲以观其徼。此两者，同出而异名，同谓之玄。玄之又玄，众妙之门。",
	})
}

// Obsidian 公开文档功能
func ObsidianPublicFiles(c *gin.Context) {
	fileName := c.Param("fileName")
	fileName = tools.NowRunConfig.OtherDataDir() + "Public" + fileName
	md := skv.GetByFileKey(fileName)
	go UpdateSkvCache(fileName)
	if md == "" {
		c.String(404, "Not Found")
		return
	}
	if c.Query("raw") == "true" {
		c.String(200, md)
	}
	c.HTML(200, "markdown.html", gin.H{
		"title":    c.Param("fileName"),
		"markdown": md,
	})
}

func UpdateSkvCache(fileName string) {
	if limitPublicPage.Allow() {
		log.Println("Updating", fileName, "...")
		err := skv.PutByFileKey(fileName)
		if err != nil {
			log.Println(err)
		}
	}
}

// WebDavServe
func WebDavServe(prefix string, rootDir string,
	validator func(c *gin.Context) bool) gin.HandlerFunc {
	logger := func(req *http.Request, err error) {
		if err != nil {
			log.Println("[WebDAV]", req.URL.Path, "Err:", err)
		}
	}
	w := webdav.Handler{
		Prefix:     prefix,
		FileSystem: webdav.Dir(rootDir),
		LockSystem: webdav.NewMemLS(),
		Logger:     logger,
	}
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, w.Prefix) {
			if validator != nil && !validator(c) {
				c.AbortWithStatus(404)
				return
			}
			c.Status(200) // 200 by default, which may be override later
			w.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

func WebDavServeAuth(c *gin.Context) bool {
	if !tools.NowRunConfig.Webdav.Server {
		return false
	}
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.Writer.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		c.Status(http.StatusUnauthorized) // 401
	}
	if username != tools.NowRunConfig.Webdav.Username || password != tools.NowRunConfig.Webdav.Password {
		http.Error(c.Writer, "WebDAV: need authorized!", http.StatusUnauthorized)
		return false
	}
	return true
}

// JWT API v1
// Hello
func JwtHello(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	claims, _ := jwt.ParseToken(auth)
	log.Println(claims)
	c.String(http.StatusOK, "hello")
}

func MailTesterHandler(c *gin.Context) {
	err := SendMail("测试邮件", "测试内容")
	if err != nil {
		c.Error(err)
		c.String(500, err.Error())
		return
	}
	c.String(200, "Successfully Send")
}

// 根据 ?fileKey=xxx.jpg 获取文件
func ObFileHanlder(c *gin.Context) {
	// 处理验证
	accessToken := c.Query("accessToken")
	if accessToken != tools.ObFileAccessToken() {
		c.Status(401)
		return
	}

	// 获取文件部分
	fileKey := c.Query("fileKey")
	if fileKey == "" {
		c.Status(404)
		return
	}
	data, err := dao.GetObject(fileKey)
	if err != nil || data == nil {
		log.Println(err)
		c.Status(404)
		return
	}
	c.Writer.Write(data)
}

type WeChatInfoStruct struct {
	Content string `json:"content"`
}

func WeChatMpInfoHandler(c *gin.Context) {
	var weChatInfoStruct WeChatInfoStruct
	if c.ShouldBindJSON(&weChatInfoStruct) != nil {
		c.String(400, "参数错误")
		return
	}
	if weChatInfoStruct.Content == "" {
		c.String(400, "参数错误")
		return
	}
	err := WeChatTemplateMesseng(weChatInfoStruct.Content)
	if err != nil {
		c.Status(500)
		c.Error(err)
		return
	}
	c.Status(200)
}

type SendMailStruct struct {
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func SendMailHandler(c *gin.Context) {
	var sendMailStruct SendMailStruct
	if c.ShouldBindJSON(&sendMailStruct) != nil {
		c.String(400, "参数错误")
		return
	}
	if sendMailStruct.Content == "" || sendMailStruct.Subject == "" {
		c.String(400, "参数错误")
		return
	}
	err := SendMail(sendMailStruct.Subject, sendMailStruct.Content)
	if err != nil {
		c.Status(500)
		c.Error(err)
		return
	}
	c.Status(200)
}

func GetMentionHandler(c *gin.Context) {
	mention := tools.NowRunConfig.Mention
	tools.Debug(mention)
	c.JSON(200, mention)
}
