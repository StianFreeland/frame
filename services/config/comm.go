package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const moduleName = "config"

func Init() bool {
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(moduleName, "init, error:", err)
		return false
	}
	return true
}
