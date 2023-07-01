package main

import (
	"log"
	"obcsapi-go/gr"
	"obcsapi-go/tools"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 限制短时间多次请求
func LimitMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			gr.TooManyRequests(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

func LimitLoginMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		if !loginLimter.Allow() {
			gr.TooManyRequests(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

// 配置文件可覆盖默认传入 defaultTokenFilePath
//
// 会根据 token.json 中配置的验证方式验证
func TokenAuthMiddleware(defaultTokenFilePath string) func(c *gin.Context) {
	return func(c *gin.Context) {
		// 处理配置默认为空的情况
		tokenFilePath := tools.RequestURIDefineOrDefalutTokenFilePath(c.Request.RequestURI, defaultTokenFilePath)
		// 读取 正确的token
		rightToken, err := tools.GetToken(tokenFilePath)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "服务器验证配置错误",
			})
			c.Abort()
			return
		}
		tools.Debug("Right Token: ", rightToken)
		var inToken string
		switch rightToken.VerifyMode {
		case "Headers-Authorization":
			inToken = c.Request.Header.Get("Authorization")
		case "Headers-Token":
			inToken = c.Request.Header.Get("Token")
		case "Query-token":
			inToken = c.Query("token")
		default:
			inToken = c.Request.Header.Get("Token")
		}
		tools.Debug("InTokenVerifyMode", rightToken.VerifyMode, "InTokenStr", inToken)

		allow := false
		if inToken == rightToken.TokenString && rightToken.IsLiving() {
			// 相符且存活
			allow = true
		} else if inToken == rightToken.TokenString {
			// 相符且不存活
			rightToken.Update()
			tools.ModTokenFile(rightToken, tokenFilePath)
		}
		// 不通过
		if !allow {
			c.JSON(401, gin.H{
				"code": 401,
				"msg":  "验证错误",
			})
			c.Abort()
			return
		}
		// 通过了
		c.Next()
	}
}

// 从 Authorization 头获取 ，只验证格式 Token {{tokenxxxStr}}
func StandardTokenAuthMiddleware(defaultTokenFilePath string) func(c *gin.Context) {
	return func(c *gin.Context) {
		// 处理配置默认为空的情况
		tokenFilePath := tools.RequestURIDefineOrDefalutTokenFilePath(c.Request.RequestURI, defaultTokenFilePath)

		// 标准头格式 Authorization: Token xxxxstr 会去除 `Token ` 这个前缀
		// 不通过
		if !tools.VerifyTokenByFilePath(tokenFilePath, strings.Replace(c.Request.Header.Get("Authorization"), "Token ", "", 1)) {
			tools.Debug("Authorization:", c.Request.Header.Get("Authorization"))
			c.JSON(401, gin.H{
				"code": 401,
				"msg":  "验证错误",
			})
			c.Abort()
			return
		}
		// 通过了
		c.Next()
	}
}

// 把 token 传给下一个处理者
func SpecialTokenMiddleware(defaultTokenFilePath string) func(c *gin.Context) {
	return func(c *gin.Context) {
		// 处理配置默认为空的情况
		tokenFilePath := tools.RequestURIDefineOrDefalutTokenFilePath(c.Request.RequestURI, defaultTokenFilePath)
		tools.Debug("[Middleware] Right Token Path: ", tokenFilePath)
		// 找到 Token 了 传递值
		c.Set("tokenfilepath", tokenFilePath)
		c.Next()
	}
}

func AllowOPTIONS() func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Status(200)
			c.Abort()
			return
		}
		c.Next()
	}
}

// 实验性质功能开关
func ExperimentalFeatures() func(c *gin.Context) {
	return func(c *gin.Context) {
		if !viper.GetBool("experimental_features") {
			gr.ErrNotFound(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

// configName 不配置默认 flase rev 翻转
func ConfigAllow(configName string, rev bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		allow := viper.GetBool(configName) // 配置文件是否允许
		if rev {                           // 如果翻转
			allow = !allow
		}
		if !allow {
			gr.ErrNotFound(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
