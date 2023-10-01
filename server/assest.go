package main

import (
	"fmt"
	"log"
	"math/rand"
	"obcsapi-go/dao"
	"obcsapi-go/skv"
	"obcsapi-go/tools"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// 一些杂七杂八的函数 又不能放到 tools 里的

var version string = "v4.2.8"

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
	log.Println("Run on", viper.GetString("host"))

	// 显示 Token
	token1, err := tools.GetToken("./token/token1.json")
	if err != nil {
		panic(fmt.Errorf("error: Fatal error Get Token file: %s \n ", err))
	}
	token2, err := tools.GetToken("./token/token2.json")
	if err != nil {
		panic(fmt.Errorf("error: Fatal error Get Token file: %s \n ", err))
	}
	log.Println("Your token is:")
	log.Println("Token1", token1.TokenString, "GenerateTime", token1.GenerateTime)
	log.Println("Token2", token2.TokenString, "GenerateTime", token2.GenerateTime)
	log.Println("Token1, Token2 自动生成，用于第三方 API 调用。也可以在对应的文件中设置很长时间的有效期")
	log.Println("你可以设置更多 token，在配置文件中使用")

}

// fileKey, memosText
func RandomMemos() (string, string) {
	var usefullList []string
	list, err := dao.ListObject(tools.NowRunConfig.DailyDir())
	if err != nil {
		log.Println(err)
	}
	for _, item := range list {
		if strings.HasSuffix(item, ".md") {
			usefullList = append(usefullList, item)
		}
	}
	rand.Seed(time.Now().UnixNano())
	// random FileKey -> get file -> Spilter
	fileKey := usefullList[rand.Intn(len(usefullList))]
	memosList := MarkdownSpilter(skv.GetByFileKey(fileKey))
	// random Memos
	memosText := memosList[rand.Intn(len(memosList))]
	return fileKey, memosText
}
