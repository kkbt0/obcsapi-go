// 可用变量 inputText
// 结果变量 result


// 前端传入字符串
console.log(inputText)

// 转为 json
inputJson = JSON.parse(inputText)
console.log(inputJson.text1)
console.log(inputJson.text2)

// 返回的字符串
result = "输入的字符串: " + inputText + "\n返回: \ntext1: " + inputJson.text1 + "\ntext2: " + inputJson.text2

// 增加到今日
obcsapi.daily_append(result)