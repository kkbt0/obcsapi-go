package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

func Greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func BaseHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("./template/404.html")
	if err != nil {
		log.Panicln("Template File Error:", err)
		return
	}
	indexInfo := IndexInfo{Title: "404", Content: "404 Not Found"}
	tpl.Execute(w, indexInfo)
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
	if VerifyToken1(tokenFromJSON.TokenString) {
		fmt.Fprintf(w, "a right Token")
	} else {
		fmt.Fprintf(w, "a error Token")
	}
}
