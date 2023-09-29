package command

import (
	"obcsapi-go/dao"

	lua "github.com/yuin/gopher-lua"
)

var appExports = map[string]lua.LGFunction{
	"DailyTextAppend":      DailyTextAppend,
	"DailyTextAppendMemos": DailyTextAppendMemos,
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
