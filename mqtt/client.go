package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"sailing.cn/iot"
	"strings"
	"time"
)

type Config struct {
	ClientId string `yaml:"client_id"`
	Password string `yaml:"password"`
	Server   string `yaml:"server"`
	Qos      byte   `yaml:"qos"`
}

type Client interface {
	ClientId() string
	Init() bool
	DisConnect()
	IsConnected() bool
	Subscribe(topic string, handler func(msg Message))
	Subscribe2(topic string, handler func(c mqtt.Client, msg mqtt.Message)) error
	Publish(topic string, body interface{}) error
}

var (
	instance Client
)

type client struct {
	Id         string
	Password   string
	Server     string
	ServerCert []byte
	Client     mqtt.Client
	qos        byte
}

func (c *client) DisConnect() {
	if c.Client != nil {
		c.Client.Disconnect(0)
	}
}

func (c *client) IsConnected() bool {
	if c.Client != nil {
		return c.Client.IsConnectionOpen()
	}
	return false
}

func (c *client) Subscribe(topic string, handler func(msg Message)) {
	token := c.Client.Subscribe(topic, c.qos, func(c1 mqtt.Client, message mqtt.Message) {
		handler(Message{message})
		message.Ack()
	})
	if token.Wait() && token.Error() != nil {
		log.Errorf("订阅主题 %s 失败：%s", topic, token.Error().Error())
	}
}

func (c *client) Subscribe2(topic string, handler func(c mqtt.Client, msg mqtt.Message)) error {
	token := c.Client.Subscribe(topic, c.qos, handler)
	if token.Wait() && token.Error() != nil {
		log.Errorf("订阅主题 %s 失败：%s", topic, token.Error().Error())
		return token.Error()
	}
	return nil
}

func (c *client) Publish(topic string, body interface{}) error {
	var bytes []byte
	if v, ok := body.([]byte); ok {
		bytes = v
	} else if v, ok := body.(string); ok {
		bytes = []byte(v)
	} else {
		bytes, _ = json.Marshal(body)
	}
	token := c.Client.Publish(topic, c.qos, false, bytes)
	if token.Wait() && token.Error() != nil {
		log.Errorf("发送消息 %s 失败：%s", topic, token.Error().Error())
		return token.Error()
	}
	return nil
}

func Get(cnf *Config) Client {
	if instance == nil || instance.IsConnected() == false {
		instance = create(cnf.ClientId, cnf.Password, cnf.Server)
	}
	return instance
}

func (c *client) Init() bool {
	options := mqtt.NewClientOptions()
	options.AddBroker(c.Server)
	options.SetClientID(assembleClientId(c))
	options.SetUsername(c.Id)
	options.SetPassword(iot.HmacSha256(c.Password, iot.TimestampString()))
	options.SetKeepAlive(60 * 2 * time.Second)
	options.SetAutoReconnect(true)
	options.SetConnectRetry(true)
	options.SetConnectTimeout(2 * time.Second)
	//options.SetCleanSession(false)
	options.SetCleanSession(true)
	options.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Errorf("与服务器断开连接")
	}
	options.OnReconnecting = func(c mqtt.Client, options *mqtt.ClientOptions) {
		log.Errorf("正在重新连接")
	}
	if strings.Contains(c.Server, "tls") || strings.Contains(c.Server, "ssl") {
		log.Infof("server support tls connection")
		if c.ServerCert != nil {
			certPool := x509.NewCertPool()
			certPool.AppendCertsFromPEM(c.ServerCert)
			options.SetTLSConfig(&tls.Config{
				RootCAs:            certPool,
				InsecureSkipVerify: false,
			})
		} else {
			options.SetTLSConfig(&tls.Config{
				InsecureSkipVerify: true,
			})
		}
	} else {
		options.SetTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
		})
	}

	c.Client = mqtt.NewClient(options)
	if token := c.Client.Connect(); token.Wait() && token.Error() != nil {
		log.Warningf("mqtt client %s init failed,error = %v", c.Id, token.Error())
		return false
	}
	return true
}

func (c *client) ClientId() string {
	return c.Id
}

func create(user, password, servers string) Client {
	config := Config{
		ClientId: user,
		Password: password,
		Server:   servers,
	}
	return createWithConfig(config)
}

func createWithConfig(cnf Config) Client {
	c := client{}
	c.Id = cnf.ClientId
	c.Password = cnf.Password
	c.Server = cnf.Server
	c.qos = cnf.Qos
	c.Init()
	return &c
}

func assembleClientId(c *client) string {
	segments := make([]string, 4)
	segments[0] = c.Id
	segments[1] = "0"
	segments[2] = "0"
	segments[3] = iot.TimestampString()
	return strings.Join(segments, "_")
}
