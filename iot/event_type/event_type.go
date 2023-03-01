package event_type

import "sailing.cn/iot"

const (
	//ADD_SUB_DEVICE_NOTIFY 添加子设备通知
	ADD_SUB_DEVICE_NOTIFY iot.EventType = "add_sub_device_notify"
	// VERSION_REPORT 版本上报
	VERSION_REPORT iot.EventType = "version_report"

	//FIRMWARE_UPGRADE 固件升级
	FIRMWARE_UPGRADE iot.EventType = "firmware_upgrade"
	//SOFTWARE_UPGRADE 软件升级
	SOFTWARE_UPGRADE iot.EventType = "software_upgrade"
	//CONFIG_UPGRADE 配置升级
	CONFIG_UPGRADE          iot.EventType = "config_upgrade"
	VERSION_QUERY           iot.EventType = "version_query"           //版本查询
	UPGRADE_PROGRESS_REPORT iot.EventType = "upgrade_progress_report" //设备升级状态上报
	MESSAGE_RECEUVED        iot.EventType = "message_received"        //消息状态变更
)
