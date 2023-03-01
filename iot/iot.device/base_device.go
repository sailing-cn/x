package iot_device

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
	log "github.com/sirupsen/logrus"
	"os"
	"sailing.cn/iot"
	"sailing.cn/iot/event_type"
	"sailing.cn/iot/iot.device/topics"
	"strings"
	"time"
)

type DeviceConfig struct {
	Id                 string
	Password           string
	Servers            string
	Qos                byte
	BatchSubDeviceSize int
	AuthType           uint8
	ServerCaPath       string
	CertFilePath       string
	CertKeyFilePath    string
}

const (
	AUTH_TYPE_PASSWORD uint8 = 0
	AUTH_TYPE_X509     uint8 = 1
)

type BaseDevice interface {
	Init() bool
	DisConnect()
	IsConnected() bool

	AddMessageHandler(handler MessageHandler)
	AddCommandHandler(handler CommandHandler)
	AddPropertiesSetHandler(handler DevicePropertiesSetHandler)
	SetPropertyQueryHandler(handler DevicePropertyQueryHandler)
	SetVersionReporter(handler VersionReporter)
	SetDeviceUpgradeHandler(handler DeviceUpgradeHandler)

	SetDeviceStatusLogCollector(collector DeviceStatusLogCollector)
	SetDevicePropertyLogCollector(collector DevicePropertyLogCollector)
	SetDeviceMessageLogCollector(collector DeviceMessageLogCollector)
	SetDeviceCommandLogCollector(collector DeviceCommandLogCollector)
}

type baseIotDevice struct {
	Id              string // 设备Id，平台又称为deviceId
	Password        string // 设备密码
	AuthType        uint8  // 鉴权类型，0：密码认证；1：x.509证书认证
	ServerCaPath    string // 平台CA证书
	CertFilePath    string // 设备证书路径
	CertKeyFilePath string // 设备证书key路径
	Servers         string
	Client          mqtt.Client

	commandHandlers                []CommandHandler
	messageHandlers                []MessageHandler
	propertiesSetHandlers          []DevicePropertiesSetHandler
	propertyQueryHandler           DevicePropertyQueryHandler
	propertiesQueryResponseHandler DevicePropertyQueryResponseHandler
	subDevicesAddHandler           SubDevicesAddHandler
	subDevicesDeleteHandler        SubDevicesDeleteHandler
	versionReporter                VersionReporter
	deviceUpgradeHandler           DeviceUpgradeHandler
	//lcc                            *LogCollectionConfig
	deviceStatusLogCollector   DeviceStatusLogCollector
	devicePropertyLogCollector DevicePropertyLogCollector
	deviceMessageLogCollector  DeviceMessageLogCollector
	deviceCommandLogCollector  DeviceCommandLogCollector
	fileUrls                   map[string]string
	qos                        byte
	batchSubDeviceSize         int
}

func (device *baseIotDevice) DisConnect() {
	if device.Client != nil {
		device.Client.Disconnect(0)
	}
}
func (device *baseIotDevice) IsConnected() bool {
	if device.Client != nil {
		return device.Client.IsConnectionOpen()
	}

	return false
}

