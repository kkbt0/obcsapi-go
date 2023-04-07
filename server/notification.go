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
	token1, err := tools.GetToken("token1")
	if err != nil {
		log.Println(err)
	}
	login_address := fmt.Sprintf("%s?backend_address=%s&token=%s", tools.ConfigGetString("front_url"), tools.ConfigGetString("backend_url"), token1.TokenString)
	var ansList []string
	ansList = append(ansList, "登录链接")
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
	return sendMail("Obcsapi 每日邮件提醒", strings.Join(ansList, "<br>"))
}

// 每分钟查询 发送到微信提醒
func WechatMpReminder() error {
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
		err = WeChatTemplateMesseng(strings.Join(ansList, "\n"))
	} else {
		return nil
	}

	if err != nil {
		return err
	}
	err = dao.TextAppend(tools.ConfigGetString("ob_daily_other_dir")+"WeChatSended/"+tools.TimeFmt("200601")+".md", "\n"+strings.Join(ansList, "\n"))
	if err != nil {
		return err
	}
	err = dao.MdTextStore("提醒任务.md", strings.Join(otherList, "\n"))
	return err
}
