package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"sailing.cn/v2/encrypt"
	"sailing.cn/v2/utils/timestamp"
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
	SetConnectionLostHandler(handler ConnectionLostHandler)
}

var (
	instance Client
)

type ConnectionLostHandler func(client Client, err error)

type client struct {
	Id                    string
	Password              string
	Server                string
	ServerCert            []byte
	Client                mqtt.Client
	qos                   byte
	connectionLostHandler ConnectionLostHandler
}

func (c *client) SetConnectionLostHandler(handler ConnectionLostHandler) {
	c.connectionLostHandler = handler
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
	options.SetClientID(AssembleClientId(c.Id))
	options.SetUsername(c.Id)
	options.SetPassword(encrypt.HmacSha256(c.Password, timestamp.TimestampString()[:10]))
	options.SetKeepAlive(60 * 2 * time.Second)
	options.SetAutoReconnect(true)
	options.SetConnectRetry(true)
	options.SetConnectTimeout(2 * time.Second)
	//options.SetCleanSession(false)
	options.SetCleanSession(true)
	options.OnConnectionLost = func(cl mqtt.Client, err error) {
		log.Errorf("与服务器断开连接")
		if c.connectionLostHandler != nil {
			cl.Disconnect(0)
			createClient(c, options)
			c.connectionLostHandler(c, err)
		}
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

	return createClient(c, options)
}

func createClient(c *client, options *mqtt.ClientOptions) bool {
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

func AssembleClientId(id string) string {
	segments := make([]string, 4)
	segments[0] = id
	segments[1] = "0"
	segments[2] = "0"
	segments[3] = timestamp.TimestampString()
	return strings.Join(segments, "_")
}

func FormatTopic(topic, deviceId string) string {
	return strings.ReplaceAll(topic, "{device_id}", deviceId)
}

func FormatTopicWithRequest(topic, deviceId string, requestId string) string {
	str := strings.ReplaceAll(topic, "{device_id}", deviceId)
	return strings.ReplaceAll(str, "#", "request_id ="+requestId)
}

func GetTopicDeviceId(topic string) string {
	//$oc/devices/1458634485842055168_1/sys/properties/report
	return strings.Split(topic, "/")[2]
}
