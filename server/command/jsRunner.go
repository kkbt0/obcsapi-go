package command

import (
	"fmt"
	"obcsapi-go/tools"
	"os"

	"github.com/dop251/goja"
)

// Example: command.RunJsByFile("script/demo.js", "{\"a\":1,\"b\":2}")
func RunJsByFile(filePath string, inputText string) (string, error) {

	// 读取 js 文件
	jsText, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("js run err: ReadFile failed")
	}
	// 使用解释器
	vm := goja.New()
	// 使用的函数
	// console.log("xxx")
	console := vm.NewObject()
	console.Set("log", js_log)
	// obcsapi.daily_append("xxx")
	obcsapi := vm.NewObject()
	obcsapi.Set("daily_append", js_obcsapi_daily_append)

	// 导入函数
	vm.Set("console", console)
	vm.Set("obcsapi", obcsapi)
	// 导入变量
	vm.Set("inputText", inputText)

	// 运行
	v, err := vm.RunString(string(jsText))
	if err != nil {
		return "", err
	}
	// 将结果转换成字符串返回
	result, ok := v.Export().(string)
	if !ok {
		return "", fmt.Errorf("js run err: Export failed , result type is not string or not defined")
	}
	tools.Debug("inputText:", inputText, "result:", result, "返回内容格式化为字符串 ok:", ok)
	return result, nil
}
