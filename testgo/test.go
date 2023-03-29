package main

import (
	"fmt"
	"obcsapi-go/dao"
	"regexp"
	"strings"
	"testing"
)

func main() {
	fmt.Println("Run")
	//estRegexp2()
	fmt.Println(dao.CheckObject("3.md"))
	fmt.Println(dao.GetTextObject("2.md"))
	fmt.Println(dao.ObjectStore("object.md", []byte("a new object")))
	fmt.Println(dao.MdTextStore("text.md", "A new object"))
	fmt.Println(dao.TextAppend("text.md", "\nNew line"))
	fmt.Println(dao.DailyTextAppendMemos("今日测试"))
}

func testRegexp2() {
	s := "-?!*<>[#^|]"
	fmt.Println(ReplaceUnAllowedChars(s))
}

//  obsidian 文件名非法字符 * " \ / < > : | ? 链接失效 # ^ [ ] | 替换为 _
func ReplaceUnAllowedChars(s string) string {
	unAllowedChars := "*\"\\/<>:|?#^[]|"
	fmt.Println(unAllowedChars)
	for _, c := range unAllowedChars {
		s = strings.ReplaceAll(s, string(c), "_")
	}
	return s
}

func testRegexp(t *testing.T) {
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
