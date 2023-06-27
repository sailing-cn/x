package conf

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	WebapiConf = &WebapiConfig{}
	GrpcConf   = &GrpcConfig{}
)

// GrpcConfig GRPC服务配置
type GrpcConfig struct {
	Grpc struct {
		Name    string `json:"name"`                     //服务名称
		Addr    string `json:"addr"`                     //监听地址
		Port    int    `json:"port"`                     //监听端口
		Version string `json:"version"`                  //版本
		MaxSend int    `json:"max_send" yaml:"max_send"` //最大发送字节
		MaxRecv int    `json:"max_recv" yaml:"max_recv"` //最大接收字节
	} `json:"grpc"` //grpc 配置项
	Services map[string]string `json:"services"` //依赖服务配置项
}

// WebapiConfig WEBAPI配置
type WebapiConfig struct {
	Webapi struct {
		Name    string `json:"name"`    //网关名称
		Addr    string `json:"addr"`    //监听地址
		Port    int    `json:"port"`    //监听端口
		Version string `json:"version"` //版本
		Prefix  string `json:"prefix"`  //路由前缀
		Mode    string `json:"mode"`    //模式 dev:开发模式
	} `json:"webapi"` //api配置项
	Services map[string]string `json:"services"` //依赖服务配置项
}

// NewWebapiConfig 创建一个WEBAPI配置,从配置文件读取,默认路径 ./conf.d/conf.yml
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
	WebapiConf = cfg
	return cfg
}

// NewGrpcConfig 创建一个GRPC服务配置,从配置文件读取,默认路径 ./conf.d/conf.yml
func NewGrpcConfig(paths ...string) *GrpcConfig {
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
	cfg := &GrpcConfig{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		log.Panicf("解析配置文件错误:%s", err.Error())
		panic(1)
	}
	GrpcConf = cfg
	return cfg
}
