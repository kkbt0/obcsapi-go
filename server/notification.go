package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"obcsapi-go/dao"
	"obcsapi-go/tools"
	"strings"
	"time"
)

func WeChatTemplateMesseng(text string) error {
	accessToken, err := mp.AccessToken.Fresh()
	if err != nil {
		log.Println(err)
		return err
	}
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + accessToken
	content := WeChatTemplateMessengDataVCPart{Value: text, Color: "#173177"}
	dataPart := WeChatTemplateMessengDataPart{Content: content}
	mainPart := WeChatTemplateMessengMainPart{
		Touser:     tools.ConfigGetString("wechat_openid"),
		TemplateId: tools.ConfigGetString("wechat_template_id"),
		TopColor:   "#FF0000",
		Data:       dataPart,
	}
	jsonBytes, err := json.Marshal(mainPart)
	if err != nil {
		log.Println(err)
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Send Wechat Template Message:", string(body))
	return nil
}

type WeChatTemplateMessengMainPart struct {
	Touser     string                        `json:"touser"`
	TemplateId string                        `json:"template_id"`
	TopColor   string                        `json:"topcolor"`
	Data       WeChatTemplateMessengDataPart `json:"data"`
}
type WeChatTemplateMessengDataPart struct {
	Content WeChatTemplateMessengDataVCPart `json:"content"`
}
type WeChatTemplateMessengDataVCPart struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

// 每日邮件
func DailyEmailReminder() error {
	login_address := "https://note.ftls.xyz/web/"
	var ansList []string
	ansList = append(ansList, "默认笔记前端链接")
	ansList = append(ansList, fmt.Sprintf("<a href=\"%s\">%s</a>", login_address, login_address))
	ansList = append(ansList, "<h3>每日提醒</h3>\n")

	// 获取 今天 昨天 前天的日记 和 每日提醒.md 的未完成任务
	md1, err := dao.GetTextObject("每日提醒.md")
	if err != nil {
		log.Println(err)
	} else {
		md1List := strings.Split(md1, "\n")
		for i := 0; i < len(md1List); i++ {
			md1List[i] = strings.ReplaceAll(md1List[i], "\n", "")
			md1List[i] = strings.ReplaceAll(md1List[i], "\t", "")
			if strings.HasPrefix(md1List[i], "- [ ]") {
				ansList = append(ansList, md1List[i])
			}
		}
	}

	md3days := dao.Get3DaysList()
	for i := 0; i < len(md3days); i++ { // 取每天 md
		date := fmt.Sprintf("<h3>%s</h3>", time.Now().AddDate(0, 0, i-2).In(time.FixedZone("CST", 8*3600)).Format("2006-01-02"))
		ansList = append(ansList, date)
		dayMemosList := strings.Split(md3days[i], "\n")
		for j := 0; j < len(dayMemosList); j++ {
			dayMemosList[j] = strings.ReplaceAll(dayMemosList[j], "\n", "")
			dayMemosList[j] = strings.ReplaceAll(dayMemosList[j], "\t", "")
			if strings.HasPrefix(dayMemosList[j], "- [ ]") {
				ansList = append(ansList, dayMemosList[j])
			}
		}
	}

	for i := 0; i < len(ansList); i++ {
		ansList[i] = strings.ReplaceAll(ansList[i], "[ ]", "")
	}
	return tools.SendMail(fmt.Sprintf("Obcsapi 每日邮件提醒 (%d)", len(ansList)-6), strings.Join(ansList, "<br>"))
}

// 每分钟查询 发送到微信提醒 or Mail
func EveryMinReminder() error {
	md0, err := dao.GetTextObject("提醒任务.md")
	if err != nil {
		return err
	}
	var otherList []string // 保存其余未使用到的
	var ansList []string   // 要发送的行

	rawTodoList := strings.Split(md0, "\n")
	for i := 0; i < len(rawTodoList); i++ {
		rawTodoList[i] = strings.ReplaceAll(rawTodoList[i], "\n", "")
		rawTodoList[i] = strings.ReplaceAll(rawTodoList[i], "\t", "")
		if strings.HasPrefix(rawTodoList[i], tools.TimeFmt("20060102 1504")) {
			ansList = append(ansList, rawTodoList[i])
		} else {
			otherList = append(otherList, rawTodoList[i])
		}
	}
	for i := 0; i < len(ansList); i++ {
		ansList[i] = strings.ReplaceAll(ansList[i], "[ ]", "")
	}

	if len(ansList) != 0 { // 如果列表不空发送消息
		log.Println("Wechat Reminder", len(ansList))
		err = WeChatTemplateMesseng(strings.Join(ansList, "\n"))
		// 邮件提醒
		var emailAns []string

		for _, iter := range ansList {
			if strings.Contains(iter, "发邮件提醒我") {
				emailAns = append(emailAns, iter)
			}
		}

		if len(emailAns) != 0 {
			log.Println("Email Reminder", len(emailAns))
			if len(emailAns) == 1 {
				err = tools.SendMail(fmt.Sprintf("Obcsapi 提醒: %s", emailAns[0]), "一个提醒<br>"+emailAns[0])
			} else {
				err = tools.SendMail(fmt.Sprintf("Obcsapi 提醒 (%d)", len(emailAns)), strings.Join(emailAns, "<br>"))
			}
			if err != nil {
				log.Println(err)
			}
		}
		// 邮件提醒结束
	} else {
		return nil
	}
	if err != nil {
		return err
	}

	err = dao.TextAppend(tools.NowRunConfig.OtherDataDir()+"WeChatSended/"+tools.TimeFmt("200601")+".md", "\n"+strings.Join(ansList, "\n"))
	if err != nil {
		return err
	}
	err = dao.MdTextStore("提醒任务.md", strings.Join(otherList, "\n"))
	return err
}
