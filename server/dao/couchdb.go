package dao

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"obcsapi-go/tools"
	"strings"
	"time"

	"github.com/go-kivik/kivik/v3"
)

// 获取节点数据
func CouchDbGetLeftData(db *kivik.DB, left_id string) (string, error) {
	var leftNode LeftDoc
	err := db.Get(context.TODO(), left_id).ScanDoc(&leftNode)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return leftNode.Data, nil
}

// 根据文件 key 获取所有 子节点，然后拼接返回 不可对图片类使用
func CouchDbGetTextObject(db *kivik.DB, text_file_key string) (string, error) {
	var fileNodeDoc FileDoc
	err := db.Get(context.TODO(), text_file_key).ScanDoc(&fileNodeDoc)
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			return "No such file: " + text_file_key, nil
		}
		return "", err
	}
	if fileNodeDoc.Deleted {
		return "No such file: " + text_file_key, nil
	}
	var ans []string
	for _, v := range fileNodeDoc.Children {
		tem, _ := CouchDbGetLeftData(db, v)
		ans = append(ans, tem)
	}
	return strings.Join(ans, ""), nil
}

// 判断文件是否存在 第一个是存在（在 Ob 内），第二个是否有删除标记（Ob内不在，但 Db 可能有删除标记）
func CouchDbCheckObject(db *kivik.DB, text_file_key string) (bool, bool) {
	var fileNodeDoc FileDoc
	err := db.Get(context.TODO(), text_file_key).ScanDoc(&fileNodeDoc)
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			return false, false // 数据库没这条记录
		}
		log.Println(err)
		return false, false // 错误
	}
	if fileNodeDoc.Deleted {
		return false, true // 数据库有记录，但有删除标记
	}
	return true, false // 数据库有记录，正常情况
}

func CouchDbFileStorage(db *kivik.DB, file_key string, file_bytes []byte) error {
	now := time.Now().UnixMilli()
	leftId := fmt.Sprintf("h:%d", now)
	fmt.Println(leftId)
	leftData := base64.StdEncoding.EncodeToString(file_bytes)
	leftDoc := LeftDoc{
		ID:   leftId,
		Data: leftData,
		Type: "leaf",
	}
	_, err := db.Put(context.TODO(), leftId, leftDoc)
	if err != nil {
		return err
	}
	exist, delSign := CouchDbCheckObject(db, file_key)
	if !exist && !delSign { // 数据库没这条记录
		fmt.Println("数据库没这条记录")
		fileDoc := FileDoc{
			ID:       file_key,
			Children: []string{leftId},
			Ctime:    now,
			Mtime:    now,
			Size:     strings.Count(leftData, ""),
			Type:     "newnote",
			Deleted:  false,
		}
		_, err = db.Put(context.TODO(), file_key, fileDoc)
		if err != nil {
			return err
		}
	} else if (!exist && delSign) || (exist && !delSign) { // 数据库有记录
		var fileNodeDoc FileDoc
		db.Get(context.TODO(), file_key).ScanDoc(&fileNodeDoc)
		fmt.Println("数据库有记录")
		fileDoc := FileDoc{
			Rev:      fileNodeDoc.Rev,
			ID:       file_key,
			Children: []string{leftId},
			Ctime:    now,
			Mtime:    now,
			Size:     strings.Count(leftData, ""),
			Type:     "newnote",
			Deleted:  false,
		}
		_, err := db.Put(context.TODO(), file_key, fileDoc)
		if err != nil {
			return err
		}
	}
	return nil
}

func CouchDbMdFiletorage(db *kivik.DB, file_key string, text string) error {
	now := time.Now().UnixMilli()
	leftId := fmt.Sprintf("h:%d", now)
	fmt.Println(leftId)
	leftDoc := LeftDoc{
		ID:   leftId,
		Data: text,
		Type: "leaf",
	}
	_, err := db.Put(context.TODO(), leftId, leftDoc)
	if err != nil {
		return err
	}
	exist, delSign := CouchDbCheckObject(db, file_key)
	if !exist && !delSign { // 数据库没这条记录
		fmt.Println("数据库没这条记录")
		fileDoc := FileDoc{
			ID:       file_key,
			Children: []string{leftId},
			Ctime:    now,
			Mtime:    now,
			Size:     strings.Count(text, ""),
			Type:     "plain",
			Deleted:  false,
		}
		_, err = db.Put(context.TODO(), file_key, fileDoc)
		if err != nil {
			return err
		}
	} else if (!exist && delSign) || (exist && !delSign) { // 数据库有记录
		var fileNodeDoc FileDoc
		db.Get(context.TODO(), file_key).ScanDoc(&fileNodeDoc)
		fmt.Println("数据库有记录")
		fileDoc := FileDoc{
			Rev:      fileNodeDoc.Rev,
			ID:       file_key,
			Children: []string{leftId},
			Ctime:    now,
			Mtime:    now,
			Size:     strings.Count(text, ""),
			Type:     "plain",
			Deleted:  false,
		}
		_, err := db.Put(context.TODO(), file_key, fileDoc)
		if err != nil {
			return err
		}
	}
	return nil
}

