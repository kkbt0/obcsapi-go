package tools

import (
	"bufio"
	"crypto/md5"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

//go:embed config.example.yaml
var configExample string

var YamlConfigMd5 string = GenerateMd5()

func CheckFiles() {
	log.Println("Check Need Files")
	_, err := os.Stat("config.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			err := os.WriteFile("config.yaml", []byte(configExample), 0666)
			if err != nil {
				log.Println("Write File Err", err)
			}
		} else {
			log.Panicln("Error: Stat config.yaml")
		}
	}
	_, err = os.Stat("config.run.json")

	if err != nil {
		if os.IsNotExist(err) {
			data, _ := json.Marshal(ExampleRunconfig())
			os.WriteFile("config.run.json", data, 0666)
		} else {
			log.Panicln("Error: Stat config.run.json")
		}
	}
	_, err = os.Stat("token/")
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("token/", 0777)
			os.Chmod("token/", 0777)
		} else {
			log.Panicln("Error: Stat token/")
		}
	}
	_, err = os.Stat("token/token1.json")
	if err != nil {
		if os.IsNotExist(err) {
			ModTokenFile(Token{TokenString: GengerateTokenString(32), GenerateTime: TimeFmt("2006-01-02 15:04:05"), LiveTime: "30s", VerifyMode: "Headers-Token"}, "./token/token1.json")
		} else {
			log.Println(err)
			log.Panicln("Error: Stat token/token1.json")
		}
	}
	_, err = os.Stat("token/token2.json")
	if err != nil {
		if os.IsNotExist(err) {
			time.Sleep(time.Duration(3) * time.Millisecond)
			ModTokenFile(Token{TokenString: GengerateTokenString(32), GenerateTime: TimeFmt("2006-01-02 15:04:05"), LiveTime: "876000h", VerifyMode: "Headers-Token"}, "./token/token2.json")
		} else {
			log.Println(err)
			log.Panicln("Error: Stat token/token2.json")
		}
	}
}

// 从配置中获取 参数
func ConfigGetString(parm string) string {
	return viper.GetString(parm)
}

// 从配置中获取 参数
func ConfigGetInt(parm string) int {
	return viper.GetInt(parm)
}

// Time fmt eg 2006-01-02 15:04:05
func TimeFmt(fmt string) string {
	// fmt.Println(time.Now().In(cstZone).Format("2006-01-02 15:04:05"))
	var cstZone = time.FixedZone("CST", 8*3600)
	return time.Now().In(cstZone).Format(fmt)
}

// obsidian 文件名非法字符 * " \ / < > : | ? 链接失效 # ^ [ ] | 替换为 _
func ReplaceUnAllowedChars(s string) string {
	unAllowedChars := "*\"\\/<>:|?#^[]|"
	for _, c := range unAllowedChars {
		s = strings.ReplaceAll(s, string(c), "_")
	}
	return s
}

// 和 downloads 除了保存文件名不一样 剩下都一样
func Downloader(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	reader := bufio.NewReaderSize(res.Body, 32*1024)
	file, err := os.Create("tem.file")
	writer := bufio.NewWriter(file)
	io.Copy(writer, reader)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf, err := os.ReadFile("tem.file")
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)[0:n]
}

// jwt secret + 日期 进行 MD5  保证来自服务器签发
func ObFileAccessToken() string {
	md5Str := md5.New()
	md5Str.Write([]byte(YamlConfigMd5 + TimeFmt("2006-01-02")))
	return hex.EncodeToString(md5Str.Sum(nil))
}

func InitViper() {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error: Fatal error config file: %s \n ", err))
	}
	// 配置	swagger 开关
	if viper.GetBool("swagger") {
		Debug("OBCSAPI_SWAGGER_DISABLE 删除,已开启 swagger")
		os.Unsetenv("OBCSAPI_SWAGGER_DISABLE")
	} else {
		Debug("OBCSAPI_SWAGGER_DISABLE 设置为 1,已关闭 swagger")
		os.Setenv("OBCSAPI_SWAGGER_DISABLE", "1")
	}
}

func UpdateViper() error {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	return viper.ReadInConfig()
}

func GenerateMd5() string {
	CheckFiles() // 程序最开始要执行的部分
	InitViper()
	md5 := md5.New()
	fileStr, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	md5.Write(fileStr)
	MD5Str := hex.EncodeToString(md5.Sum(nil))
	return MD5Str
}

func ExampleRunconfig() RunConfig {
	var exampleRunconfig RunConfig
	exampleRunconfig.Basic.DisableLogin = false
	exampleRunconfig.ObDaily = ObsidianDailyConfig{
		ObDailyDir:                        "日记/",
		ObDaily:                           "2006-01-02",
		ObDailyAttachmentDir:              "附件/",
		ObDailyAttachmentDirUnderDailyDir: true,
		ObOtherDataDir:                    "其他/",
	}
	exampleRunconfig.WeChatMp.ReturnStr = "📩 已保存，<a href='https://kkbt.gitee.io/web/'>点击查看今日笔记</a>"
	exampleRunconfig.Webdav = WebDavConfig{
		Server:     false,
		Username:   "testuser",
		Password:   "testpassword",
		ObLocalDir: "note/",
	}
	exampleRunconfig.ImageHosting = ImageHostingConfig{
		BaseURL:          "http://localhost:8900/images/",
		Prefix:           "200601/kkbt_",
		UseRawName:       true,
		RandomCharLength: 5,
	}
	exampleRunconfig.Reminder = ReminderConfig{
		DailyEmailRemderTime: "9999",
		ReminderDicionary:    "dictionary-200k.txt",
	}
	exampleRunconfig.Mention = MentionConfig{
		Tags: []string{"收藏"},
	}
	return exampleRunconfig
}
