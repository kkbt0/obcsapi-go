-- 常用和工具函数
-- utils.table_to_string(table)

local M = {}

-- 测试用例
function M.add(a, b)
    return a + b
end

-- 打印 table 为字符串
function M.table_to_string(t)
    local str = ""
    for k, v in pairs(t) do
        str = str .. string.format("%s: %s\n", k, tostring(v))
    end
    return str
end

return M