func (device *baseIotDevice) Init() bool {

	options := mqtt.NewClientOptions()
	options.AddBroker(device.Servers)
	options.SetClientID(assembleClientId(device))
	options.SetUsername(device.Id)
	options.SetPassword(iot.HmacSha256(device.Password, iot.TimestampString()))
	options.SetKeepAlive(250 * time.Second)
	options.SetAutoReconnect(true)
	options.SetConnectRetry(true)
	options.SetConnectTimeout(2 * time.Second)

	if strings.Contains(device.Servers, "tls") || strings.Contains(device.Servers, "ssl") {
		glog.Infof("server support tls connection")

		// 设备使用x.509 证书认证
		if device.AuthType == AUTH_TYPE_X509 {
			if len(device.ServerCaPath) == 0 || len(device.CertFilePath) == 0 || len(device.CertKeyFilePath) == 0 {
				glog.Error("device use x.509 auth but not set cert")
				panic("not set cert")
			}

			ca, err := os.ReadFile(device.ServerCaPath)
			if err != nil {
				glog.Error("load server ca failed\n")
				panic(err)
			}
			serverCaPool := x509.NewCertPool()
			serverCaPool.AppendCertsFromPEM(ca)

			deviceCert, err := tls.LoadX509KeyPair(device.CertFilePath, device.CertKeyFilePath)
			if err != nil {
				glog.Error("load device cert failed")
				panic("load device cert failed")
			}
			var clientCerts []tls.Certificate
			clientCerts = append(clientCerts, deviceCert)

			cipherSuites := []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			}
			tlsConfig := &tls.Config{
				RootCAs:            serverCaPool,
				Certificates:       clientCerts,
				InsecureSkipVerify: true,
				MaxVersion:         tls.VersionTLS12,
				MinVersion:         tls.VersionTLS12,
				CipherSuites:       cipherSuites,
			}
			options.SetTLSConfig(tlsConfig)

		}

		if device.AuthType == 0 {
			options.SetTLSConfig(&tls.Config{
				InsecureSkipVerify: true,
			})
		}

	} else {
		options.SetTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
		})
	}
	device.Client = mqtt.NewClient(options)
	token := device.Client.Connect()
	token.Wait()
	if token.Error() != nil {
		glog.Warningf("device %s init failed,error = %v", device.Id, token.Error())
		return false
	}
	device.subscribeDefaultTopics()
	go logFlush()
	return true
}

func (device *baseIotDevice) AddMessageHandler(handler MessageHandler) {
	if handler == nil {
		return
	}
	device.messageHandlers = append(device.messageHandlers, handler)
}
func (device *baseIotDevice) AddCommandHandler(handler CommandHandler) {
	if handler == nil {
		return
	}

	device.commandHandlers = append(device.commandHandlers, handler)
}
func (device *baseIotDevice) AddPropertiesSetHandler(handler DevicePropertiesSetHandler) {
	if handler == nil {
		return
	}
	device.propertiesSetHandlers = append(device.propertiesSetHandlers, handler)
}
func (device *baseIotDevice) SetVersionReporter(handler VersionReporter) {
	device.versionReporter = handler
}

func (device *baseIotDevice) SetDeviceUpgradeHandler(handler DeviceUpgradeHandler) {
	device.deviceUpgradeHandler = handler
}

func (device *baseIotDevice) SetPropertyQueryHandler(handler DevicePropertyQueryHandler) {
	device.propertyQueryHandler = handler
}

func (device *baseIotDevice) SetDeviceStatusLogCollector(collector DeviceStatusLogCollector) {
	device.deviceStatusLogCollector = collector
}

func (device *baseIotDevice) SetDevicePropertyLogCollector(collector DevicePropertyLogCollector) {
	device.devicePropertyLogCollector = collector
}

func (device *baseIotDevice) SetDeviceMessageLogCollector(collector DeviceMessageLogCollector) {
	device.deviceMessageLogCollector = collector
}

func (device *baseIotDevice) SetDeviceCommandLogCollector(collector DeviceCommandLogCollector) {
	device.deviceCommandLogCollector = collector
}

func (device *baseIotDevice) createCommandMqttHandler() func(client mqtt.Client, message mqtt.Message) {
	commandHandler := func(client mqtt.Client, message mqtt.Message) {
		go func() {
			command := &iot.Command{}
			if json.Unmarshal(message.Payload(), command) != nil {
				glog.Warningf("unmarshal platform command failed,device id = %s，message = %s", device.Id, message)
			}

			handleFlag := true
			for _, handler := range device.commandHandlers {
				handleFlag = handleFlag && handler(*command)
			}
			var res string
			if handleFlag {
				glog.Infof("device %s handle command success", device.Id)
				res = iot.Interface2JsonString(iot.CommandResponse{
					ResultCode: 0,
					RequestId:  command.CommandId,
					Paras:      command.Paras,
				})
			} else {
				glog.Warningf("device %s handle command failed", device.Id)
				res = iot.Interface2JsonString(iot.CommandResponse{
					ResultCode: 1,
					//todo 下面只是显示作用
					RequestId: command.CommandId,
					Paras:     command.Paras,
				})
			}

			if token := device.Client.Publish(iot.FormatTopic(topics.CommandResponseTopic, device.Id)+iot.GetTopicRequestId(message.Topic()), 1, false, res); token.Wait() && token.Error() != nil {
				glog.Infof("device %s send command response failed", device.Id)
			}
		}()

	}

	return commandHandler
}

