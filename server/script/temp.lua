--- 一个前端渲染模板的例子

local json_str = [[
{
    "code": 200,
    "data": {
        "type": "message-tbb",
        "command": "#temp 鹅鹅鹅是谁写的",
        "parts": [
            {
                "type": "md",
                "text": "\"obcsapiv4-json\"是由未知的作者编写的。\n"
            },
            {
                "type": "button",
                "text": "Resend",
                "command": "@web-resend"
            },
            {
                "type": "button",
                "text": "Save",
                "command": "@web-save-memos"
            }
        ]
    }
}
]]

return json_str