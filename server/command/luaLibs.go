package command

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func LuaRunGoError(L *lua.LState, err error) error {
	errNew := fmt.Errorf("[Lua Run Go Func panic] %v", err)
	L.Push(lua.LString(errNew.Error()))
	return errNew
}
