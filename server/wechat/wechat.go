package wechat

import (
	"bufio"
	"fmt"
	"os"

	"log"
	. "obcsapi-go/dao"
	"obcsapi-go/talk"
	"obcsapi-go/tools"

	"github.com/gin-gonic/gin"
	"github.com/sidbusy/weixinmp"
)

var Mp = weixinmp.New(tools.ConfigGetString("wechat_token"), tools.ConfigGetString("wechat_appid"), tools.ConfigGetString("wechat_secret"))

var WeChatMode = 1 // default 0 = 对话/指令模式 ; 1 = 输入模式

func WeChatMpHandlers(c *gin.Context) {
	log.Println("WeChat MP Run")
	openid := tools.ConfigGetString("wechat_openid") // OpenID
	if !Mp.Request.IsValid(c.Writer, c.Request) {
		return
	}
	if Mp.Request.FromUserName != openid {
		Mp.ReplyTextMsg(c.Writer, "你好陌生人")
		log.Println("陌生人:", Mp.Request.FromUserName)
		return
	}
	r_str := tools.NowRunConfig.WeChatMp.ReturnStr
	if r_str == "" {
		r_str = "📩 已保存"
	}
	var err error
	switch Mp.Request.MsgType {
	case weixinmp.MsgTypeText: // 文字消息
		r_str, err = WeChatTextAndVoice(Mp.Request.Content)
	case weixinmp.MsgTypeImage: // 图片消息
		fileby, _ := PicDownloader(Mp.Request.PicUrl)
		file_key := fmt.Sprintf("%s%s.jpg", tools.NowRunConfig.DailyAttachmentDir(), tools.TimeFmt("20060102150405"))
		ObjectStore(file_key, fileby)
		// 前端会监测 ![https://..](..) 将 http:// 放到 后面 ![..](https://..)
		// append_memos_in_daily(client, fmt.Sprintf("![%s](%s)", mp.Request.PicUrl, file_key))
		err = DailyTextAppendMemos(fmt.Sprintf("![](%s)", file_key))
	case weixinmp.MsgTypeVoice: // 语言消息
		if Mp.Request.Recognition != "" {
			r_str, err = WeChatTextAndVoice(Mp.Request.Recognition)
		} else {
			r_str = "没有识别到文字"
		}
	case weixinmp.MsgTypeLocation: // 位置消息
		err = DailyTextAppendMemos(fmt.Sprintf("位置信息: 位置 %s <br>经纬度( %f , %f )", Mp.Request.Label, Mp.Request.LocationX, Mp.Request.LocationY))
	case weixinmp.MsgTypeLink: // 链接消息
		err = DailyTextAppendMemos(fmt.Sprintf("[%s](%s)<br>%s...", Mp.Request.Title, Mp.Request.Url, Mp.Request.Description))
	case weixinmp.MsgTypeVideo:
		r_str = "不支持的视频消息"
	default:
		r_str = "未知消息"
	}
	if err != nil {
		log.Println(err)
		r_str = "Error"
	}
	Mp.ReplyTextMsg(c.Writer, r_str)
}

func WeChatTextAndVoice(text string) (string, error) {
	if WeChatMode == 0 { // 对话指令模式
		return WeChatTalk(text)
	} else if text == "对话模式" || text == "指令模式" || text == "命令模式" || text == "对话模式。" || text == "指令模式。" || text == "Talk" {
		WeChatMode = 0
		return "对话模式，输入 退出 返回输入模式", nil
	} else {
		return talk.GetReminderFromString(text)
	}
}

// 指令/对话模式 预设处理 如返回今日待办
func WeChatTalk(input string) (string, error) {
	//打开对话日志文件，如果不存在则创建
	date := tools.TimeFmt("20060102")
	file, err := os.OpenFile("./log/dialogues."+date+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("打开文件失败！", err)
		return "", err
	}
	defer file.Close()
	//写入数据
	writerInput := bufio.NewWriter(file)
	writerInput.WriteString(fmt.Sprintf("I: %s\n", input))
	writerInput.Flush()

	// 根据输入添加自定义逻辑，生成适当的回复
	// todo 返回今日待办
	var output string
	if input == "输入模式" || input == "退出" || input == "exit" || input == "Exit" || input == "q" {
		WeChatMode = 1
		output = "输入模式"
	} else {
		output = talk.GetResponse(input)
	}

	writerOutput := bufio.NewWriter(file)
	writerOutput.WriteString(fmt.Sprintf("O: %s\n", output))
	writerOutput.Flush()

	return output, nil
}
