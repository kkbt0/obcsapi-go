package tools

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// 主要用于运行时可修改配置
// 运行时可修改配置

var NowRunConfig RunConfig

type WebDavConfig struct {
	Server     bool   `json:"server"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	ObLocalDir string `json:"ob_local_dir"`
}

type RunConfig struct {
	Webdav WebDavConfig `json:"webdav"`
}

func GetRunConfigHandler(c *gin.Context) {
	ReloadRunConfig()
	c.JSON(200, NowRunConfig)
}

func PostConfigHandler(c *gin.Context) {
	var config RunConfig
	if c.ShouldBindJSON(&config) != nil {
		c.String(400, "参数错误")
		return
	}
	data, err := json.Marshal(&config)
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	err = os.WriteFile("config.run.json", data, 0666)
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	err = ReloadRunConfig()
	if err != nil {
		c.Error(err)
		c.Status(500)
		return
	}
	c.String(200, "Success")

}

func ReloadRunConfig() error {
	log.Println("Reload Run Config")
	configByte, err := os.ReadFile("./config.run.json")
	if err != nil {
		log.Println(err)
		return err
	}
	config := RunConfig{}
	err = json.Unmarshal(configByte, &config)
	if err != nil {
		log.Println(err)
		return err
	}
	NowRunConfig = config
	return nil
}
