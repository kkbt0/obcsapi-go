package dao

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"obcsapi-go/tools"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	_ "github.com/go-kivik/couchdb/v3"
	"github.com/go-kivik/kivik/v3"
)

type DataSource int

var dataSource DataSource
var sess *session.Session
var couchDb *kivik.DB
var webDavPath string

const (
	Unknown DataSource = iota
	S3
	CouchDb
	LocalStorage
)

func init() {
	tools.CheckFiles() // 监测文件是否存在，不存在则创建文件
	tools.ReloadRunConfig()
	// 初始化 dataSource
	switch tools.ConfigGetInt("data_source") {
	case 1:
		log.Println("Data Source: S3")
		dataSource = S3
		var err error
		sess, err = NewS3Session()
		if err != nil {
			log.Panicln("S3 Session Create Error,Please Check your Config file")
		}
	case 2:
		log.Println("Data Source: CouchDB")
		dataSource = CouchDb
		var err error
		client, err := kivik.New("couch", tools.ConfigGetString("couchdb_url"))
		if err != nil {
			log.Panicln("CouchDB Client Create Error,Please Check your Config file")
		}
		couchDb = client.DB(context.TODO(), tools.ConfigGetString("couchdb_db_name"))
	case 3:
		log.Println("Data Source: Local (Use Webdav)")
		dataSource = LocalStorage
		webDavPath = "./webdav/" + tools.NowRunConfig.Webdav.ObLocalDir
	default:
		log.Panicln("Data Source: Unknown")
		dataSource = Unknown
	}
}

// 注意不要循环 import 使用方法

// 用于获取日记目录
func GetDailyFileKey() string {
	return tools.NowRunConfig.DailyDir() + tools.TimeFmt("2006-01-02") + ".md"
}

func GetMoreDailyFileKey(addDateDay int) string {
	return tools.NowRunConfig.DailyDir() + time.Now().AddDate(0, 0, addDateDay).In(time.FixedZone("CST", 8*3600)).Format("2006-01-02") + ".md"
}

// 获取指定位置文件 并读取为 Str
func GetTextObject(text_file_key string) (string, error) {
	switch dataSource {
	case S3:
		return S3GetTextObject(sess, text_file_key)
	case CouchDb:
		return CouchDbGetTextObject(couchDb, text_file_key)
	case LocalStorage:
		return LocalStorageGetTextObject(webDavPath, text_file_key)
	}
	return "", fmt.Errorf("err GetTextObject Data Source: Unknown")
}

func CheckObject(file_key string) (bool, error) {
	switch dataSource {
	case S3:
		file, _ := S3GetObject(sess, file_key)
		if file != nil {
			return true, nil
		} else {
			return false, nil
		}
	case CouchDb:
		exist, _ := CouchDbCheckObject(couchDb, file_key)
		return exist, nil
	case LocalStorage:
		return LocalStorageCheckObject(webDavPath, file_key)
	}
	return false, fmt.Errorf("err CheckObject Data Source: Unknown")
}

// 覆盖指定位置文件 二进制使用 不可和 TextAppend 混用
func ObjectStore(file_key string, file_bytes []byte) error {
	switch dataSource {
	case S3:
		return S3ObjectStore(sess, file_key, file_bytes)
	case CouchDb:
		return CouchDbFileStorage(couchDb, file_key, file_bytes)
	case LocalStorage:
		return LocalStorageObjectStore(webDavPath, file_key, file_bytes)
	}
	return fmt.Errorf("err Data Source: Unknown")
}

// 覆盖指定位置文件 纯文本使用
func MdTextStore(file_key string, text string) error {
	switch dataSource {
	case S3:
		return S3ObjectStore(sess, file_key, []byte(text))
	case CouchDb:
		return CouchDbMdFiletorage(couchDb, file_key, text)
	case LocalStorage:
		return LocalStorageObjectStore(webDavPath, file_key, []byte(text))
	}
	return fmt.Errorf("err Data Source: Unknown")
}

// 指定位置文件 文本增加内容
func TextAppend(file_key string, text string) error {
	switch dataSource {
	case S3:
		return S3TextAppend(sess, file_key, text)
	case CouchDb:
		return CouchDbTextAppend(couchDb, file_key, text)
	case LocalStorage:
		return LocalStorageTextAppend(webDavPath, file_key, text)
	}
	return fmt.Errorf("err TextAppend Data Source: Unknown")
}

// 为今日日记 增加文本
func DailyTextAppend(text string) error {
	switch dataSource {
	case S3:
		return S3DailyTextAppend(sess, text)
	case CouchDb:
		return CouchDbDailyTextAppend(couchDb, text)
	case LocalStorage:
		return LocalStorageDailyTextAppend(webDavPath, text)
	}
	return fmt.Errorf("err DailyTextAppend Data Source: Unknown")
}

