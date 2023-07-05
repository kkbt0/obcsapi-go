#!/bin/bash

content=$1

data='{
  "model": "gpt-3.5-turbo",
  "messages": [{"role": "user", "content": '"$content"'}],
  "temperature": 0.7
}'

curl -X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer sk-xxxxxxx" \
-d "$data" \
https://api.openai.com/v1/chat/chat/completions
