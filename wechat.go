package main

import (
	"fmt"

	"log"
	. "obcsapi-go/dao"
	"obcsapi-go/tools"

	"github.com/gin-gonic/gin"
	"github.com/sidbusy/weixinmp"
)

func WeChatMpHandlers(c *gin.Context) {
	log.Println("WeChat MP Run")
	openid := tools.ConfigGetString("wechat_openid") // OpenID
	mp := weixinmp.New(tools.ConfigGetString("wechat_token"), tools.ConfigGetString("wechat_appid"), tools.ConfigGetString("wechat_secret"))
	if !mp.Request.IsValid(c.Writer, c.Request) {
		return
	}
	if mp.Request.FromUserName != openid {
		mp.ReplyTextMsg(c.Writer, "你不是恐咖兵糖")
		log.Println("陌生人:", mp.Request.FromUserName)
		return
	}
	r_str := tools.ConfigGetString("wechat_return_str")
	var err error
	switch mp.Request.MsgType {
	case weixinmp.MsgTypeText: // 文字消息
		err = DailyTextAppendMemos(mp.Request.Content) //
	case weixinmp.MsgTypeImage: // 图片消息
		fileby, _ := PicDownloader(mp.Request.PicUrl)
		file_key := fmt.Sprintf("%s%s/%s.jpg", tools.ConfigGetString("ob_daily_attachment_dir"), tools.TimeFmt("200601"), tools.TimeFmt("20060102150405"))
		ObjectStore(file_key, fileby)
		// 前端会监测 ![https://..](..) 将 http:// 放到 后面 ![..](https://..)
		// append_memos_in_daily(client, fmt.Sprintf("![%s](%s)", mp.Request.PicUrl, file_key))
		err = DailyTextAppendMemos(fmt.Sprintf("![](%s)", file_key))
	case weixinmp.MsgTypeVoice: // 语言消息
		err = DailyTextAppendMemos(fmt.Sprintf("语音: %s", mp.Request.Recognition))
	case weixinmp.MsgTypeLocation: // 位置消息
		err = DailyTextAppendMemos(fmt.Sprintf("位置信息: 位置 %s <br>经纬度( %f , %f )", mp.Request.Label, mp.Request.LocationX, mp.Request.LocationY))
	case weixinmp.MsgTypeLink: // 链接消息
		err = DailyTextAppendMemos(fmt.Sprintf("[%s](%s)<br>%s...", mp.Request.Title, mp.Request.Url, mp.Request.Description))
	case weixinmp.MsgTypeVideo:
		r_str = "不支持的视频消息"
	default:
		r_str = "未知消息"
	}
	if err != nil {
		log.Println(err)
		r_str = "Error"
	}
	mp.ReplyTextMsg(c.Writer, r_str)
}
