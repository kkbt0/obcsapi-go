package tools

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// 用于存储 Token
//
// VerifyMode: Authentication Token Query(token) (可以从对话模式编写bash脚本获取本地token文件内容)
//
// LiveHours == "" or error ==> 0s
type Token struct {
	TokenString  string `json:"token"`
	GenerateTime string `json:"generate_time"`
	LiveTime     string `json:"live_time"`
	VerifyMode   string `json:"verify_mode"`
}

func (token *Token) IsLiving() bool {
	// 现在时间 < 生成时间 + 存活时间
	nowTime, _ := time.Parse("2006-01-02 15:04:05", TimeFmt("2006-01-02 15:04:05"))
	rightTokenTime, err := time.Parse("2006-01-02 15:04:05", token.GenerateTime)
	if err != nil {
		log.Println(err)
	}
	liveTime, _ := time.ParseDuration(token.LiveTime)
	if !nowTime.Before(rightTokenTime.Add(liveTime)) {
		Debug("Live Time:", liveTime)
		Debug("Now: ", nowTime, "AllowLiveTime", rightTokenTime.Add(liveTime))
		return false
	}
	return true
}

func (token *Token) Update() {
	// Debug("This Is A Expired Token. It will be updated!")
	token.TokenString = GengerateTokenString(32)
	token.GenerateTime = TimeFmt("2006-01-02 15:04:05")
}

// 用于 http json 格式解析
type TokenFromJSON struct {
	TokenString string `json:"token"`
}

const allowChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 生成随机 token
func GengerateTokenString(n int) string {
	rand.Seed(time.Now().UnixMilli()) // 保证每秒生成不同的随机 token  , unix 时间戳，ms
	ans := make([]byte, n)
	for i := range ans {
		ans[i] = allowChars[rand.Intn(len(allowChars))]
	}
	return string(ans)
}

// 更新 Token File
func ModTokenFile(new_token Token, tokenFilePath string) error {
	data, err := json.Marshal(&new_token)
	if err != nil {
		return err
	}
	return os.WriteFile(tokenFilePath, data, 0666)
}

// 获取 token tokenFilePath 传入 token1(全权限，有效期) or token2（只能发送） 从而获取本地存储的 token 文件内容
func GetToken(tokenFilePath string) (Token, error) {
	tokenBytes, err := os.ReadFile(tokenFilePath)
	if err != nil {
		return Token{}, err
	}
	token := Token{}
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		return Token{}, err
	} else {
		return token, nil
	}
}

// rightToken 正确的 Token
// inToken 要验证的 token 字符串
func VerifyToken(rightToken Token, inToken string) bool {
	Debug("InToken: ", inToken)
	Debug("Right Token: ", rightToken)
	// 验证 Token 相符合 且未过期
	return inToken == rightToken.TokenString && rightToken.IsLiving()
}

// tokenFilePath Token文件位置 并且自动更新
func VerifyTokenByFilePath(tokenFilePath string, inToken string) bool {
	rightToken, err := GetToken(tokenFilePath)
	if err != nil {
		log.Println("Token Get Error:", err)
		return false
	}
	Debug("InToken: ", inToken)
	Debug("Right Token: ", rightToken.TokenString)
	Debug("Right Token: ", rightToken)
	// 验证 Token 相符合
	if inToken == rightToken.TokenString {
		//如果在符合的情况下但是过期了
		if !rightToken.IsLiving() {
			rightToken.Update()
			ModTokenFile(rightToken, tokenFilePath)
			return false
		}
		return true
	}
	return false
}

// 判断 requestURI 为空 则使用 defaultTokenFilePath
func RequestURIDefineOrDefalutTokenFilePath(requestURI string, defaultTokenFilePath string) string {
	Debug("RequestURI: ", requestURI)
	// ?xxx=xxx
	index := strings.Index(requestURI, "?")
	if index != -1 {
		requestURI = requestURI[:index]
	}
	tokenFilePath := ConfigGetString(requestURI)
	if tokenFilePath == "" {
		Debug("config.yaml does not define which token to use")
		Debug("Use Middleware Define Default Token: ", defaultTokenFilePath)
		tokenFilePath = defaultTokenFilePath
	}
	return tokenFilePath
}
