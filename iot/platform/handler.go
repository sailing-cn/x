package platform

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"sailing.cn/v2/iot"
)

// PropertiesReportHandler 属性上报处理器
type PropertiesReportHandler func(request iot.DeviceService)

// GatewayBatchReportSubDeviceHandler 网关批量上报子设备处理器
type GatewayBatchReportSubDeviceHandler func(request iot.DevicesService)

// DeviceInfoReportHandler 设备信息上报事件处理
type DeviceInfoReportHandler func(deviceId string, paras iot.ReportDeviceInfoPara)

// CommandHandler 命令处理器
type CommandHandler func(command iot.CommandResponse) bool

// MessageHandler 消息处理器
type MessageHandler func(message iot.Message) bool

// VersionReportHandler 版本上报处理器
type VersionReportHandler func(deviceId string, version iot.VersionInfo)

// UpgradeProcessReportHandler 远程升级处理器
type UpgradeProcessReportHandler func(deviceId string, result iot.UpgradeProgress)

// DeviceOnlineHandler 设备上线处理器
type DeviceOnlineHandler func(info *iot.DeviceOnlineInfo)

// DeviceOfflineHandler 设备离线处理器
type DeviceOfflineHandler func(info *iot.DeviceOfflineInfo)

// ConnectLostHandler 断开连接处理器
type ConnectLostHandler func(client mqtt.Client, err error)
