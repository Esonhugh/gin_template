package config

import (
	"github.com/sirupsen/logrus"

	"gin_template/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var log = logrus.WithField("module", "config").WithField("server", "internal")

type Config struct {
	Viper *viper.Viper
	*ConfigContent
}

// GlobalConfig 默认全局配置
var GlobalConfig *Config

// Init 使用 ./application.yaml 初始化全局配置

func Init() {
	GlobalConfig = &Config{
		Viper: viper.New(),
	}
	GlobalConfig.Viper.SetConfigName("app")
	GlobalConfig.Viper.SetConfigType("yaml")
	GlobalConfig.Viper.AddConfigPath(".")
	GlobalConfig.Viper.AddConfigPath("../")    // For Debug
	GlobalConfig.Viper.AddConfigPath("../../") // For Debug
	GlobalConfig.Viper.AddConfigPath("/etc/" + utils.PackageName())

	err := GlobalConfig.Viper.ReadInConfig()
	if err != nil {
		log.Panic("Config Reading Error", err)
	}
	viper.OnConfigChange(func(in fsnotify.Event) {
		updateConfig()
	})
	// hot Change.
	viper.WatchConfig()
	updateConfig()
}

func updateConfig() {
	// ToDo: When New Config Added.
	err := GlobalConfig.Viper.Unmarshal(&GlobalConfig.ConfigContent)
	if err != nil {
		log.Warnf("Config Reading Error", err)
	}
	// database.reInit()
}