func (device *baseIotDevice) createPropertiesSetMqttHandler() func(client mqtt.Client, message mqtt.Message) {
	propertiesSetHandler := func(client mqtt.Client, message mqtt.Message) {
		go func() {
			propertiesSetRequest := &iot.DevicePropertyDownRequest{}
			if json.Unmarshal(message.Payload(), propertiesSetRequest) != nil {
				glog.Warningf("unmarshal platform properties set request failed,device id = %s，message = %s", device.Id, message)
			}

			handleFlag := true
			for _, handler := range device.propertiesSetHandlers {
				handleFlag = handleFlag && handler(*propertiesSetRequest)
			}

			var res string
			response := struct {
				ResultCode byte   `json:"result_code"`
				ResultDesc string `json:"result_desc"`
			}{}
			if handleFlag {
				response.ResultCode = 0
				response.ResultDesc = "Set property success."
				res = iot.Interface2JsonString(response)
			} else {
				response.ResultCode = 1
				response.ResultDesc = "Set properties failed."
				res = iot.Interface2JsonString(response)
			}
			if token := device.Client.Publish(iot.FormatTopic(topics.PropertiesSetResponseTopic, device.Id)+iot.GetTopicRequestId(message.Topic()), device.qos, false, res); token.Wait() && token.Error() != nil {
				glog.Warningf("unmarshal platform properties set request failed,device id = %s，message = %s", device.Id, message)
			}
		}()
	}

	return propertiesSetHandler
}

func (device *baseIotDevice) createMessageMqttHandler() func(client mqtt.Client, message mqtt.Message) {
	messageHandler := func(client mqtt.Client, message mqtt.Message) {
		go func() {
			msg := &iot.Message{}
			if json.Unmarshal(message.Payload(), msg) != nil {
				glog.Warningf("unmarshal device message failed,device id = %s,message = %s", device.Id, message)
			}

			for _, handler := range device.messageHandlers {
				handler(*msg)
			}
		}()
	}

	return messageHandler
}

func (device *baseIotDevice) createPropertiesQueryMqttHandler() func(client mqtt.Client, message mqtt.Message) {
	propertiesQueryHandler := func(client mqtt.Client, message mqtt.Message) {
		go func() {
			propertiesQueryRequest := &iot.DevicePropertyQueryRequest{}
			if json.Unmarshal(message.Payload(), propertiesQueryRequest) != nil {
				glog.Warningf("device %s unmarshal properties query request failed %s", device.Id, message)
			}

			queryResult := device.propertyQueryHandler(*propertiesQueryRequest)
			responseToPlatform := iot.Interface2JsonString(queryResult)
			if token := device.Client.Publish(iot.FormatTopic(topics.PropertiesQueryResponseTopic, device.Id)+iot.GetTopicRequestId(message.Topic()), device.qos, false, responseToPlatform); token.Wait() && token.Error() != nil {
				glog.Warningf("device %s send properties query response failed.", device.Id)
			}
		}()
	}

	return propertiesQueryHandler
}

func (device *baseIotDevice) createPropertiesQueryResponseMqttHandler() func(client mqtt.Client, message mqtt.Message) {
	propertiesQueryResponseHandler := func(client mqtt.Client, message mqtt.Message) {
		propertiesQueryResponse := &iot.DevicePropertyQueryResponse{}
		if json.Unmarshal(message.Payload(), propertiesQueryResponse) != nil {
			glog.Warningf("device %s unmarshal property response failed,message %s", device.Id, iot.Interface2JsonString(message))
		}
		device.propertiesQueryResponseHandler(*propertiesQueryResponse)
	}

	return propertiesQueryResponseHandler
}

