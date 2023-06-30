#!/bin/bash
# 需要自建后端 PocketBase
# 用于创建 PocketBase 一条记录 代码仅供参考

# text: I: #发送 xxxx
input=`bash sh/lastinput.sh`
text=$(echo "${input:7}" | jq -Rs '.')

send2whisper() {
  url="https://pocketbase.com/api/collections/xxxxx/records"
  dt="{\"text\":$text,\"createdAt\":\"$(date -u +'%Y-%m-%dT%H:%M:%SZ')\"}"
  key=$(cat pktoken)
  hd="authorization: $key"
  
  response=$(curl -s -X POST -H "Content-Type: application/json" -H "$hd" -d "$dt" -w "%{http_code}" "$url")
  status_code=$(echo "$response" | tail -n1)
  
  if [[ "$status_code" != "200" ]]; then
    adminurl="https://pocketbase.com/api/admins/auth-with-password"
    admindt="{\"identity\":\"useremail\",\"password\":\"password\"}"
    
    response=$(curl -s -X POST -H "Content-Type: application/json" -d "$admindt" "$adminurl")
    token=$(echo "$response" | jq -r '.token')
    echo "$token" > pktoken
    
    key=$(cat pktoken)
    hd="authorization: $key"
    
    response=$(curl -s -X POST -H "Content-Type: application/json" -H "$hd" -d "$dt" -w "%{http_code}" "$url")
    status_code=$(echo "$response" | tail -n1)
    
    if [[ "$status_code" != "200" ]]; then
      echo "Failed 发送说说"
      return 1
    fi
  fi
  
  echo "Success 发送说说"
}

send2whisper "$text"
