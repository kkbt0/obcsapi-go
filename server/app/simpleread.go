package app

import (
	"encoding/json"
	"fmt"
	"obcsapi-go/dao"
	"obcsapi-go/gr"
	"obcsapi-go/tools"

	"github.com/gin-gonic/gin"
)

// SimpRead WebHook Used
type SimpReadWebHookStruct struct {
	Title       string `json:"title"`
	Url         string `json:"url"`
	Description string `json:"desc"`
	Tags        string `json:"tags"`
	Content     string `json:"content"`
	Note        string `json:"note"`
}

// SimpRead WebHook Used Token2
// @Summary 简悦 WebHook 保存文章
// @Description SimpRead 简悦 WebHook POST 简悦 WebHook 保存文章
// @Tags Ob
// @Security Token
// @Accept json
// @Param json body SimpReadWebHookStruct true "SimpRead 简悦 POST"
// @Router /ob/sr/webhook [post]
func SRWebHook(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var simpReadJson SimpReadWebHookStruct
	var err error
	if err = decoder.Decode(&simpReadJson); err != nil {
		gr.ErrServerError(c, err)
		return
	}
	serverTime := tools.TimeFmt("200601021504")
	yaml := fmt.Sprintf("---\ntitle: %s\nsctime: %s\n---\n", simpReadJson.Title, serverTime)
	file_str := fmt.Sprintf("%s[[简悦WebHook生成]]\n生成时间: %s\n原文: %s\n标题: %s\n描述: %s\n标签: %s\n内容: \n%s", yaml, serverTime, simpReadJson.Url, simpReadJson.Title, simpReadJson.Description, simpReadJson.Tags, simpReadJson.Content)
	file_key := fmt.Sprintf("%sSimpRead/%s %s.md", tools.NowRunConfig.OtherDataDir(), tools.ReplaceUnAllowedChars(simpReadJson.Title), serverTime)
	if err = dao.MdTextStore(file_key, file_str); err != nil {
		gr.ErrServerError(c, err)
		return
	}
	gr.Success(c)
}