func (device *baseIotDevice) subscribeDefaultTopics() {
	// 订阅平台命令下发topic
	topic := iot.FormatTopic(topics.CommandDownTopic, device.Id)
	if token := device.Client.Subscribe(topic, device.qos, device.createCommandMqttHandler()); token.Wait() && token.Error() != nil {
		glog.Warningf("device %s subscribe platform send command topic %s failed", device.Id, topic)
		panic(0)
	}

	// 订阅平台消息下发的topic
	topic = iot.FormatTopic(topics.MessageDownTopic, device.Id)
	if token := device.Client.Subscribe(topic, device.qos, device.createMessageMqttHandler()); token.Wait() && token.Error() != nil {
		glog.Warningf("device % subscribe platform send message topic %s failed.", device.Id, topic)
		panic(0)
	}

	// 订阅平台设置设备属性的topic
	topic = iot.FormatTopic(topics.PropertiesSetRequestTopic, device.Id)
	if token := device.Client.Subscribe(topic, device.qos, device.createPropertiesSetMqttHandler()); token.Wait() && token.Error() != nil {
		glog.Warningf("device %s subscribe platform set properties topic %s failed", device.Id, topic)
		panic(0)
	}

	// 订阅平台查询设备属性的topic
	topic = iot.FormatTopic(topics.PropertiesQueryRequestTopic, device.Id)
	if token := device.Client.Subscribe(topic, device.qos, device.createPropertiesQueryMqttHandler()); token.Wait() && token.Error() != nil {
		glog.Warningf("device %s subscriber platform query device properties topic failed %s", device.Id, topic)
		panic(0)
	}

	// 订阅查询设备影子响应的topic
	topic = iot.FormatTopic(topics.DeviceShadowQueryResponseTopic, device.Id)
	if token := device.Client.Subscribe(topic, device.qos, device.createPropertiesQueryResponseMqttHandler()); token.Wait() && token.Error() != nil {
		glog.Warningf("device %s subscribe query device shadow topic %s failed", device.Id, topic)
		panic(0)
	}

	// 订阅平台下发到设备的事件
	topic = iot.FormatTopic(topics.PlatformEventToDeviceTopic, device.Id)
	if token := device.Client.Subscribe(topic, device.qos, device.handlePlatformToDeviceData()); token.Wait() && token.Error() != nil {
		glog.Warningf("device %s subscribe query device shadow topic %s failed", device.Id, topic)
		panic(0)
	}
}

