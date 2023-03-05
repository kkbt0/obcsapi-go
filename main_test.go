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
	fmt.Println(mail_config["smtp_host"])
	fmt.Println(mail_config["port"])
	fmt.Println(mail_config["username"])
	fmt.Println(mail_config["password"])
	// TODO: 发送邮件
}

func TestRomanTokne(t *testing.T) {
	token := GengerateToken(99)
	fmt.Println(len(token))
	fmt.Println(token)
}

func TestModToken(t *testing.T) {
	token1 := Token{TokenString: GengerateToken(10), GenerateTime: GengerateToken(20)}
	ModTokenFile(token1, "./token/", "token1")
	token2 := Token{TokenString: GengerateToken(10), GenerateTime: GengerateToken(20)}
	ModTokenFile(token2, "./token/", "token2")
}
