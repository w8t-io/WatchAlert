package initialize

import (
	"github.com/spf13/viper"
	"log"
	"watchAlert/globals"
)

var (
	configFile = "config/config.yaml"
)

func InitConfig() {

	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatal("配置读取失败:", err)
	}
	if err := v.Unmarshal(&globals.Config); err != nil {
		log.Fatal("配置解析失败:", err)
	}
}
