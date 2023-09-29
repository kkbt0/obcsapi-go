package command

import (
	"fmt"
	"obcsapi-go/dao"

	lua "github.com/yuin/gopher-lua"
)

func LuaRunGoError(L *lua.LState, err error) error {
	errNew := fmt.Errorf("[Lua Run Go Func panic] %v", err)
	L.Push(lua.LString(errNew.Error()))
	return errNew
}

// func Double(L *lua.LState) int {
// 	lv := L.ToInt(1)            /* get argument */
// 	L.Push(lua.LNumber(lv * 2)) /* push result */
// 	return 1                    /* number of results */
// }

func LuaModuleAppLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.SetField(mod, "name", lua.LString("obcsapi-app"))
	L.Push(mod)
	return 1
}

var exports = map[string]lua.LGFunction{
	"DailyTextAppend":      DailyTextAppend,
	"DailyTextAppendMemos": DailyTextAppendMemos,
	"ErrorTest":            ErrorTest,
	"ErrorTestSuccess":     ErrorTestSuccess,
}

func DailyTextAppend(L *lua.LState) int {
	text := L.ToString(1) // 读取参数 2
	err := dao.DailyTextAppend(text)
	if err != nil {
		panic(LuaRunGoError(L, err))
	}
	return 0
}

func DailyTextAppendMemos(L *lua.LState) int {
	text := L.ToString(1) // 读取参数
	err := dao.DailyTextAppendMemos(text)
	if err != nil {
		panic(LuaRunGoError(L, err))
	}

	return 0
}

func ErrorTest(L *lua.LState) int {
	arg1 := L.ToString(1) // 读取参数
	err := fmt.Errorf("读取 %v 错误", arg1)
	if err != nil {
		panic(LuaRunGoError(L, err))
	}
	return 0
}

func ErrorTestSuccess(L *lua.LState) int {
	return 0
}
