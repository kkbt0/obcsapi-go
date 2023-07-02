package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"obcsapi-go/gr"
	"obcsapi-go/models"
	"obcsapi-go/tools"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var stateRandom string = tools.RandomString(8)
var oauth2UserInfo = "https://gitee.com/api/v5/user"
var oauth2Config = &oauth2.Config{
	ClientID:     viper.GetString("oauth2_gitee_ClientID"),
	ClientSecret: viper.GetString("oauth2_gitee_ClientSecret"),
	Scopes:       []string{"user_info"}, // Adjust the scopes based on your requirements
	Endpoint: oauth2.Endpoint{
		AuthURL:   "https://api.gitee.com/oauth/authorize",
		TokenURL:  "https://api.gitee.com/oauth/token",
		AuthStyle: 0,
	},
	RedirectURL: viper.GetString("backend_url_full") + "/auth/oauth2-callback",
}

// 一次性授权码
var authorizationCode *AuthorizationCode

// 请求授权
func InitiateOAuthFlow(c *gin.Context) {
	stateRandom = tools.RandomString(8)
	c.Redirect(http.StatusTemporaryRedirect, oauth2Config.AuthCodeURL(stateRandom))
}

// 请求注册授权信息 返回注册 URL
func InitiateOAuthFlowAndSetState(c *gin.Context) {
	stateRandom = "set_" + tools.RandomString(4)
	c.String(200, oauth2Config.AuthCodeURL(stateRandom))
	// c.Redirect(http.StatusTemporaryRedirect, oauth2Config.AuthCodeURL(stateRandom))
}

func HandleCallback(c *gin.Context) {
	if c.Query("state") != stateRandom {
		gr.RJSON(c, nil, 400, 400, "Invalid State", gr.H{})
		return
	}
	code := c.Query("code")
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		gr.RJSON(c, nil, 400, 400, "Failed to exchange code for token", gr.H{})
		return
	}
	client := oauth2Config.Client(context.Background(), token)
	// log.Println(token)
	resp, err := client.Get(oauth2UserInfo)
	if err != nil {
		gr.RJSON(c, err, 500, 500, "Failed to retrieve user information", gr.H{})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		gr.RJSON(c, err, 500, 500, "Failed to retrieve user information", gr.H{})
		return
	}
	// Read the json
	decoder := json.NewDecoder(resp.Body)
	var giteeUserInfo models.GiteeUserInfo
	if err = decoder.Decode(&giteeUserInfo); err != nil {
		gr.ErrServerError(c, err)
		return
	}
	// 请求注册授权信息 设置信息
	if strings.HasPrefix(c.Query("state"), "set_") {
		tools.NowRunConfig.OAuth2UserInfo.GiteeUserInfo = giteeUserInfo
		if err = tools.UpdateConfig(tools.NowRunConfig); err != nil {
			gr.ErrServerError(c, err)
			return
		}
		c.Redirect(302, viper.GetString("web_url_full"))
		return
	}
	// 登录
	if giteeUserInfo == tools.NowRunConfig.OAuth2UserInfo.GiteeUserInfo {
		// 前端带 一次性 授权码 的 URL
		authorizationCode, _ = GenerateAuthorizationCode(time.Now().Add(30 * time.Second))
		c.Redirect(302, viper.GetString("web_url_full")+"?code="+authorizationCode.Code)
		return
	}
	gr.RJSON(c, nil, 500, 500, "No Connect OAuth2.0", gr.H{})
}

type CodeStuct struct {
	Code string `json:"code"`
}

// 根据授权码分发 jwt token
func LoginByCodeHandler(c *gin.Context) {
	var codeStuct CodeStuct
	if c.ShouldBindJSON(&codeStuct) != nil {
		gr.ErrBindJSONErr(c)
		return
	}
	if codeStuct.Code == "" {
		gr.ErrEmpty(c)
		return
	}
	tools.Debug(authorizationCode)
	if authorizationCode.IsValid(codeStuct.Code) {
		tokenString, err := GetTokenString()
		if err != nil {
			gr.ErrServerError(c, err)
			return
		}
		c.JSON(200, gin.H{
			"code":  200,
			"token": tokenString,
			"msg":   "登录成功",
		})
	} else {
		gr.RJSON(c, nil, 400, 400, "Invalid Code", gr.H{})
	}
}
