package tools

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

// 写入日志文件
func writeLog(content string) error {
	// 创建日志目录
	if err := os.MkdirAll("logs", 0755); err != nil {
		return err
	}

	// 生成日志文件名
	filename := fmt.Sprintf("logs/url_test_%s.log", time.Now().Format("20060102_150405"))

	// 写入文件
	return os.WriteFile(filename, []byte(content), 0644)
}

// 提取标题的函数
func extractTitle(html string) string {
	logContent := "开始提取标题...\n"
	
	// 1. 首先尝试从 meta 标签中获取（微信文章通常在这里）
	metaTitleRegex := regexp.MustCompile(`<meta[^>]*property="og:title"[^>]*content="([^"]*)"`)
	if matches := metaTitleRegex.FindStringSubmatch(html); len(matches) > 1 {
		logContent += "从 meta 标签中找到标题\n"
		logContent += fmt.Sprintf("标题内容: %s\n", matches[1])
		return strings.TrimSpace(matches[1])
	}

	// 2. 尝试从 <title> 标签中获取
	titleRegex := regexp.MustCompile(`<title[^>]*>(.*?)</title>`)
	if matches := titleRegex.FindStringSubmatch(html); len(matches) > 1 {
		logContent += "从 <title> 标签中找到标题\n"
		logContent += fmt.Sprintf("标题内容: %s\n", matches[1])
		return strings.TrimSpace(matches[1])
	}

	// 3. 尝试从 h1 标签中获取
	h1Regex := regexp.MustCompile(`<h1[^>]*>(.*?)</h1>`)
	if matches := h1Regex.FindStringSubmatch(html); len(matches) > 1 {
		logContent += "从 h1 标签中找到标题\n"
		logContent += fmt.Sprintf("标题内容: %s\n", matches[1])
		return strings.TrimSpace(matches[1])
	}

	// 4. 尝试从 #activity-name 中获取（微信文章特有的标题类）
	activityNameRegex := regexp.MustCompile(`<h1 class="rich_media_title"[^>]*>(.*?)</h1>`)
	if matches := activityNameRegex.FindStringSubmatch(html); len(matches) > 1 {
		logContent += "从 rich_media_title 中找到标题\n"
		logContent += fmt.Sprintf("标题内容: %s\n", matches[1])
		return strings.TrimSpace(matches[1])
	}

	logContent += "未找到标题\n"
	writeLog(logContent)
	return ""
}

// 提取作者信息
func extractAuthor(html string) string {
	// 1. 尝试从 meta 标签中获取
	metaAuthorRegex := regexp.MustCompile(`<meta[^>]*property="og:article:author"[^>]*content="([^"]*)"`)
	if matches := metaAuthorRegex.FindStringSubmatch(html); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// 2. 尝试从 meta name="author" 中获取
	authorRegex := regexp.MustCompile(`<meta[^>]*name="author"[^>]*content="([^"]*)"`)
	if matches := authorRegex.FindStringSubmatch(html); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	return ""
}

// 提取描述信息
func extractDescription(html string) string {
	// 1. 尝试从 meta 标签中获取
	metaDescRegex := regexp.MustCompile(`<meta[^>]*property="og:description"[^>]*content="([^"]*)"`)
	if matches := metaDescRegex.FindStringSubmatch(html); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// 2. 尝试从 meta name="description" 中获取
	descRegex := regexp.MustCompile(`<meta[^>]*name="description"[^>]*content="([^"]*)"`)
	if matches := descRegex.FindStringSubmatch(html); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	return ""
}

// 提取正文内容
func extractContent(html string) string {
	logContent := "开始提取正文...\n"
	
	// 1. 尝试从 js_content 中获取
	contentRegex := regexp.MustCompile(`<div[^>]*id="js_content"[^>]*>(.*?)</div>`)
	if matches := contentRegex.FindStringSubmatch(html); len(matches) > 1 {
		content := matches[1]
		// 移除 HTML 标签
		content = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(content, "")
		// 移除多余空白
		content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ")
		content = strings.TrimSpace(content)
		
		logContent += "从 js_content 中找到正文\n"
		logContent += fmt.Sprintf("正文内容: %s\n", content)
		writeLog(logContent)
		return content
	}

	// 2. 尝试从 article 标签中获取
	articleRegex := regexp.MustCompile(`<article[^>]*>(.*?)</article>`)
	if matches := articleRegex.FindStringSubmatch(html); len(matches) > 1 {
		content := matches[1]
		// 移除 HTML 标签
		content = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(content, "")
		// 移除多余空白
		content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ")
		content = strings.TrimSpace(content)
		
		logContent += "从 article 标签中找到正文\n"
		logContent += fmt.Sprintf("正文内容: %s\n", content)
		writeLog(logContent)
		return content
	}

	// 3. 尝试从 main 标签中获取
	mainRegex := regexp.MustCompile(`<main[^>]*>(.*?)</main>`)
	if matches := mainRegex.FindStringSubmatch(html); len(matches) > 1 {
		content := matches[1]
		// 移除 HTML 标签
		content = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(content, "")
		// 移除多余空白
		content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ")
		content = strings.TrimSpace(content)
		
		logContent += "从 main 标签中找到正文\n"
		logContent += fmt.Sprintf("正文内容: %s\n", content)
		writeLog(logContent)
		return content
	}

	logContent += "未找到正文内容\n"
	writeLog(logContent)
	return ""
}

func TestURLFetch(t *testing.T) {
	logContent := "开始测试 URL 抓取...\n"
	
	// 使用微信文章进行测试
	url := "https://mp.weixin.qq.com/s/1qoCSFfApzwfe2kjkcnbFA"
	logContent += fmt.Sprintf("测试 URL: %s\n", url)
	
	// 创建 HTTP 客户端
	client := &http.Client{}
	
	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logContent += fmt.Sprintf("创建请求失败: %v\n", err)
		writeLog(logContent)
		t.Fatalf("创建请求失败: %v", err)
	}

	// 添加请求头
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Add("Cookie", "wxuin=1234567890; wxsid=abcdefghijklmnopqrstuvwxyz;")
	
	// 发送请求
	logContent += "发送 HTTP 请求...\n"
	resp, err := client.Do(req)
	if err != nil {
		logContent += fmt.Sprintf("获取 URL 失败: %v\n", err)
		writeLog(logContent)
		t.Fatalf("获取 URL 失败: %v", err)
	}
	defer resp.Body.Close()

	logContent += fmt.Sprintf("收到响应，状态码: %d\n", resp.StatusCode)

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logContent += fmt.Sprintf("读取响应内容失败: %v\n", err)
		writeLog(logContent)
		t.Fatalf("读取响应内容失败: %v", err)
	}

	// 将字节数组转换为字符串
	html := string(body)
	logContent += fmt.Sprintf("获取到 HTML 内容，长度: %d 字节\n", len(html))

	// 记录完整的 HTML 内容
	logContent += "原始 HTML 内容:\n"
	logContent += html
	logContent += "\n\n"

	// 提取标题
	title := extractTitle(html)
	logContent += fmt.Sprintf("提取到的标题: %s\n", title)

	// 提取作者
	author := extractAuthor(html)
	logContent += fmt.Sprintf("提取到的作者: %s\n", author)

	// 提取描述
	description := extractDescription(html)
	logContent += fmt.Sprintf("提取到的描述: %s\n", description)

	// 提取正文内容
	content := extractContent(html)
	logContent += fmt.Sprintf("提取到的正文: %s\n", content)

	// 写入日志文件
	if err := writeLog(logContent); err != nil {
		t.Fatalf("写入日志文件失败: %v", err)
	}
} 