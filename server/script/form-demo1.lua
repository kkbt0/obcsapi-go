-- !!! 需要命名为 form.lua !!!
-- 用于处理 前端生成的表单 通用处理方法 会添加到日志中
-- @param app_input string 前端输入的内容，是一个字符串，内容是 json
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

if result["type"] == "AI" then
    local http = require("http")
    local client = http.client()
    -- 参数
    -- https://cloud.baidu.com/doc/WENXINWORKSHOP/s/Ilkkrb0i5
    local access_token = "YOUR_ACCESS_TOKEN"
    local url = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/eb-instant?access_token="..access_token
    -- 构造请求
    local tb = { messages = {{role = "user",content = "书面化文本: "..result["记录"]}}}
    local dt, err = json.encode(tb)
    if err then error(err) end
    --local dt = "{\"status\": \""..text.."\",\"visibility\": \"public\"}"
    local request = http.request("POST", url ,dt)
    request:header_set("Content-Type", "application/json")
    -- 发送请求
    local result, err = client:do_request(request)
    if err then error(err) end
    if not (result.code == 200) then
      print(dt)
      error(result.body)
    end
    print(dt)
    print(result)
    print(result.body)
    ret = "\n请求结果 JSON 字符串:\n"..result.body
    local result, err = json.decode(result.body)
    if err then error(err) end
    -- 保存到笔记
    local err = app.AppendDailyMemos(result["result"])
    if err then error(err) end
    return "\nAI 使用 "..result["usage"]["total_tokens"].."个 tokens 返回结果:\n"..result["result"].."\n"..ret
end

-- 使用工具库 table_to_string 将 table 转为字符串
table_str = utils.table_to_string(result)
print(table_str)

ret = "后端返回: 前端输入的内容是 "..app_input -- Lua 拼接字符串
ret = ret.."\n解析 JSON 最终转为字符串结果:\n"..table_str
--- 将结果保存到 memos

local err1 = app.AppendDailyMemos(table_str)
if err1 then
    error(err1)
end

return ret.."\n并且保存到 memos 中" -- 返回结果，供给前端显示


-- 使用的表单 JSON Schema 在 前端-齿轮-Mention-表单可选列表-第二个空格
-- 表单可选列表 第一个空是 表单名 第二个空是 JSON Schema 内容
--[[

{
    "type": "object",
    "displayType": "row",
    "properties": {
        "记录": {
            "title": "记录",
            "type": "string",
            "ui:options": {
                "type": "textarea",
                "rows": 10
            }
        },
        "type": {
            "title": "类型",
            "type": "string",
            "enum": [
                "AI"
            ],
            "default": "AI"
        }
    }
}

]]--