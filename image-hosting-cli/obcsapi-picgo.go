package main

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"

	"net/http"

	"time"

	"golang.design/x/clipboard"
)

var version string = "v1.0"

//go:embed obcsapi-picgo.ini
var configExample string

var config []string

// https://picgo.github.io/PicGo-Core-Doc/zh/guide/commands.html#use
// upload|u
func main() {
	runPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Panicln(err)
	}
	// 初始化配置文件
	configByte, err := os.ReadFile(path.Join(runPath, "obcsapi-picgo.ini"))
	if err != nil {
		log.Println("配置文件读取失败，生成配置文件，初始化完成")
		file, _ := os.Create(path.Join(runPath, "obcsapi-picgo.ini"))
		file.WriteString(configExample)
		defer file.Close()
		return
	}
	configStr := string(configByte)
	config = strings.Split(configStr, "\n")
	for i := 0; i < len(config); i++ {
		config[i] = strings.TrimSpace(config[i])
	}
	if len(config) < 3 {
		log.Println("配置文件参数不足")
	}

	if len(os.Args) == 1 {
		fmt.Printf("Obcsapi-ImagesHost-PicGo-Core %s  Help:\n", version)
		fmt.Println("upload    ->    Upload WindowsClipboard Picture(only png)")
		fmt.Println("upload test1.jpg test2.jpg    ->    Upload Pictures")
		log.Println("缺少参数")
		return
	}
	// debug use
	file, _ := os.Create(path.Join(runPath, "log.txt"))
	file.WriteString(strings.Join(os.Args, " ") + "\n")
	defer file.Close()
	// debug use

	fmt.Printf("[PicGo INFO]: Obcsapi-ImagesHost-PicGo-Core %s by kkbt\n[PicGo INFO]: %s \n[PicGo INFO]: Your Command is: %s\n[PicGo INFO]: Uploading...\n", version, os.Args[0], strings.Join(os.Args, " "))

	switch os.Args[1] {
	case "upload", "u":
		if len(os.Args) == 2 {
			GetFromWindowsClipboard()
			return
		}
		// 读取 upload 后面的参数
		fmt.Println("[PicGo SUCCESS]:")
		for i := 2; i < len(os.Args); i++ {
			if strings.HasPrefix(os.Args[i], "http") {
				res, err := http.Get(os.Args[i])
				if err != nil {
					log.Println("Read File Err,Skip This One", err)
				} else {
					body, err := io.ReadAll(res.Body)
					res.Body.Close()
					if res.StatusCode > 299 {
						log.Fatalf("Download File Err,Skip This One .Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
					}
					if err != nil {
						log.Println("Download File Err,Skip This One", err)
					}
					// fileName := TimeFmt("20060102_") + RandomString(6) + ".jpg"
					fileName := filepath.Base(os.Args[i])
					fmt.Println(PostObcsapiImageHost(fileName, body))
				}
			} else {
				fileName := path.Base(filepath.ToSlash(os.Args[i])) // windows path 风格转 unix 然后获取文件名
				file, err := os.ReadFile(os.Args[i])
				if err != nil {
					log.Println("Read File Err,Skip This One", err)
				} else {
					fmt.Println(PostObcsapiImageHost(fileName, file))
				}
			}
		}
	default:
		fmt.Printf("未考虑 %s", os.Args[1])
	}
}

func GetFromWindowsClipboard() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	img := clipboard.Read(clipboard.FmtImage)
	if len(img) == 0 {
		fmt.Println("[PicGo ERROR]: Fail Read from clipboard")
	} else {
		fmt.Printf("[PicGo SUCCESS]:\n%s", PostObcsapiImageHost(TimeFmt("20060102_")+RandomString(6)+".png", img))
	}
}

func PostObcsapiImageHost(fileName string, file []byte) string {
	// 随机名字

	body := &bytes.Buffer{}
	//创建一个multipart类型的写文件
	writer := multipart.NewWriter(body)
	//使用给出的属性名paramName和文件名filePath创建一个新的form-data头 // "file", RandomString(6)+".png"
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		fmt.Println(" post err=", err)
	}
	_, err = part.Write(file)
	if err != nil {
		fmt.Println("Write err=", err)
	}
	err = writer.Close()
	if err != nil {
		fmt.Println("post err=", err)
	}

	req, err := http.NewRequest( // stored to context
		"POST",
		config[0],
		body,
	)
	if err != nil {
		fmt.Println("err=", err)
	}
	req.Header.Set("Token", config[1])
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err=", err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var hostJson ObcsapiImagesHostJson
	err = decoder.Decode(&hostJson)
	if err != nil {
		fmt.Println("err=", err)
	}
	if config[2] == "url2" {
		return hostJson.Data.Url2
	}
	return hostJson.Data.Url
}

type ObcsapiImagesHostJson struct {
	Data ObcsapiImagesHostJsonUrls `json:"data"`
}
type ObcsapiImagesHostJsonUrls struct {
	Url  string `json:"url"`
	Url2 string `json:"url2"`
}

// Format("2006-01-02 15:04:05")
func TimeFmt(fmt string) string {
	var cstZone = time.FixedZone("CST", 8*3600)
	return time.Now().In(cstZone).Format(fmt)
}

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}
