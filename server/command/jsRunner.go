package command

import (
	"fmt"
	"log"
	"os"

	"github.com/dop251/goja"
)

// Example: command.RunJsByFile("script/demo.js", "{\"a\":1,\"b\":2}")
func RunJsByFile(filePath string, inputText string) {

	// 读取 js 文件
	jsText, err := os.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	// 使用解释器
	vm := goja.New()
	// 导入函数
	// console.log("xxx")
	console := vm.NewObject()
	console.Set("log", js_log)

	vm.Set("console", console)
	// 导入变量
	vm.Set("inputText", inputText)

	// 运行
	v, err := vm.RunString(string(jsText))
	if err != nil {
		log.Println(err)
		return
	}
	// Export：将结果转换成go基础类型
	result, ok := v.Export().(string)
	if !ok {
		log.Println("Export failed , result type is not string or not defined")
		return
	}
	fmt.Printf("result:%s,ok:%t\n", result, ok)
}
