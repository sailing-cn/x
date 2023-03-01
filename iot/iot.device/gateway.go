package iot_device

import "sailing.cn/iot"

type baseGateway interface {
	// 设置平台添加子设备回调函数
	SetSubDevicesAddHandler(handler SubDevicesAddHandler)

	// 设置平台删除子设备回调函数
	SetSubDevicesDeleteHandler(handler SubDevicesDeleteHandler)
}

type Gateway interface {
	baseGateway

	// 网关更新子设备状态
	UpdateSubDeviceState(subDevicesStatus iot.SubDevicesStatus) bool

	// 网关删除子设备
	DeleteSubDevices(deviceIds []string) bool

	// 网关添加子设备
	AddSubDevices(deviceInfos []iot.SubDeviceInfo) bool

	// 网关同步子设备列表,默认实现不指定版本
	SyncAllVersionSubDevices()

	// 网关同步特定版本子设备列表
	SyncSubDevices(version int)
}

type AsyncGateway interface {
	baseGateway

	// 网关更新子设备状态
	UpdateSubDeviceState(subDevicesStatus iot.SubDevicesStatus) AsyncResult

	// 网关删除子设备
	DeleteSubDevices(deviceIds []string) AsyncResult

	// 网关添加子设备
	AddSubDevices(deviceInfos []iot.SubDeviceInfo) AsyncResult

	// 网关同步子设备列表,默认实现不指定版本
	SyncAllVersionSubDevices() AsyncResult

	// 网关同步特定版本子设备列表
	SyncSubDevices(version int) AsyncResult
}
