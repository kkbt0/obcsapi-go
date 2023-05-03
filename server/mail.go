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
	token1, err := tools.GetToken("token1")
	if err != nil {
		return err
	}
	token2, err := tools.GetToken("token2")
	if err != nil {
		return err
	}

	// 发送邮件
	login_address := fmt.Sprintf("%s?backend_address=%s&token=%s", tools.ConfigGetString("front_url"), tools.ConfigGetString("backend_url"), token1.TokenString)
	content1 := "这是 Obsidian Cloud Storeage API 后台发送的登录链接，请谨慎保管。这是 token1 ，会在设定的时间后失效。<br> 登录链接: "
	content2 := fmt.Sprintf("<a href=\"%s\">%s</a>", login_address, login_address)
	token1_info := fmt.Sprintf("<br>Token1 权限1: %s ,生成时间 %s , 设定有效期: %s <br>", token1.TokenString, token1.GenerateTime, viper.GetString("token1_live_time"))
	token2_info := fmt.Sprintf("<br>Token2 只发送: %s ,生成时间 %s , 设定有效期: 无限 <br>", token2.TokenString, token2.GenerateTime)
	update_url := fmt.Sprintf("更新 token 链接 <a href=\"http://%s/api/sendtoken2mail\">http://%s/api/sendtoken2mail</a><br>或<a href=\"https://%s/api/sendtoken2mail\">https://%s/api/sendtoken2mail</a><br>", tools.ConfigGetString("backend_url"), tools.ConfigGetString("backend_url"), tools.ConfigGetString("backend_url"), tools.ConfigGetString("backend_url"))
	sendMail("ObCSAPI 登录链接", content1+content2+"<br>"+token1_info+token2_info+update_url)
	return nil
}

// 发送邮件 需要传入 主题 和 内容 ，其余配置选项均从配置中读取使用
func sendMail(subjct string, content string) error {

	config := tools.NowRunConfig.Mail

	// 配置邮件
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	em := email.NewEmail()
	em.From = fmt.Sprintf("%s <%s>", config.SenderName, config.SenderEmail)
	em.To = []string{config.ReceiverEmail}
	em.Subject = subjct
	// em.Text = content
	em.HTML = []byte(subjct + ":<br>" + content)

	// 发送邮件
	sender_addr := fmt.Sprintf("%s:%d", config.SmtpHost, config.Port) // smtpdm.aliyun.com:80"
	err := em.Send(sender_addr, smtp.PlainAuth("", config.UserName, config.Password, config.SmtpHost))
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Send successfully ... Subject: ", subjct)
	return nil
}
