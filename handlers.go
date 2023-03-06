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

// NewCaptcha 生成 token 生成验证码 邮件发送 or 邮件发送登录链接 直接附带 token
func SendTokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.RequestURI)
	log.Println("Succeed Send Token")
	fmt.Fprintf(w, "Succeed Send Token")
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
