package wechat

import (
	"log"

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
