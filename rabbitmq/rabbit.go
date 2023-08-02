package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sailing.cn/utils"
	"sync"
)

var conf = struct {
	Conf *ConnConfig `yaml:"amqp"`
}{}

type rabbit struct {
	conn       *amqp.Connection
	chConsumer *amqp.Channel
	chProducer *amqp.Channel
	wgChannel  *sync.WaitGroup
	consumers  map[string]*amqp.Channel
}

func NewRabbitMQ() RabbitMQ {
	if conf.Conf == nil {
		log.Errorf("rabbit 配置信息不能为空")
		return nil
	}
	return &rabbit{
		wgChannel: &sync.WaitGroup{},
	}
}

func Init() {
	path := filepath.Join(utils.GetExecPath(), "conf.d", "conf.yml")
	file, err := os.ReadFile(path)
	if err != nil {
		log.Errorf("读取配置文件出错:%s", err)
	}
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		log.Errorf("解析配置文件出错:%s", err)
	}
}
