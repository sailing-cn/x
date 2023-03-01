package topics

const (
	DEVICE_ONLINE                         = "$sailing.iot/device/online"                                       //设备上线
	DEVICE_OFFLINE                        = "$sailing.iot/device/offline"                                      //设备离线
	PLATFORM_PROPERTIES_UP                = "$sailing.iot/devices/+/sys/properties/report"                     //平台订阅设备属性上报
	BATCH_REPORT_SUB_DEVICE               = "$sailing.iot/devices/+/sys/gateway/sub_devices/properties/report" //平台订阅网关设备批量上报子设备属性
	DEVICE_EVENT_TO_PLATFORM_TOPIC string = "$sailing.iot/devices/+/sys/events/up"                             //设备事件上报
	DEVICE_MESSAGE_UP              string = "$sailing.iot/devices/+/sys/messages/up"                           //设备发送消息
	COMMAND_RESPONSE               string = "$sailing.iot/devices/+/sys/commands/response/#"                   //命令响应

	// CommandDownTopic 平台下发命令topic
	CommandDownTopic string = "$sailing.iot/devices/{device_id}/sys/commands/down/request_id="
)
