package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

type IndexInfo struct {
	Title   string
	Content string
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.RequestURI)
	tpl, err := template.ParseFiles("./template/index.html")
	if err != nil {
		log.Panicln("Template File Error:", err)
		return
	}
	indexInfo := IndexInfo{Title: "404", Content: "404 Not Found"}
	tpl.Execute(w, indexInfo)
}

func Config() {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error: Fatal error config file: %s \n ", err))
	}
}

func main() {
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

	local_ip, _ := LocalIPv4s()
	log.Printf("LocalIp http://%s:%s\n", local_ip[0], viper.GetString("port"))

	http.HandleFunc("/", HelloHandler)
	http.Handle("/api/sendtoken2mail", limit(http.HandlerFunc(SendTokenHandler))) // 对请求将 token发送到 email 速率进行限制
	http.ListenAndServe(fmt.Sprintf("%s:%s", viper.GetString("host"), viper.GetString("port")), nil)
}

var limiter = rate.NewLimiter(0.3, 1)

// 短时间多次请求限制
func limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