// 平台向设备下发的事件callback
func (device *baseIotDevice) handlePlatformToDeviceData() func(client mqtt.Client, message mqtt.Message) {
	handler := func(client mqtt.Client, message mqtt.Message) {
		data := &iot.EventData{}
		err := json.Unmarshal(message.Payload(), data)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, entry := range data.Services {
			eventType := entry.EventType
			switch eventType {
			case "add_sub_device_notify":
				// 子设备添加
				subDeviceInfo := &iot.AddSubDeviceParas{}
				if json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), subDeviceInfo) != nil {
					continue
				}
				device.subDevicesAddHandler(*subDeviceInfo)
			case "delete_sub_device_notify":
				subDeviceInfo := &iot.DeleteSubDeviceParas{}
				if json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), subDeviceInfo) != nil {
					continue
				}
				device.subDevicesDeleteHandler(*subDeviceInfo)

			case "get_upload_url_response":
				//获取文件上传URL
				fileResponse := &iot.FileResponseServiceEventParas{}
				if json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), fileResponse) != nil {
					continue
				}
				device.fileUrls[fileResponse.ObjectName+topics.FileActionUpload] = fileResponse.Url
			case "get_download_url_response":
				fileResponse := &iot.FileResponseServiceEventParas{}
				if json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), fileResponse) != nil {
					continue
				}
				device.fileUrls[fileResponse.ObjectName+topics.FileActionDownload] = fileResponse.Url
			case event_type.VERSION_QUERY:
				// 查询软固件版本
				log.Infof("查询设备版本:%s", time.Now().String())
				device.reportVersion()

			case event_type.FIRMWARE_UPGRADE:
				upgradeInfo := &iot.UpgradePara{}
				if json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), upgradeInfo) != nil {
					continue
				}
				device.upgradeDevice(event_type.FIRMWARE_UPGRADE, upgradeInfo)

			case event_type.SOFTWARE_UPGRADE:
				upgradeInfo := &iot.UpgradePara{}
				if json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), upgradeInfo) != nil {
					continue
				}
				device.upgradeDevice(event_type.SOFTWARE_UPGRADE, upgradeInfo)
			case event_type.CONFIG_UPGRADE:
				upgradeInfo := &iot.UpgradePara{}
				err := json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), upgradeInfo)
				if err != nil {
					log.Errorf("解析平台事件参数失败:%s", err.Error())
					continue
				}
				device.upgradeDevice(event_type.CONFIG_UPGRADE, upgradeInfo)
			case event_type.MESSAGE_RECEUVED:
				messageResponse := &iot.Message{}
				err := json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), messageResponse)
				if err != nil {
					log.Errorf("消息接受失败:%s", err.Error())
					continue
				}
				device.updateMessageStatus("message_received", messageResponse)
			case "log_config":
				// 平台下发日志收集通知
				//fmt.Println("platform send log collect command")
				//logConfig := &LogCollectionConfig{}
				//if json.Unmarshal([]byte(iot.Interface2JsonString(entry.Paras)), logConfig) != nil {
				//	continue
				//}
				//
				//lcc := &LogCollectionConfig{
				//	logCollectSwitch: logConfig.logCollectSwitch,
				//	endTime:          logConfig.endTime,
				//}
				//device.lcc = lcc
				//device.reportLogsWorker()
			}
		}

	}

	return handler
}

func (device *baseIotDevice) reportVersion() {
	version := device.versionReporter()
	dataEntry := iot.ServiceEvent{
		ServiceId: iot.OTA,
		EventType: event_type.VERSION_REPORT,
		EventTime: iot.GetEventTimeStamp(),
		Paras:     version,
	}
	data := iot.EventData{
		DeviceId: device.Id,
		Services: []iot.ServiceEvent{dataEntry},
	}

	device.Client.Publish(iot.FormatTopic(topics.DeviceToPlatformTopic, device.Id), device.qos, false, iot.Interface2JsonString(data))
}
func (device *baseIotDevice) updateMessageStatus(upgradeType iot.EventType, upgradeInfo *iot.Message) {

}
func (device *baseIotDevice) upgradeDevice(upgradeType iot.EventType, upgradeInfo *iot.UpgradePara) {
	progress := device.deviceUpgradeHandler(upgradeType, *upgradeInfo)
	dataEntry := iot.ServiceEvent{
		ServiceId: iot.OTA,
		EventType: event_type.UPGRADE_PROGRESS_REPORT,
		EventTime: iot.GetEventTimeStamp(),
		Paras:     progress,
	}
	data := iot.EventData{
		DeviceId: device.Id,
		Services: []iot.ServiceEvent{dataEntry},
	}

	if token := device.Client.Publish(iot.FormatTopic(topics.DeviceToPlatformTopic, device.Id), device.qos, false, iot.Interface2JsonString(data)); token.Wait() && token.Error() != nil {
		glog.Errorf("device %s upgrade failed,type %d", device.Id, upgradeType)
	}
}

func assembleClientId(device *baseIotDevice) string {
	segments := make([]string, 4)
	segments[0] = device.Id
	segments[1] = "0"
	segments[2] = "0"
	segments[3] = iot.TimestampString()

	return strings.Join(segments, "_")
}
func logFlush() {
	ticker := time.Tick(5 * time.Second)
	for {
		select {
		case <-ticker:
			glog.Flush()
		}
	}
}
