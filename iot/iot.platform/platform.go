package iot_platform

import (
	"github.com/golang/glog"
	log "github.com/sirupsen/logrus"
	"sailing.cn/iot"
	"sailing.cn/iot/event_type"
	"sailing.cn/iot/iot.device/topics"
	"sailing.cn/utils"
)

var instance Platform

type iotPlatform struct {
	base basePlatform
}

// Client MQTT客户端
type Client struct {
}

func (platform *iotPlatform) Init() bool {
	result := platform.base.Init()
	platform.base.subscribeDefaultTopics()
	return result
}
func (platform *iotPlatform) DisConnect() {
	platform.base.DisConnect()
}

func (platform *iotPlatform) IsConnected() bool {
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
	client := basePlatform{}
	client.Id = config.Id
	client.Password = config.Password
	client.Servers = config.Servers
	client.messageHandlers = []MessageHandler{}
	client.commandHandlers = []CommandHandler{}

	client.fileUrls = map[string]string{}

	client.qos = config.Qos
	client.batchSubDeviceSize = 100

	result := &iotPlatform{
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

func (platform *iotPlatform) AddMessageHandler(handler MessageHandler) {
	platform.base.addMessageHandler(handler)
}
func (platform *iotPlatform) AddCommandHandler(handler CommandHandler) {
	platform.base.addCommandHandler(handler)
}
func (platform *iotPlatform) SetDeviceInfoReportHandler(handler DeviceInfoReportHandler) {
	platform.base.deviceInfoReportHandler = handler
}

func (platform *iotPlatform) SetDeviceInfoReport(handler DeviceInfoReportHandler) {
	//TODO implement me
	panic("implement me")
}

// SetDeviceOnlineHandler 设备上线处理事件
func (platform *iotPlatform) SetDeviceOnlineHandler(handler DeviceOnlineHandler) {
	platform.base.deviceOnlineHandler = handler
}

// SetDeviceOfflineHandler 设备上线处理事件
func (platform *iotPlatform) SetDeviceOfflineHandler(handler DeviceOfflineHandler) {
	platform.base.deviceOfflineHandler = handler
}

func (platform *iotPlatform) EventDown(event iot.EventData) bool {
	if token := platform.base.Client.Publish(iot.FormatTopic(topics.PlatformEventToDeviceTopic, event.DeviceId), platform.base.qos, false, iot.Interface2JsonString(event)); token.Wait() && token.Error() != nil {
		log.Errorf("平台下发事件失败,device:%s", event.DeviceId)
		return false
	}
	return true
}

func (platform *iotPlatform) SetGatewayBatchReportSubDeviceHandler(handler GatewayBatchReportSubDeviceHandler) {
	platform.base.batchReportHandler = handler
}

func (platform *iotPlatform) SetPropertiesReportHandler(handler PropertiesReportHandler) {
	platform.base.propertiesReportHandler = handler
}

func (platform *iotPlatform) SetVersionReportHandler(handler VersionReportHandler) {
	platform.base.versionReportHandler = handler
}

func (platform *iotPlatform) QueryVersion(device string) bool {
	event := iot.EventData{
		DeviceId: device,
		Services: []iot.ServiceEvent{{
			ServiceId: iot.OTA,
			EventType: event_type.VERSION_QUERY,
			EventTime: utils.Timestamp(),
			Paras:     nil,
		}},
	}
	return platform.EventDown(event)
}

func (platform *iotPlatform) SendUpgradeMessage(device string, upgradeType iot.EventType, data iot.UpgradePara) bool {
	event := iot.EventData{
		DeviceId: device,
		Services: []iot.ServiceEvent{{
			ServiceId: iot.OTA,
			EventType: upgradeType,
			EventTime: utils.Timestamp(),
			Paras:     data,
		}},
	}
	return platform.EventDown(event)
}

func (platform *iotPlatform) SetUpgradeProcessReportHandler(handler UpgradeProcessReportHandler) {
	platform.base.upgradeProcessReportHandler = handler
}

func (platform *iotPlatform) SendMessage(message iot.Message) bool {
	msg := utils.Interface2JsonString(message)
	if token := platform.base.Client.Publish(iot.FormatTopic(topics.MessageDownTopic, message.DeviceId), platform.base.qos, false, msg); token.Wait() && token.Error() != nil {
		glog.Errorf("发送设备[%s]消息失败", message.DeviceId)
		return false
	}
	return true
}
func (platform *iotPlatform) SendCommand(command iot.Command) bool {
	cmd := utils.Interface2JsonString(command)
	if token := platform.base.Client.Publish(iot.FormatTopic1(topics.CommandDownTopic, command.DeviceId, command.CommandId)+command.CommandId, platform.base.qos, false, cmd); token.Wait() && token.Error() != nil {
		glog.Errorf("发送设备[%s]消息失败", command.CommandId)
		return false
	}
	return true
}
