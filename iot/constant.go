package iot

const (
	// SUB_DEVICE_MANAGER 子设备管理
	SUB_DEVICE_MANAGER = "$sub_device_manager"
	//OTA 固件
	OTA = "$ota"
)

const (
	// TASK_CREATED 任务已创建
	TASK_CREATED = "TASK_CREATED"
	//QUERY_VERSION 版本查询
	QUERY_VERSION = "QUERY_VERSION"
	// DEVICE_VERSION_REPORT 设备上报版本
	DEVICE_VERSION_REPORT = "DEVICE_VERSION_REPORT"
	//UPGRADE 升级
	UPGRADE = "UPGRADE"
	//UPGRADE_PROCESS_REPORT 升级进度上报
	UPGRADE_PROCESS_REPORT = "UPGRADE_PROCESS_REPORT"
	//UPDATE_SHADOW_VALUE 更新设备影子数据
	UPDATE_SHADOW_VALUE = "UPDATE_SHADOW_VALUE"
	//SEND_MESSAGE 发送消息
	SEND_MESSAGE = "ON_SEND_MESSAGE"
	//UPDATE_MESSAGE_STATUS 更新消息状态
	UPDATE_MESSAGE_STATUS = "UPDATE_MESSAGE_STATUS"
	//UPDATE_COMMAND_STATUS 更新命令状态
	UPDATE_COMMAND_STATUS = "UPDATE_COMMAND_STATUS"
	//DOWN_COMMAND 下发命令
	DOWN_COMMAND = "DOWN_COMMAND"
)

type AuthType string
type NodeType string
type DeviceStatus string

// OperatingMode 操作方式
type OperatingMode string

const (
	SECRET       AuthType = "secret"      //密码
	CERTIFICATES AuthType = "certificate" //证书

	GATEWAY  NodeType = "gateway"  //网关
	ENDPOINT NodeType = "endpoint" //子设备

	INACTIVE DeviceStatus = "inactive" //未激活
	OFFLINE  DeviceStatus = "offline"  //离线
	ONLINE   DeviceStatus = "online"   //在线
	FROZEN   DeviceStatus = "frozen"   //冻结

	READ  OperatingMode = "read"  //读
	WRITE OperatingMode = "write" //写
)

type TaskStatus string

const (
	//WAITING 等待中
	WAITING TaskStatus = "waiting"
	//FAIL 失败
	FAIL TaskStatus = "fail"
	//PROCESSING 执行中
	PROCESSING TaskStatus = "processing"
	//SUCCESS 成功
	SUCCESS TaskStatus = "success"
	//STOPPED 停止
	STOPPED TaskStatus = "stopped"
)

const (

	// ON_DEVICE_CONNECTED 设备连接
	ON_DEVICE_CONNECTED = "DEVICE_CONNECTED"

	// ON_DEVICE_DISCONNECTED 设备断开连接
	ON_DEVICE_DISCONNECTED = "DEVICE_DISCONNECTED"

	//ON_PRODUCT_UPDATE 产品更新
	ON_PRODUCT_UPDATE = "PRODUCT_UPDATE"

	// ON_SUB_DEVICE_CREATE  添加子设备
	ON_SUB_DEVICE_CREATE = "SUB_DEVICE_CREATE"

	//ON_SUB_DEVICE_UPDATE 更新子设备
	ON_SUB_DEVICE_UPDATE = "SUB_DEVICE_UPDATE"

	//ON_SUB_DEVICE_DELETE 删除子设备
	ON_SUB_DEVICE_DELETE = "SUB_DEVICE_DELETE"

	//ON_TASK_CREATE 创建任务
	ON_TASK_CREATE = "ON_TASK_CREATE"
	// ON_TASk_STOP 停止任务
	ON_TASk_STOP = "TASk_STOP"
	// ON_SEND_MESSAGE 发送消息
	ON_SEND_MESSAGE = "ON_SEND_MESSAGE"
	//ON_DOWN_COMMAND 下发命令
	ON_DOWN_COMMAND = "ON_DOWN_COMMAND"
	//UPDATE_PROCESS_SUCCESS 更新成功
	UPDATE_PROCESS_SUCCESS = "UPDATE_PROCESS_SUCCESS"
)

// 数据类型
const (
	Int     = "int"
	Decimal = "decimal"
	Bool    = "bool"
	String  = "string"
)

const (
	//ADD_SUB_DEVICE_NOTIFY 添加子设备通知
	ADD_SUB_DEVICE_NOTIFY EventType = "add_sub_device_notify"
	// VERSION_REPORT 版本上报
	VERSION_REPORT EventType = "version_report"
	//FIRMWARE_UPGRADE 固件升级
	FIRMWARE_UPGRADE EventType = "firmware_upgrade"
	//SOFTWARE_UPGRADE 软件升级
	SOFTWARE_UPGRADE EventType = "software_upgrade"
	//CONFIG_UPGRADE 配置升级
	CONFIG_UPGRADE          EventType = "config_upgrade"
	VERSION_QUERY           EventType = "version_query"           //版本查询
	UPGRADE_PROGRESS_REPORT EventType = "upgrade_progress_report" //设备升级状态上报
	MESSAGE_RECEUVED        EventType = "message_received"        //消息状态变更
)
