--- app 为 obcsapi 模块
--- 本示例为 使用 app 调用 obcsapi 的添加到 memos 的示例

local app = require("app") -- 导入 obcsapi 模块
local test = require("test")

print(app_input) -- 从前端传来的参数
-- 带有异常处理的 obcsapi app 调用
print("DailyTextAppend")
local err1 = app.DailyTextAppend(app_input)
if err1 then
    error(err1)
end

print("DailyTextAppendMemos")
local err2 = app.DailyTextAppendMemos(app_input)
if err2 then
    error(err2)
end


-- local err3 = test.ErrorTest()
-- print(err3)
-- if err3 then -- 等于  if err6 ~= nil then
--     error(err3) -- 打印错误信息并终止
-- end

--- 引入自定义 lua 模块写法
local utils = require("script/utils")
local ans = utils.add(2,3)
---

result = "This is a lua result to obcsapi-web , 输入的内容是 "..app_input -- Lua 拼接字符串

return result
