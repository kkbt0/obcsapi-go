package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("../config.yaml")
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
