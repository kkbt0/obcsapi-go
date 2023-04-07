package main

import (
	"log"
	"obcsapi-go/tools"

	"github.com/robfig/cron"
)

// 定时任务
func RunCronJob() {
	c := cron.New()
	c.AddFunc("1/60 * * * * ?", func() { // 每分钟执行一次
		log.Println("Run a scheduled task...")
		// 要执行的代码
		err := MessagesSend()
		if err != nil {
			log.Println(err)
		}
	})
	c.Start()
}

func MessagesSend() error {
	var err error = nil
	if tools.ConfigGetString("email_reminder_time") == tools.TimeFmt("1504") {
		err = DailyEmailReminder()
	}
	if err != nil {
		log.Println(err)
	}
	err = WechatMpReminder()
	return err
}
