package main

import (
	"log"
	"obcsapi-go/command"
	"obcsapi-go/tools"
	"os"

	"github.com/robfig/cron"
)

// 定时任务
func RunCronJob() {
	log.Println("Start scheduled tasks...")
	c := cron.New()
	c.AddFunc(tools.ConfigGetString("cron"), func() { // 每分钟执行一次
		// 定时执行 Lua 脚本
		err1 := LuaCronJob()
		if err1 != nil {
			log.Println(err1)
		}
		// 要执行的代码 提醒任务
		err2 := MessagesSend()
		if err2 != nil {
			log.Println(err2)
		}
	})
	c.Start()
}

func LuaCronJob() error {
	tools.Debug("Start Lua CronJob...")
	_, err := os.Stat("script/cron.lua")
	if err != nil {
		if os.IsNotExist(err) {
			tools.Debug("Error: Stat script/cron.lua")
		}
		return err
	}
	_, err = command.LuaRunner("script/cron.lua", "")
	return err
}

func MessagesSend() error {
	var err error = nil
	if tools.NowRunConfig.Reminder.DailyEmailRemderTime == tools.TimeFmt("1504") {
		err = DailyEmailReminder()
	}
	if err != nil {
		log.Println(err)
	}
	err = EveryMinReminder()
	return err
}
