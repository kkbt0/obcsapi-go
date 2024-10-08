package dao

import (
	"bufio"
	"log"
	"obcsapi-go/tools"
	"os"
	"path/filepath"
	"strings"
)

// TODO: No such file 早晚得改 暂时将就着用吧
func LocalStorageGetFileText(webDavPath string, text_file_key string) (string, error) {
	_, err := os.Stat(webDavPath + text_file_key)
	if err != nil && os.IsNotExist(err) {
		return "No such file: " + text_file_key, nil
	} else if err != nil {
		return "", err
	}
	buf, err := os.ReadFile(webDavPath + text_file_key)
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
	buf, err := os.ReadFile(webDavPath + text_file_key)
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

func LocalStorageStoreObject(webDavPath string, file_key string, file_bytes []byte) error {
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

func LocalStorageAppendText(webDavPath string, file_key string, text string) error {
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
		return LocalStorageStoreObject(webDavPath, file_key, []byte(text))
	}
}

func LocalStorageAppendDailyText(webDavPath string, text string) error {
	return LocalStorageAppendText(webDavPath, GetDailyFileKey(), text)
}

func LocalStorageGetTodayDaily(webDavPath string) string {
	ans, _ := LocalStorageGetFileText(webDavPath, GetDailyFileKey())
	return ans
}

func LocalStorageGet3DaysList(webDavPath string) [3]string {
	var ans [3]string
	for i := 0; i < 3; i++ { // 0 1 2 -> -2 -1 0
		day, err := LocalStorageGetFileText(webDavPath, tools.NowRunConfig.DailyFileKeyMore(i-2))
		if err != nil {
			log.Println(err)
		}
		ans[i] = day
	}
	return ans
}

func LocalStorageGetMoreDaliyMdText(webDavPath string, addDateDay int) (string, error) {
	day, err := LocalStorageGetFileText(webDavPath, tools.NowRunConfig.DailyFileKeyMore(addDateDay))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return day, nil
}

func LocalStorageListObject(webDavPath string, prefix string) ([]string, error) {
	var result []string
	err := filepath.Walk(webDavPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(path, filepath.Clean(webDavPath)+"/"+prefix) {
			result = append(result, strings.Replace(path, filepath.Clean(webDavPath)+"/", "", 1))
		}
		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
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
