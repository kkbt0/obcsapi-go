package main

import (
	"fmt"
	"log"
	"net/smtp"
	"obcsapi-go/tools"

	"github.com/jordan-wright/email"
)

// 发送邮件 需要传入 主题 和 内容 ，其余配置选项均从配置中读取使用
func SendMail(subjct string, content string) error {

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
