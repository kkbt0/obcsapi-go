local app = require("app") -- 导入 obcsapi 模块

-- 打印输入的参数
print(app_input)

-- 尝试 JSON 解析
local json = require("json")
local result, err = json.decode(app_input)
if err then
    error(err)
end

print(result['text1'])

ret = "后端返回: 前端输入的内容是 "..app_input -- Lua 拼接字符串
ret = ret.."\n解析 JSON 结果: text1 = "..result['text1'].." text2: "..result['text2']
return ret