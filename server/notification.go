package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"obcsapi-go/tools"
)

func WeChatTemplateMesseng(text string) {
	accessToken, err := mp.AccessToken.Fresh()
	if err != nil {
		log.Println(err)
	}
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + accessToken
	content := WeChatTemplateMessengDataVCPart{Value: text, Color: "#173177"}
	dataPart := WeChatTemplateMessengDataPart{Content: content}
	mainPart := WeChatTemplateMessengMainPart{
		Touser:     tools.ConfigGetString("wechat_openid"),
		TemplateId: tools.ConfigGetString("wechat_template_id"),
		Url:        "https://kkbt.gitee.io/obweb/#/Memos",
		TopColor:   "#FF0000",
		Data:       dataPart,
	}
	jsonBytes, err := json.Marshal(mainPart)
	if err != nil {
		log.Println(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println("Send Wechat Template Message:", string(body))

}

type WeChatTemplateMessengMainPart struct {
	Touser     string                        `json:"touser"`
	TemplateId string                        `json:"template_id"`
	Url        string                        `json:"url"`
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
