package topics

const (
	// MessageDownTopic 平台下发消息topic
	MessageDownTopic string = "$sailing.iot/devices/{device_id}/sys/messages/down"

	// MessageUpTopic 设备上报消息topic
	MessageUpTopic string = "$sailing.iot/devices/{device_id}/sys/messages/up"

	// CommandDownTopic 平台下发命令topic
	CommandDownTopic string = "$sailing.iot/devices/{device_id}/sys/commands/down/#"

	// CommandResponseTopic 设备响应平台命令
	CommandResponseTopic string = "$sailing.iot/devices/{device_id}/sys/commands/response/request_id="

	// PropertiesUpTopic 设备上报属性
	PropertiesUpTopic string = "$sailing.iot/devices/{device_id}/sys/properties/report"

	// PropertiesSetRequestTopic 平台设置属性topic
	PropertiesSetRequestTopic string = "$sailing.iot/devices/{device_id}/sys/properties/set/#"

	// PropertiesSetResponseTopic 设备响应平台属性设置topic
	PropertiesSetResponseTopic string = "$sailing.iot/devices/{device_id}/sys/properties/set/response/request_id="

	// PropertiesQueryRequestTopic 平台查询设备属性
	PropertiesQueryRequestTopic string = "$sailing.iot/devices/{device_id}/sys/properties/get/#"

	// PropertiesQueryResponseTopic 设备响应平台属性查询
	PropertiesQueryResponseTopic string = "$sailing.iot/devices/{device_id}/sys/properties/get/response/request_id="

	// DeviceShadowQueryRequestTopic 设备侧获取平台的设备影子数据
	DeviceShadowQueryRequestTopic string = "$sailing.iot/devices/{device_id}/sys/shadow/get/request_id="

	// DeviceShadowQueryResponseTopic 设备侧响应获取平台设备影子
	DeviceShadowQueryResponseTopic string = "$sailing.iot/devices/{device_id}/sys/shadow/get/response/#"

	// GatewayBatchReportSubDeviceTopic 网关批量上报子设备属性
	GatewayBatchReportSubDeviceTopic string = "$sailing.iot/devices/{device_id}/sys/gateway/sub_devices/properties/report"

	// FileActionUpload 平台下发文件上传和下载URL
	FileActionUpload   string = "upload"
	FileActionDownload string = "download"

	// DeviceToPlatformTopic 设备或网关向平台发送请求
	DeviceToPlatformTopic string = "$sailing.iot/devices/{device_id}/sys/events/up"

	// PlatformEventToDeviceTopic 平台向设备下发事件topic
	PlatformEventToDeviceTopic string = "$sailing.iot/devices/{device_id}/sys/events/down"
)
