#!/bin/bash
# 获取 token 值 -r 代表输出原始字符串
result=`cat token/token1.json | jq -r .token` 
echo "${result}"