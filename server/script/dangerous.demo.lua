--- 运行前端传入的 Lua 脚本
--- 非常危险 

local func1 = loadstring(app_input)
local result = func1()
return result