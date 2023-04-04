package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"obcsapi-go/dao"
	"obcsapi-go/tools"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html")
	c.String(200, indeHtml)
}

func Greet(c *gin.Context) {
	c.String(200, "Hello World! %s", time.Now())
}

func BaseHandler(c *gin.Context) {
	c.String(404, "404")
}

// NewCaptcha 生成或更新 token 邮件发送登录链接 直接附带 token
func SendTokenHandler(c *gin.Context) {
	log.Println("Succeed Send Token")
	// 修改 Token1
	tools.ModTokenFile(tools.Token{TokenString: tools.GengerateToken(32), GenerateTime: tools.TimeFmt("2006-01-02 15:04:05")}, "token1")
	// 发送所有 Token 消息
	emailSendToken()
	c.String(200, "Succeed Send Token")
}

// 验证 Token 1 有效性

func VerifyToken1Handler(c *gin.Context) {
	// 解析 token json {"token":"sometoken1"}
	decoder := json.NewDecoder(c.Request.Body)
	var tokenFromJSON tools.TokenFromJSON
	err := decoder.Decode(&tokenFromJSON)
	if err != nil {
		fmt.Println("JSON Decoder Error:", err)
	}
	if tools.VerifyToken1(tokenFromJSON.TokenString) {
		c.String(200, "a right Token")
	} else {
		c.String(200, "a error Token")
	}
}

// 一个简易图床
func ImagesHostUplaodHanler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	log.Println("ImagesHost Upload:", file.Filename, "=>", tools.ReplaceUnAllowedChars(file.Filename))
	// filePath: /images/202303/test.jpg
	typeName := path.Ext(file.Filename)
	filePath := []string{tools.TimeFmt(tools.ConfigGetString("images_hosted_fmt"))}
	// filePath := fmt.Sprintf("%s%s", tools.TimeFmt(tools.ConfigGetString("images_hosted_fmt")), file.Filename)
	if tools.ConfigGetString("images_hosted_use_raw_name") == "true" {
		filePath = append(filePath, strings.TrimSuffix(tools.ReplaceUnAllowedChars(file.Filename), typeName)) // 200601/test
	}
	if tools.ConfigGetInt("images_hosted_random_name_length") != 0 {
		filePath = append(filePath, "_"+tools.RandomString(tools.ConfigGetInt("images_hosted_random_name_length"))) // 200601/test_e5md1
	}
	filePath = append(filePath, typeName) // 200601/test_e5md1.jpg
	c.SaveUploadedFile(file, "./images/"+strings.Join(filePath, ""))
	// Bd ocr
	if tools.ConfigGetString("bd_ocr_access_token") != "" {
		switch typeName {
		case ".jpg", ".jpeg", ".png", ".bmp":
			ans, err := tools.BdGeneralBasicOcr("./images/" + strings.Join(filePath, ""))
			if err != nil {
				c.Error(err)
			}
			var textList []string
			for _, v := range ans {
				textList = append(textList, v.Words)
			}
			inMdText := fmt.Sprintf("%s%s", tools.ConfigGetString("images_hosted_url"), strings.Join(filePath, ""))
			inMdText += fmt.Sprintf("\n%s\n\n---\n", strings.Join(textList, "\n"))
			dao.TextAppend(tools.ConfigGetString("ob_daily_other_dir")+"OcrData/bdocr.md", inMdText)
			fmt.Println(ans)
		default:
			log.Println("UnSupported file type: ", typeName)
		}
	}
	// END ocr
	c.JSON(200, gin.H{
		"data": gin.H{
			"url":  fmt.Sprintf("%s%s", tools.ConfigGetString("images_hosted_url"), strings.Join(filePath, "")),
			"url2": fmt.Sprintf("%s%s", tools.ConfigGetString("images_hosted_url"), strings.Join(filePath, "")),
		},
	})
}

func ObsidianPublicFiles(c *gin.Context) {
	fileName := c.Param("fileName")
	fileName = tools.ConfigGetString("ob_daily_other_dir") + "Public" + fileName
	md, err := dao.GetTextObject(fileName)
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	c.HTML(200, "markdown.html", gin.H{
		"title":    c.Param("fileName"),
		"markdown": md,
	})
}
