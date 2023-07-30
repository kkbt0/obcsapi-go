package wechat

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"obcsapi-go/dao"
	"obcsapi-go/gr"
	"obcsapi-go/tools"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func WorkWechatVerifyURL(c *gin.Context) {
	wxcpt := NewWXBizMsgCrypt(
		viper.GetString("work_wechat_token"),
		viper.GetString("work_wechat_encoding_aeskey"),
		viper.GetString("work_wechat_receiverid"),
		XmlType)
	echoStr, cryptErr := wxcpt.VerifyURL(c.Query("msg_signature"), c.Query("timestamp"), c.Query("nonce"), c.Query("echostr"))
	if cryptErr != nil {
		log.Println("verifyUrl fail", cryptErr)
		return
	}
	log.Println("verifyUrl success echoStr")
	c.String(200, string(echoStr))
}

type MsgData struct {
	ToUserName   CDATA `xml:"ToUserName"`
	FromUserName CDATA `xml:"FromUserName"`
	CreateTime   int   `xml:"CreateTime"`
	MsgType      CDATA `xml:"MsgType"`
	MsgId        int   `xml:"MsgId"`
	AgentID      int   `xml:"AgentID"`
	// text
	Content CDATA `xml:"Content"`
	// image
	PicUrl CDATA `xml:"PicUrl"`
	// link
	Title       CDATA `xml:"Title"`
	Description CDATA `xml:"Description"`
	Url         CDATA `xml:"Url"`
}

type WorkWechatXMLResponse struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATA    `xml:"ToUserName"`
	FromUserName CDATA    `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      CDATA    `xml:"MsgType"`
	Content      CDATA    `xml:"Content"`
}

func WorkWechatMsgHandler(c *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}
	wxcpt := NewWXBizMsgCrypt(
		viper.GetString("work_wechat_token"),
		viper.GetString("work_wechat_encoding_aeskey"),
		viper.GetString("work_wechat_receiverid"),
		XmlType)
	xmlMsg, cryptErr := wxcpt.DecryptMsg(c.Query("msg_signature"), c.Query("timestamp"), c.Query("nonce"), bodyBytes)
	if cryptErr != nil {
		gr.ErrServerError(c, err)
		return
	}

	var msgContent MsgData
	err = xml.Unmarshal(xmlMsg, &msgContent)
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}
	// 处理
	r_str, _ := WorkWechatMsg2Obsidian(msgContent)
	workWechatXMLResponse := WorkWechatXMLResponse{
		ToUserName:   CDATA{Value: viper.GetString("work_wechat_user_id")},
		FromUserName: CDATA{Value: viper.GetString("work_wechat_corpid")},
		CreateTime:   time.Now().Unix(),
		MsgType:      CDATA{Value: "text"},
		Content:      CDATA{Value: r_str},
	}
	xmlBytes, err := xml.MarshalIndent(workWechatXMLResponse, "", "   ")
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}

	encryptMsg, cryptErr := wxcpt.EncryptMsg(string(xmlBytes), string(rune(time.Now().Unix())), c.Query("nonce"))
	if cryptErr != nil {
		log.Println("DecryptMsg fail", cryptErr)
		return
	}

	c.String(200, string(encryptMsg))
}

func WorkWechatMsg2Obsidian(msg MsgData) (string, error) {
	r_str := tools.NowRunConfig.WeChatMp.ReturnStr
	var err error
	if r_str == "" {
		r_str = "📩 已保存"
	}
	switch msg.MsgType.Value {
	case "text":
		r_str, err = WeChatTextAndVoice(msg.Content.Value)
	case "image":
		fileby, _ := dao.PicDownloader(msg.PicUrl.Value)
		file_key := fmt.Sprintf("%s%s.jpg", tools.NowRunConfig.DailyAttachmentDir(), tools.TimeFmt("20060102150405"))
		dao.ObjectStore(file_key, fileby)
		// 前端会监测 ![https://..](..) 将 http:// 放到 后面 ![..](https://..)
		// append_memos_in_daily(client, fmt.Sprintf("![%s](%s)", mp.Request.PicUrl, file_key))
		err = dao.DailyTextAppendMemos(fmt.Sprintf("![](%s)", file_key))
	case "link":
		err = dao.DailyTextAppendMemos(fmt.Sprintf("[%s](%s)<br>%s...", msg.Title.Value, msg.Url.Value, msg.Description.Value))
	default:
		r_str = "不支持的消息类型: " + msg.MsgType.Value
	}
	if err != nil {
		log.Println(err)
		r_str = "Error"
	}
	return r_str, nil
}

/////// Send Message

type Message struct {
	ToUser  string             `json:"touser"`
	MsgType string             `json:"msgtype"`
	AgentID int                `json:"agentid"`
	Text    MessageTextContent `json:"text"`
}
type MessageTextContent struct {
	Content string `json:"content"`
}

type MessageResponse struct {
	Errcode int `json:"errcode"`
}

func WorkWechatdSendMessage(message Message, twice bool) error {
	accessToken, err := GetWorkWechatAccessToken()
	if err != nil {
		return err
	}
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token="+accessToken, bytes.NewBuffer(messageJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 更新 access_token
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var messageResponse MessageResponse
	if err := json.Unmarshal(body, &messageResponse); err != nil {
		return err
	}
	// 42001  token 过期
	if messageResponse.Errcode == 42001 && !twice {
		_, err = FetchNewAccessToken()
		if err != nil {
			return err
		}
		return WorkWechatdSendMessage(message, true)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %s", resp.Status)
	}

	return nil
}

func WorkWechatSendText(text string) error {
	message := Message{
		ToUser:  viper.GetString("work_wechat_user_id"),
		MsgType: "text",
		Text:    MessageTextContent{Content: text},
		AgentID: viper.GetInt("work_wechat_agentid"),
	}
	err := WorkWechatdSendMessage(message, false)
	if err != nil {
		log.Println("Error sending message:", err)
		return err
	}
	return nil
}
