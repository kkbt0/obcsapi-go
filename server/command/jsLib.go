package command

import (
	"fmt"

	"github.com/dop251/goja"
)

func js_log(call goja.FunctionCall) goja.Value {
	str := call.Argument(0)
	fmt.Println(str.String())
	return str
}
