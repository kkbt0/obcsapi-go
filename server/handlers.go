package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"obcsapi-go/auth"
	"obcsapi-go/dao"
	"obcsapi-go/gr"
	"obcsapi-go/skv"
	"obcsapi-go/tools"
	"obcsapi-go/wechat"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/net/webdav"
)

func IndexHandler(c *gin.Context) {
	if viper.GetString("web_url_full") != "" {
		c.Redirect(307, viper.GetString("web_url_full"))
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/web/")
}

func NotFoundHandler(c *gin.Context) {
	gr.ErrNotFound(c)
}

// 一个简易图床
func ImagesHostUplaodHanler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		gr.ErrServerError(c, err)
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

	// choice mode

	// read the file into buffer
	uploadedFile, _ := file.Open()
	defer uploadedFile.Close()
	buffer := make([]byte, file.Size)
	_, _ = uploadedFile.Read(buffer)

	// Bd ocr
	if config.UsBdOcr {
		switch typeName {
		case ".jpg", ".jpeg", ".png", ".bmp":
			ans, err := tools.BdGeneralBasicOcr(buffer)
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
			tools.Debug(ans)
		default:
			log.Println("UnSupported file type: ", typeName)
		}
	}
	// END ocr

	// store
	rUrl := fmt.Sprintf("%s%s", config.BaseURL, strings.Join(filePath, ""))
	switch tools.NowRunConfig.ImageHosting.StorageMode {
	case "local":
		err = c.SaveUploadedFile(file, "./webdav/images/"+strings.Join(filePath, ""))
	case "obsidian":
		err = dao.ObjectStore(strings.Join(filePath, ""), buffer)
		rUrl = strings.Join(filePath, "")
	case "s3":
		rUrl, err = tools.S3FileStore("images/"+strings.Join(filePath, ""), buffer)
	default:
		err = c.SaveUploadedFile(file, "./webdav/images/"+strings.Join(filePath, ""))
	}
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}

	gr.RJSON(c, nil, 200, 200, "", gin.H{
		"url":  rUrl,
		"url2": rUrl,
	})
}

type VersionResponseJosn struct {
	Code          int       `json:"code"`
	ServerVersion string    `json:"server_version"`
	ConfigVersion string    `json:"config_version"`
	ServerTime    time.Time `json:"server_time"`
	Msg           string    `json:"msg"`
}

// @Summary 服务器信息与测试接口
// @Produce json
// @Router /info [get]
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
		gr.ErrNotFound(c)
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
// @Summary JWT 测试接口
// @Tags 前端
// @Security JWT
// @Accept plain,octet-stream
// @Produce json
// @Router /api/v1/sayHello [get]
func JwtHello(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	claims, _ := auth.ParseToken(authToken)
	log.Println(claims)
	gr.RJSON(c, nil, 200, 200, "hello", gr.H{})
}

func MailTesterHandler(c *gin.Context) {
	err := tools.SendMail("测试邮件", "测试内容")
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}
	gr.Success(c)
}

// 根据 ?fileKey=xxx.jpg 获取文件
func ObFileHanlder(c *gin.Context) {
	// 处理验证
	accessToken := c.Query("accessToken")
	if accessToken != tools.ObFileAccessToken() {
		gr.ErrAuth(c)
		return
	}

	// 获取文件部分
	fileKey := c.Query("fileKey")
	if fileKey == "" {
		gr.ErrNotFound(c)
		return
	}
	data, err := dao.GetObject(fileKey)
	if err != nil || data == nil {
		log.Println(err)
		gr.ErrNotFound(c)
		return
	}
	c.Writer.Write(data)
}

type WeChatInfoStruct struct {
	Content string `json:"content"`
}

// @Summary 微信通知
// @Tags 通知
// @Security Token
// @Accept json
// @Produce json
// @Param json body WeChatInfoStruct true "WeChatInfoStruct"
// @Router /api/wechatmpmsg [post]
func WeChatMpInfoHandler(c *gin.Context) {
	var weChatInfoStruct WeChatInfoStruct
	if c.ShouldBindJSON(&weChatInfoStruct) != nil {
		gr.ErrBindJSONErr(c)
		return
	}
	if weChatInfoStruct.Content == "" {
		gr.ErrEmpty(c)
		return
	}
	err := wechat.WeChatTemplateMesseng(weChatInfoStruct.Content)
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}
	gr.Success(c)
}

