package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"obcsapi-go/dao"
	. "obcsapi-go/dao"
	"obcsapi-go/tools"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gin-gonic/gin"
)

// Token2 静读天下使用的 API
// @Summary 静读天下使用的 API
// @Description 静读天下使用的 API，标注-设置-ReadWise 设置该路径和 token 值即可
// @Tags Ob
// @Security AuthorizationToken
// @Accept json
// @Produce plain
// @Param 划线和标注 body MoodReader true "MoodReader"
// @Router /ob/moonreader [post]
func MoodReaderHandler(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var moodReader MoodReader
	err := decoder.Decode(&moodReader)
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	fmt.Println(moodReader.Highlights[0])
	title := moodReader.Highlights[0].Title
	text := moodReader.Highlights[0].Text
	author := moodReader.Highlights[0].Author
	note := moodReader.Highlights[0].Note
	file_key := fmt.Sprintf("%sMoonReader/%s.md", tools.NowRunConfig.OtherDataDir(), tools.ReplaceUnAllowedChars(title))
	append_text := fmt.Sprintf("文: %s\n批: %s\n于: %s\n\n---\n", text, note, tools.TimeFmt("2006-01-02 15:04"))
	exist, _ := CheckObject(file_key)
	if exist {
		err = TextAppend(file_key, append_text)
	} else {
		yaml := fmt.Sprintf("---\ntitle: %s\nauthor: %s\n---\n书名: %s\n作者: %s\n简介: \n评价: \n\n---\n", title, author, title, author)
		err = TextAppend(file_key, yaml+append_text)
	}
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	c.String(200, "Success")
}

// @Summary fv 悬浮球使用的 API
// @Description 安卓软件 fv 悬浮球使用的 API 用于自定义任务的 图片、文字
// @Tags Ob
// @Security Token
// @Accept plain,octet-stream
// @Produce plain
// @Param 內容 body string true "fv payload 內容"
// @Router /ob/fv [post]
func fvHandler(c *gin.Context) {
	if c.GetHeader("Content-Type") == "text/plain" {
		content, _ := ioutil.ReadAll(c.Request.Body)
		DailyTextAppendMemos(string(content))
		c.String(200, "Success")
		return
	} else if c.GetHeader("Content-Type") == "application/octet-stream" {
		content, _ := ioutil.ReadAll(c.Request.Body)
		file_key := fmt.Sprintf("%s%s.jpg", tools.NowRunConfig.DailyAttachmentDir(), tools.TimeFmt("20060102150405"))
		ObjectStore(file_key, content)
		DailyTextAppendMemos(fmt.Sprintf("![](%s)", file_key))
		c.String(200, "Success")
		return
	}
	c.String(404, "Error")
}

// SimpRead WebHook Used
type SimpReadWebHookStruct struct {
	Title       string `json:"title"`
	Url         string `json:"url"`
	Description string `json:"desc"`
	Tags        string `json:"tags"`
	Content     string `json:"content"`
	Note        string `json:"note"`
}

// SimpRead WebHook Used Token2
// @Summary 简悦 WebHook 保存文章
// @Description SimpRead 简悦 WebHook POST 简悦 WebHook 保存文章
// @Tags Ob
// @Security Token
// @Accept json
// @Param json body SimpReadWebHookStruct true "SimpRead 简悦 POST"
// @Router /ob/sr/webhook [post]
func SRWebHook(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var simpReadJson SimpReadWebHookStruct
	err := decoder.Decode(&simpReadJson)
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	serverTime := tools.TimeFmt("200601021504")
	yaml := fmt.Sprintf("---\ntitle: %s\nsctime: %s\n---\n", simpReadJson.Title, serverTime)
	file_str := fmt.Sprintf("%s[[简悦WebHook生成]]\n生成时间: %s\n原文: %s\n标题: %s\n描述: %s\n标签: %s\n内容: \n%s", yaml, serverTime, simpReadJson.Url, simpReadJson.Title, simpReadJson.Description, simpReadJson.Tags, simpReadJson.Content)
	file_key := fmt.Sprintf("%sSimpRead/%s %s.md", tools.NowRunConfig.OtherDataDir(), tools.ReplaceUnAllowedChars(simpReadJson.Title), serverTime)
	MdTextStore(file_key, file_str)
}

