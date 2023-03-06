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
		fmt.Println(err)
	}
	if mp.Request.MsgType == weixinmp.MsgTypeText { // 文字消息
		append_memos_in_daily(client, mp.Request.Content) //
		mp.ReplyTextMsg(w, r_str)
	} else if mp.Request.MsgType == weixinmp.MsgTypeImage { // 图片消息
		fileby, _ := downloader(mp.Request.PicUrl)
		file_key := fmt.Sprintf("日志/附件/%s/%s.jpg", timeFmt("200601"), timeFmt("200601021504"))
		store(client, file_key, fileby)
		append_memos_in_daily(client, fmt.Sprintf("![%s](%s)", mp.Request.PicUrl, file_key))
		mp.ReplyTextMsg(w, r_str)
	} else if mp.Request.MsgType == weixinmp.MsgTypeVoice { // 语言消息
		append_memos_in_daily(client, fmt.Sprintf("语音: %s", mp.Request.Recognition))
		mp.ReplyTextMsg(w, r_str)
	} else if mp.Request.MsgType == weixinmp.MsgTypeLocation { // 位置消息
		append_memos_in_daily(client, fmt.Sprintf("位置信息: 位置 %s <br>经纬度( %f , %f )", mp.Request.Label, mp.Request.LocationX, mp.Request.LocationY))
		mp.ReplyTextMsg(w, r_str)
	} else if mp.Request.MsgType == weixinmp.MsgTypeLink { // 链接消息
		append_memos_in_daily(client, fmt.Sprintf("[%s](%s)<br>%s...", mp.Request.Title, mp.Request.Url, mp.Request.Description))
		mp.ReplyTextMsg(w, r_str)
	} else if mp.Request.MsgType == weixinmp.MsgTypeVideo {
		mp.ReplyTextMsg(w, "不支持的视频消息")
	} else {
		mp.ReplyTextMsg(w, "未知消息")
	}
}
