package dao

import (
	"fmt"
	"log"
	"obcsapi-go/tools"
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

func NewS3Session() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(tools.ConfigGetString("access_key"), tools.ConfigGetString("secret_key"), ""),
		Endpoint:    aws.String(tools.ConfigGetString("end_point")),
		Region:      aws.String(tools.ConfigGetString("region")),
	})
	return sess, err
}

// get text used
func S3GetTextObject(sess *session.Session, text_file_key string) (string, error) {
	tem, err := S3GetObject(sess, text_file_key)
	if tem == nil {
		return "No such file: " + text_file_key, nil
	}
	return string(tem), err
}

// get_object
func S3GetObject(sess *session.Session, file_key string) ([]byte, error) {
	file, err := os.Create("tem.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(
		file,
		&s3.GetObjectInput{
			Bucket: aws.String(tools.ConfigGetString("bucket")),
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
	return buf, nil
}

// 直接上传存储,覆盖
func S3ObjectStore(sess *session.Session, file_key string, file_bytes []byte) error {
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
		Bucket: aws.String(tools.ConfigGetString("bucket")),
		Key:    aws.String(file_key),
		Body:   fp,
	})
	if err != nil {
		return err
	}
	return nil
}
func S3TextAppend(sess *session.Session, file_key string, text string) error {
	try_get_file, err := S3GetObject(sess, file_key)
	if try_get_file == nil && err != nil {
		err = S3ObjectStore(sess, file_key, []byte(text))
	} else {
		err = S3ObjectStore(sess, file_key, []byte(string(try_get_file)+text))
	}
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func S3DailyTextAppend(sess *session.Session, text string) error {
	return S3TextAppend(sess, GetDailyFileKey(), text)
}

func S3GetTodayDaily(sess *session.Session) string {
	tem, err := S3GetTextObject(sess, GetDailyFileKey())
	if err != nil {
		log.Println(err)
		return "Have Error!"
	}
	return tem
}

func S3Get3DaysList(sess *session.Session) [3]string {
	var ans [3]string
	for i := 0; i < 3; i++ { // 0 1 2 -> -2 -1 0
		day, err := S3GetTextObject(sess, tools.NowRunConfig.DailyFileKeyMore(i-2))
		if err != nil {
			log.Println(err)
		}
		ans[i] = day
	}
	return ans
}

func S3GetMoreDaliyMdText(sess *session.Session, addDateDay int) (string, error) {
	day, err := S3GetTextObject(sess, tools.NowRunConfig.DailyFileKeyMore(addDateDay))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return day, nil
}

// 获取文件预先签名 5 min 有效期 即使 file 不存在也会返回 URL
func S3GetPreSignURL(sess *session.Session, file_key string) (string, error) {
	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(tools.ConfigGetString("bucket")),
		Key:    aws.String(file_key),
	})
	urlStr, err := req.Presign(5 * time.Minute)
	if err != nil {
		return "", err
	}
	return urlStr, nil
}

// md text img url to preSigned url ![](a.jpg) -> ![](a.jpg&signed)
func S3ReplaceMdUrl2PreSignedUrl(in_md []byte) []byte {
	sess, err := NewS3Session()
	if err != nil {
		log.Println(err)
	}
	pattern := regexp.MustCompile(`!\[(.*?)\]\(([^http:].*)\)`)
	//pattern := regexp.MustCompile(`!\[(.*?)\]\(\s*([^)"'\s]+)\s*\)`)
	replaceFunc := func(match []byte) []byte {
		// 获取匹配到的链接 并转为 预签名 url
		description := pattern.ReplaceAllString(string(match), "$1")
		link := pattern.ReplaceAllString(string(match), "$2")
		link2 := link
		// 若请求 以 .md 结尾，则拒绝，避免文本泄露
		if strings.HasSuffix(link, ".md") {
			link2 = link
		} else {
			link2, err = S3GetPreSignURL(sess, link)
			if err != nil {
				log.Println(err)
				return []byte(fmt.Sprintf("![%s](%s)", description, link))
			}
		}
		// fmt.Println(link2)
		// 替换链接为临时鉴权链接
		return []byte(fmt.Sprintf("![%s](%s)", description, link2))
	}
	return pattern.ReplaceAllFunc(in_md, replaceFunc)
}

func S3ListObject(sess *session.Session, prefix string) ([]string, error) {
	var result []string
	svc := s3.New(sess)

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(tools.ConfigGetString("bucket")),
	}
	resultList, err := svc.ListObjectsV2(input)
	if err != nil {
		fmt.Println("Error", err)
		return nil, err
	}
	for _, object := range resultList.Contents {
		if strings.HasPrefix(*object.Key, prefix) && strings.Replace(*object.Key, prefix, "", 1) != "" {
			result = append(result, *object.Key)
		}
	}
	return result, nil
}
