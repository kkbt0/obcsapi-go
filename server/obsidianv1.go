package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"obcsapi-go/dao"
	. "obcsapi-go/dao"
	"obcsapi-go/gr"
	"obcsapi-go/skv"
	"obcsapi-go/tools"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ObDailyV1 struct {
	Date       string   `json:"date"`
	MdText     []string `json:"md_text"`
	MdShowText []string `json:"md_show_text"`
	// MdElements [][]md.Element `json:"md_elements"`
}

// 分割 md 文本 便于根据行号修改
func MarkdownSpilter(text string) []string {
	arr := strings.Split(text, "\n")
	var last string
	var result []string
	for _, s := range arr {
		if strings.HasPrefix(s, "- ") {
			result = append(result, s)
		} else {
			if last == "" {
				last = s
			} else {
				last += s
			}
			if len(result) > 0 {
				result[len(result)-1] += "\n" + last
				last = ""
			} else {
				result = append(result, last)
				last = ""
			}
		}
	}
	return result
}

// ?day=-1 yesterday daily
// @Summary Memos 请求
// @Description 默认一周以前的查找缓存返回 即 <= -7 且不允许请求 一年之前的日记
// @Tags 前端
// @Security JWT
// @Produce json
// @Param day query integer  false  "请求几天前的"
// @Router /api/v1/daily [get]
func ObV1GetDailyHandler(c *gin.Context) {
	addData := c.Query("day")
	addDataInt := 0
	var err error
	if addData != "" { // 有参数
		addDataInt, err = strconv.Atoi(addData)
		if err != nil {
			gr.ErrServerError(c, err)
			return
		}
	}
	var text string
	if addDataInt <= -366 { // 不允许请求 一年之前的日记
		gr.RJSON(c, nil, 400, 400, "超出允许范围", gr.H{})
		return
	} else if addDataInt <= -7 { // 一周之前的 使用缓存
		text = skv.GetByFileKey(GetMoreDailyFileKey(addDataInt))

	} else { // 一周内 请求并缓存
		if err = skv.PutByFileKey(GetMoreDailyFileKey(addDataInt)); err != nil {
			gr.ErrServerError(c, err)
			return
		}
		text = skv.GetByFileKey(GetMoreDailyFileKey(addDataInt))
	}
	md_show_text := MarkdownSpilter(MdShowText(
		tools.NowRunConfig.DailyDir()+tools.NowRunConfig.DailyDateKeyMore(addDataInt)+".md", text))
	c.JSON(200, ObDailyV1{
		MdText:     MarkdownSpilter(text),
		MdShowText: md_show_text, // TODO 显示图像
		// MdElements: md.ParseMemos(md_show_text),
		Date: tools.NowRunConfig.DailyDateKeyMore(addDataInt),
	})

}

// ?day=-1 yesterday daily 请求日记
// @Summary Memos 请求 (无缓存)
// @Tags 前端
// @Security JWT
// @Produce json
// @Param day query integer  false  "请求几天前的"
// @Router /api/v1/daily/nocache [get]
func ObV1GetDailyNoCacheHandler(c *gin.Context) {
	addData := c.Query("day")
	addDataInt := 0
	var err error
	if addData != "" { // 有参数
		addDataInt, err = strconv.Atoi(addData)
		if err != nil {
			gr.ErrServerError(c, err)
			return
		}
	}
	var text string
	if addDataInt <= -366 { // 不允许请求 一年之前的日记
		gr.RJSON(c, nil, 400, 400, "超出允许范围", gr.H{})
		return
	} else { // 请求并缓存
		err = skv.PutByFileKey(GetMoreDailyFileKey(addDataInt))
		if err != nil {
			gr.ErrServerError(c, err)
			return
		}
		text = skv.GetByFileKey(GetMoreDailyFileKey(addDataInt))
	}
	md_show_text := MarkdownSpilter(MdShowText(
		tools.NowRunConfig.DailyDir()+tools.NowRunConfig.DailyDateKeyMore(addDataInt)+".md", text))
	c.JSON(200, ObDailyV1{
		MdText:     MarkdownSpilter(text),
		MdShowText: md_show_text,
		// MdElements: md.ParseMemos(md_show_text),
		Date: tools.NowRunConfig.DailyDateKeyMore(addDataInt),
	})

}

