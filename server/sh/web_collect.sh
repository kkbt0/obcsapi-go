#!/bin/bash

# 功能: 收藏网页
# 触发: 需要以 http 开头的链接单独一行

host="http://localhost:8900"
token2=`cat token/token2.json | jq -r .token` 
# 注意不能有不允许的文件名字符 如 : 

# 获取最后的输入
input=`bash sh/lastinput.sh`

# 提取最后从第四位开始的内容 去除 I: 
# 如果想，可以从输入中提取更多参数，或者利用 Bash 对字符串进行处理

echo "输入内容: ${input:3}"

content=$(echo "${input:3}" | grep '^http' | tr -d '\n' | jq -Rs '.')

# 构造请求体的 JSON 字符串
data="{\"url\": $content }"

echo -n "裁剪网页结果: "
# 发送 POST 请求

response=$(curl -s -X POST -H "Token: $token2" -H "Content-Type: application/json" -d "$data" -w "\n%{http_code}" "$host/ob/url")
status_code=$(echo "$response" | tail -n1)
if [ "$status_code" -eq 200 ]; then
    echo "请求成功"
    file_key=$(echo "$response" | head -n 1 | jq -r '.data.file_key')
    title=$(echo "$response"| head -n 1  | jq -r '.data.title')
    echo "Title: $title"
    echo "FileKey: $title"
else
    echo "请求失败"
fi
