package platform

import (
	"github.com/golang/glog"
	log "github.com/sirupsen/logrus"
	"sailing.cn/v2/convert"
	"sailing.cn/v2/iot"
	mqtt2 "sailing.cn/v2/mqtt"
	"sailing.cn/v2/utils/timestamp"
)

var instance Platform

type platform struct {
	base base
}

// Client MQTT客户端
type Client struct {
}

func (platform *platform) Init() bool {
	result := platform.base.Init()
	platform.base.subscribeDefaultTopics()
	return result
}
func (platform *platform) DisConnect() {
	platform.base.DisConnect()
}

func (platform *platform) IsConnected() bool {
	return platform.base.IsConnected()
}

func Get(cnf *Config) Platform {
	if instance == nil || instance.IsConnected() == false {
		instance = CreatePlatform(cnf.Id, cnf.Password, cnf.Servers)
	}
	return instance
}

func CreatePlatform(user, password, servers string) Platform {
	config := Config{
		Id:                 user,
		Password:           password,
		Servers:            servers,
		BatchSubDeviceSize: 0,
	}
	return CreatePlatformWithConfig(config)
}
func CreatePlatformWithConfig(config Config) Platform {
	client := base{}
	client.Id = config.Id
	client.Password = config.Password
	client.Servers = config.Servers
	client.messageHandlers = []MessageHandler{}
	client.commandHandlers = []CommandHandler{}

	client.fileUrls = map[string]string{}

	client.qos = config.Qos
	client.batchSubDeviceSize = 100

	result := &platform{
		base: client,
	}
	return result
}

type Platform interface {
	BasePlatform
	EventDown(event iot.EventData) bool
	SendMessage(message iot.Message) bool
	SendCommand(command iot.Command) bool
	//QueryVersion 查询设备版本
	QueryVersion(device string) bool
	SendUpgradeMessage(device string, upgradeType iot.EventType, data iot.UpgradePara) bool
	SetPropertiesReportHandler(handler PropertiesReportHandler)
	SetGatewayBatchReportSubDeviceHandler(handler GatewayBatchReportSubDeviceHandler)
	SetDeviceInfoReportHandler(handler DeviceInfoReportHandler)
	SetDeviceOnlineHandler(handler DeviceOnlineHandler)
	SetDeviceOfflineHandler(handler DeviceOfflineHandler)
	AddMessageHandler(handler MessageHandler)
	AddCommandHandler(handler CommandHandler)
	SetVersionReportHandler(handler VersionReportHandler)
	SetUpgradeProcessReportHandler(handler UpgradeProcessReportHandler)
}

func (platform *platform) AddMessageHandler(handler MessageHandler) {
	platform.base.addMessageHandler(handler)
}
func (platform *platform) AddCommandHandler(handler CommandHandler) {
	platform.base.addCommandHandler(handler)
}
func (platform *platform) SetDeviceInfoReportHandler(handler DeviceInfoReportHandler) {
	platform.base.deviceInfoReportHandler = handler
}

func (platform *platform) SetDeviceInfoReport(handler DeviceInfoReportHandler) {
	//TODO implement me
	panic("implement me")
}

// SetDeviceOnlineHandler 设备上线处理事件
func (platform *platform) SetDeviceOnlineHandler(handler DeviceOnlineHandler) {
	platform.base.deviceOnlineHandler = handler
}

// SetDeviceOfflineHandler 设备上线处理事件
func (platform *platform) SetDeviceOfflineHandler(handler DeviceOfflineHandler) {
	platform.base.deviceOfflineHandler = handler
}

func (platform *platform) EventDown(event iot.EventData) bool {
	if token := platform.base.Client.Publish(mqtt2.FormatTopic(iot.PlatformEventToDeviceTopic, event.DeviceId), platform.base.qos, false, convert.ToString(event)); token.Wait() && token.Error() != nil {
		log.Errorf("平台下发事件失败,device:%s", event.DeviceId)
		return false
	}
	return true
}

func (platform *platform) SetGatewayBatchReportSubDeviceHandler(handler GatewayBatchReportSubDeviceHandler) {
	platform.base.batchReportHandler = handler
}

func (platform *platform) SetPropertiesReportHandler(handler PropertiesReportHandler) {
	platform.base.propertiesReportHandler = handler
}

func (platform *platform) SetVersionReportHandler(handler VersionReportHandler) {
	platform.base.versionReportHandler = handler
}

func (platform *platform) QueryVersion(device string) bool {
	event := iot.EventData{
		DeviceId: device,
		Services: []iot.ServiceEvent{{
			ServiceId: iot.OTA,
			EventType: iot.VERSION_QUERY,
			EventTime: timestamp.Timestamp(),
			Paras:     nil,
		}},
	}
	return platform.EventDown(event)
}

func (platform *platform) SendUpgradeMessage(device string, upgradeType iot.EventType, data iot.UpgradePara) bool {
	event := iot.EventData{
		DeviceId: device,
		Services: []iot.ServiceEvent{{
			ServiceId: iot.OTA,
			EventType: upgradeType,
			EventTime: timestamp.Timestamp(),
			Paras:     data,
		}},
	}
	return platform.EventDown(event)
}

func (platform *platform) SetUpgradeProcessReportHandler(handler UpgradeProcessReportHandler) {
	platform.base.upgradeProcessReportHandler = handler
}

func (platform *platform) SendMessage(message iot.Message) bool {
	msg := convert.ToString(message)
	if token := platform.base.Client.Publish(mqtt2.FormatTopic(iot.MessageDownTopic, message.DeviceId), platform.base.qos, false, msg); token.Wait() && token.Error() != nil {
		glog.Errorf("发送设备[%s]消息失败", message.DeviceId)
		return false
	}
	return true
}
func (platform *platform) SendCommand(command iot.Command) bool {
	cmd := convert.ToString(command)
	if token := platform.base.Client.Publish(mqtt2.FormatTopicWithRequest(iot.CommandDownTopic, command.DeviceId, command.CommandId)+command.CommandId, platform.base.qos, false, cmd); token.Wait() && token.Error() != nil {
		glog.Errorf("发送设备[%s]消息失败", command.CommandId)
		return false
	}
	return true
}

// SetConnectLostHandler 设置断开连接处理器
func (platform *platform) SetConnectLostHandler(handler ConnectLostHandler) {
	platform.base.connectLostHandler = handler
}
