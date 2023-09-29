local app = require("app") -- 导入 obcsapi 模块


-- 带有异常处理的 obcsapi app 调用
local success, result = pcall(app.DailyTextAppend, "Some Str1")
if not success then
    print("An error occurred:", result)
else
    print("Function executed successfully. Result:", result)
end

local success2, result2 = pcall(app.DailyTextAppendMemos, "Some Str2")
if not success2 then
    print("An error occurred:", result2)
else
    print("Function executed successfully. Result:", result)
end

