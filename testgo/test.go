package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"obcsapi-go/tools"
	"strings"
)

func main() {
	fmt.Println("Run")
	file, err := ioutil.ReadFile("1.txt")
	if err != nil {
		log.Println(err)
	}
	rawTodoList := strings.Split(string(file), "\n")
	var todoList []string
	for i := 0; i < len(rawTodoList); i++ {
		rawTodoList[i] = strings.ReplaceAll(rawTodoList[i], "\n", "")
		rawTodoList[i] = strings.ReplaceAll(rawTodoList[i], "\t", "")
		if strings.HasPrefix(rawTodoList[i], "20") {
			todoList = append(todoList, rawTodoList[i])
		}
	}
	fmt.Println("RawFile", string(file))
	fmt.Println("---todoList---")
	fmt.Println(todoList)
	fmt.Println("---checkList---")
	for _, v := range todoList {
		if strings.HasPrefix(v, tools.TimeFmt("20060102 1504")) {
			fmt.Println(v)
		}
	}
}
