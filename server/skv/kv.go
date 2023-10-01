package skv

import (
	"log"
	"obcsapi-go/dao"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/spf13/viper"
)

var db, _ = bolt.Open(GetKvDbFile(), 0600, &bolt.Options{Timeout: 1 * time.Second})

func InitKv() error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBlocks"))
		if b == nil {
			_, err := tx.CreateBucket([]byte("MyBlocks"))
			if err != nil {
				log.Println(err)
			}
		}
		return nil
	})
	return err
}

func GetKv(key string) ([]byte, error) {
	var result []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBlocks"))
		if b != nil {
			result = b.Get([]byte(key))
		}
		return nil
	})
	return result, err
}

func UpdateKv(key string, val string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBlocks"))
		if b != nil {
			err := b.Put([]byte(key), []byte(val))
			if err != nil {
				log.Println(err)
			}
		}
		return nil
	})
	return err
}

// Only Text
func GetByFileKey(filekey string) string {
	var result string
	err := InitKv()
	if err != nil {
		log.Println(err)
	}

	tem, err := GetKv(filekey)
	if err != nil {
		log.Println(err)
	}
	result = string(tem)

	if result == "" { // nothing
		text, err := dao.GetFileText(filekey)
		if err != nil {
			log.Println(err)
		}
		if text != "No such file: "+filekey {
			err = UpdateKv(filekey, text)
			if err != nil {
				log.Println(err)
			}
			result = text
		}
	}
	return result // "" means nothing
}

func PutByFileKey(filekey string) error {
	err := InitKv()
	if err != nil {
		log.Println(err)
	}
	text, err := dao.GetFileText(filekey)
	if err != nil {
		log.Println(err)
	}
	if text != "No such file: "+filekey {
		err = UpdateKv(filekey, text)
		if err != nil {
			log.Println(err)
		}
	}
	return err
}

func PutFile(filekey string, val string) error {
	err := InitKv()
	if err != nil {
		log.Println(err)
	}
	err = UpdateKv(filekey, val)
	if err != nil {
		log.Println(err)
	}
	return err
}

type KvSerchResult struct {
	Key string `json:"filekey"`
	Val string `json:"content"`
}

func KvSerch(key string) ([]KvSerchResult, error) {
	var ans []KvSerchResult = []KvSerchResult{}
	err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("MyBlocks"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if strings.Contains(string(k), key) || strings.Contains(string(v), key) {
				ans = append(ans, KvSerchResult{Key: string(k), Val: string(v)})
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ans, nil

}

func GetKvDbFile() string {
	name := viper.GetString("kvdb")
	if name == "" {
		name = "mykv.db"
	}
	return name
}
