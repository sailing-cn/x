package iot_platform

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"sailing.cn/iot"
	"sailing.cn/iot/event_type"
	"sailing.cn/iot/iot.platform/topics"
	"strings"
	"time"
)

type Config struct {
	Id                 string
	Password           string
	Servers            string
	Qos                byte
	BatchSubDeviceSize int
}

type BasePlatform interface {
	Init() bool
	DisConnect()
	IsConnected() bool
}

type basePlatform struct {
	Id                          string
	Password                    string
	Servers                     string
	ServerCert                  []byte
	Client                      mqtt.Client
	fileUrls                    map[string]string
	qos                         byte
	batchSubDeviceSize          int
	commandHandlers             []CommandHandler
	messageHandlers             []MessageHandler
	propertiesReportHandler     PropertiesReportHandler
	batchReportHandler          GatewayBatchReportSubDeviceHandler
	deviceInfoReportHandler     DeviceInfoReportHandler
	deviceOnlineHandler         DeviceOnlineHandler
	deviceOfflineHandler        DeviceOfflineHandler
	versionReportHandler        VersionReportHandler
	upgradeProcessReportHandler UpgradeProcessReportHandler
}

func (p *basePlatform) Init() bool {
	options := mqtt.NewClientOptions()
	options.AddBroker(p.Servers)
	options.SetClientID(assembleClientId(p))
	options.SetUsername(p.Id)
	options.SetPassword(iot.HmacSha256(p.Password, iot.TimestampString()))
	options.SetKeepAlive(250 * time.Second)
	options.SetAutoReconnect(true)
	options.SetConnectRetry(true)
	options.SetConnectTimeout(2 * time.Second)
	if strings.Contains(p.Servers, "tls") || strings.Contains(p.Servers, "ssl") {
		log.Infof("server support tls connection")
		if p.ServerCert != nil {
			certPool := x509.NewCertPool()
			certPool.AppendCertsFromPEM(p.ServerCert)
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

	p.Client = mqtt.NewClient(options)
	if token := p.Client.Connect(); token.Wait() && token.Error() != nil {
		log.Warningf("device %s init failed,error = %v", p.Id, token.Error())
		return false
	}
	return true
}
func (p *basePlatform) DisConnect() {
	if p.Client != nil {
		p.Client.Disconnect(0)
	}
}

func (p *basePlatform) IsConnected() bool {
	if p.Client != nil {
		return p.Client.IsConnectionOpen()
	}
	return false
}

func (p *basePlatform) SendMessage(message iot.Message) bool {
	panic("implement me" + message.Name)
}

func (p *basePlatform) subscribeDefaultTopics() {
	//设备上线
	if token := p.Client.Subscribe(topics.DEVICE_ONLINE, p.qos, p.createDeviceOnlineHandler()); token.Wait() && token.Error() != nil {
		panic(0)
	}
	//设备离线
	if token := p.Client.Subscribe(topics.DEVICE_OFFLINE, p.qos, p.createDeviceOfflineHandler()); token.Wait() && token.Error() != nil {
		panic(0)
	}
	if token := p.Client.Subscribe(topics.PLATFORM_PROPERTIES_UP, p.qos, p.createPropertiesReportHandler()); token.Wait() && token.Error() != nil {
		//log.Warningf("设备 %s 订阅平台发送命令topic %s 失败", device.Id, topic)
		panic(0)
	}
	if token := p.Client.Subscribe(topics.BATCH_REPORT_SUB_DEVICE, p.qos, p.createBatchSubDeviceReportHandler()); token.Wait() && token.Error() != nil {
		panic(0)
	}
	if token := p.Client.Subscribe(topics.DEVICE_EVENT_TO_PLATFORM_TOPIC, p.qos, p.handleDeviceToPlatformEvent()); token.Wait() && token.Error() != nil {
		log.Errorf("订阅设备事件上报失败: %s", token.Error().Error())
		panic(0)
	}

	// 订阅设备消息的topic
	if token := p.Client.Subscribe(iot.FormatTopic(topics.DEVICE_MESSAGE_UP, p.Id), p.qos, p.createMessageMqttHandler()); token.Wait() && token.Error() != nil {
		log.Errorf("平台订阅设备消息失败：%s", token.Error().Error())
		panic(0)
	}
	//订阅设备命令响应topic
	if token := p.Client.Subscribe(iot.FormatTopic(topics.COMMAND_RESPONSE, p.Id), p.qos, p.createCommandMqttHandler()); token.Wait() && token.Error() != nil {
		log.Errorf("平台订阅命令响应失败：%s", token.Error().Error())
		panic(0)
	}
}

func (p *basePlatform) addMessageHandler(handler MessageHandler) {
	p.messageHandlers = append(p.messageHandlers, handler)
}
func (p *basePlatform) addCommandHandler(handler CommandHandler) {
	p.commandHandlers = append(p.commandHandlers, handler)
}
func (p *basePlatform) createPropertiesReportHandler() func(client mqtt.Client, message mqtt.Message) {
	propertiesReportHandler := func(client mqtt.Client, message mqtt.Message) {
		go func() {
			request := &iot.DeviceService{}
			if err := json.Unmarshal(message.Payload(), request); err != nil {
				log.Errorf("解析设备上报属性失败:%s ERROR:%v", message, err)
			}
			request.DeviceId = iot.GetTopicDeviceId(message.Topic())
			p.propertiesReportHandler(*request)
		}()
	}
	return propertiesReportHandler
}

func (p *basePlatform) createBatchSubDeviceReportHandler() func(client mqtt.Client, message mqtt.Message) {
	handler := func(client mqtt.Client, message mqtt.Message) {
		go func() {
			request := &iot.DevicesService{}
			if err := json.Unmarshal(message.Payload(), request); err != nil {
				log.Warningf("解析网关批量上报子设备属性失败,\r\n错误:%e\r\n报文:%s", err, message)
			}
			p.batchReportHandler(*request)
		}()
	}
	return handler
}

func (p *basePlatform) createMessageMqttHandler() func(client mqtt.Client, message mqtt.Message) {
	messageHandler := func(client mqtt.Client, message mqtt.Message) {
		go func() {
			msg := &iot.Message{}
			if json.Unmarshal(message.Payload(), msg) != nil {
				log.Warningf("unmarshal device message failed,device id = %s,message = %s", p.Id, message)
			}

			for _, handler := range p.messageHandlers {
				handler(*msg)
			}
		}()
	}

	return messageHandler
}

func (p *basePlatform) createDeviceOnlineHandler() func(client mqtt.Client, message mqtt.Message) {
	handler := func(client mqtt.Client, message mqtt.Message) {
		device := &iot.DeviceOnlineInfo{}
		payload := message.Payload()
		if err := json.Unmarshal(payload, device); err != nil {
			log.Errorf("解析设备上线失败:%s", err)
		} else {
			p.deviceOnlineHandler(device)
		}
	}
	return handler
}
func (p *basePlatform) createDeviceOfflineHandler() func(client mqtt.Client, message mqtt.Message) {
	handler := func(client mqtt.Client, message mqtt.Message) {
		device := &iot.DeviceOfflineInfo{}
		if err := json.Unmarshal(message.Payload(), device); err != nil {
			log.Errorf("解析设备上线失败:%s", err)
		} else {
			p.deviceOfflineHandler(device)
		}
	}
	return handler
}

// todo 2022/4/15添加
func (p *basePlatform) createCommandMqttHandler() func(client mqtt.Client, message mqtt.Message) {
	commandHandler := func(client mqtt.Client, message mqtt.Message) {
		go func() {
			request := &iot.CommandResponse{}
			if json.Unmarshal(message.Payload(), request) != nil {
				log.Warningf("解析命令响应失败", message)
			}
			for _, handler := range p.commandHandlers {
				handler(*request)
			}
		}()

	}

	return commandHandler
}

func (p *basePlatform) handleDeviceToPlatformEvent() func(client mqtt.Client, message mqtt.Message) {
	handler := func(client mqtt.Client, message mqtt.Message) {
		event := &iot.EventData{}
		err := json.Unmarshal(message.Payload(), event)
		if err != nil {
			log.Errorf(err.Error())
			return
		}
		for _, service := range event.Services {
			switch service.EventType {
			case event_type.VERSION_REPORT:
				version := &iot.VersionInfo{}
				if json.Unmarshal([]byte(iot.Interface2JsonString(service.Paras)), version) != nil {
					continue
				}
				p.versionReportHandler(event.DeviceId, *version)
				break
			case event_type.UPGRADE_PROGRESS_REPORT:
				process := &iot.UpgradeProgress{}
				if json.Unmarshal([]byte(iot.Interface2JsonString(service.Paras)), process) != nil {
					continue
				}
				p.upgradeProcessReportHandler(event.DeviceId, *process)
				break
			}
		}
	}
	return handler
}

func assembleClientId(device *basePlatform) string {
	segments := make([]string, 4)
	segments[0] = device.Id
	segments[1] = "0"
	segments[2] = "0"
	segments[3] = iot.TimestampString()
	return strings.Join(segments, "_")
}
