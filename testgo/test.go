package main

import (
	"log"
	"os"

	"github.com/studio-b12/gowebdav"
)

func main() {
	c := gowebdav.NewClient("http://localhost:8900/webdav", "kkbt", "xxxxx")
	err := c.Write("testdb/xxx.md", []byte("New xxx"), 0644)
	if err != nil {
		log.Println(err)
	}
	err = c.Write("testdb/xxx1.md", []byte("New xxx1"), 0644)
	if err != nil {
		log.Println(err)
	}
	err = c.Write("testdb/xxx/xxx/xxx1.md", []byte("New xxx1"), 0644)
	if err != nil {
		log.Println(err)
	}
}

func ReadTalkLog(path string) (string, error) {
	file, err := os.ReadFile("io.txt")
	if err != nil {
		return "", err
	}
	return string(file), nil
}
