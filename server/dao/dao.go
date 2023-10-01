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

	"github.com/aws/aws-sdk-go/aws/session"
	_ "github.com/go-kivik/couchdb/v3"
	"github.com/go-kivik/kivik/v3"
	"github.com/spf13/viper"
	"github.com/studio-b12/gowebdav"
)

type DataSource int

var dataSource DataSource
var sess *session.Session
var couchDb *kivik.DB
var webDavDirPath string
var webDavClient *gowebdav.Client
var webDavPrePath string

const (
	Unknown DataSource = iota
	S3
	CouchDb
	LocalStorage
	WebDav
)

func init() {
	// tools.CheckFiles() // 监测文件是否存在，不存在则创建文件
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
		// webDavPath is a LocalDir
		webDavDirPath = "./webdav/" + tools.NowRunConfig.Webdav.ObLocalDir
	case 4:
		log.Println("Data Source: WebDav")
		dataSource = WebDav
		webDavClient = gowebdav.NewClient(tools.ConfigGetString("web_dav_url"), tools.ConfigGetString("web_dav_username"), tools.ConfigGetString("web_dav_password"))
		webDavPrePath = tools.ConfigGetString("web_dav_dir")
	default:
		log.Panicln("Data Source: Unknown")
		dataSource = Unknown
	}
}

// 注意不要循环 import 使用方法

// 用于获取日记目录
func GetDailyFileKey() string {
	return tools.NowRunConfig.DailyFileKey()
}

func GetMoreDailyFileKey(addDateDay int) string {
	return tools.NowRunConfig.DailyFileKeyMore(addDateDay)
}

// 获取指定位置文件 并读取为 Str
func GetFileText(text_file_key string) (string, error) {
	switch dataSource {
	case S3:
		return S3GetFileText(sess, text_file_key)
	case CouchDb:
		return CouchDbGetFileText(couchDb, text_file_key)
	case LocalStorage:
		return LocalStorageGetFileText(webDavDirPath, text_file_key)
	case WebDav:
		return WebDavGetFileText(webDavClient, webDavPrePath, text_file_key)
	}
	return "", fmt.Errorf("err GetFileText Data Source: Unknown")
}

// 获取指定位置文件 can nil , nil means no such object
func GetObject(fileKey string) ([]byte, error) {
	switch dataSource {
	case S3:
		return S3GetObject(sess, fileKey)
	case CouchDb:
		return CouchDbGetObject(couchDb, fileKey)
	case LocalStorage:
		return LocalStorageGetObject(webDavDirPath, fileKey)
	case WebDav:
		return WebDavGetObject(webDavClient, webDavPrePath, fileKey)
	}
	return nil, fmt.Errorf("err GetFileText Data Source: Unknown")
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
		return LocalStorageCheckObject(webDavDirPath, file_key)
	case WebDav:
		return WebDavCheckObject(webDavClient, webDavPrePath, file_key)
	}
	return false, fmt.Errorf("err CheckObject Data Source: Unknown")
}

// 覆盖指定位置文件 二进制使用 不可和 AppendText 混用
func StoreObject(file_key string, file_bytes []byte) error {
	switch dataSource {
	case S3:
		return S3StoreObject(sess, file_key, file_bytes)
	case CouchDb:
		return CouchDbFileStorage(couchDb, file_key, file_bytes)
	case LocalStorage:
		return LocalStorageStoreObject(webDavDirPath, file_key, file_bytes)
	case WebDav:
		return WebDavObjectStorage(webDavClient, webDavPrePath, file_key, file_bytes)
	}
	return fmt.Errorf("err Data Source: Unknown")
}

// 覆盖指定位置文件 纯文本使用
func CoverStoreTextFile(file_key string, text string) error {
	switch dataSource {
	case S3:
		return S3StoreObject(sess, file_key, []byte(text))
	case CouchDb:
		return CouchDbMdFiletorage(couchDb, file_key, text)
	case LocalStorage:
		return LocalStorageStoreObject(webDavDirPath, file_key, []byte(text))
	case WebDav:
		return WebDavObjectStorage(webDavClient, webDavPrePath, file_key, []byte(text))
	}
	return fmt.Errorf("err Data Source: Unknown")
}

// 指定位置文件 文本增加内容
func AppendText(file_key string, text string) error {
	switch dataSource {
	case S3:
		return S3AppendText(sess, file_key, text)
	case CouchDb:
		return CouchDbAppendText(couchDb, file_key, text)
	case LocalStorage:
		return LocalStorageAppendText(webDavDirPath, file_key, text)
	case WebDav:
		return WebDavAppendText(webDavClient, webDavPrePath, file_key, text)
	}
	return fmt.Errorf("err AppendText Data Source: Unknown")
}

// 为今日日记 增加文本
func AppendDailyText(text string) error {
	switch dataSource {
	case S3:
		return S3AppendDailyText(sess, text)
	case CouchDb:
		return CouchDbAppendDailyText(couchDb, text)
	case LocalStorage:
		return LocalStorageAppendDailyText(webDavDirPath, text)
	case WebDav:
		return WebDavAppendDailyText(webDavClient, webDavPrePath, text)
	}
	return fmt.Errorf("err AppendDailyText Data Source: Unknown")
}

