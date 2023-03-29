package main

import (
	"fmt"
	"log"
	"net/smtp"
	"obcsapi-go/tools"

	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
)

// 邮件发送 Tokne 包含两种 Token
func emailSendToken() error {
	// get token
	token1, err := GetToken("token1")
	if err != nil {
		return err
	}
	token2, err := GetToken("token2")
	if err != nil {
		return err
	}

	// 发送邮件
	login_address := fmt.Sprintf("%s?backend_address=%s&token=%s", tools.ConfigGetString("front_url"), tools.ConfigGetString("backend_url"), token1.TokenString)
	content1 := "这是 Obsidian Cloud Storeage API 后台发送的登录链接，请谨慎保管。这是全权限 token ，会在设定的时间后失效。<br> 登录链接: "
	content2 := fmt.Sprintf("<a href=\"%s\">%s</a>", login_address, login_address)
	token1_info := fmt.Sprintf("<br>Token1 全权限: %s ,生成时间 %s , 设定有效期: %s <br>", token1.TokenString, token1.GenerateTime, viper.GetString("token1_live_time"))
	token2_info := fmt.Sprintf("<br>Token2 只发送: %s ,生成时间 %s , 设定有效期: 无限 <br>", token2.TokenString, token2.GenerateTime)
	update_url := fmt.Sprintf("更新 token 链接 <a href=\"http://%s/api/sendtoken2mail\">http://%s/api/sendtoken2mail</a><br>或<a href=\"https://%s/api/sendtoken2mail\">https://%s/api/sendtoken2mail</a><br>", tools.ConfigGetString("backend_url"), tools.ConfigGetString("backend_url"), tools.ConfigGetString("backend_url"), tools.ConfigGetString("backend_url"))
	sendMail("ObCSAPI 登录链接", content1+content2+"<br>"+token1_info+token2_info+update_url)
	return nil
}

// 发送邮件 需要传入 主题 和 内容 ，其余配置选项均从配置中读取使用
func sendMail(subjct string, content string) error {
	// read configuration
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error: Fatal error config file: %s \n ", err))
	}
	mail_config := viper.GetStringMap("smtp_mail")

	// 配置邮件
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	em := email.NewEmail()
	em.From = fmt.Sprintf("%s <%s>", mail_config["mail_sender_name"], mail_config["mail_sender_address"])
	em.To = []string{mail_config["mail_send_to"].(string)}
	em.Subject = subjct
	// em.Text = content
	em.HTML = []byte(subjct + ":<br>" + content)

	// 发送邮件
	sender_addr := fmt.Sprintf("%s:%d", mail_config["smtp_host"], mail_config["port"]) // smtpdm.aliyun.com:80"
	err = em.Send(sender_addr, smtp.PlainAuth("", mail_config["username"].(string), mail_config["password"].(string), mail_config["smtp_host"].(string)))
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println("Send successfully ... Subject: ", subjct)
	return nil
}
