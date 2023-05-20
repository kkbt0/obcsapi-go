package talk

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
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
				if strings.HasPrefix(response, "Command ") {
					cmd := strings.TrimPrefix(response, "Command ")
					output, err := exec.Command("bash", "-c", cmd).Output()
					if err != nil {
						return fmt.Sprintf("执行命令时出错：%v", err)
					}
					if len(output) == 0 {
						return "命令已执行，无输出"
					}
					return string(output)
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
