package iot_device

//import (
//	"iot.sailing.cn/iot"
//	"iot.sailing.cn/iot/iot.device/topics"
//)
//
//// 子设备添加回调函数
//type SubDevicesAddHandler func(devices iot.AddSubDeviceParas)
//
////子设备删除糊掉函数
//type SubDevicesDeleteHandler func(devices iot.DeleteSubDeviceParas)
//
//// 处理平台下发的命令
//type CommandHandler func(iot.Command) bool
//
//// 设备消息
//type MessageHandler func(message iot.Message) bool
//
//// 平台设置设备属性
//type DevicePropertiesSetHandler func(message DevicePropertyDownRequest) bool
//
//// 平台查询设备属性
//type DevicePropertyQueryHandler func(query DevicePropertyQueryRequest) iot.DevicePropertyEntity
//
//// 设备执行软件/固件升级.upgradeType = 0 软件升级，upgradeType = 1 固件升级
//type DeviceUpgradeHandler func(upgradeType iot.EventType, info iot.UpgradePara) iot.UpgradeProgress
//
//// 设备上报版本
//type VersionReporter func() *iot.VersionInfo
//
////SubDevicesStatus 网关更新子设备状态
//type SubDevicesStatus struct {
//	DeviceStatuses []DeviceStatus `json:"device_statuses"`
//}
//
//type DeviceStatus struct {
//	DeviceId string `json:"device_id"`
//	Status   string `json:"status"` // 子设备状态。 OFFLINE：设备离线 ONLINE：设备上线
//}
//
//// 平台设置设备属性==================================================
//type DevicePropertyDownRequest struct {
//	ObjectDeviceId string                           `json:"object_device_id"`
//	Services       []DevicePropertyDownRequestEntry `json:"services"`
//}
//
//type DevicePropertyDownRequestEntry struct {
//	ServiceId  string      `json:"service_id"`
//	Properties interface{} `json:"properties"`
//}
//
//// 平台设置设备属性==================================================
//type DevicePropertyQueryRequest struct {
//	ObjectDeviceId string `json:"object_device_id"`
//	ServiceId      string `json:"service_id"`
//}
//
//// 设备获取设备影子数据
//type DevicePropertyQueryResponseHandler func(response DevicePropertyQueryResponse)
//
//type DevicePropertyQueryResponse struct {
//	ObjectDeviceId string             `json:"object_device_id"`
//	Shadow         []DeviceShadowData `json:"shadow"`
//}
//
//type DeviceShadowData struct {
//	ServiceId string                     `json:"service_id"`
//	Desired   DeviceShadowPropertiesData `json:"desired"`
//	Reported  DeviceShadowPropertiesData `json:"reported"`
//	Version   int                        `json:"version"`
//}
//type DeviceShadowPropertiesData struct {
//	Properties interface{} `json:"properties"`
//	EventTime  string      `json:"event_time"`
//}
//
//// 文件上传下载管理
//func CreateFileUploadDownLoadResultResponse(filename, action string, result bool) FileResultResponse {
//	code := 0
//	if !result {
//		code = 1
//	}
//
//	paras := FileResultServiceEventParas{
//		ObjectName: filename,
//		ResultCode: code,
//	}
//
//	serviceEvent := FileResultResponseServiceEvent{
//		Paras: paras,
//	}
//	serviceEvent.ServiceId = "$file_manager"
//	if action == topics.FileActionDownload {
//		serviceEvent.EventType = "download_result_report"
//	}
//	if action == topics.FileActionUpload {
//		serviceEvent.EventType = "upload_result_report"
//	}
//	serviceEvent.EventTime = iot.GetEventTimeStamp()
//
//	var services []FileResultResponseServiceEvent
//	services = append(services, serviceEvent)
//
//	response := FileResultResponse{
//		Services: services,
//	}
//
//	return response
//}
//
//// 设备获取文件上传下载请求体
//type FileRequest struct {
//	ObjectDeviceId string                    `json:"object_device_id"`
//	Services       []FileRequestServiceEvent `json:"services"`
//}
//
//// 平台下发文件上传和下载URL响应
//type FileResponse struct {
//	ObjectDeviceId string                     `json:"object_device_id"`
//	Services       []FileResponseServiceEvent `json:"services"`
//}
//
//type FileResultResponse struct {
//	ObjectDeviceId string                           `json:"object_device_id"`
//	Services       []FileResultResponseServiceEvent `json:"services"`
//}
//
//type BaseServiceEvent struct {
//	ServiceId string `json:"service_id"`
//	EventType string `json:"event_type"`
//	EventTime int64  `json:"event_time,omitempty"`
//}
//
//type FileRequestServiceEvent struct {
//	BaseServiceEvent
//	Paras FileRequestServiceEventParas `json:"paras"`
//}
//
//type FileResponseServiceEvent struct {
//	BaseServiceEvent
//	Paras FileResponseServiceEventParas `json:"paras"`
//}
//
//type FileResultResponseServiceEvent struct {
//	BaseServiceEvent
//	Paras FileResultServiceEventParas `json:"paras"`
//}
//
//// 设备获取文件上传下载URL参数
//type FileRequestServiceEventParas struct {
//	FileName       string      `json:"file_name"`
//	FileAttributes interface{} `json:"file_attributes"`
//}
//
//// 平台下发响应参数
//type FileResponseServiceEventParas struct {
//	Url            string      `json:"url"`
//	BucketName     string      `json:"bucket_name"`
//	ObjectName     string      `json:"object_name"`
//	Expire         int         `json:"expire"`
//	FileAttributes interface{} `json:"file_attributes"`
//}
//
//// 上报文件上传下载结果参数
//type FileResultServiceEventParas struct {
//	ObjectName        string `json:"object_name"`
//	ResultCode        int    `json:"result_code"`
//	StatusCode        int    `json:"status_code"`
//	StatusDescription string `json:"status_description"`
//}
//
//// 上报设备信息请求
//type ReportDeviceInfoRequest struct {
//	ObjectDeviceId string                         `json:"object_device_id,omitempty"`
//	Services       []ReportDeviceInfoServiceEvent `json:"services,omitempty"`
//}
//
//type ReportDeviceInfoServiceEvent struct {
//	BaseServiceEvent
//	Paras ReportDeviceInfoEventParas `json:"paras,omitempty"`
//}
//
//// 设备信息上报请求参数
//type ReportDeviceInfoEventParas struct {
//	DeviceSdkVersion string `json:"device_sdk_version,omitempty"`
//	SwVersion        string `json:"sw_version,omitempty"`
//	FwVersion        string `json:"fw_version,omitempty"`
//}
//
//// 上报设备日志请求
//type ReportDeviceLogRequest struct {
//	Services []ReportDeviceLogServiceEvent `json:"services,omitempty"`
//}
//
//type ReportDeviceLogServiceEvent struct {
//	BaseServiceEvent
//	Paras DeviceLogEntry `json:"paras,omitempty"`
//}
//
//// 设备状态日志收集器
//type DeviceStatusLogCollector func(endTime string) []DeviceLogEntry
//
//// 设备属性日志收集器
//type DevicePropertyLogCollector func(endTime string) []DeviceLogEntry
//
//// 设备消息日志收集器
//type DeviceMessageLogCollector func(endTime string) []DeviceLogEntry
//
//// 设备命令日志收集器
//type DeviceCommandLogCollector func(endTime string) []DeviceLogEntry
//
//type DeviceLogEntry struct {
//	Timestamp string `json:"timestamp"` // 日志产生时间
//	Type      string `json:"type"`      // 日志类型：DEVICE_STATUS，DEVICE_PROPERTY ，DEVICE_MESSAGE ，DEVICE_COMMAND
//	Content   string `json:"content"`   // 日志内容
//}