func CouchDbTextAppend(db *kivik.DB, file_key string, text string) error {
	now := time.Now().UnixMilli()
	leftId := fmt.Sprintf("h:%d", now)
	fmt.Println(leftId)
	leftDoc := LeftDoc{
		ID:   leftId,
		Data: text,
		Type: "leaf",
	}
	_, err := db.Put(context.TODO(), leftId, leftDoc)
	if err != nil {
		return err
	}
	exist, delSign := CouchDbCheckObject(db, file_key)
	if !exist && !delSign { // 文件不存在，新建文件
		fmt.Println("文件不存在，新建文件")
		fileDoc := FileDoc{
			ID:       file_key,
			Children: []string{leftId},
			Ctime:    now,
			Mtime:    now,
			Size:     strings.Count(text, ""),
			Type:     "plain",
			Deleted:  false,
		}
		_, err = db.Put(context.TODO(), file_key, fileDoc)
		if err != nil {
			return err
		}
		return nil
	} else if !exist && delSign { // 数据库存在，但是有删除标记
		var fileNodeDoc FileDoc
		db.Get(context.TODO(), file_key).ScanDoc(&fileNodeDoc)
		fmt.Println("数据库存在，但是有删除标记")
		fileDoc := FileDoc{
			Rev:      fileNodeDoc.Rev,
			ID:       file_key,
			Children: []string{leftId},
			Ctime:    now,
			Mtime:    now,
			Size:     strings.Count(text, ""),
			Type:     "plain",
			Deleted:  false,
		}
		_, err := db.Put(context.TODO(), file_key, fileDoc)
		if err != nil {
			return err
		}
		return nil
	} else if exist && !delSign { // 数据库存在，正常情况
		var fileNodeDoc FileDoc
		db.Get(context.TODO(), file_key).ScanDoc(&fileNodeDoc)
		fmt.Println("数据库存在，正常情况")
		fileDoc := FileDoc{
			Rev:      fileNodeDoc.Rev,
			ID:       file_key,
			Children: append(fileNodeDoc.Children, leftId),
			Ctime:    now,
			Mtime:    now,
			Size:     strings.Count(text, ""),
			Type:     "plain",
			Deleted:  false,
		}
		_, err := db.Put(context.TODO(), file_key, fileDoc)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("没有预料的的情况")
}

func CouchDbDailyTextAppend(db *kivik.DB, text string) error {
	return CouchDbTextAppend(db, GetDailyFileKey(), text)
}

func CouchDbGetTodayDaily(db *kivik.DB) string {
	ans, _ := CouchDbGetTextObject(couchDb, GetDailyFileKey())
	return ans
}

func CouchDbGetTodayDailyList(db *kivik.DB) []Daily {
	day := GetTodayDaily()
	return []Daily{{
		Date:       tools.TimeFmt("2006-01-02"),
		ServerTime: tools.TimeFmt("200601021504"),
		Data:       day,
		MdShowData: day, // TODO 图像 用 Base64 代替
	}}
}

func CouchDbGet3DaysDailyList(db *kivik.DB) [3]Daily {
	var ans [3]Daily
	for i := 0; i < 3; i++ { // 0 1 2 -> -2 -1 0
		date := time.Now().AddDate(0, 0, i-2).In(time.FixedZone("CST", 8*3600)).Format("2006-01-02")
		day, err := CouchDbGetTextObject(db, fmt.Sprintf("%s%s.md", tools.ConfigGetString("ob_daily_dir"), date))
		if err != nil {
			log.Println(err)
		}
		ans[i] = Daily{
			Data:       day,
			MdShowData: day, // TODO 图像 用 Base64 代替
			Date:       date,
			ServerTime: tools.TimeFmt("200601021504"),
		}
	}
	return ans
}

func CouchDbGet3DaysList(db *kivik.DB) [3]string {
	var ans [3]string
	for i := 0; i < 3; i++ { // 0 1 2 -> -2 -1 0
		date := time.Now().AddDate(0, 0, i-2).In(time.FixedZone("CST", 8*3600)).Format("2006-01-02")
		day, err := CouchDbGetTextObject(db, fmt.Sprintf("%s%s.md", tools.ConfigGetString("ob_daily_dir"), date))
		if err != nil {
			log.Println(err)
		}
		ans[i] = day
	}
	return ans
}

func CouchDbGetMoreDaliyMdText(db *kivik.DB, addDateDay int) (string, error) {
	date := time.Now().AddDate(0, 0, addDateDay).In(time.FixedZone("CST", 8*3600)).Format("2006-01-02")
	day, err := CouchDbGetTextObject(db, fmt.Sprintf("%s%s.md", tools.ConfigGetString("ob_daily_dir"), date))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return day, nil
}
