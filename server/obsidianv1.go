package main

import (
	"fmt"
	"log"
	. "obcsapi-go/dao"
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
func ObV1GetDailyHandler(c *gin.Context) {
	addData := c.Query("day")
	addDataInt := 0
	var err error
	if addData != "" { // 有参数
		addDataInt, err = strconv.Atoi(addData)
		if err != nil {
			c.Error(err)
			c.Status(500)
			return
		}
	}
	var text string
	if addDataInt <= -366 { // 不允许请求 一年之前的日记
		c.String(400, "超出允许范围")
		return
	} else if addDataInt <= -7 { // 一周之前的 使用缓存
		text = skv.GetByFileKey(GetMoreDailyFileKey(addDataInt))

	} else { // 一周内 请求并缓存
		err = skv.PutByFileKey(GetMoreDailyFileKey(addDataInt))
		if err != nil {
			c.Error(err)
			c.Status(500)
			return
		}
		text = skv.GetByFileKey(GetMoreDailyFileKey(addDataInt))
	}
	c.JSON(200, ObDailyV1{
		MdText:     MarkdownSpilter(text),
		MdShowText: MarkdownSpilter(MdShowText(text)), // TODO 显示图像
		Date:       tools.NowRunConfig.DailyDateKeyMore(addDataInt),
	})

}

// ?day=-1 yesterday daily 请求日记
func ObV1GetDailyNoCacheHandler(c *gin.Context) {
	addData := c.Query("day")
	addDataInt := 0
	var err error
	if addData != "" { // 有参数
		addDataInt, err = strconv.Atoi(addData)
		if err != nil {
			c.Error(err)
			c.Status(500)
			return
		}
	}
	var text string
	if addDataInt <= -366 { // 不允许请求 一年之前的日记
		c.String(400, "超出允许范围")
		return
	} else { // 请求并缓存
		err = skv.PutByFileKey(GetMoreDailyFileKey(addDataInt))
		if err != nil {
			c.Error(err)
			c.Status(500)
			return
		}
		text = skv.GetByFileKey(GetMoreDailyFileKey(addDataInt))
	}
	c.JSON(200, ObDailyV1{
		MdText:     MarkdownSpilter(text),
		MdShowText: MarkdownSpilter(MdShowText(text)), // TODO 显示图像
		Date:       tools.NowRunConfig.DailyDateKeyMore(addDataInt),
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
func ObV1PostLineHandler(c *gin.Context) {
	var modText ObV1ModMdText
	if c.ShouldBindJSON(&modText) != nil {
		c.String(400, "参数错误")
		return
	}
	if modText.DayFileKey == "" {
		modText.DayFileKey = tools.NowRunConfig.DailyDateKeyMore(0)
	}
	fileKey := tools.NowRunConfig.DailyDir() + modText.DayFileKey + ".md"
	dailyText, err := GetTextObject(fileKey)
	// ailyText, err := GetMoreDaliyMdText(modText.DayBefore)
	if err != nil {
		c.Error(err)
		c.Status(500)
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
			c.String(400, "原来数据，参数错误")
			return
		}
	}
	newText := strings.Join(textList, "\n")
	MdTextStore(fileKey, newText) // 存入数据源
	skv.PutFile(fileKey, newText) // 存入缓存
	c.JSON(200, ObDailyV1{
		MdText:     MarkdownSpilter(newText),
		MdShowText: MarkdownSpilter(MdShowText(newText)),
		Date:       modText.DayFileKey,
	})
}

func ObV1UpdateCacheHandler(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.Status(400)
		return
	}
	err := skv.PutByFileKey(key)
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	c.Status(200)
}
