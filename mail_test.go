package main

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

func TestSmtpMain(t *testing.T) {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error: Fatal error config file: %s \n ", err))
	}
	mail_config := viper.GetStringMap("smtp_mail")
	fmt.Println(mail_config)
	err = emailSendToken()
	if err != nil {
		fmt.Println(err)
	}
	// sendMail("ObCSAPI 登录链接", "这是 Obsidian Cloud Storeage API 后台发送的登录链接，请谨慎保管。这是全权限 token ，会在设定的时间后失效。<br> <a href=\"https://www.ftls.xyz/index.html?address=api.domain.com&token=xxx\">https://www.ftls.xyz/index.html?address=api.domain.com&token=xxx</a>")
}
