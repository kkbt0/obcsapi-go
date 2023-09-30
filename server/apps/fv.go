package apps

import (
	"fmt"
	"io"
	"obcsapi-go/dao"
	"obcsapi-go/gr"
	"obcsapi-go/tools"

	"github.com/gin-gonic/gin"
)

// @Summary fv 悬浮球使用的 API
// @Description 安卓软件 fv 悬浮球使用的 API 用于自定义任务的 图片、文字
// @Tags Ob
// @Security Token
// @Accept plain,octet-stream
// @Produce json
// @Param 內容 body string true "fv payload 內容"
// @Router /ob/fv [post]
func FvHandler(c *gin.Context) {
	if c.GetHeader("Content-Type") == "text/plain" {
		content, _ := io.ReadAll(c.Request.Body)
		dao.DailyTextAppendMemos(string(content))
		gr.Success(c)
		return
	} else if c.GetHeader("Content-Type") == "application/octet-stream" {
		content, _ := io.ReadAll(c.Request.Body)
		file_key := fmt.Sprintf("%s%s.jpg", tools.NowRunConfig.DailyAttachmentDir(), tools.TimeFmt("20060102150405"))
		dao.ObjectStore(file_key, content)
		dao.DailyTextAppendMemos(fmt.Sprintf("![](%s)", file_key))
		gr.Success(c)
		return
	}
	gr.RJSON(c, nil, 404, 404, "Error Content Type", gin.H{})
}
