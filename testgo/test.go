package main

import (
	"fmt"
	"path"
)

func main() {
	fmt.Println(GetAbsoluteFileKey("1/2/3.md", "4.md"))
	fmt.Println(GetAbsoluteFileKey("1/2/3.md", "../4.md"))
	fmt.Println(GetAbsoluteFileKey("1/2/3.md", "../.4.md"))
	fmt.Println(GetAbsoluteFileKey("1/2/3.md", "../../4.md"))
}

func GetAbsoluteFileKey(filekey string, relativePath string) string {
	fmt.Println("------")
	// 获取当前文件所在的目录
	dir := path.Dir(filekey)
	fmt.Println(dir)
	// 拼接相对路径和当前目录，得到新的filekey
	newFilekey := path.Join(dir, relativePath)
	return newFilekey
}
