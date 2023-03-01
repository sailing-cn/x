package logs

import (
	"context"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sailing.cn/logs/amqp"
	"sailing.cn/rabbitmq"
	"sailing.cn/utils"
)

var (
	conf = &Config{}
)

type LogConfig struct {
	AMQP *rabbitmq.ConnConfig `json:"amqp" yaml:"amqp"`
	Type string               `json:"type" yaml:"type"`
}

type Config struct {
	Log *LogConfig `json:"logs" yaml:"logs"`
}

func UseRabbitMQ(ctx context.Context, app string) {
	Init()
	hook, err := amqp.NewHook(ctx, app)
	if hook == nil {
		logrus.Errorf("初始化 amqp hook 错误:%s", err)
		return
	}
	logrus.AddHook(hook)
}

func Init() {
	path := filepath.Join(utils.GetExecPath(), "conf.d", "conf.yml")
	file, err := os.ReadFile(path)
	if err != nil {
		logrus.Errorf("读取配置文件出错:%s", err)
	}
	err = yaml.Unmarshal(file, conf)
	if err != nil {
		logrus.Errorf("解析配置文件出错:%s", err)
	}
}
