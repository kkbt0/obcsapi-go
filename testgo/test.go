package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

func main() {
	url := "https://www.ftls.xyz/posts/anki-sync-server-rs-docker/" // 替换为您要下载的网页URL

	// 发送HTTP GET请求并获取响应
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("HTTP GET error:", err)
		return
	}
	defer response.Body.Close()

	// 读取响应的HTML内容
	htmlBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("HTML read error:", err)
		return
	}

	// 将HTML内容转换为Markdown
	markdown := convertToMarkdown(htmlBytes)

	fmt.Println(markdown)
}

func convertToMarkdown(html []byte) string {
	// TODO: 在这里编写将HTML转换为Markdown的代码
	// 您可以使用golang.org/x/net/html包来解析HTML并将其转换为Markdown格式

	// 示例代码中暂未包含转换逻辑，您需要根据自己的需求实现该功能

	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(string(html))

	if err != nil {
		log.Println(err)
	}

	return markdown // 返回Markdown字符串
}
