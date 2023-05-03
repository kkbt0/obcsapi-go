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

type MailSmtpConfig struct {
	SmtpHost      string `json:"smtp_host"`
	Port          int    `json:"smtp_port"`
	UserName      string `json:"user_name"`
	Password      string `json:"password"`
	SenderEmail   string `json:"sender_email"`
	SenderName    string `json:"sender_name"`
	ReceiverEmail string `json:"receiver_email"`
}

type ImageHostingConfig struct {
	BaseURL          string `json:"base_url"`
	Prefix           string `json:"prefix"`
	UseRawName       bool   `json:"use_raw_name"`
	RandomCharLength int    `json:"random_char_length"`
	BdOcrAccessToken string `json:"bd_ocr_access_token"`
}

type WeChatMpConfig struct {
	ReturnStr string `json:"return_str"`
}

type ObsidianDailyConfig struct {
	ObDailyDir           string `json:"ob_daily_dir"`
	ObDailyAttachmentDir string `json:"ob_daily_attachment_dir"`
	ObOtherDataDir       string `json:"ob_other_data_dir"`
}

type ReminderConfig struct {
	DailyEmailRemderTime string `json:"daily_email_remder_time"`
	ReminderDicionary    string `json:"reminder_dicionary"`
}

type RunConfig struct {
	Webdav       WebDavConfig        `json:"webdav"`
	Mail         MailSmtpConfig      `json:"mail"`
	ImageHosting ImageHostingConfig  `json:"image_hosting"`
	WeChatMp     WeChatMpConfig      `json:"wechat_mp"`
	ObDaily      ObsidianDailyConfig `json:"ob_daily_config"`
	Reminder     ReminderConfig      `json:"reminder"`
}

func GetRunConfigHandler(c *gin.Context) {
	ReloadRunConfig()
	c.JSON(200, NowRunConfig)
}

func PostConfigHandler(c *gin.Context) {
	var config RunConfig
	err := c.ShouldBindJSON(&config)
	if err != nil {
		c.Error(err)
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

// 支持格式化时间
func (runConfig *RunConfig) DailyDir() string {
	return TimeFmt(runConfig.ObDaily.ObDailyDir)
}
func (runConfig *RunConfig) DailyAttachmentDir() string {
	return TimeFmt(runConfig.ObDaily.ObDailyAttachmentDir)
}

func (runConfig *RunConfig) OtherDataDir() string {
	return TimeFmt(runConfig.ObDaily.ObOtherDataDir)
}
