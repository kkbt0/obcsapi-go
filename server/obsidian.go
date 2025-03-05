package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"obcsapi-go/gr"
	"obcsapi-go/skv"
	"obcsapi-go/tools"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"obcsapi-go/dao"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 提取标题的函数
func extractTitle(html string) string {
	// 1. 首先尝试从 meta 标签中获取（微信文章通常在这里）
	metaTitleRegex := regexp.MustCompile(`<meta[^>]*property="og:title"[^>]*content="([^"]*)"`)
	if matches := metaTitleRegex.FindStringSubmatch(html); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// 2. 尝试从 <title> 标签中获取
	titleRegex := regexp.MustCompile(`<title[^>]*>(.*?)</title>`)
	if matches := titleRegex.FindStringSubmatch(html); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// 3. 尝试从 h1 标签中获取
	h1Regex := regexp.MustCompile(`<h1[^>]*>(.*?)</h1>`)
	if matches := h1Regex.FindStringSubmatch(html); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// 4. 尝试从 #activity-name 中获取（微信文章特有的标题类）
	activityNameRegex := regexp.MustCompile(`<h1 class="rich_media_title"[^>]*>(.*?)</h1>`)
	if matches := activityNameRegex.FindStringSubmatch(html); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	return ""
}

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
	var memosData dao.MemosData
	var err error
	if err = decoder.Decode(&memosData); err != nil {
		gr.ErrServerError(c, err)
		return
	}
	if err = dao.AppendDailyMemos(memosData.Content); err != nil {
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
	var memosData dao.MemosData
	var err error
	if err = decoder.Decode(&memosData); err != nil {
		gr.ErrServerError(c, err)
		return
	}
	if err = dao.AppendDailyMemos(memosData.Content); err != nil {
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
	var urlStruct dao.UrlStruct
	if err := decoder.Decode(&urlStruct); err != nil {
		gr.ErrServerError(c, err)
		return
	}
	var text []byte

	if text, err = tools.Downloader(urlStruct.Url); err != nil {
		gr.ErrServerError(c, err)
		return
	}

	// 提取所有图片链接并修改 HTML
	imgRegex := regexp.MustCompile(`data-src="([^"]*)"`)
	matches := imgRegex.FindAllStringSubmatch(string(text), -1)
	modifiedHTML := imgRegex.ReplaceAllStringFunc(string(text), func(match string) string {
		return strings.Replace(match, "data-src", "src", 1)
	})

	// 保存调试文件
	if viper.GetBool("debug") {
		debugDir := "debug"
		if err := os.MkdirAll(debugDir, 0755); err != nil {
			log.Printf("创建调试目录失败: %v", err)
		} else {
			if err = os.WriteFile(filepath.Join(debugDir, "debug_original.html"), text, 0644); err != nil {
				log.Printf("保存原始HTML失败: %v", err)
			}

			// 保存图片链接
			var imgLinks []string
			for _, match := range matches {
				if len(match) > 1 {
					imgLinks = append(imgLinks, match[1])
				}
			}
			if err = os.WriteFile(filepath.Join(debugDir, "debug_images.txt"), []byte(strings.Join(imgLinks, "\n")), 0644); err != nil {
				log.Printf("保存图片链接失败: %v", err)
			}

			// 保存修改后的 HTML
			if err = os.WriteFile(filepath.Join(debugDir, "debug_modified.html"), []byte(modifiedHTML), 0644); err != nil {
				log.Printf("保存修改后HTML失败: %v", err)
			}
		}
	}

	// 配置 HTML 到 Markdown 的转换器
	converter := md.NewConverter("", true, &md.Options{
		// 使用默认配置，但确保图片链接被保留
		GetAbsoluteURL: func(selec *goquery.Selection, rawURL string, domain string) string {
			return rawURL
		},
	})

	// 转换为 Markdown
	markdown, err := converter.ConvertString(modifiedHTML)
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}

	// 清理微信公众号文章末尾的无用内容
	if strings.Contains(urlStruct.Url, "mp.weixin.qq.com") {
		markdown = cleanWechatArticle(markdown)
	}

	// 保存转换后的 Markdown
	if viper.GetBool("debug") {
		if err = os.WriteFile(filepath.Join("debug", "debug_result.md"), []byte(markdown), 0644); err != nil {
			log.Printf("保存Markdown失败: %v", err)
		}
	}

	// 使用更强大的标题提取方法
	title := extractTitle(string(text))
	if title == "" {
		// 如果提取失败，回退到使用 Markdown 第一行
		converter := md.NewConverter("", true, nil)
		var markdown string
		if markdown, err = converter.ConvertString(string(text)); err != nil {
			gr.ErrServerError(c, err)
			return
		}
		title = strings.Split(markdown, "\n")[0]
	}

	serverTime := tools.TimeFmt("200601021504")
	yaml := fmt.Sprintf("---\nurl: %s\ntitle: %s\nsctime: %s\n---\n[[ObSavePage]]\n",
		urlStruct.Url,
		tools.ReplaceUnAllowedChars(strings.TrimSpace(title)),
		serverTime)

	file_key := fmt.Sprintf("%sHtmlPages/%s %s.md",
		tools.NowRunConfig.OtherDataDir(),
		serverTime,
		tools.ReplaceUnAllowedChars(strings.TrimSpace(title)))
	if err = dao.CoverStoreTextFile(file_key, yaml+markdown); err != nil {
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
		err = dao.CoverStoreTextFile(fileKey, generalJson.Content)
	} else {
		err = dao.AppendText(fileKey, generalJson.Content)
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
	if mdText, err := dao.GetFileText(tools.NowRunConfig.DailyFileKeyMore(ObTodayAddDateNum())); err != nil {
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
	if err = dao.CoverStoreTextFile(fileKey, string(bodyBytes)); err != nil {
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
	if err = dao.AppendText(fileKey, string(bodyBytes)); err != nil {
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

// 清理微信公众号文章末尾的无用内容
func cleanWechatArticle(content string) string {
	// 查找截断词的位置
	cutPoint := strings.Index(content, "预览时标签不可点")
	if cutPoint == -1 {
		return content
	}

	// 截取有效内容
	return strings.TrimSpace(content[:cutPoint])
}
