package dao

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"obcsapi-go/tools"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	_ "github.com/go-kivik/couchdb/v3"
	"github.com/go-kivik/kivik/v3"
)

type DataSource int

var dataSource DataSource
var sess *session.Session
var couchDb *kivik.DB

const (
	Unknown DataSource = iota
	S3
	CouchDb
)

func init() {
	tools.CheckFiles() // 监测文件是否存在，不存在则创建文件
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
	default:
		log.Panicln("Data Source: Unknown")
		dataSource = Unknown
	}
}

// 注意不要循环 import 使用方法

// 用于获取日记目录
func GetDailyFileKey() string {
	return tools.ConfigGetString("ob_daily_dir") + tools.TimeFmt("2006-01-02") + ".md"
}

// 获取指定位置文件 并读取为 Str
func GetTextObject(text_file_key string) (string, error) {
	switch dataSource {
	case S3:
		return S3GetTextObject(sess, text_file_key)
	case CouchDb:
		return CouchDbGetTextObject(couchDb, text_file_key)
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
	}
	return fmt.Errorf("err DailyTextAppend Data Source: Unknown")
}

// 为今日日记 增加一行 Memos 格式内容
func DailyTextAppendMemos(text string) error {
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
	}
	var ans [3]string
	return ans
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
