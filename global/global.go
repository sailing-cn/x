package global

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sailing.cn/utils"
)

var (
	WebapiConf = &WebapiConfig{}
	GrpcConf   = &GrpcConfig{}
)

type ServerConfig struct {
	Name    string `json:"name"`
	Addr    string `json:"addr"`
	Port    int    `json:"port"`
	Version string `json:"version"`
}

type GrpcConfig struct {
	Grpc struct {
		Name    string `json:"name"`
		Addr    string `json:"addr"`
		Port    int    `json:"port"`
		Version string `json:"version"`
		MaxSend int    `json:"max_send" yaml:"max_send"`
		MaxRecv int    `json:"max_recv" yaml:"max_recv"`
	} `json:"grpc"`
	Services map[string]string `json:"services"`
}
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

type WebsocketConfig struct {
	Websocket struct {
		Addr string `yaml:"addr"`
		Port int    `yaml:"port"`
	}
}

func (c *GrpcConfig) Init(path string) {
	if len(path) == 0 {
		path = filepath.Join(utils.GetExecPath(), "conf.d", "conf.yml")
	}
	file, err := os.ReadFile(path)
	if err != nil {
		log.Errorf("读取配置文件出错:%s", err)
	}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		log.Errorf("解析配置文件出错:%s", err)
	}
}

func (c *WebapiConfig) Init(path string) {
	if len(path) == 0 {
		path = filepath.Join(utils.GetExecPath(), "conf.d", "conf.yml")
	}
	file, err := os.ReadFile(path)
	if err != nil {
		log.Errorf("读取配置文件出错:%s", err)
	}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		log.Errorf("解析配置文件出错:%s", err)
	}
}

func (c *WebsocketConfig) Init(path string) {
	if len(path) == 0 {
		path = filepath.Join(utils.GetExecPath(), "conf.d", "conf.yml")
	}
	file, err := os.ReadFile(path)
	if err != nil {
		log.Errorf("读取配置文件出错:%s", err)
	}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		log.Errorf("解析配置文件出错:%s", err)
	}
}

type Services struct {
	//Grpc struct {
	//	Name    string `json:"name"`
	//	Addr    string `json:"addr"`
	//	Port    int    `json:"port"`
	//	Version string `json:"version"`
	//	MaxSend int    `json:"max_send" yaml:"max_send"`
	//	MaxRecv int    `json:"max_recv" yaml:"max_recv"`
	//} `json:"grpc"`
	Services map[string]string `json:"services"`
}

func (c *Services) Init(path string) {
	if len(path) == 0 {
		path = filepath.Join(utils.GetExecPath(), "conf.d", "edge.yml")
	}
	file, err := os.ReadFile(path)
	if err != nil {
		log.Errorf("读取配置文件出错:%s", err)
	}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		log.Errorf("解析配置文件出错:%s", err)
	}
}
