package command

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// 返回值 int 指的是返回值的个数
// 0 不反回任何东西 异常或正常退出
// 2 反回错误和错误信息 nil,err
// 1 返回错误或值 err,nil
var testExports = map[string]lua.LGFunction{
	"ErrorTest": ErrorTest,
}

func LuaModuleTestLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), testExports)
	L.SetField(mod, "name", lua.LString("lua-test"))
	L.Push(mod)
	return 1
}

func ErrorTest(L *lua.LState) int {
	err := fmt.Errorf("错误")
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}
