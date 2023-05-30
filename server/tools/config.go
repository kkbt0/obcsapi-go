package tools

import (
	"encoding/json"
	"log"
	"os"
	"time"

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
	ObDaily              string `json:"ob_daily"`
	ObDailyAttachmentDir string `json:"ob_daily_attachment_dir"`
	ObOtherDataDir       string `json:"ob_other_data_dir"`
}

type ReminderConfig struct {
	DailyEmailRemderTime string `json:"daily_email_remder_time"`
	ReminderDicionary    string `json:"reminder_dicionary"`
}

type MentionConfig struct {
	Tags []string `json:"tags"`
}

type RunConfig struct {
	ObDaily      ObsidianDailyConfig `json:"ob_daily_config"`
	WeChatMp     WeChatMpConfig      `json:"wechat_mp"`
	Webdav       WebDavConfig        `json:"webdav"`
	Mail         MailSmtpConfig      `json:"mail"`
	ImageHosting ImageHostingConfig  `json:"image_hosting"`
	Reminder     ReminderConfig      `json:"reminder"`
	Mention      MentionConfig       `json:"mention"`
}

func GetRunConfigHandler(c *gin.Context) {
	ReloadRunConfig()
	c.JSON(200, NowRunConfig)
}

func PostConfigHandler(c *gin.Context) {
	var config RunConfig = NowRunConfig
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

// 支持格式化时间 06日志/ 202303/2023-03-01

// 06日志/
func (runConfig *RunConfig) DailyDir() string {
	return runConfig.ObDaily.ObDailyDir
}

// 06日志/  + 202303/2023-03-01.md
func (runConfig *RunConfig) DailyFileKey() string {
	return runConfig.DailyDir() + TimeFmt(runConfig.ObDaily.ObDaily) + ".md"
}

// 06日志/  + 202302/2023-02-01.md
func (runConfig *RunConfig) DailyFileKeyMore(addDateDay int) string {
	return runConfig.DailyDir() + time.Now().AddDate(0, 0, addDateDay).In(time.FixedZone("CST", 8*3600)).Format(runConfig.ObDaily.ObDaily) + ".md"
}

// 202302/2023-02-01
func (runConfig *RunConfig) DailyDateKeyMore(addDateDay int) string {
	return time.Now().AddDate(0, 0, addDateDay).In(time.FixedZone("CST", 8*3600)).Format(runConfig.ObDaily.ObDaily)

}

func (runConfig *RunConfig) DailyFileKeyTime(inTime time.Time) string {
	diff := time.Until(inTime)
	return runConfig.DailyDir() + time.Now().Add(time.Hour*time.Duration(diff.Hours())).In(time.FixedZone("CST", 8*3600)).Format(runConfig.ObDaily.ObDaily) + ".md"
}

// 06日志/  + 附件/202302/
func (runConfig *RunConfig) DailyAttachmentDir() string {
	return runConfig.DailyDir() + TimeFmt(runConfig.ObDaily.ObDailyAttachmentDir)
}

func (runConfig *RunConfig) OtherDataDir() string {
	return TimeFmt(runConfig.ObDaily.ObOtherDataDir)
}
