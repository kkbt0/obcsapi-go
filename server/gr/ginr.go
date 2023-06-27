package gr

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type RJson struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
type H map[string]any

func RJSON(c *gin.Context, err error, status_code int, code int, msg string, data any) {
	if err != nil {
		log.Println(err)
	}
	switch c.Query("field") {
	case "code":
		c.String(status_code, fmt.Sprintf("%d", code))
	case "msg":
		c.String(status_code, msg)
	case "data":
		c.JSON(status_code, data)
	default:
		c.JSON(status_code, RJson{
			Code: code,
			Msg:  msg,
			Data: data,
		})
	}

}

func ErrBindJSONErr(c *gin.Context) {
	RJSON(c, nil, 400, 400, "JSON 格式错误", H{})
}

func ErrServerError(c *gin.Context, err error) {
	RJSON(c, err, 500, 500, "Server Error", H{})
}
func Success(c *gin.Context) {
	RJSON(c, nil, 200, 200, "Success", H{})
}

func ErrEmpty(c *gin.Context) {
	RJSON(c, nil, 400, 400, "内容为空", H{})
}

func ErrAuth(c *gin.Context) {
	RJSON(c, nil, 401, 401, "Auth Error", H{})
}

func ErrNotFound(c *gin.Context) {
	RJSON(c, nil, 404, 404, "Not Found", H{})
}
