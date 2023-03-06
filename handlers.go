package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Captcha struct {
	Captcha string
}

// NewCaptcha 生成或更新 token 邮件发送登录链接 直接附带 token
func SendTokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.RequestURI)
	log.Println("Succeed Send Token")

	// 修改 Token1
	ModTokenFile(Token{TokenString: GengerateToken(32), GenerateTime: timeFmt("2006-01-02 15:04:05")}, "token1")
	// 发送所有 Token 消息
	emailSendToken()
	fmt.Fprintf(w, "Succeed Send Token")
}

// 验证 Token 1 有效性

func VerifyToken1Handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.RequestURI)
	// 解析 token json {"token":"sometoken1"}
	decoder := json.NewDecoder(r.Body)
	var tokenFromJSON TokenFromJSON
	err := decoder.Decode(&tokenFromJSON)
	if err != nil {
		fmt.Println("JSON Decoder Error:", err)
	}
	if VerifyToken1(tokenFromJSON) {
		fmt.Fprintf(w, "a right Token")
	} else {
		fmt.Fprintf(w, "a error Token")
	}
}

// !!!deprecated!!! 弃用 邮件获取的验证码 前端输入 提交处理
func VerifyCaptchaHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.RequestURI)
	// 对传入的 token json 解析判断
	decoder := json.NewDecoder(r.Body)
	var captcha Captcha
	err := decoder.Decode(&captcha)
	if err != nil {
		fmt.Println(err)
	}
	if captcha.Captcha == "right_Captcha" {
		fmt.Fprintf(w, "a Tem Token")
	}
}
