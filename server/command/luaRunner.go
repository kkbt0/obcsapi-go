package command

import (
	"flag"
	"fmt"
	"obcsapi-go/tools"
	"os"

	libs "github.com/vadv/gopher-lua-libs"
	lua "github.com/yuin/gopher-lua"
)

// Lua Libs Docs https://gitee.com/kkbt/gopher-lua-libs

func LuaRunner(scriptFilePath string, inputText string) (string, error) {
	// flag
	exec := flag.String("execute", scriptFilePath, "execute lua script")
	flag.Parse()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	// 初始化
	luaVm := lua.NewState()
	defer luaVm.Close()

	// 加载 lib
	libs.Preload(luaVm)
	// 加载 obcsapi app
	luaVm.PreloadModule("app", LuaModuleAppLoader)
	luaVm.PreloadModule("test", LuaModuleTestLoader)
	// 加载变量
	luaVm.SetGlobal("app_input", lua.LString(inputText))
	if *exec != `` {

		if err := luaVm.DoFile(*exec); err != nil {
			return "", fmt.Errorf("[Lua ERROR] Error executing file: %v", err)
		}
	} else {
		return "", fmt.Errorf("[Lua ERROR]: Target file was not given")
	}
	result := luaVm.Get(-1) // get the value at the top of the stack
	tools.Debug("[Lua Runner Result]", result.String())

	return result.String(), nil
}