// 为今日日记 增加一行 Memos 格式内容
func DailyTextAppendMemos(text string) error {
	// zk 判断
	if len(text) > 30 && strings.HasPrefix(text, "zk ") {
		text = text[3:]
		fileKey := tools.NowRunConfig.DailyAttachmentDir() + tools.TimeFmt("20060102150405.md")
		err := MdTextStore(fileKey, text)
		if err != nil {
			return err
		}
		return DailyTextAppendMemos(fmt.Sprintf("![[%s]]", fileKey))
	}

	var todo = "todo"
	if strings.Contains(text, todo) {
		text = fmt.Sprintf("\n- [ ] %s %s", tools.TimeFmt("15:04"), strings.Replace(text, "todo", "", 1))
	} else {
		text = fmt.Sprintf("\n- %s %s", tools.TimeFmt("15:04"), text)
	}
	switch dataSource {
	case S3:
		return S3DailyTextAppend(sess, text)
	case CouchDb:
		return CouchDbDailyTextAppend(couchDb, text)
	case LocalStorage:
		return LocalStorageDailyTextAppend(webDavPath, text)
	}
	return fmt.Errorf("err DailyTextAppendMemos Data Source: Unknown")
}

// 获取今日日记
func GetTodayDaily() string {
	switch dataSource {
	case S3:
		return S3GetTodayDaily(sess)
	case CouchDb:
		return CouchDbGetTodayDaily(couchDb)
	case LocalStorage:
		return LocalStorageGetTodayDaily(webDavPath)
	}
	return "没有预料的情况，可能是数据源出现了问题"
}

// 获取今日日记列表，只有一个元素
func GetTodayDailyList() []Daily {
	switch dataSource {
	case S3:
		return S3GetTodayDailyList(sess)
	case CouchDb:
		return CouchDbGetTodayDailyList(couchDb)
	case LocalStorage:
		return LocalStorageGetTodayDailyList(webDavPath)
	}
	var ans []Daily
	return ans
}

// 获取今日日记列表，只有一个元素，对 url 进行了相应处理，可在前端显示
func Get3DaysDailyList() [3]Daily {
	switch dataSource {
	case S3:
		return S3Get3DaysDailyList(sess)
	case CouchDb:
		return CouchDbGet3DaysDailyList(couchDb)
	case LocalStorage:
		return LocalStorageGet3DaysDailyList(webDavPath)
	}
	var ans [3]Daily
	return ans
}

func Get3DaysList() [3]string {
	switch dataSource {
	case S3:
		return S3Get3DaysList(sess)
	case CouchDb:
		return CouchDbGet3DaysList(couchDb)
	case LocalStorage:
		return LocalStorageGet3DaysList(webDavPath)
	}
	var ans [3]string
	return ans
}

// 0 today 1 tomorrow -1 yesterday
func GetMoreDaliyMdText(addDateDay int) (string, error) {
	switch dataSource {
	case S3:
		return S3GetMoreDaliyMdText(sess, addDateDay)
	case CouchDb:
		return CouchDbGetMoreDaliyMdText(couchDb, addDateDay)
	case LocalStorage:
		return LocalStorageGetMoreDaliyMdText(webDavPath, addDateDay)
	}
	return "", errors.New("没有预料的情况，可能是数据源出现了问题")
}

func MdShowText(text string) string {
	text = MdShowTextDailyZk(text)
	switch dataSource {
	case S3:
		return string(S3ReplaceMdUrl2PreSignedUrl([]byte(text)))
	case CouchDb:
		return text
	case LocalStorage:
		return text
	}
	return text
}

// ------Tools--------

// 下载图片到 tem.jpg
func PicDownloader(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	reader := bufio.NewReaderSize(res.Body, 32*1024)
	file, err := os.Create("tem.jpg")
	writer := bufio.NewWriter(file)
	io.Copy(writer, reader)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf, err := os.ReadFile("tem.jpg")
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// 读取 ![[支持类文件/xxx.md]]
func MdShowTextDailyZk(text string) string {
	regexp.MustCompile(`!\[(.*?)\]\(([^http:].*)\)`)
	pattern := regexp.MustCompile(`!\[\[(.*?)\]\]`)
	replaceFunc := func(match []byte) []byte {
		// 获取匹配到的链接 并转为 预签名 url
		ans := text
		fileKey := pattern.ReplaceAllString(string(match), "$1")
		if strings.HasPrefix(fileKey, tools.NowRunConfig.OtherDataDir()) && strings.HasSuffix(fileKey, ".md") {
			ans, _ = GetTextObject(fileKey)
		}
		return []byte(fmt.Sprintf("zk %s", ans))
	}
	return string(pattern.ReplaceAllFunc([]byte(text), replaceFunc))

}
