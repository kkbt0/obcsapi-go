package wechat

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"obcsapi-go/gr"
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

type WorkWechatDecryptXMLData struct {
	ToUserName string `xml:"ToUserName"`
	AgentID    string `xml:"AgentID"`
	Encrypt    string `xml:"Encrypt"`
}

type MsgData struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	MsgId        int64  `xml:"MsgId"`
	AgentID      int    `xml:"AgentID"`
}

type WorkWechatXMLResponse struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATA    `xml:"ToUserName"`
	FromUserName CDATA    `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      CDATA    `xml:"MsgType"`
	Content      CDATA    `xml:"Content"`
}

type WorkWechatEncryptXMLResponse struct {
	Encrypt      string `xml:"Encrypt"`
	MsgSignature string `xml:"MsgSignature"`
	TimeStamp    string `xml:"TimeStamp"`
	Nonce        string `xml:"Nonce"`
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
	workWechatXMLResponse := WorkWechatXMLResponse{
		ToUserName:   CDATA{Value: viper.GetString("work_wechat_user_id")},
		FromUserName: CDATA{Value: viper.GetString("work_wechat_corpid")},
		CreateTime:   time.Now().Unix(),
		MsgType:      CDATA{Value: "text"},
		Content:      CDATA{Value: "Echo: " + msgContent.Content},
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
