package main

import (
	"net/http"
	"obcsapi-go/tools"

	"github.com/gin-gonic/gin"
)

// 限制短时间多次请求
func LimitMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			http.Error(c.Writer, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			c.Abort()
			return
		}
		c.Next()
	}
}

func Token1AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		if !tools.VerifyToken1(c.Request.Header.Get("Token")) {
			c.JSON(401, gin.H{
				"code": 401,
				"msg":  "验证错误",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func Token2AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		if !tools.VerifyToken2(c.Request.Header.Get("Token")) {
			c.JSON(401, gin.H{
				"code": 401,
				"msg":  "验证错误",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
