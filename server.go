package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/spf13/viper"
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
	http.ListenAndServe(fmt.Sprintf("%s:%s", viper.GetString("host"), viper.GetString("port")), nil)
}
