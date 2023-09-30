package command

import (
	"obcsapi-go/dao"
	"obcsapi-go/tools"
	"time"

	lua "github.com/yuin/gopher-lua"
)

var appExports = map[string]lua.LGFunction{
	"DailyTextAppend":      DailyTextAppend,      // 日志新增文本
	"DailyTextAppendMemos": DailyTextAppendMemos, // 日志新增 memos
	"TextAppend":           TextAppend,           // 指定文件添加文本
	"GetTextObject":        GetTextObject,        // 获取指定文件内容，返回字符串
	"MdTextStore":          MdTextStore,          // 覆盖指定位置文件 纯文本使用
	"TodayGet":             TodayGet,             // 获取今日日记 md 文件字符串 // 每天凌晨 00:00 - 03:59  判断为 today daily 为 昨天的日志
}

func LuaModuleAppLoader(L *lua.LState) int {

	modApp := L.SetFuncs(L.NewTable(), appExports)
	L.SetField(modApp, "name", lua.LString("obcsapi-app"))
	L.Push(modApp)

	return 1
}

func DailyTextAppend(L *lua.LState) int {
	text := L.ToString(1) // 读取参数
	err := dao.DailyTextAppend(text)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

func DailyTextAppendMemos(L *lua.LState) int {
	text := L.ToString(1) // 读取参数
	err := dao.DailyTextAppendMemos(text)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

func TextAppend(L *lua.LState) int {
	file_key := L.ToString(1)
	text := L.ToString(2)
	err := dao.TextAppend(file_key, text)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

func GetTextObject(L *lua.LState) int {
	text_file_key := L.ToString(1) // 读取参数
	text, err := dao.GetTextObject(text_file_key)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(text))
	return 1
}

func MdTextStore(L *lua.LState) int {
	file_key := L.ToString(1)
	text := L.ToString(2)
	err := dao.MdTextStore(file_key, text)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

func TodayGet(L *lua.LState) int {
	mdText, err := dao.GetTextObject(tools.NowRunConfig.DailyFileKeyMore(ObTodayAddDateNum()))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(mdText))
	return 1
}

// ----------- Tools --------------
// 每天凌晨 00:00 - 03:59  判断为 today daily 为 昨天的日志
func ObTodayAddDateNum() int {
	hour := time.Now().Hour()
	if hour >= 0 && hour <= 3 {
		return -1
	}
	return 0
}
