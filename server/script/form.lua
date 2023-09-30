--- 用于处理 前端生成的表单
--- @param app_input string 前端输入的内容，是一个字符串，内容是 json
-- @return string Lua 拼接字符串，最终给前端显示

local app = require("app") -- 导入 obcsapi 模块
local utils = require("script/utils") -- 导入自定义函数工具

-- 打印输入的参数
print(app_input)

-- 尝试 JSON 解析 将 json 解析为 table
local json = require("json")
local result, err = json.decode(app_input)
if err then
    error(err)
end


-- 此处可以做一些事情 
-- 比如 取出 json 中的一些字段做判断，执行不同的函数，或者脚本

-- 使用工具库 table_to_string 将 table 转为字符串
table_str = utils.table_to_string(result)
print(table_str)

ret = "后端返回: 前端输入的内容是 "..app_input -- Lua 拼接字符串
ret = ret.."\n解析 JSON 最终转为字符串结果:\n"..table_str
--- 将结果保存到 memos

local err1 = app.DailyTextAppend(app_input)
if err1 then
    error(err1)
end

return ret.."\n并且保存到 memos 中" -- 返回结果，供给前端显示