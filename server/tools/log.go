package tools

import (
	"log"

	"github.com/spf13/viper"
)

func Debug(v ...interface{}) {
	if viper.GetBool("debug") {
		debugMsg := append([]interface{}{"[ObC-Debug]"}, v...)
		log.Println(debugMsg...)
	}
}
