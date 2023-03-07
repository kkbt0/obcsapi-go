package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func get_client() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(ConfigGetString("access_key"), ConfigGetString("secret_key"), ""),
		Endpoint:    aws.String(ConfigGetString("end_point")),
		Region:      aws.String(ConfigGetString("region")),
	})
	return sess, err
}

// get text used
func get(sess *session.Session, text_file_key string) (string, error) {
	tem, err := get_object(sess, text_file_key)
	if tem == nil {
		return "No such file: " + text_file_key, nil
	}
	return string(tem), err
}

// get_object
func get_object(sess *session.Session, file_key string) ([]byte, error) {
	file, err := os.Create("tem.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(
		file,
		&s3.GetObjectInput{
			Bucket: aws.String(ConfigGetString("bucket")),
			Key:    aws.String(file_key),
		})
	if err != nil {
		return nil, err
	}
	if numBytes == 0 {
		return nil, nil
	}
	buf, err := os.ReadFile("tem.txt")
	if err != nil {
		return nil, err
	}
	// 对图片 url 进行预签名
	return replace_md_url2pre_url(sess, buf), nil
}

// """直接上传存储,可能覆盖"""
func store(sess *session.Session, file_key string, file_bytes []byte) error {
	file, err := os.Create("tem.txt")
	if err != nil {
		return err
	}
	_, err = file.Write(file_bytes)
	if err != nil {
		return err
	}
	defer file.Close()
	fp, err := os.Open("tem.txt")
	if err != nil {
		return err
	}
	defer fp.Close()
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(ConfigGetString("bucket")),
		Key:    aws.String(file_key),
		Body:   fp,
	})
	if err != nil {
		return err
	}
	return nil
}
func append(sess *session.Session, file_key string, text string) error {
	try_get_file, err := get_object(sess, file_key)
	if try_get_file == nil && err != nil {
		err = store(sess, file_key, []byte(text))
	} else {
		err = store(sess, file_key, []byte(string(try_get_file)+text))
	}
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func daily_file_key() string {
	return "日志/" + timeFmt("2006-01-02") + ".md"
}
func append_daily(sess *session.Session, text string) error {
	err := append(sess, daily_file_key(), text)
	if err != nil {
		return err
	}
	return nil
}
func append_memos_in_daily(sess *session.Session, text string) error {
	var todo = "todo"
	if strings.Contains(text, todo) {
		text = fmt.Sprintf("\n- [ ] %s %s", timeFmt("15:04"), text)
	} else {
		text = fmt.Sprintf("\n- %s %s", timeFmt("15:04"), text)
	}
	err := append_daily(sess, text)
	if err != nil {
		return err
	}
	return nil
}
func today_daily(sess *session.Session) string {
	tem, err := get(sess, daily_file_key())
	if err != nil {
		return "Have Error!"
	}
	return tem
}
func get_today_daily_list(sess *session.Session) []string {
	return []string{today_daily(sess)}
}
func get_3_daily_list(sess *session.Session) [3]Daily {
	// fmt.Println(time.Now().In(cstZone).Format("2006-01-02 15:04:05"))
	var cstZone = time.FixedZone("CST", 8*3600)
	var ans [3]Daily
	for i := 0; i < 3; i++ { // 0 1 2 -> -2 -1 0
		date := time.Now().AddDate(0, 0, i-2).In(cstZone).Format("2006-01-02")
		day, err := get(sess, fmt.Sprintf("日志/%s.md", date))
		if err != nil {
			fmt.Println(err)
		}
		ans[i] = Daily{Data: day, Date: date, ServerTime: timeFmt("200601021504")}
	}
	return ans
}

func downloader(url string) ([]byte, error) {
	fmt.Println("Hello World!")
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

// 获取文件预先签名 5 min 有效期 即使 file 不存在也会返回 URL
func getPreSignURL(sess *session.Session, file_key string) (string, error) {
	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(ConfigGetString("bucket")),
		Key:    aws.String(file_key),
	})
	urlStr, err := req.Presign(5 * time.Minute)
	if err != nil {
		return "", err
	}
	return urlStr, nil
}

// md text img url to preSigned url ![](a.jpg) -> ![](a.jpg&signed)
func replace_md_url2pre_url(sess *session.Session, in_md []byte) []byte {
	pattern := regexp.MustCompile(`!\[(.*?)\]\(\s*([^)"'\s]+)\s*\)`)
	replaceFunc := func(match []byte) []byte {
		// 获取匹配到的链接 并转为 预签名 url
		description := pattern.ReplaceAllString(string(match), "$1")
		link := pattern.ReplaceAllString(string(match), "$2")
		link2, err := getPreSignURL(sess, link)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(link2)
		// 替换链接为临时鉴权链接
		return []byte(fmt.Sprintf("![%s](%s)", description, link2))
	}
	return pattern.ReplaceAllFunc(in_md, replaceFunc)
}