// 为今日日记 增加一行 Memos 格式内容
func AppendDailyMemos(text string) error {
	// zk 判断
	if len(text) > 30 && strings.HasPrefix(text, "zk ") {
		text = text[3:]
		fileKey := tools.NowRunConfig.DailyAttachmentDir() + tools.TimeFmt("20060102150405.md")
		err := CoverStoreTextFile(fileKey, text)
		if err != nil {
			return err
		}
		return AppendDailyMemos(fmt.Sprintf("![[%s]]", fileKey))
	}

	var todo = "todo"
	if strings.Contains(text, todo) {
		text = fmt.Sprintf("\n- [ ] %s %s", tools.TimeFmt("15:04"), strings.Replace(text, "todo", "", 1))
	} else {
		text = fmt.Sprintf("\n- %s %s", tools.TimeFmt("15:04"), text)
	}
	return AppendDailyText(text)
}

// 获取今日日记 废弃
func GetTodayDaily() string {
	switch dataSource {
	case S3:
		return S3GetTodayDaily(sess)
	case CouchDb:
		return CouchDbGetTodayDaily(couchDb)
	case LocalStorage:
		return LocalStorageGetTodayDaily(webDavDirPath)
	case WebDav:
		return WebDavGetTodayDaily(webDavClient, webDavPrePath)
	}
	return "没有预料的情况，可能是数据源出现了问题"
}

func Get3DaysList() [3]string {
	switch dataSource {
	case S3:
		return S3Get3DaysList(sess)
	case CouchDb:
		return CouchDbGet3DaysList(couchDb)
	case LocalStorage:
		return LocalStorageGet3DaysList(webDavDirPath)
	case WebDav:
		return WebDavGet3DaysList(webDavClient, webDavPrePath)
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
		return LocalStorageGetMoreDaliyMdText(webDavDirPath, addDateDay)
	case WebDav:
		return WebDavGetMoreDaliyMdText(webDavClient, webDavPrePath, addDateDay)
	}
	return "", errors.New("没有预料的情况，可能是数据源出现了问题")
}

func MdShowText(fileKey string, text string) string {
	tools.Debug("From File:", fileKey)
	// 先替换一遍 .md 结尾的
	text = MdShowTextDailyZk(fileKey, text)
	switch dataSource {
	case S3:
		if viper.GetBool("s3_wiki_link_use_presign") {
			return string(S3ReplaceMdUrl2PreSignedUrl([]byte(text)))
		} else {
			return ObFileUrl(fileKey, text)
		}
	case CouchDb:
		return ObFileUrl(fileKey, text)
	case LocalStorage:
		return ObFileUrl(fileKey, text)
	case WebDav:
		return ObFileUrl(fileKey, text) // 其实 Basic Auth 也可以获取，不过 markdown 预览不支持
	}
	return text
}

func ListObject(prefix string) ([]string, error) {
	switch dataSource {
	case S3:
		// all objects filtered -r
		return S3ListObject(sess, prefix)
	case CouchDb:
		// all objects filtered -r
		return CouchDbListObject(couchDb, prefix)
	case LocalStorage:
		// all objects filtered -r
		return LocalStorageListObject(webDavDirPath, prefix)
	case WebDav:
		// only one dir objects
		return WebDavListObject(webDavClient, webDavPrePath, prefix)
	}
	return nil, errors.New("没有预料的情况，可能是数据源出现了问题")
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

// 读取 ![[日志/附件/202305/xxx.md]]
func MdShowTextDailyZk(fromFileKey string, text string) string {
	pattern := regexp.MustCompile(`!\[\[(.*?)\]\]`)
	replaceFunc := func(match []byte) []byte {
		// 获取匹配到的链接 并转为 预签名 url
		ans := "There is something wrong !"
		var err error
		fileKey := pattern.ReplaceAllString(string(match), "$1")
		if strings.HasPrefix(fileKey, "../") {
			fileKey = GetAbsoluteFileKey(fromFileKey, fileKey)
		}
		if strings.HasPrefix(fileKey, tools.NowRunConfig.DailyAttachmentDir()) && strings.HasSuffix(fileKey, ".md") {
			ans, err = GetFileText(fileKey)
		} else if strings.HasSuffix(fileKey, ".md") && tools.ConfigGetString("allow_wiki_link_all") == "true" {
			ans, err = GetFileText(fileKey)
		}
		if err != nil {
			ans = "Found Error: " + err.Error()
		}
		ans = strings.ReplaceAll(ans, "\n-", "\n -") // 避免分割错误
		return []byte(fmt.Sprintf("[zk] %s", ans))
	}
	return string(pattern.ReplaceAllFunc([]byte(text), replaceFunc))

}

func ObFileUrl(fromFileKey string, text string) string {
	pattern := regexp.MustCompile(`!\[(.*?)\]\(([^http:].*)\)`)
	//pattern := regexp.MustCompile(`!\[(.*?)\]\(\s*([^)"'\s]+)\s*\)`)
	replaceFunc := func(match []byte) []byte {
		description := pattern.ReplaceAllString(string(match), "$1")
		link := pattern.ReplaceAllString(string(match), "$2")
		if strings.HasPrefix(link, "../") {
			link = GetAbsoluteFileKey(fromFileKey, link)
		}
		link2 := link
		// 若请求 以 .md 结尾，则拒绝，避免文本泄露
		if !strings.HasSuffix(link, ".md") {
			link2 = fmt.Sprintf("%s/ob/file?accessToken=%s&fileKey=%s", tools.ConfigGetString("backend_url_full"), tools.ObFileAccessToken(), link)
		}
		// fmt.Println(link2)
		return []byte(fmt.Sprintf("![%s](%s)", description, link2))
	}
	return string(pattern.ReplaceAllFunc([]byte(text), replaceFunc))
}
