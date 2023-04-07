package main

import (
	"fmt"
	"log"
	"obcsapi-go/tools"

	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

// 一些杂七杂八的函数 又不能放到 tools 里的

var version string = "v4.0.8"

func ShowConfig() {

	// Read configuration
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error: Fatal error config file: %s \n ", err))
	}

	// output configuration
	log.Println(viper.GetString("name"), version, "ConfigVersion", viper.GetString("version"), viper.GetString("description"))
	log.Println("https://gitee.com/kkbt/obcsapi-go")
	log.Println("Server Time:", tools.TimeFmt("2006-01-02 15:04"))
	log.Println("Tokne File Path:", viper.GetString("token_path"))
	log.Println("Run on", viper.GetString("host"))

	// 显示 Token
	token1, err := tools.GetToken("token1")
	if err != nil {
		panic(fmt.Errorf("error: Fatal error Get Token file: %s \n ", err))
	}
	token2, err := tools.GetToken("token2")
	if err != nil {
		panic(fmt.Errorf("error: Fatal error Get Token file: %s \n ", err))
	}
	log.Println("Your token is:")
	log.Println("Token1", token1.TokenString, "GenerateTime", token1.GenerateTime)
	log.Println("Token2", token2.TokenString)
	log.Println("Token1 用于前端，有有效期概念。Token2 用于第三方 API 调用，无限时间。也可以在配置文件中设置很长时间的 Token1有效期")
	log.Println("你可以访问 /api/sendtoken2mail 路径更新 Token1 ,如果你配置了邮箱服务，程序会将 Token 相关信息发送到指定邮箱")

}

// 定时任务
func RunCronJob() {
	log.Println("Starting cron job...")
	c := cron.New()
	c.AddFunc("1/60 * * * * ?", func() { // 每分钟执行一次
		log.Println("Run a scheduled task...")
		// 要执行的代码
		WeChatTemplateMesseng("abc")
	})
	c.Start()
}
