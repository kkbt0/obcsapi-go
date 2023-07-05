#!/bin/bash
input=`bash sh/lastinput.sh`
command=$(echo "${input:3}" | jq -Rs '.') # remove I: 

ai_result=`bash sh/lib/ai.sh "${command}"`
# echo $ai_result

deal=$(echo ${ai_result} | jq -r ".choices[] | .message.content" | jq -sR 'gsub("\n"; "\n")')

json_payload='{"code":200,"data":{"type":"message-tbb","command":'${command}',"parts":[{"type":"md","text":'${deal}'},{"type":"button","text":"Resend","command":"@web-resend"},{"type":"button","text":"Save","command":"@web-save-memos"}]}}'

# Use jq to pretty-print the JSON payload
formatted_json=$(echo "$json_payload" | jq '.')

echo $formatted_json