package iot

const (
	DEVICE_ONLINE                         = "$sailing.iot/device/online"                                       //设备上线
	DEVICE_OFFLINE                        = "$sailing.iot/device/offline"                                      //设备离线
	PLATFORM_PROPERTIES_UP                = "$sailing.iot/devices/+/sys/properties/report"                     //平台订阅设备属性上报
	BATCH_REPORT_SUB_DEVICE               = "$sailing.iot/devices/+/sys/gateway/sub_devices/properties/report" //平台订阅网关设备批量上报子设备属性
	DEVICE_EVENT_TO_PLATFORM_TOPIC string = "$sailing.iot/devices/+/sys/events/up"                             //设备事件上报
	DEVICE_MESSAGE_UP              string = "$sailing.iot/devices/+/sys/messages/up"                           //设备发送消息
	COMMAND_RESPONSE               string = "$sailing.iot/devices/+/sys/commands/response/#"                   //命令响应

	// CommandDownTopic 平台下发命令topic
	//CommandDownTopic string = "$sailing.iot/devices/{device_id}/sys/commands/down/request_id="
)

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
