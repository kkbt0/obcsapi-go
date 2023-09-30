package wechat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/viper"
)

type AccessTokenResponse struct {
	Errcode     int    `json:"errcode"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func FetchNewAccessToken() (string, error) {
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", viper.GetString("work_wechat_corpid"), viper.GetString("work_wechat_corpsecret"))
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var tokenResp AccessTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	if err := os.WriteFile("workwechat_access_token.txt", []byte(tokenResp.AccessToken), 0644); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func readAccessTokenFromFile() (string, error) {
	// Read access token from the local file
	token, err := os.ReadFile("workwechat_access_token.txt")
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func GetWorkWechatAccessToken() (string, error) {
	var token string
	var err error
	if isTokenExpired() {
		// Token has expired, make a new request
		token, err = FetchNewAccessToken()
		if err != nil {
			return "", err
		}
		return token, nil
	}
	token, err = readAccessTokenFromFile()
	if err != nil {
		return "", err
	}
	return token, nil
}

func isTokenExpired() bool {
	fileInfo, err := os.Stat("workwechat_access_token.txt")
	if err != nil {
		return true // Return true to fetch a new token if there is any error with file info
	}
	// Get the last modified time of the file
	lastModifiedTime := fileInfo.ModTime()
	// Calculate the time difference
	duration := time.Since(lastModifiedTime)
	// Compare with 2 hours (2 hours * 60 minutes * 60 seconds)
	return duration.Seconds() > 2*60*60
}
