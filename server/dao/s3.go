package dao

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"obcsapi-go/tools"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetS3Client() (*s3.Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: tools.ConfigGetString("end_point"),
		}, nil
	})
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(tools.ConfigGetString("region")),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(tools.ConfigGetString("access_key"), tools.ConfigGetString("secret_key"), "")),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}

// get text used
func S3GetFileText(s3Client *s3.Client, fileKey string) (string, error) {
	tem, err := S3GetObject(s3Client, fileKey)
	if tem == nil {
		return "No such file: " + fileKey, nil
	}
	return string(tem), err
}

// get_object
func S3GetObject(s3Client *s3.Client, objectKey string) ([]byte, error) {
	buffer := manager.NewWriteAtBuffer([]byte{})
	downloader := manager.NewDownloader(s3Client)
	numBytes, err := downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
		Bucket: aws.String(tools.ConfigGetString("bucket")),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, err
	}
	if numBytes < 1 {
		return nil, errors.New("zero bytes written to memory")
	}
	return buffer.Bytes(), nil
}

// 直接上传存储,覆盖
func S3StoreObject(s3Client *s3.Client, fileKey string, fileBytes []byte) error {
	uploader := manager.NewUploader(s3Client)
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(tools.ConfigGetString("bucket")),
		Key:    aws.String(fileKey),
		Body:   bytes.NewReader(fileBytes),
	})
	return err
}

func S3AppendText(s3Client *s3.Client, fileKey string, text string) error {
	try_get_file, err := S3GetObject(s3Client, fileKey)
	if try_get_file == nil && err != nil {
		err = S3StoreObject(s3Client, fileKey, []byte(text))
	} else {
		err = S3StoreObject(s3Client, fileKey, []byte(string(try_get_file)+text))
	}
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func S3GetTodayDaily(s3Client *s3.Client) string {
	tem, err := S3GetFileText(s3Client, GetDailyFileKey())
	if err != nil {
		log.Println(err)
		return "Have Error!"
	}
	return tem
}
func S3AppendDailyText(s3Client *s3.Client, text string) error {
	return S3AppendText(s3Client, GetDailyFileKey(), text)
}
func S3GetTodayDaily2(s3Client *s3.Client) string {
	tem, err := S3GetFileText(s3Client, GetDailyFileKey())
	if err != nil {
		log.Println(err)
		return "Have Error!"
	}
	return tem
}
func S3Get3DaysList(s3Client *s3.Client) [3]string {
	var ans [3]string
	for i := 0; i < 3; i++ { // 0 1 2 -> -2 -1 0
		day, err := S3GetFileText(s3Client, tools.NowRunConfig.DailyFileKeyMore(i-2))
		if err != nil {
			log.Println(err)
		}
		ans[i] = day
	}
	return ans
}

func S3GetMoreDaliyMdText(s3Client *s3.Client, addDateDay int) (string, error) {
	day, err := S3GetFileText(s3Client, tools.NowRunConfig.DailyFileKeyMore(addDateDay))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return day, nil
}

// 获取文件预先签名 5 min 有效期 即使 file 不存在也会返回 URL
func S3GetPreSignURL(s3Client *s3.Client, fileKey string) (string, error) {
	client := s3.NewPresignClient(s3Client)
	req, err := client.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(tools.ConfigGetString("bucket")),
		Key:    aws.String(fileKey),
	}, s3.WithPresignExpires(time.Minute*5)) // 生成 5 分钟的预签名 URL
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

// md text img url to preSigned url ![](a.jpg) -> ![](a.jpg&signed)
func S3ReplaceMdUrl2PreSignedUrl(in_md []byte) []byte {
	client, err := GetS3Client()
	if err != nil {
		log.Println(err)
	}
	pattern := regexp.MustCompile(`!\[(.*?)\]\(([^http:].*)\)`)
	//pattern := regexp.MustCompile(`!\[(.*?)\]\(\s*([^)"'\s]+)\s*\)`)
	replaceFunc := func(match []byte) []byte {
		// 获取匹配到的链接 并转为 预签名 url
		description := pattern.ReplaceAllString(string(match), "$1")
		link := pattern.ReplaceAllString(string(match), "$2")
		tem, err := url.QueryUnescape(link)
		if err != nil {
			log.Println(err)
		} else {
			tools.Debug(link, "->", tem)
			link = tem
		}
		link2 := link
		// 若请求 以 .md 结尾，则拒绝，避免文本泄露
		if strings.HasSuffix(link, ".md") {
			link2 = link
		} else {
			link2, err = S3GetPreSignURL(client, link)
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

func S3ListObject(s3Client *s3.Client, prefix string) ([]string, error) {
	var result []string
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(tools.ConfigGetString("bucket")),
		Prefix: aws.String(prefix),
	}
	resultList, err := s3Client.ListObjectsV2(context.Background(), input)
	if err != nil {
		log.Fatal(err)
	}

	for _, object := range resultList.Contents {
		result = append(result, *object.Key)
	}
	return result, nil
}
