package main

import (
	"encoding/json"
	"fmt"
	"io"
	"obcsapi-go/dao"
	. "obcsapi-go/dao"
	"obcsapi-go/gr"
	"obcsapi-go/skv"
	"obcsapi-go/tools"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// @Summary 通用 API 接口 Memos
// @Description 通用 API 接口,添加 Memos
// @Tags Ob
// @Security Token
// @Accept json
// @Produce json
// @Param json body MemosData true "MemosData"
// @Router /ob/general [post]
func GeneralHeader(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var memosData MemosData
	var err error
	if err = decoder.Decode(&memosData); err != nil {
		gr.ErrServerError(c, err)
		return
	}
	if err = AppendDailyMemos(memosData.Content); err != nil {
		gr.ErrServerError(c, err)
		return
	}
	gr.Success(c)
}

// @Summary 通用 API 接口 (Memos Flomo Like API)
// @Description 通用 API 接口,添加 Memos
// @Tags Ob
// @Accept json
// @Produce json
// @Param token path string true "设定的 token 值"
// @Param json body MemosData true "MemosData"
// @Router /ob/general/{token} [post]
func GeneralHeader2(c *gin.Context) {
	fromMiddlewareTokenFilePath, exist := c.Get("tokenfilepath")
	if !exist {
		gr.ErrServerError(c, nil)
		return
	}
	rightTokenFilePath, ok := fromMiddlewareTokenFilePath.(string)
	if !ok {
		gr.ErrServerError(c, nil)
		return
	}
	tools.Debug("RightToken FilePath: ", rightTokenFilePath)

	if !tools.VerifyTokenByFilePath(rightTokenFilePath, c.Param("paramtoken")[1:]) {
		tools.Debug("ParamToken: ", c.Param("paramtoken"))
		gr.ErrAuth(c)
		return
	}
	decoder := json.NewDecoder(c.Request.Body)
	var memosData MemosData
	var err error
	if err = decoder.Decode(&memosData); err != nil {
		gr.ErrServerError(c, err)
		return
	}
	if err = AppendDailyMemos(memosData.Content); err != nil {
		gr.ErrServerError(c, err)
		return
	}
	gr.Success(c)
}

// Token2
// @Summary 裁剪网页
// @Description 裁剪网页
// @Tags Ob
// @Accept json
// @Produce json
// @Security Token
// @Param json body UrlStruct true "MemosData"
// @Router /ob/url [post]
func Url2MdHandler(c *gin.Context) {
	var err error
	decoder := json.NewDecoder(c.Request.Body)
	var urlStruct UrlStruct
	if err := decoder.Decode(&urlStruct); err != nil {
		gr.ErrServerError(c, err)
		return
	}
	var text []byte

	if text, err = tools.Downloader(urlStruct.Url); err != nil {
		gr.ErrServerError(c, err)
		return
	}
	converter := md.NewConverter("", true, nil)
	var markdown string
	if markdown, err = converter.ConvertString(string(text)); err != nil {
		gr.ErrServerError(c, err)
		return
	}

	title := strings.Split(markdown, "\n")[0]
	serverTime := tools.TimeFmt("200601021504")
	yaml := fmt.Sprintf("---\nurl: %s\ntitle: %s\nsctime: %s\n---\n[[ObSavePage]]\n", urlStruct.Url, tools.ReplaceUnAllowedChars(strings.TrimSpace(title)), serverTime)
	file_key := fmt.Sprintf("%sHtmlPages/%s %s.md", tools.NowRunConfig.OtherDataDir(), serverTime, tools.ReplaceUnAllowedChars(strings.TrimSpace(title)))
	if err = CoverStoreTextFile(file_key, yaml+markdown); err != nil {
		gr.ErrServerError(c, err)
		return
	}
	gr.RJSON(c, nil, 200, 200, "", gin.H{
		"file_key": file_key,
		"title":    title,
	})
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
// @Produce json
// @Param json body GeneralAllStruct true "GeneralAllStruct"
// @Router /ob/generalall [post]
func GeneralPostAllHandler(c *gin.Context) {
	if !viper.GetBool("allow_general_all_post") {
		gr.ErrNotFound(c)
		return
	}
	var generalJson GeneralAllStruct
	if c.ShouldBindJSON(&generalJson) != nil {
		gr.ErrBindJSONErr(c)
		return
	}
	if generalJson.Content == "" {
		gr.ErrEmpty(c)
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
		err = CoverStoreTextFile(fileKey, generalJson.Content)
	} else {
		err = AppendText(fileKey, generalJson.Content)
	}
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}
	gr.Success(c)
}

// @Summary 通用 API 接口 All
// @Description 通用 API 接口，获取所有文件。需要配置声明允许使用该接口
// @Tags Ob
// @Security Token
// @Produce plain
// @Param filekey query string true "文件名，有路径，如 dir/text.md"
// @Router /ob/generalall [get]
func GeneralGetAllHandler(c *gin.Context) {
	if !viper.GetBool("allow_general_all_get") {
		gr.ErrNotFound(c)
		return
	}
	filekey := c.Query("filekey")
	text, err := dao.GetFileText(filekey)
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}
	c.String(200, string(text))
}

// @Summary Today Daily Get 今日日志获取
// @Description Today Daily Get 今日日志获取 注意：每天凌晨 00:00 - 03:59  判断为 today daily 为 昨天的日志
// @Tags Ob
// @Security Token
// @Accept plain
// @Produce plain
// @Router /ob/today [get]
func ObGetTodayDailyHandler(c *gin.Context) {
	if mdText, err := GetFileText(tools.NowRunConfig.DailyFileKeyMore(ObTodayAddDateNum())); err != nil {
		gr.ErrServerError(c, err)
	} else {
		c.String(200, mdText)
	}
}

// @Summary Today Daily Put 今日日志覆写
// @Description Today Daily Put 完全覆盖内容 注意：每天凌晨 00:00 - 03:59  判断为 today daily 为 昨天的日志
// @Tags Ob
// @Security Token
// @Accept plain
// @Produce plain
// @Param 內容 body string true "完全覆盖 内容"
// @Router /ob/today [put]
func ObTodayPutHandler(c *gin.Context) {
	fileKey := tools.NowRunConfig.DailyFileKeyMore(ObTodayAddDateNum())
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}
	skv.PutFile(fileKey, string(bodyBytes))
	if err = CoverStoreTextFile(fileKey, string(bodyBytes)); err != nil {
		gr.ErrServerError(c, err)
		return
	}
	gr.Success(c)
}

// @Summary Today Daily Post 今日日志新增
// @Description Today Daily Post 新增内容，末尾添加 注意：每天凌晨 00:00 - 03:59  判断为 today daily 为 昨天的日志
// @Tags Ob
// @Security Token
// @Accept plain
// @Produce plain
// @Param 內容 body string true "新增内容，末尾添加"
// @Router /ob/today [post]
func ObTodayPostHandler(c *gin.Context) {
	fileKey := tools.NowRunConfig.DailyFileKeyMore(ObTodayAddDateNum())
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}
	if err = AppendText(fileKey, string(bodyBytes)); err != nil {
		gr.ErrServerError(c, err)
		return
	}
	gr.Success(c)
}

// 每天凌晨 00:00 - 03:59  判断为 today daily 为 昨天的日志
func ObTodayAddDateNum() int {
	hour := time.Now().Hour()
	if hour >= 0 && hour <= 3 {
		return -1
	}
	return 0
}
