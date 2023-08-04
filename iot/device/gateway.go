package device

import "sailing.cn/v2/iot"

type baseGateway interface {
	// SetSubDevicesAddHandler 设置平台添加子设备回调函数
	SetSubDevicesAddHandler(handler SubDevicesAddHandler)

	// SetSubDevicesDeleteHandler 设置平台删除子设备回调函数
	SetSubDevicesDeleteHandler(handler SubDevicesDeleteHandler)
}

type Gateway interface {
	baseGateway

	// UpdateSubDeviceState 网关更新子设备状态
	UpdateSubDeviceState(subDevicesStatus iot.SubDevicesStatus) bool

	// DeleteSubDevices 网关删除子设备
	DeleteSubDevices(deviceIds []string) bool

	// AddSubDevices 网关添加子设备
	AddSubDevices(deviceInfos []iot.SubDeviceInfo) bool

	// SyncAllVersionSubDevices 网关同步子设备列表,默认实现不指定版本
	SyncAllVersionSubDevices()

	// SyncSubDevices 网关同步特定版本子设备列表
	SyncSubDevices(version int)
}

type AsyncGateway interface {
	baseGateway

	// UpdateSubDeviceState 网关更新子设备状态
	UpdateSubDeviceState(subDevicesStatus iot.SubDevicesStatus) AsyncResult

	// DeleteSubDevices 网关删除子设备
	DeleteSubDevices(deviceIds []string) AsyncResult

	// AddSubDevices 网关添加子设备
	AddSubDevices(deviceInfos []iot.SubDeviceInfo) AsyncResult

	// SyncAllVersionSubDevices 网关同步子设备列表,默认实现不指定版本
	SyncAllVersionSubDevices() AsyncResult

	// SyncSubDevices 网关同步特定版本子设备列表
	SyncSubDevices(version int) AsyncResult
}
