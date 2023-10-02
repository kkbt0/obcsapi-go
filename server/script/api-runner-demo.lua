local json = require("json")
local utils = require("script/utils") -- 导入自定义函数工具

local result, err = json.decode(app_input)
if err then error(err) end

local table_str = utils.table_to_string(result)

return table_str