package main

import (
	"fmt"
	"strings"

	"log"
	. "obcsapi-go/dao"
	"obcsapi-go/tools"

	"github.com/DanPlayer/timefinder"
	"github.com/gin-gonic/gin"
	"github.com/sidbusy/weixinmp"
)

var mp = weixinmp.New(tools.ConfigGetString("wechat_token"), tools.ConfigGetString("wechat_appid"), tools.ConfigGetString("wechat_secret"))

func WeChatMpHandlers(c *gin.Context) {
	log.Println("WeChat MP Run")
	openid := tools.ConfigGetString("wechat_openid") // OpenID
	if !mp.Request.IsValid(c.Writer, c.Request) {
		return
	}
	if mp.Request.FromUserName != openid {
		mp.ReplyTextMsg(c.Writer, "你好陌生人")
		log.Println("陌生人:", mp.Request.FromUserName)
		return
	}
	r_str := tools.NowRunConfig.WeChatMp.ReturnStr
	if r_str == "" {
		r_str = "📩 已保存"
	}
	var err error
	switch mp.Request.MsgType {
	case weixinmp.MsgTypeText: // 文字消息
		// 提醒任务判断
		// 初始化timefinder 对自然语言（中文）提取时间
		var segmenter = timefinder.New("./static/jieba_dict.txt,./static/" + tools.NowRunConfig.Reminder.ReminderDicionary)
		extract := segmenter.TimeExtract(mp.Request.Content) // 如果提取出了时间
		if strings.Contains(mp.Request.Content, "提醒我") && len(extract) != 0 {
			err = TextAppend("提醒任务.md", "\n"+extract[0].Format("20060102 1504 ")+mp.Request.Content)
			if err != nil {
				log.Println(err)
			}
			err = TextAppend(tools.NowRunConfig.DailyFileKeyTime(extract[0]), "\n- [ ] "+mp.Request.Content+" ⏳ "+extract[0].Format("2006-01-02 15:04"))
			r_str = "已添加至提醒任务:" + extract[0].Format("20060102 1504")
		} else {
			err = DailyTextAppendMemos(mp.Request.Content) //
		}
	case weixinmp.MsgTypeImage: // 图片消息
		fileby, _ := PicDownloader(mp.Request.PicUrl)
		file_key := fmt.Sprintf("%s%s.jpg", tools.NowRunConfig.DailyAttachmentDir(), tools.TimeFmt("20060102150405"))
		ObjectStore(file_key, fileby)
		// 前端会监测 ![https://..](..) 将 http:// 放到 后面 ![..](https://..)
		// append_memos_in_daily(client, fmt.Sprintf("![%s](%s)", mp.Request.PicUrl, file_key))
		err = DailyTextAppendMemos(fmt.Sprintf("![](%s)", file_key))
	case weixinmp.MsgTypeVoice: // 语言消息
		// 提醒任务判断
		// 初始化timefinder 对自然语言（中文）提取时间
		var segmenter = timefinder.New("./static/jieba_dict.txt,./static/" + tools.NowRunConfig.Reminder.ReminderDicionary)
		extract := segmenter.TimeExtract(mp.Request.Recognition)
		if strings.Contains(mp.Request.Recognition, "提醒我") && len(extract) != 0 {
			err = TextAppend("提醒任务.md", "\n"+extract[0].Format("20060102 1504 ")+mp.Request.Recognition)
			if err != nil {
				log.Println(err)
			}
			err = TextAppend(tools.NowRunConfig.DailyFileKeyTime(extract[0]), "\n- [ ] "+mp.Request.Recognition+" ⏳ "+extract[0].Format("2006-01-02 15:04"))
			r_str = "已添加至提醒任务:" + extract[0].Format("20060102 1504")
		} else if mp.Request.Recognition != "" {
			err = DailyTextAppendMemos("语音: " + mp.Request.Recognition) //
		} else {
			r_str = "没有识别到文字"
		}
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
