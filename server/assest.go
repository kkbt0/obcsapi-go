package main

import (
	"fmt"
	"log"
	"obcsapi-go/tools"
	"time"

	"github.com/spf13/viper"
)

// 一些杂七杂八的函数 又不能放到 tools 里的

var version string = "v4.1.0"

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

func ChineseSegmenterTest() {
	var msg string
	var extract []time.Time

	msg = " 6月9日有一场show要去观看"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "后天早上10:35的会议，需要及时参与"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "明天下午三点的飞机，提醒我坐车"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "一个小时后提醒我喝水"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "明天早上8:00喊我起床"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "明天早上8点喊我起床"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "明早十点喊我喝水"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "明天早上十点喊我喝水"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "明天下午三点提醒我喝水"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "一天后提醒我喝水"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "一年后提醒我喝水"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "一个月后提醒我喝水"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "一月后提醒我喝水"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "我要住到大后天"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "我要住到明天"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "下个月到上个月再到这个月"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "我要住到明天下午三点十分"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "帮我预定明天凌晨3点的飞机"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "今天13:00的飞机"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "3月15号的飞机"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "昨天凌晨2点"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "十分钟后提醒我喝水"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)
}
