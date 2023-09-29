package command

import (
	"flag"
	"fmt"

	libs "github.com/vadv/gopher-lua-libs"
	lua "github.com/yuin/gopher-lua"
)

var (
	exec = flag.String("execute", "script/demo.lua", "execute lua script")
)

// Lua Libs Docs https://gitee.com/kkbt/gopher-lua-libs

// err = command.LuaRunner()
//
//	if err != nil {
//		log.Println("[LuaRunner Error]", err)
//	}
func LuaRunner() error {
	flag.Parse()
	luaVm := lua.NewState()
	defer luaVm.Close()
	// 加载 lib
	libs.Preload(luaVm)
	// 加载 obcsapi app
	luaVm.PreloadModule("app", LuaModuleAppLoader)
	if *exec != `` {
		if err := luaVm.DoFile(*exec); err != nil {
			return fmt.Errorf("[Lua ERROR] Error executing file: %v", err)
		}
	} else {
		return fmt.Errorf("[Lua ERROR]: Target file was not given")
	}
	return nil
}