type ObV1ModMdText struct {
	DayFileKey string `json:"day"`
	LineNum    int    `json:"line_num"`
	Content    string `json:"content"`
	OldContent string `json:"old"`
}

// 根据行号 修改内容
// eg: {"line_num":99,"content":"new Memos","day": "2023-01-01","old":""}
// @Summary 根据行号修改内容
// @Description 根据行号修改内容，line_num 大于原文件行数，如 9999 新增 Memos 。需要原文件不完整 FileKey 和原来的行的内容进行校验。成功后返回更新后的内容。
// @Tags 前端
// @Security JWT
// @Accept json
// @Produce json
// @Param json body ObV1ModMdText true "根据行号修改内容"
// @Router /api/v1/line [post]
func ObV1PostLineHandler(c *gin.Context) {
	var modText ObV1ModMdText
	if c.ShouldBindJSON(&modText) != nil {
		gr.ErrBindJSONErr(c)
		return
	}
	if modText.DayFileKey == "" {
		modText.DayFileKey = tools.NowRunConfig.DailyDateKeyMore(0)
	}
	fileKey := tools.NowRunConfig.DailyDir() + modText.DayFileKey + ".md"
	dailyText, err := GetTextObject(fileKey)
	// ailyText, err := GetMoreDaliyMdText(modText.DayBefore)
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}
	if dailyText == "No such file: "+fileKey { // 防止提示语被 append 进入文件
		dailyText = ""
	}
	textList := MarkdownSpilter(dailyText)
	if modText.LineNum >= len(textList) { // 如果超出行数则认定为 Memos
		if len(modText.Content) > 30 && strings.HasPrefix(modText.Content, "zk ") { // zk 判读
			modText.Content = modText.Content[3:]
			fileKey := tools.NowRunConfig.DailyAttachmentDir() + tools.TimeFmt("20060102150405.md")
			err := MdTextStore(fileKey, modText.Content)
			if err != nil {
				log.Println(err)
			}
			modText.Content = fmt.Sprintf("![[%s]]", fileKey)
		}
		textList = append(textList, tools.TimeFmt("- 15:04 ")+modText.Content)
	} else {
		if textList[modText.LineNum] == modText.OldContent {
			textList[modText.LineNum] = modText.Content // 行数已经有内容并校验成功 认定为覆写
		} else {
			gr.RJSON(c, nil, 400, 400, "原来数据错误,校验失败", gr.H{})
			return
		}
	}
	newText := strings.Join(textList, "\n")
	MdTextStore(fileKey, newText) // 存入数据源
	skv.PutFile(fileKey, newText) // 存入缓存
	md_show_text := MarkdownSpilter(MdShowText(fileKey, newText))
	c.JSON(200, ObDailyV1{
		MdText:     MarkdownSpilter(newText),
		MdShowText: md_show_text,
		// MdElements: md.ParseMemos(md_show_text),
		Date: modText.DayFileKey,
	})
}

// @Summary 更新文件的缓存
// @Tags 前端
// @Security JWT
// @Param key query string  true  "更新文件 FileKey 完整的"
// @Router /api/v1/cacheupdate [post]
func ObV1UpdateCacheHandler(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		gr.ErrEmpty(c)
		return
	}
	err := skv.PutByFileKey(key)
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}
	gr.Success(c)
}

func LisFileHandler(c *gin.Context) {
	list, err := dao.ListObject("")
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}
	c.JSON(200, list)
}

func TextGetHandler(c *gin.Context) {
	fileKey := c.Query("fileKey")
	if err := skv.PutByFileKey(fileKey); err != nil {
		gr.ErrServerError(c, err)
		return
	}
	text := skv.GetByFileKey(fileKey)
	c.String(200, text)
}

func TextPostHandler(c *gin.Context) {
	fileKey := c.Query("fileKey")
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}
	bodyString := string(bodyBytes)
	skv.PutFile(fileKey, bodyString)
	err = MdTextStore(fileKey, bodyString)
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}
	gr.Success(c)
}
