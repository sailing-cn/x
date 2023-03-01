package iot_platform

import "sailing.cn/iot"

type PropertiesReportHandler func(request iot.DeviceService)
type GatewayBatchReportSubDeviceHandler func(request iot.DevicesService)

// DeviceInfoReportHandler 设备信息上报事件处理
type DeviceInfoReportHandler func(deviceId string, paras iot.ReportDeviceInfoPara)

type CommandHandler func(command iot.CommandResponse) bool

type MessageHandler func(message iot.Message) bool

type VersionReportHandler func(deviceId string, version iot.VersionInfo)

type UpgradeProcessReportHandler func(deviceId string, result iot.UpgradeProgress)

type DeviceOnlineHandler func(info *iot.DeviceOnlineInfo)

type DeviceOfflineHandler func(info *iot.DeviceOfflineInfo)
