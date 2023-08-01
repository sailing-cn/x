package platform

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"sailing.cn/v2/convert"
	"sailing.cn/v2/encrypt"
	"sailing.cn/v2/iot"
	mqtt2 "sailing.cn/v2/mqtt"
	"sailing.cn/v2/utils/timestamp"
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

type base struct {
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

func (p *base) Init() bool {
	options := mqtt.NewClientOptions()
	options.AddBroker(p.Servers)
	options.SetClientID(mqtt2.AssembleClientId(p.Id))
	options.SetUsername(p.Id)
	options.SetPassword(encrypt.HmacSha256(p.Password, timestamp.TimestampString()[:10]))
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
func (p *base) DisConnect() {
	if p.Client != nil {
		p.Client.Disconnect(0)
	}
}

func (p *base) IsConnected() bool {
	if p.Client != nil {
		return p.Client.IsConnectionOpen()
	}
	return false
}

func (p *base) SendMessage(message iot.Message) bool {
	panic("implement me" + message.Name)
}

func (p *base) subscribeDefaultTopics() {
	//设备上线
	if token := p.Client.Subscribe(iot.DEVICE_ONLINE, p.qos, p.createDeviceOnlineHandler()); token.Wait() && token.Error() != nil {
		panic(0)
	}
	//设备离线
	if token := p.Client.Subscribe(iot.DEVICE_OFFLINE, p.qos, p.createDeviceOfflineHandler()); token.Wait() && token.Error() != nil {
		panic(0)
	}
	if token := p.Client.Subscribe(iot.PLATFORM_PROPERTIES_UP, p.qos, p.createPropertiesReportHandler()); token.Wait() && token.Error() != nil {
		panic(0)
	}
	if token := p.Client.Subscribe(iot.BATCH_REPORT_SUB_DEVICE, p.qos, p.createBatchSubDeviceReportHandler()); token.Wait() && token.Error() != nil {
		panic(0)
	}
	if token := p.Client.Subscribe(iot.DEVICE_EVENT_TO_PLATFORM_TOPIC, p.qos, p.handleDeviceToPlatformEvent()); token.Wait() && token.Error() != nil {
		log.Errorf("订阅设备事件上报失败: %s", token.Error().Error())
		panic(0)
	}

	// 订阅设备消息的topic
	if token := p.Client.Subscribe(mqtt2.FormatTopic(iot.DEVICE_MESSAGE_UP, p.Id), p.qos, p.createMessageMqttHandler()); token.Wait() && token.Error() != nil {
		log.Errorf("平台订阅设备消息失败：%s", token.Error().Error())
		panic(0)
	}
	//订阅设备命令响应topic
	if token := p.Client.Subscribe(mqtt2.FormatTopic(iot.COMMAND_RESPONSE, p.Id), p.qos, p.createCommandMqttHandler()); token.Wait() && token.Error() != nil {
		log.Errorf("平台订阅命令响应失败：%s", token.Error().Error())
		panic(0)
	}
}

func (p *base) addMessageHandler(handler MessageHandler) {
	p.messageHandlers = append(p.messageHandlers, handler)
}
func (p *base) addCommandHandler(handler CommandHandler) {
	p.commandHandlers = append(p.commandHandlers, handler)
}
func (p *base) createPropertiesReportHandler() func(client mqtt.Client, message mqtt.Message) {
	propertiesReportHandler := func(client mqtt.Client, message mqtt.Message) {
		go func() {
			request := &iot.DeviceService{}
			if err := json.Unmarshal(message.Payload(), request); err != nil {
				log.Errorf("解析设备上报属性失败:%s ERROR:%v", message, err)
			}
			request.DeviceId = mqtt2.GetTopicDeviceId(message.Topic())
			p.propertiesReportHandler(*request)
		}()
	}
	return propertiesReportHandler
}

func (p *base) createBatchSubDeviceReportHandler() func(client mqtt.Client, message mqtt.Message) {
	handler := func(client mqtt.Client, message mqtt.Message) {
		go func() {
			request := &iot.DevicesService{}
			if err := json.Unmarshal(message.Payload(), request); err != nil {
				log.Warningf("解析网关批量上报子设备属性失败,\r\n错误:%s\r\n报文:%s", err.Error(), message)
			}
			p.batchReportHandler(*request)
		}()
	}
	return handler
}

func (p *base) createMessageMqttHandler() func(client mqtt.Client, message mqtt.Message) {
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

func (p *base) createDeviceOnlineHandler() func(client mqtt.Client, message mqtt.Message) {
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
func (p *base) createDeviceOfflineHandler() func(client mqtt.Client, message mqtt.Message) {
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
func (p *base) createCommandMqttHandler() func(client mqtt.Client, message mqtt.Message) {
	commandHandler := func(client mqtt.Client, message mqtt.Message) {
		go func() {
			request := &iot.CommandResponse{}
			if json.Unmarshal(message.Payload(), request) != nil {
				log.Warningf("解析命令响应失败:%s", message)
			}
			for _, handler := range p.commandHandlers {
				handler(*request)
			}
		}()

	}

	return commandHandler
}

func (p *base) handleDeviceToPlatformEvent() func(client mqtt.Client, message mqtt.Message) {
	handler := func(client mqtt.Client, message mqtt.Message) {
		event := &iot.EventData{}
		err := json.Unmarshal(message.Payload(), event)
		if err != nil {
			log.Errorf(err.Error())
			return
		}
		for _, service := range event.Services {
			switch service.EventType {
			case iot.VERSION_REPORT:
				version := &iot.VersionInfo{}
				if json.Unmarshal([]byte(convert.ToString(service.Paras)), version) != nil {
					continue
				}
				p.versionReportHandler(event.DeviceId, *version)
				break
			case iot.UPGRADE_PROGRESS_REPORT:
				process := &iot.UpgradeProgress{}
				if json.Unmarshal([]byte(convert.ToString(service.Paras)), process) != nil {
					continue
				}
				p.upgradeProcessReportHandler(event.DeviceId, *process)
				break
			}
		}
	}
	return handler
}
