#!/bin/bash

# 该脚本需要在配置中 allow_general_all_post: true 
# 前端显示 allow_wiki_link_all: true
# 功能： 收藏文字到指定文件夹
# 定义配置 一般来说 host token2 无需更改  file_key 可根据需要自定义

host="http://localhost:8900"
token2=`cat token/token2.json | jq -r .token` 
# 注意不能有不允许的文件名字符 如 : 
current_time=$(date +"%Y-%m-%d-%H%M%S")
# 文件名，可以从 input 中获取 目前只是取了秒级时间
file_key="收藏文件夹/$current_time.md"


# 获取最后的输入
input=`bash sh/lastinput.sh`

# 提取最后从第四位开始的内容 去除 I: 
# 如果想，可以从输入中提取更多参数，或者利用 Bash 对字符串进行处理
# content=${input:3}
# text: I: #收藏 https:/xxx.com/index.html
content=$(echo "${input:7}" | jq -Rs '.')

# 构造请求体的 JSON 字符串
data="{\"content\": $content, \"mod\": \"cover\", \"file_key\": \"$file_key\"}"

echo "覆盖 ${file_key} 结果: "
# 发送 POST 请求
curl -X POST \
     -H "Token: $token2" \
     -H "Content-Type: application/json" \
     -d "$data" \
     "$host/ob/generalall"

memos_data="{\"content\": \"![[$file_key]]\"}"
echo -e "\n添加 Memos 结果: "
curl -X POST \
     -H "Token: $token2" \
     -H "Content-Type: application/json" \
     -d "$memos_data" \
     "$host/ob/general"
