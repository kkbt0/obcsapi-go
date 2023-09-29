local app = require("app") -- 导入 obcsapi 模块

local err = app.DailyTextAppend("xxxx")
local err2 = app.DailyTextAppendMemos("xxxx222")

print(err)
print(err2)
print("Hello World")
print(app.name)




-- 异常处理


local success, errMsg = pcall(app.ErrorTest,"This is a test error from lua file")
print(success, errMsg)
local success2, errMsg2 = pcall(app.ErrorTestSuccess)
print(success2, errMsg2)


local success, result = pcall(app.ErrorTest, "some_argument")
-- 检查是否出现异常
if not success then
    -- 处理异常，result 变量包含了异常信息
    print("An error occurred:", result)
else
    -- 函数执行成功，result 变量包含了函数的返回值
    print("Function executed successfully. Result:", result)
end