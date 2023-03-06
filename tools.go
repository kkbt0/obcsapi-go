package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/spf13/viper"
)

func ShowConfig() {
	// Read configuration
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error: Fatal error config file: %s \n ", err))
	}

	// output configuration
	log.Println(viper.GetString("name"), viper.GetString("version"), viper.GetString("description"))
	log.Println("Server Time:", timeFmt("2006-01-02 15:04"))
	log.Println("Tokne File Path:", viper.GetString("token_path"))
	log.Println("Run on", viper.GetString("host"))
}

// 从配置中获取 参数
func ConfigGetString(parm string) string {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error: Fatal error config file: %s \n ", err))
	}
	return viper.GetString(parm)
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
