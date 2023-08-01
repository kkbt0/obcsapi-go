package app

import (
	"encoding/json"
	"fmt"
	"obcsapi-go/dao"
	"obcsapi-go/gr"
	"obcsapi-go/tools"

	"github.com/gin-gonic/gin"
)

type MoodReaderHighlights struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	Author string `json:"author"`
	Note   string `json:"note"`
}

type MoodReader struct {
	Highlights []MoodReaderHighlights `json:"highlights"`
}

// Token2 静读天下使用的 API
// @Summary 静读天下使用的 API
// @Description 静读天下使用的 API，标注-设置-ReadWise 设置该路径和 token 值即可
// @Tags Ob
// @Security AuthorizationToken
// @Accept json
// @Produce json
// @Param 划线和标注 body MoodReader true "MoodReader"
// @Router /ob/moonreader [post]
func MoodReaderHandler(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var moodReader MoodReader
	err := decoder.Decode(&moodReader)
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}
	fmt.Println(moodReader.Highlights[0])
	title := moodReader.Highlights[0].Title
	text := moodReader.Highlights[0].Text
	author := moodReader.Highlights[0].Author
	note := moodReader.Highlights[0].Note
	file_key := fmt.Sprintf("%sMoonReader/%s.md", tools.NowRunConfig.OtherDataDir(), tools.ReplaceUnAllowedChars(title))
	append_text := fmt.Sprintf("文: %s\n批: %s\n于: %s\n\n---\n", text, note, tools.TimeFmt("2006-01-02 15:04"))
	exist, _ := dao.CheckObject(file_key)
	if exist {
		err = dao.TextAppend(file_key, append_text)
	} else {
		yaml := fmt.Sprintf("---\ntitle: %s\nauthor: %s\n---\n书名: %s\n作者: %s\n简介: \n评价: \n\n---\n", title, author, title, author)
		err = dao.TextAppend(file_key, yaml+append_text)
	}
	if err != nil {
		gr.ErrServerError(c, err)
		return
	}
	gr.Success(c)
}
