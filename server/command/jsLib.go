package command

import (
	"fmt"
	"log"
	"obcsapi-go/dao"

	"github.com/dop251/goja"
)

// console.log("xxx")
func js_log(call goja.FunctionCall) goja.Value {
	str := call.Argument(0)
	fmt.Println(str.String())
	return str
}

// obcsapi.daily_append("xxx")
func js_obcsapi_daily_append(call goja.FunctionCall) goja.Value {
	str := call.Argument(0)
	fmt.Println(str.String())
	err := dao.DailyTextAppend(str.String())
	if err != nil {
		log.Println(err)
	}
	// TODO: 可以返回异常
	return str
}
