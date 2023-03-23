package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"
)

// 用于存储
type Token struct {
	TokenString  string `json:"token"`
	GenerateTime string `json:"generate_time"`
}

// 用于 http json 格式解析
type TokenFromJSON struct {
	TokenString string `json:"token"`
}

const allowChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 生成随机 token
func GengerateToken(n int) string {
	rand.Seed(time.Now().Unix()) // 保证每秒生成不同的随机 token  , unix 时间戳，秒
	ans := make([]byte, n)
	for i := range ans {
		ans[i] = allowChars[rand.Intn(len(allowChars))]
	}
	return string(ans)
}

// 更新 Token File
func ModTokenFile(new_token Token, token_class string) error {
	data, err := json.Marshal(&new_token)
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigGetString("token_path")+token_class, data, 0666)
}

// 获取 token token_class 传入 token1(全权限，有效期) or token2（只能发送） 从而获取本地存储的 token 文件内容
func GetToken(token_class string) (Token, error) {
	tokenBytes, err := os.ReadFile(ConfigGetString("token_path") + token_class)
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

func VerifyToken1(inToken string) bool {
	rightToken, err := GetToken("token1")
	if err != nil {
		log.Println("Token Get Error:", err)
		return false
	}
	nowTime, err := time.Parse("2006-01-02 15:04:05", timeFmt("2006-01-02 15:04:05"))
	if err != nil {
		log.Println(err)
		return false
	}
	rightTokenTime, err := time.Parse("2006-01-02 15:04:05", rightToken.GenerateTime)
	if err != nil {
		log.Println(err)
		return false
	}
	liveTime, _ := time.ParseDuration(ConfigGetString("token1_live_time"))
	// log.Println(nowTime, rightTokenTime.Add(liveTime))
	// 验证 Token 相符合 且 现在时间 < 生成时间 + 存活时间
	if inToken == rightToken.TokenString && nowTime.Before(rightTokenTime.Add(liveTime)) {
		return true
	}
	return false
}

func VerifyToken2(inToken string) bool {
	rightToken, err := GetToken("token2")
	if err != nil {
		log.Println("Token Get Error:", err)
		return false
	}
	// 验证 Token 相符合
	if inToken == rightToken.TokenString {
		return true
	}
	return false
}
