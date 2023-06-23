package dao

import (
	"obcsapi-go/tools"
	"path"
)

// dir1/dir2/3.md ../4.md -> dir1/4.md 注意歧义：
// 4.md -> dir1/dir2/4.md  but it also can be 4.md
func GetAbsoluteFileKey(fromFilekey string, relativePath string) string {
	// 获取当前文件所在的目录
	dir := path.Dir(fromFilekey)
	// 拼接相对路径和当前目录，得到新的filekey
	newFilekey := path.Join(dir, relativePath)
	tools.Debug("GetAbsoluteFileKey: in", fromFilekey, relativePath, "->", newFilekey)
	return newFilekey
}
