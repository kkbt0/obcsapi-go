package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func ShowConfig() {
	// Read configuration
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error: Fatal error config file: %s \n ", err))
	}

	// output configuration
	log.Println(viper.GetString("name"), viper.GetString("version"), viper.GetString("description"))
	log.Println("Server Time:", timeFmt("2006-01-02 15:04"))
	log.Println("Tokne File Path:", viper.GetString("token_path"))
	log.Println("Run on", viper.GetString("host"))
}

// 从配置中获取 参数
func ConfigGetString(parm string) string {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error: Fatal error config file: %s \n ", err))
	}
	return viper.GetString(parm)
}

// Time fmt eg 2006-01-02 15:04:05
func timeFmt(fmt string) string {
	// fmt.Println(time.Now().In(cstZone).Format("2006-01-02 15:04:05"))
	var cstZone = time.FixedZone("CST", 8*3600)
	return time.Now().In(cstZone).Format(fmt)
}

//  obsidian 文件名非法字符 * " \ / < > : | ? 链接失效 # ^ [ ] | 替换为 _
func ReplaceUnAllowedChars(s string) string {
	unAllowedChars := "*\"\\/<>:|?#^[]|"
	fmt.Println(unAllowedChars)
	for _, c := range unAllowedChars {
		s = strings.ReplaceAll(s, string(c), "_")
	}
	return s
}
