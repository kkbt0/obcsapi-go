package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gin-gonic/gin"
)

type Daily struct {
	Data       string `json:"data"`
	MdShowData string `json:"md_show_data"`
	Date       string `json:"date"`
	ServerTime string `json:"serverTime"`
}
type MemosData struct {
	Content string `json:"content"`
}
type MoodReader struct {
	Highlights []MoodReaderHighlights `json:"highlights"`
}
type MoodReaderHighlights struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	Author string `json:"author"`
	Note   string `json:"note"`
}

// Token1
func ObTodayHandler(c *gin.Context) {
	client, err := get_client()
	if err != nil {
		c.Error(err)
	}
	switch c.Request.Method {
	case "OPTIONS":
		c.Status(200)
	case "GET":
		// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json_data := get_today_daily_list(client)[0]
		daily_data := Daily{Date: timeFmt("2006-01-02"), ServerTime: timeFmt("200601021504"), Data: json_data, MdShowData: string(replace_md_url2pre_url([]byte(json_data)))}
		data, _ := json.Marshal([]Daily{daily_data})
		c.String(200, string(data))
	case "POST":
		decoder := json.NewDecoder(c.Request.Body)
		var memosData MemosData
		err := decoder.Decode(&memosData)
		if err != nil {
			log.Println(err)
		}
		append_memos_in_daily(client, memosData.Content)
		c.String(200, "Success")
	default:
		c.Status(404)
	}
}

// Token1
func ObPostTodayAllHandler(c *gin.Context) {
	client, err := get_client()
	if err != nil {
		c.Error(err)
	}
	decoder := json.NewDecoder(c.Request.Body)
	var memosData MemosData
	err = decoder.Decode(&memosData)
	if err != nil {
		c.Error(err)
	} else {
		store(client, daily_file_key(), []byte(memosData.Content))
	}
	c.String(200, "Success")
}

// Tokne1
func ObGet3DaysHandler(c *gin.Context) {
	client, err := get_client()
	if err != nil {
		c.Status(500)
		c.Error(err)
		return
	}
	three_list := get_3_daily_list(client)
	data, err := json.Marshal(three_list)
	if err != nil {
		c.Status(500)
		c.Error(err)
		return
	}
	c.String(200, string(data))
}

// Token2 静读天下使用的 API
func MoodReaderHandler(c *gin.Context) {
	right_token2, _ := GetToken("token2")
	if c.Request.Header.Get("Token") != "Token "+right_token2.TokenString {
		c.Status(401)
		return
	}
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
	client, _ := get_client()
	file_key := fmt.Sprintf("支持类文件/MoonReader/%s.md", ReplaceUnAllowedChars(title))
	append_text := fmt.Sprintf("文: %s\n批: %s\n于: %s\n\n---\n", text, note, timeFmt("2006-01-02 15:04"))
	md, _ := get_object(client, file_key)
	if md != nil {
		err = append(client, file_key, append_text)
	} else {
		yaml := fmt.Sprintf("---\ntitle: %s\nauthor: %s\n---\n书名: %s\n作者: %s\n简介: \n评价: \n\n---\n", title, author, title, author)
		err = append(client, file_key, yaml+append_text)
	}
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	c.String(200, "Success")
}

// 安卓软件 fv 悬浮球使用的 API 用于自定义任务的 图片、文字
func fvHandler(c *gin.Context) {
	client, _ := get_client()
	if c.GetHeader("Content-Type") == "text/plain" {
		content, _ := ioutil.ReadAll(c.Request.Body)
		append_memos_in_daily(client, string(content))
		c.String(200, "Success")
		return
	} else if c.GetHeader("Content-Type") == "application/octet-stream" {
		content, _ := ioutil.ReadAll(c.Request.Body)
		file_key := fmt.Sprintf("日志/附件/%s/%s.jpg", timeFmt("200601"), timeFmt("20060102150405"))
		store(client, file_key, content)
		append_memos_in_daily(client, fmt.Sprintf("![](%s)", file_key))
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
func SRWebHook(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var simpReadJson SimpReadWebHookStruct
	err := decoder.Decode(&simpReadJson)
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	serverTime := timeFmt("200601021504")
	yaml := fmt.Sprintf("---\ntitle: %s\nsctime: %s\n---\n", simpReadJson.Title, serverTime)
	file_str := fmt.Sprintf("%s[[简悦WebHook生成]]\n生成时间: %s\n原文: %s\n标题: %s\n描述: %s\n标签: %s\n内容: \n%s", yaml, serverTime, simpReadJson.Url, simpReadJson.Title, simpReadJson.Description, simpReadJson.Tags, simpReadJson.Content)
	file_key := fmt.Sprintf("支持类文件/SimpRead/%s %s.md", ReplaceUnAllowedChars(simpReadJson.Title), serverTime)
	client, err := get_client()
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	store(client, file_key, []byte(file_str))
}

// 通用 API 接口 使用 Token2 验证
func GeneralHeader(c *gin.Context) {
	switch c.Request.Method {
	case "OPTIONS":
		c.Status(200)
	case "POST":
		decoder := json.NewDecoder(c.Request.Body)
		var memosData MemosData
		err := decoder.Decode(&memosData)
		if err != nil {
			c.Error(err)
			c.Status(500)
			return
		}
		client, err := get_client()
		if err != nil {
			c.Error(err)
			c.Status(500)
			return
		}
		err = append_memos_in_daily(client, memosData.Content)
		if err != nil {
			c.Error(err)
			c.Status(500)
			return
		}
		c.String(200, "Success")
	default:
		c.Status(404)
	}
}

type UrlStruct struct {
	Url string `json:"url"`
}

// Token2
func Url2MdHandler(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var urlStruct UrlStruct
	err := decoder.Decode(&urlStruct)
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	text, err := Downloader(urlStruct.Url)
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
	client, err := get_client()
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	serverTime := timeFmt("200601021504")
	yaml := fmt.Sprintf("---\nurl: %s\nsctime: %s\n---\n", urlStruct.Url, serverTime)
	file_key := fmt.Sprintf("支持类文件/HtmlPages/%s.md", serverTime)
	store(client, file_key, []byte(yaml+markdown))
	c.Status(200)
}
