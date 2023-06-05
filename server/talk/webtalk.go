package talk

import (
	"bufio"
	"fmt"
	"log"
	"obcsapi-go/tools"
	"os"

	"github.com/gin-gonic/gin"
)

var ChatMode = 0 // default 0 = 对话/指令模式 ; 1 = 输入模式

func ChatText(text string) (string, error) {
	if ChatMode == 0 { // 对话指令模式
		return ChatTalk(text)
	} else if text == "对话模式" || text == "指令模式" || text == "命令模式" || text == "对话模式。" || text == "指令模式。" || text == "Talk" {
		ChatMode = 0
		return "对话模式，输入 退出 返回输入模式", nil
	} else {
		return GetReminderFromString(text)
	}
}

// 指令/对话模式 预设处理 如返回今日待办
func ChatTalk(input string) (string, error) {
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
	var output string
	if input == "输入模式" || input == "退出" || input == "exit" || input == "Exit" || input == "q" {
		ChatMode = 1
		output = "输入模式"
	} else {
		output = GetResponse(input)
	}

	writerOutput := bufio.NewWriter(file)
	writerOutput.WriteString(fmt.Sprintf("O: %s\n", output))
	writerOutput.Flush()

	return output, nil
}

type TalkStruct struct {
	Content string `json:"content"`
}

// @Summary 指令模式接口
// @Tags 前端
// @Security JWT
// @Accept json
// @Produce plain
// @Param json body TalkStruct true "TalkStruct"
// @Router /api/v1/talk [post]
func TalkHandler(c *gin.Context) {
	var talkStruct TalkStruct
	if c.ShouldBindJSON(&talkStruct) != nil {
		c.String(400, "参数错误")
		return
	}
	if talkStruct.Content == "" {
		c.String(400, "参数错误")
		return
	}
	r_str, err := ChatText(talkStruct.Content)
	if err != nil {
		c.Status(500)
		return
	}
	c.String(200, r_str)
}
