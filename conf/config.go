package conf

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type WebapiConfig struct {
	Webapi struct {
		Name    string `json:"name"`
		Addr    string `json:"addr"`
		Port    int    `json:"port"`
		Version string `json:"version"`
		Prefix  string `json:"prefix"`
		Mode    string `json:"mode"`
	} `json:"webapi"`
	Services map[string]string `json:"services"`
}

func NewWebapiConfig(paths ...string) *WebapiConfig {
	if len(paths) == 0 {
		viper.AddConfigPath("./conf.d/")
		viper.SetConfigName("conf")
		viper.SetConfigType("yaml")
	} else {
		for _, s := range paths {
			viper.AddConfigPath(s)
		}
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Panicf("找不到配置文件,请检查配置文件路径是否正确,默认路径 ./conf.d/conf.yml")
			panic(1)
		} else {
			log.Panicf("读取配置文件错误:%s", err.Error())
		}
	}
	cfg := &WebapiConfig{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		log.Panicf("解析配置文件错误:%s", err.Error())
		panic(1)
	}
	return cfg
}
