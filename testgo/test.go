package main

import (
	"fmt"
	"regexp"
)

func main() {
	fmt.Println("Run")
	olrdMd := "![](https://github.com/1.jpg)\n![xxx](https://github.com/2.jpg)\n![](xxx)xxxx\n![xx](xxx)xxxx"
	pattern := regexp.MustCompile(`!\[(.*?)\]\(([^http:].*)\)`)
	replaceFunc := func(match []byte) []byte {
		// 获取匹配到的链接
		desp := pattern.ReplaceAllString(string(match), "$1")
		link := pattern.ReplaceAllString(string(match), "$2")
		// 替换链接为临时鉴权链接
		return []byte(fmt.Sprintf("![%s](%s)", desp, dealPic(link)))
	}
	result := pattern.ReplaceAllFunc([]byte(olrdMd), replaceFunc)
	fmt.Println(string(result))
}

func dealPic(key string) string {
	return fmt.Sprintf("%s+捕获", key)
}
