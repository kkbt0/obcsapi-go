package wechat

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"obcsapi-go/tools"
)

func WeChatTemplateMesseng(text string) error {
	accessToken, err := Mp.AccessToken.Fresh()
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
