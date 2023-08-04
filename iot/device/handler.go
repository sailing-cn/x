package device

import (
	"sailing.cn/v2/iot"
	"sailing.cn/v2/utils/timestamp"
)

// DevicePropertyQueryResponseHandler 设备获取设备影子数据
type DevicePropertyQueryResponseHandler func(response iot.DevicePropertyQueryResponse)

// MessageHandler 设备消息
type MessageHandler func(message iot.Message) bool

// CommandHandler 处理平台下发的命令
type CommandHandler func(iot.Command) bool

// DevicePropertiesSetHandler 平台设置设备属性
type DevicePropertiesSetHandler func(message iot.DevicePropertyDownRequest) bool

// VersionReporter 设备上报版本
type VersionReporter func() *iot.VersionInfo

// DeviceUpgradeHandler 设备执行软件/固件升级.upgradeType = software_upgrade 软件升级，upgradeType = firmware_upgrade 固件升级，upgradeType = config_upgrade 配置升级
type DeviceUpgradeHandler func(upgradeType iot.EventType, info iot.UpgradePara) iot.UpgradeProgress

// DevicePropertyQueryHandler 平台查询设备属性
type DevicePropertyQueryHandler func(query iot.DevicePropertyQueryRequest) iot.DevicePropertyEntity

// SubDevicesAddHandler 子设备添加回调函数
type SubDevicesAddHandler func(devices iot.AddSubDeviceParas)

// SubDevicesDeleteHandler 子设备删除糊掉函数
type SubDevicesDeleteHandler func(devices iot.DeleteSubDeviceParas)

// DeviceStatusLogCollector 设备状态日志收集器
type DeviceStatusLogCollector func(endTime string) []iot.DeviceLog

// DevicePropertyLogCollector 设备属性日志收集器
type DevicePropertyLogCollector func(endTime string) []iot.DeviceLog

// DeviceMessageLogCollector 设备消息日志收集器
type DeviceMessageLogCollector func(endTime string) []iot.DeviceLog

// DeviceCommandLogCollector 设备命令日志收集器
type DeviceCommandLogCollector func(endTime string) []iot.DeviceLog

// CreateFileUploadDownLoadResultResponse 文件上传下载管理
func CreateFileUploadDownLoadResultResponse(filename, action string, result bool) iot.FileResultResponse {
	code := 0
	if !result {
		code = 1
	}

	paras := iot.FileResultServiceEventParas{
		ObjectName: filename,
		ResultCode: code,
	}

	serviceEvent := iot.FileResultResponseServiceEvent{
		Paras: paras,
	}
	serviceEvent.ServiceId = "$file_manager"
	if action == iot.FileActionDownload {
		serviceEvent.EventType = "download_result_report"
	}
	if action == iot.FileActionUpload {
		serviceEvent.EventType = "upload_result_report"
	}
	serviceEvent.EventTime = timestamp.Timestamp()

	var services []iot.FileResultResponseServiceEvent
	services = append(services, serviceEvent)

	response := iot.FileResultResponse{
		Services: services,
	}

	return response
}
