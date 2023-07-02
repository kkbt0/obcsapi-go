package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

// AuthorizationCode 结构体表示一次性授权码
type AuthorizationCode struct {
	Code       string
	Expiration time.Time
	Used       bool
	Init       bool
}

// GenerateAuthorizationCode 生成一次性授权码
func GenerateAuthorizationCode(expiration time.Time) (*AuthorizationCode, error) {
	// 生成随机字节
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	// 将字节转换为 base64 编码的字符串
	code := base64.URLEncoding.EncodeToString(bytes)

	authCode := &AuthorizationCode{
		Code:       code,
		Expiration: expiration,
		Used:       false,
		Init:       true,
	}

	return authCode, nil
}

// IsExpired 判断授权码是否过期
func (ac *AuthorizationCode) IsExpired() bool {
	return time.Now().After(ac.Expiration)
}

// IsUsed 判断授权码是否被使用过
func (ac *AuthorizationCode) IsUsed() bool {
	return ac.Used
}

func (ac *AuthorizationCode) IsValid(inCodeStr string) bool {
	if ac.Code != inCodeStr {
		return false
	}
	if !ac.Init || ac.IsExpired() || ac.IsUsed() {
		return false
	}
	ac.Used = true
	return true
}
