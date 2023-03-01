package amqp

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"sailing.cn/rabbitmq"
	"sailing.cn/wrong"
	"strings"
)

type AmqpHook struct {
	rabbit  rabbitmq.RabbitMQ
	AppName string
}

func (hook AmqpHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func NewHook(ctx context.Context, appName string) (*AmqpHook, error) {
	rabbitmq.Init()
	var rabbit = rabbitmq.NewRabbitMQ()
	if rabbit != nil {
		var setup rabbitmq.Setup = func() {
		}
		rabbitmq.KeepConnectionAndSetup(ctx, rabbit, setup)
		ctx.Value("")
		hook := AmqpHook{rabbit, appName}
		return &hook, nil
	} else {
		return nil, wrong.New("初始化日志rabbitmq失败")
	}
}

func (hook AmqpHook) Fire(entry *logrus.Entry) error {
	var msg = format(entry, hook.AppName)
	return hook.rabbit.Publish("sailing.log", "logs", msg)
}

func format(entry *logrus.Entry, appName string) map[string]interface{} {
	m := make(map[string]interface{})
	m["timestamp"] = entry.Time.UnixNano() / 1e6
	m["creation_time"] = entry.Time.Format("2006-01-02 15:04:05")
	m["target"] = findCaller(5)
	m["message"] = entry.Message
	m["app"] = appName
	m["level"] = entry.Level.String()
	return m
}

// 对caller进行递归查询, 直到找到非logrus包产生的第一个调用.
// 因为filename我获取到了上层目录名, 因此所有logrus包的调用的文件名都是 logrus/...
// 因此通过排除logrus开头的文件名, 就可以排除所有logrus包的自己的函数调用
func findCaller(skip int) string {
	file := ""
	line := 0
	for i := 0; i < 10; i++ {
		file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	//fmt.Println(file)
	//fmt.Println(line)
	if !ok {
		return "", 0
	}
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line
}