// @Summary 企业微信通知
// @Tags 通知
// @Security Token
// @Accept json
// @Produce json
// @Param json body WeChatInfoStruct true "WeChatInfoStruct"
// @Router /api/workwechatmsg [post]
func WorkWeChatMpInfoHandler(c *gin.Context) {
	var weChatInfoStruct WeChatInfoStruct
	if c.ShouldBindJSON(&weChatInfoStruct) != nil {
		gr.ErrBindJSONErr(c)
		return
	}
	if weChatInfoStruct.Content == "" {
		gr.ErrEmpty(c)
		return
	}
	err := wechat.WorkWechatSendText(weChatInfoStruct.Content)
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}
	gr.Success(c)
}

type SendMailStruct struct {
	Subject string `json:"subject"`
	Content string `json:"content"`
}

// @Summary 邮件通知
// @Tags 通知
// @Accept json
// @Security Token
// @Produce json
// @Param json body SendMailStruct true "SendMailStruct"
// @Router /api/sendmail [post]
func SendMailHandler(c *gin.Context) {
	var sendMailStruct SendMailStruct
	if c.ShouldBindJSON(&sendMailStruct) != nil {
		gr.ErrBindJSONErr(c)
		return
	}
	if sendMailStruct.Content == "" || sendMailStruct.Subject == "" {
		gr.ErrEmpty(c)
		return
	}
	err := tools.SendMail(sendMailStruct.Subject, sendMailStruct.Content)
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}
	gr.Success(c)
}

func GetMentionHandler(c *gin.Context) {
	mention := tools.NowRunConfig.Mention
	tools.Debug(mention)
	c.JSON(200, mention)
}

func UpdateBdAccessTokenHandler(c *gin.Context) {
	accessToken, err := tools.BdGetAccessToken(tools.NowRunConfig.BdOcr.ApiKey, tools.NowRunConfig.BdOcr.ApiSecret)
	if err != nil {
		gr.ErrServerError(c, err)
	}
	msg := "请求成功，但出现错误"
	if accessToken.AccessToken != "" {
		tools.NowRunConfig.ImageHosting.BdOcrAccessToken = accessToken.AccessToken
		err := tools.UpdateConfig(tools.NowRunConfig)
		if err != nil {
			gr.ErrServerError(c, err)
			return
		}
		msg = "请求并更新成功"
	}
	gr.RJSON(c, nil, 200, 200, msg, gin.H{})
}

func UpdateViperHandler(c *gin.Context) {
	err := tools.UpdateViper()
	if err != nil {
		gr.ErrServerError(c, fmt.Errorf("error: Update Error config file: %s \n ", err))
		return
	}
	gr.Success(c)
}

type RandomMemosStruct struct {
	FileKey       string `json:"file_key"`
	MemosText     string `json:"memos_text"`
	MemosShowText string `json:"memos_show_text"`
}

// @Summary 随机回顾
// @Tags 前端
// @Security JWT
// @Router /api/v1/random [get]
func RandomMemosHandler(c *gin.Context) {
	fileKey, memosText := RandomMemos()
	if memosText == "" {
		fileKey, memosText = RandomMemos()
	}
	c.JSON(200, RandomMemosStruct{
		MemosText:     memosText,
		MemosShowText: dao.MdShowText(fileKey, memosText),
		FileKey:       fileKey,
	})
}

type KvSerchPost struct {
	Key string `json:"key"`
}

func KvSerchHandler(c *gin.Context) {
	var kvSerchPost KvSerchPost
	if c.ShouldBindJSON(&kvSerchPost) != nil {
		gr.ErrBindJSONErr(c)
		return
	}
	result, err := skv.KvSerch(kvSerchPost.Key)
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}
	c.JSON(200, result)
}

// func FormPostHandler(c *gin.Context) {
// 	bodyBytes, err := io.ReadAll(c.Request.Body)
// 	if err != nil {
// 		gr.ErrServerError(c, err)
// 		return
// 	}
// 	bodyString := string(bodyBytes)
// 	tools.Debug(bodyString)
// 	result, err := command.RunJsByFile("script/demo.js", bodyString)
// 	if err != nil {
// 		gr.ErrServerError(c, err)
// 		return
// 	}
// 	c.String(200, result)
// }
