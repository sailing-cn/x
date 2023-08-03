package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

func NewConfig(paths ...string) *ConnConfig {
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
	cfg := &ConnConfig{
		Host:          viper.GetString("amqp.host"),
		Port:          viper.GetInt("amqp.port"),
		User:          viper.GetString("amqp.user"),
		Password:      viper.GetString("amqp.password"),
		Vhost:         viper.GetString("amqp.vhost"),
		PrefetchCount: viper.GetInt("amqp.prefetch_count"),
	}
	return cfg
}

func Init() {
	conf.Conf = NewConfig()
}
