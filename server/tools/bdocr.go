package tools

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type OcrResult struct {
	WordsResults    []WordResult `json:"words_result"`
	WordsResultsNum int          `json:"words_result_num"`
	LogId           int          `json:"log_id"`
}
type WordResult struct {
	Words string `json:"words"`
}
type AccessToken struct {
	RefreshToken  string `json:"refresh_token"`
	ExpiresIn     int    `json:"expires_in"`
	SessionKey    string `json:"session_key"`
	AccessToken   string `json:"access_token"`
	Scope         string `json:"scope"`
	SessionSecret string `json:"session_secret"`
}

func BdGeneralBasicOcr(filePath string) ([]WordResult, error) {
	// OCR START https://ai.baidu.com/ai-doc/OCR/zk3h7xz52
	ocrUrl := fmt.Sprintf("https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic?access_token=%s", NowRunConfig.ImageHosting.BdOcrAccessToken)
	//  Read file and post # image url pdf_file
	f, err := os.Open(filePath)
	if err != nil {
		return []WordResult{}, err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return []WordResult{}, err
	}
	base64Content := base64.StdEncoding.EncodeToString(content)
	urlEncodedContent := url.QueryEscape(base64Content)
	body := fmt.Sprintf("image=%s", urlEncodedContent)

	res, err := http.Post(ocrUrl, "application/x-www-form-urlencoded", strings.NewReader(body))
	if err != nil {
		return []WordResult{}, err
	}
	defer res.Body.Close()
	content, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return []WordResult{}, err
	}
	var result OcrResult
	err = json.Unmarshal(content, &result)
	if err != nil {
		return []WordResult{}, err
	}
	return result.WordsResults, nil
}
func BdGetAccessToken(apiKey string, apiSecret string) (AccessToken, error) {
	var accessToken AccessToken

	url := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?client_id=%s&client_secret=%s&grant_type=client_credentials", apiKey, apiSecret)

	payload := strings.NewReader(``)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		return accessToken, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return accessToken, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return accessToken, err
	}
	json.Unmarshal(body, &accessToken)
	// Debug("Body: ", string(body))
	// Debug("AccessToken: ", accessToken)
	return accessToken, nil
}
