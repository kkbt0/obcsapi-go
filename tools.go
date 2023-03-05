package main

import (
	"io/fs"
	"math/rand"
	"net"
	"os"
	"time"
)

const allowChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 生成随机 token2 只发送 token
func GengerateToken2(n int) string {
	ans := make([]byte, n)
	for i := range ans {
		ans[i] = allowChars[rand.Intn(len(allowChars))]
	}
	return string(ans)
}

// 更新 Token File
func ModTokenFile(new_token string, path string, token_class string) error {
	return os.WriteFile(path+token_class, []byte(new_token), fs.FileMode(os.O_WRONLY))
}

// 获取 token token_class 传入 token1(全权限，有效期) or token2（只能发送） 从而获取本地存储的 token 文件内容
func GetToken(path string, token_class string) ([]byte, error) {
	tokenFile, err := os.ReadFile(path + token_class)
	if err != nil {
		return nil, err
	}
	return tokenFile, nil
}

// Time fmt eg 2006-01-02 15:04:05
func timeFmt(fmt string) string {
	// fmt.Println(time.Now().In(cstZone).Format("2006-01-02 15:04:05"))
	var cstZone = time.FixedZone("CST", 8*3600)
	return time.Now().In(cstZone).Format(fmt)
}

// LocalIPs return all non-loopback IPv4 addresses
func LocalIPv4s() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	return ips, nil
}

// GetIPv4ByInterface return IPv4 address from a specific interface IPv4 addresses
func GetIPv4ByInterface(name string) ([]string, error) {
	var ips []string

	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	return ips, nil
}
