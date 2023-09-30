package talk

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"obcsapi-go/command"
	. "obcsapi-go/dao"
	"obcsapi-go/tools"
	"os"
	"os/exec"
	"strings"

	"github.com/DanPlayer/timefinder"
)

type Dialogue struct {
	Triggers  []string
	Responses []string
}

func loadDialoguesFromFile(filename string) ([]Dialogue, error) {

	file, err := os.Open(filename)
	if err != nil {
		log.Println("无法加载对话文件:", err)
		return nil, err
	}
	defer file.Close()

	var dialogues []Dialogue
	scanner := bufio.NewScanner(file)
	var dialogue Dialogue

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "I:") {
			if len(dialogue.Triggers) > 0 {
				dialogues = append(dialogues, dialogue)
			}
			dialogue = Dialogue{}
			trigger := strings.TrimSpace(strings.TrimPrefix(line, "I:"))
			dialogue.Triggers = append(dialogue.Triggers, trigger)
		} else if strings.HasPrefix(line, "O:") {
			response := strings.TrimSpace(strings.TrimPrefix(line, "O:"))
			dialogue.Responses = append(dialogue.Responses, response)
		}
	}

	if len(dialogue.Triggers) > 0 {
		dialogues = append(dialogues, dialogue)
	}

	if err := scanner.Err(); err != nil {
		log.Println("读取对话文件时发生错误:", err)
	}

	return dialogues, nil
}
func GetResponse(input string) string {
	dialogues, err := loadDialoguesFromFile("dialogues.txt")
	if err != nil {
		return "Load Dialogues From File Error"
	}
	for _, dialogue := range dialogues {
		for _, trigger := range dialogue.Triggers {
			if strings.Contains(input, trigger) {
				response := dialogue.Responses[randInt(0, len(dialogue.Responses))] // 随机
				if strings.HasPrefix(response, "Command ") {                        // Bash 运行
					cmd := strings.TrimPrefix(response, "Command ")
					output, err := exec.Command("bash", "-c", cmd).Output()
					if err != nil {
						return fmt.Sprintf("执行命令时出错：%v", err)
					}
					if len(output) == 0 {
						return "命令已执行，无输出"
					}
					return string(output)
				} else if strings.HasPrefix(response, "Lua ") { // Lua 运行
					scriptFilePath := strings.TrimPrefix(response, "Lua ")
					output, err := command.LuaRunner(scriptFilePath, input)
					if err != nil {
						return fmt.Sprintf("执行命令时出错：%v", err)
					}
					if len(output) == 0 {
						return "命令已执行，无输出"
					}
					return output
				}
				return response
			}
		}
	}
	return "抱歉，我无法理解。你可以重新表达或者问点其他的吗？或输入 退出 返回输入模式"
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// 提醒任务判断 如果没有识别出来 则保存为 Memos
func GetReminderFromString(text string) (string, error) {
	// 提醒任务判断 如果没有则保存为 Memos
	// 初始化timefinder 对自然语言（中文）提取时间
	r_str := tools.NowRunConfig.WeChatMp.ReturnStr
	if r_str == "" {
		r_str = "📩 已保存"
	}
	var err error

	if strings.Contains(text, "提醒我") {
		var segmenter = timefinder.New("./static/jieba_dict.txt,./static/" + tools.NowRunConfig.Reminder.ReminderDicionary)
		extract := segmenter.TimeExtract(text)
		tools.Debug("提取时间:", extract)
		if len(extract) != 0 {
			err = TextAppend("提醒任务.md", "\n"+extract[0].Format("20060102 1504 ")+text)
			if err != nil {
				log.Println(err)
			}
			err = TextAppend(tools.NowRunConfig.DailyFileKeyTime(extract[0]), "\n- [ ] "+text+" ⏳ "+extract[0].Format("2006-01-02 15:04"))
			r_str = "已添加至提醒任务:" + extract[0].Format("20060102 1504")
		} else {
			err = DailyTextAppendMemos(text)
			r_str = "监测到提醒任务，未能提取时间。已保存"
		}

	} else {
		err = DailyTextAppendMemos(text) //
	}
	return r_str, err
}
