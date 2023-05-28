package dao

import (
	"bufio"
	"io/ioutil"
	"log"
	"obcsapi-go/tools"
	"os"
	"path/filepath"
)

// TODO: No such file 早晚得改 暂时将就着用吧
func LocalStorageGetTextObject(webDavPath string, text_file_key string) (string, error) {
	_, err := os.Stat(webDavPath + text_file_key)
	if err != nil && os.IsNotExist(err) {
		return "No such file: " + text_file_key, nil
	} else if err != nil {
		return "", err
	}
	buf, err := ioutil.ReadFile(webDavPath + text_file_key)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func LocalStorageGetObject(webDavPath string, text_file_key string) ([]byte, error) {
	_, err := os.Stat(webDavPath + text_file_key)
	if err != nil && os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadFile(webDavPath + text_file_key)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func LocalStorageCheckObject(webDavPath string, file_key string) (bool, error) {
	_, err := os.Stat(webDavPath + file_key)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err

}

func LocalStorageObjectStore(webDavPath string, file_key string, file_bytes []byte) error {
	err := createFile(filepath.Dir(webDavPath + file_key)) // 递归创建文件夹
	if err != nil {
		return err
	}
	file, err := os.Create(webDavPath + file_key) // 清空
	if err != nil {
		log.Println("Creat Error")
		return err
	}
	defer file.Close() // 关闭文件
	_, err = file.Write(file_bytes)
	if err != nil {
		return err
	}
	return nil
}

func LocalStorageTextAppend(webDavPath string, file_key string, text string) error {
	exist, err := LocalStorageCheckObject(webDavPath, file_key)
	if err != nil {
		return err
	}
	if exist {
		file, err := os.OpenFile(webDavPath+file_key, os.O_WRONLY|os.O_APPEND, 0666) // Append
		if err != nil {
			return err
		}
		defer file.Close() // 关闭文件
		write := bufio.NewWriter(file)
		_, err = write.WriteString(text)
		if err != nil {
			return err
		}
		err = write.Flush()
		if err != nil {
			return err
		}
		return nil
	} else {
		return LocalStorageObjectStore(webDavPath, file_key, []byte(text))
	}
}

func LocalStorageDailyTextAppend(webDavPath string, text string) error {
	return LocalStorageTextAppend(webDavPath, GetDailyFileKey(), text)
}

func LocalStorageGetTodayDaily(webDavPath string) string {
	ans, _ := LocalStorageGetTextObject(webDavPath, GetDailyFileKey())
	return ans
}

func LocalStorageGet3DaysList(webDavPath string) [3]string {
	var ans [3]string
	for i := 0; i < 3; i++ { // 0 1 2 -> -2 -1 0
		day, err := LocalStorageGetTextObject(webDavPath, tools.NowRunConfig.DailyFileKeyMore(i-2))
		if err != nil {
			log.Println(err)
		}
		ans[i] = day
	}
	return ans
}

func LocalStorageGetMoreDaliyMdText(webDavPath string, addDateDay int) (string, error) {
	day, err := LocalStorageGetTextObject(webDavPath, tools.NowRunConfig.DailyFileKeyMore(addDateDay))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return day, nil
}

// ------ Tools ------

// 判断是否存在 递归创建文件夹
func createFile(filePath string) error {
	exist, err := LocalStorageCheckObject(filePath, "")
	if err != nil {
		log.Println(err)
	}
	if !exist {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}