// @Summary 通用 API 接口 Memos
// @Description 通用 API 接口,添加 Memos
// @Tags Ob
// @Security Token
// @Accept json
// @Produce plain
// @Param json body MemosData true "MemosData"
// @Router /ob/general [post]
func GeneralHeader(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var memosData MemosData
	err := decoder.Decode(&memosData)
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	err = DailyTextAppendMemos(memosData.Content)
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	c.String(200, "Success")
}

// @Summary 通用 API 接口 (Memos Flomo Like API)
// @Description 通用 API 接口,添加 Memos
// @Tags Ob
// @Accept json
// @Produce plain
// @Param token path string true "设定的 token 值"
// @Param json body MemosData true "MemosData"
// @Router /ob/general/{token} [post]
func GeneralHeader2(c *gin.Context) {
	fromMiddlewareTokenFilePath, exist := c.Get("tokenfilepath")
	if !exist {
		c.Status(500)
		return
	}
	rightTokenFilePath, ok := fromMiddlewareTokenFilePath.(string)
	if !ok {
		c.Status(500)
		return
	}
	tools.Debug("RightToken FilePath: ", rightTokenFilePath)

	if !tools.VerifyTokenByFilePath(rightTokenFilePath, c.Param("paramtoken")[1:]) {
		tools.Debug("ParamToken: ", c.Param("paramtoken"))
		c.Status(401)
		return
	}
	decoder := json.NewDecoder(c.Request.Body)
	var memosData MemosData
	err := decoder.Decode(&memosData)
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	err = DailyTextAppendMemos(memosData.Content)
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	c.String(200, "Success")

}

// Token2
// @Summary 裁剪网页
// @Description 裁剪网页
// @Tags Ob
// @Accept json
// @Produce plain
// @Security Token
// @Param json body UrlStruct true "MemosData"
// @Router /ob/url [post]
func Url2MdHandler(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var urlStruct UrlStruct
	err := decoder.Decode(&urlStruct)
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	text, err := tools.Downloader(urlStruct.Url)
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(string(text))
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	serverTime := tools.TimeFmt("200601021504")
	yaml := fmt.Sprintf("---\nurl: %s\nsctime: %s\n---\n", urlStruct.Url, serverTime)
	file_key := fmt.Sprintf("%sHtmlPages/%s.md", tools.NowRunConfig.OtherDataDir(), serverTime)
	err = MdTextStore(file_key, yaml+markdown)
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	c.Status(200)
}

// 通用接口

type GeneralAllStruct struct {
	Mod     string `json:"mod" enums:"append,cover"`
	FileKey string `json:"file_key" default:"test.md"`
	Content string `json:"content"`
}

// @Summary 通用 API 接口 All
// @Description 通用 API 接口，覆盖修改或增添所有文件。需要配置声明允许使用该接口
// @Tags Ob
// @Security Token
// @Accept json
// @Produce plain
// @Param json body GeneralAllStruct true "GeneralAllStruct"
// @Router /ob/generalall [post]
func GeneralPostAllHandler(c *gin.Context) {
	if tools.ConfigGetString("allow_general_all_post") != "true" {
		c.Status(404)
		return
	}
	var generalJson GeneralAllStruct
	if c.ShouldBindJSON(&generalJson) != nil {
		c.String(400, "参数错误")
		return
	}
	if generalJson.Content == "" {
		c.String(400, "参数错误")
		return
	}
	fileKey := tools.NowRunConfig.ObDaily.ObOtherDataDir + "通用接口/" + tools.TimeFmt("20060102150405.md")
	mod := "append" // cover append
	if generalJson.FileKey != "" {
		fileKey = generalJson.FileKey
	}
	if generalJson.Mod == "cover" {
		mod = generalJson.Mod
	}
	var err error
	if mod == "cover" {
		err = MdTextStore(fileKey, generalJson.Content)
	} else {
		err = TextAppend(fileKey, generalJson.Content)
	}
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	c.String(200, "Success")
}

// @Summary 通用 API 接口 All
// @Description 通用 API 接口，获取所有文件。需要配置声明允许使用该接口
// @Tags Ob
// @Security Token
// @Produce plain
// @Param filekey query string true "文件名，有路径，如 dir/text.md"
// @Router /ob/generalall [get]
func GeneralGetAllHandler(c *gin.Context) {
	if tools.ConfigGetString("allow_general_all_get") != "true" {
		c.Status(404)
		return
	}
	filekey := c.Query("filekey")
	text, err := dao.GetTextObject(filekey)
	if err != nil {
		c.Status(500)
		return
	}
	c.String(200, string(text))
}
