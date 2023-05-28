package dao

import (
	"fmt"
	"log"
	"obcsapi-go/tools"
	"strings"

	"github.com/studio-b12/gowebdav"
)

func WebDavGetTextObject(webDavClient *gowebdav.Client, prePath string, fileKey string) (string, error) {
	bytes, err := WebDavGetObject(webDavClient, prePath, fileKey)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func WebDavGetObject(webDavClient *gowebdav.Client, prePath string, fileKey string) ([]byte, error) {
	return webDavClient.Read(prePath + fileKey)
}

func WebDavCheckObject(webDavClient *gowebdav.Client, prePath string, fileKey string) (bool, error) {
	_, err := webDavClient.Stat(prePath + fileKey)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func WebDavObjectStorage(webDavClient *gowebdav.Client, prePath string, file_key string, file_bytes []byte) error {
	return webDavClient.Write(prePath+file_key, file_bytes, 0644)
}

func WebDavTextAppend(webDavClient *gowebdav.Client, prePath string, file_key string, text string) error {
	exits, err := WebDavCheckObject(webDavClient, prePath, file_key)
	if err != nil {
		return err
	}
	if !exits {
		return WebDavObjectStorage(webDavClient, prePath, file_key, []byte(text))
	} else {
		oldText := ""
		oldText, err = WebDavGetTextObject(webDavClient, prePath, file_key)
		fmt.Println(oldText)
		if err != nil {
			return err
		}
		return WebDavObjectStorage(webDavClient, prePath, file_key, []byte(oldText+text))
	}
}

func WebDavDailyTextAppend(webDavClient *gowebdav.Client, prePath string, text string) error {
	return WebDavTextAppend(webDavClient, prePath, GetDailyFileKey(), text)
}

func WebDavGetTodayDaily(webDavClient *gowebdav.Client, prePath string) string {
	ans, _ := WebDavGetTextObject(webDavClient, prePath, GetDailyFileKey())
	return ans
}

func WebDavGet3DaysList(webDavClient *gowebdav.Client, prePath string) [3]string {
	var ans [3]string
	for i := 0; i < 3; i++ { // 0 1 2 -> -2 -1 0
		day, err := WebDavGetTextObject(webDavClient, prePath, tools.NowRunConfig.DailyFileKeyMore(i-2))
		if err != nil {
			log.Println(err)
		}
		ans[i] = day
	}
	return ans
}

func WebDavGetMoreDaliyMdText(webDavClient *gowebdav.Client, prePath string, addDateDay int) (string, error) {
	day, err := WebDavGetTextObject(webDavClient, prePath, tools.NowRunConfig.DailyFileKeyMore(addDateDay))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return day, nil
}
