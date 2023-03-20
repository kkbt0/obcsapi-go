package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sidbusy/weixinmp"
)

func wechatmpfunc(w http.ResponseWriter, r *http.Request) {
	log.Println("WeChat MP Run")
	openid := ConfigGetString("wechat_openid") // OpenID
	mp := weixinmp.New(ConfigGetString("wechat_token"), ConfigGetString("wechat_appid"), ConfigGetString("wechat_secret"))
	if !mp.Request.IsValid(w, r) {
		return
	}
	if mp.Request.FromUserName != openid {
		mp.ReplyTextMsg(w, "你不是恐咖兵糖")
		log.Println("陌生人:", mp.Request.FromUserName)
		return
	}
	r_str := "📩 已保存，<a href='https://note.ftls.xyz/#/ZK/202209050658'>点击查看今日笔记</a>"
	client, err := get_client()
	if err != nil {
		log.Println(err)
	}
	switch mp.Request.MsgType {
	case weixinmp.MsgTypeText: // 文字消息
		err = append_memos_in_daily(client, mp.Request.Content) //
	case weixinmp.MsgTypeImage: // 图片消息
		fileby, _ := downloader(mp.Request.PicUrl)
		file_key := fmt.Sprintf("日志/附件/%s/%s.jpg", timeFmt("200601"), timeFmt("20060102150405"))
		store(client, file_key, fileby)
		// 前端会监测 ![https://..](..) 将 http:// 放到 后面 ![..](https://..)
		// append_memos_in_daily(client, fmt.Sprintf("![%s](%s)", mp.Request.PicUrl, file_key))
		err = append_memos_in_daily(client, fmt.Sprintf("![](%s)", file_key))
	case weixinmp.MsgTypeVoice: // 语言消息
		err = append_memos_in_daily(client, fmt.Sprintf("语音: %s", mp.Request.Recognition))
	case weixinmp.MsgTypeLocation: // 位置消息
		err = append_memos_in_daily(client, fmt.Sprintf("位置信息: 位置 %s <br>经纬度( %f , %f )", mp.Request.Label, mp.Request.LocationX, mp.Request.LocationY))
	case weixinmp.MsgTypeLink: // 链接消息
		err = append_memos_in_daily(client, fmt.Sprintf("[%s](%s)<br>%s...", mp.Request.Title, mp.Request.Url, mp.Request.Description))
	case weixinmp.MsgTypeVideo:
		r_str = "不支持的视频消息"
	default:
		r_str = "未知消息"
	}
	if err != nil {
		log.Println(err)
		r_str = "Error"
	}
	mp.ReplyTextMsg(w, r_str)
}
