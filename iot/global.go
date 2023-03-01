package iot

import (
	"time"
)

type EventType string
type UpgradeResult int

const (
	// SUB_DEVICE_MANAGER 子设备管理
	SUB_DEVICE_MANAGER = "$sub_device_manager"
	//OTA 固件
	OTA = "$ota"
)

const (
	SUCCESS             UpgradeResult = 0   //处理成功
	BUSY                UpgradeResult = 1   //设备使用中
	SIGNAL_QUALITY_BAD  UpgradeResult = 2   //信号质量差
	LATEST              UpgradeResult = 3   //已经是最新版本
	LOW_BATTERY         UpgradeResult = 4   //电量不足
	NO_SPACE            UpgradeResult = 5   //剩余空间不足
	DOWNLOAD_TIMEOUT    UpgradeResult = 6   //下载超时
	VERIFICATION_FAIL   UpgradeResult = 7   //升级包校验失败
	NOT_SUPPORT_PACKAGE UpgradeResult = 8   //升级包类型不支持
	OUT_OF_MEMORY       UpgradeResult = 9   //内存不足
	INSTALL_FAIL        UpgradeResult = 10  //安装升级包失败
	INTERNAL_ERROR      UpgradeResult = 255 //内部异常
)

// DeviceOnlineInfo 设备上线信息
type DeviceOnlineInfo struct {
	Username       string `json:"username"`
	Ts             int64  `json:"ts"`
	Sockport       int    `json:"sockport"`
	ProtoVer       int    `json:"proto_ver"`
	ProtoName      string `json:"proto_name"`
	Keepalive      int    `json:"keepalive"`
	Ipaddress      string `json:"ipaddress"`
	ExpiryInterval int    `json:"expiry_interval"`
	ConnectedAt    int64  `json:"connected_at"`
	Connack        int    `json:"connack"`
	Clientid       string `json:"clientid"`
	CleanStart     bool   `json:"clean_start"`
}

// DeviceOfflineInfo 设备离线信息
type DeviceOfflineInfo struct {
	Username       string `json:"username"`
	Ts             int64  `json:"ts"`
	Sockport       int    `json:"sockport"`
	Reason         string `json:"reason"`
	ProtoVer       int    `json:"proto_ver"`
	ProtoName      string `json:"proto_name"`
	Ipaddress      string `json:"ipaddress"`
	DisconnectedAt int64  `json:"disconnected_at"`
	Clientid       string `json:"clientid"`
}

// Message 消息
type Message struct {
	DeviceId  string      `json:"device_id"`
	Name      string      `json:"name"`
	MessageId string      `json:"message_id"`
	Content   string      `json:"content"`
	Data      interface{} `json:"data"`
	Type      string      `json:"type"`
}

type VersionInfo struct {
	SoftwareVersion string `json:"software_version"`
	FirmwareVersion string `json:"firmware_version"`
	CollectVersion  string `json:"collect_version"`
	ServiceVersion  string `json:"service_version"`
}

// 定义平台和设备之间的数据交换结构体

type EventData struct {
	DeviceId string         `json:"device_id,omitempty"`
	Services []ServiceEvent `json:"services"`
}

// ServiceEvent 平台和设备之间事件结构体
type ServiceEvent struct {
	ServiceId string      `json:"service_id"`
	EventType EventType   `json:"event_type"`
	EventTime int64       `json:"event_time"`
	Paras     interface{} `json:"paras"` // 不同类型的请求paras使用的结构体不同
}

// AddSubDeviceParas 子设备添加事件参数
type AddSubDeviceParas struct {
	Devices []SubDeviceInfo `json:"devices"`
	Version int64           `json:"version"`
}

type DeleteSubDeviceParas struct {
	Devices []SubDeviceInfo `json:"devices"`
	Version int64           `json:"version"`
}

// SubDeviceInfo 子设备信息
type SubDeviceInfo struct {
	ParentDeviceId string      `json:"parent_device_id,omitempty"`
	Type           string      `json:"type,omitempty"`
	NodeId         string      `json:"node_id,omitempty"`
	DeviceId       string      `json:"device_id,omitempty"`
	Name           string      `json:"name,omitempty"`
	Description    string      `json:"description,omitempty"`
	ManufacturerId string      `json:"manufacturer_id,omitempty"`
	Model          string      `json:"model,omitempty"`
	ProductId      string      `json:"product_id"`
	FwVersion      string      `json:"fw_version,omitempty"`
	SwVersion      string      `json:"sw_version,omitempty"`
	Status         string      `json:"status,omitempty"`
	ExtensionInfo  interface{} `json:"extension_info,omitempty"`
}

// DevicesService 网关批量上报子设备属性
type DevicesService struct {
	Devices []*DeviceService `json:"devices"`
}

type DeviceService struct {
	DeviceId string                 `json:"device_id"`
	Services []DevicePropertyEntity `json:"services"`
}

// 设备属性
type DeviceProperties struct {
	Services []DevicePropertyEntity `json:"services"`
}

// DevicePropertyEntity 设备的一个属性
type DevicePropertyEntity struct {
	ServiceId  string                 `json:"service_id"`
	Properties map[string]interface{} `json:"properties"`
	EventTime  time.Time              `json:"event_time"`
}

type ReportDeviceInfoPara struct {
	SoftwareVersion string `json:"software_version"`
	FirmwareVersion string `json:"firmware_version"`
}

type UpgradePara struct {
	Version  string `json:"version"`
	URL      string `json:"url"`
	FileSize int64  `json:"file_size"`
	Token    string `json:"token"`
	Expires  int64  `json:"expires"`
	Sign     string `json:"sign"`
}

// Command 设备命令
type Command struct {
	DeviceId    string      `json:"device_id"`
	ServiceId   string      `json:"service_id"`
	CommandName string      `json:"command_name"`
	Paras       interface{} `json:"paras"`
	CommandId   string      `json:"command_id"`
}

type CommandResponse struct {
	ResultCode   byte        `json:"result_code"`
	ResponseName string      `json:"response_name"`
	Paras        interface{} `json:"paras"`
	RequestId    string      `json:"request_id"`
}

// UpgradeProgress 设备升级状态响应，用于设备向平台反馈进度，错误信息等
type UpgradeProgress struct {
	ResultCode  UpgradeResult `json:"result_code"`
	Progress    int           `json:"progress"`    // 设备的升级进度，范围：0到100
	Version     string        `json:"version"`     // 设备当前版本号
	Description string        `json:"description"` // 升级状态描述信息，可以返回具体升级失败原因。
}

// DevicePropertyQueryRequest 平台设置设备属性
type DevicePropertyQueryRequest struct {
	DeviceId  string `json:"device_id"`
	ServiceId string `json:"service_id"`
}

type DevicePropertyQueryResponse struct {
	ObjectDeviceId string             `json:"object_device_id"`
	Shadow         []DeviceShadowData `json:"shadow"`
}

type DeviceShadowData struct {
	ServiceId string                     `json:"service_id"`
	Desired   DeviceShadowPropertiesData `json:"desired"`
	Reported  DeviceShadowPropertiesData `json:"reported"`
	Version   int                        `json:"version"`
}
type DeviceShadowPropertiesData struct {
	Properties interface{} `json:"properties"`
	EventTime  string      `json:"event_time"`
}

type BaseServiceEvent struct {
	ServiceId string `json:"service_id"`
	EventType string `json:"event_type"`
	EventTime int64  `json:"event_time,omitempty"`
}

type FileRequestServiceEvent struct {
	BaseServiceEvent
	Paras FileRequestServiceEventParas `json:"paras"`
}

type FileResponseServiceEvent struct {
	BaseServiceEvent
	Paras FileResponseServiceEventParas `json:"paras"`
}

type FileResultResponseServiceEvent struct {
	BaseServiceEvent
	Paras FileResultServiceEventParas `json:"paras"`
}

// FileRequestServiceEventParas 设备获取文件上传下载URL参数
type FileRequestServiceEventParas struct {
	FileName       string      `json:"file_name"`
	FileAttributes interface{} `json:"file_attributes"`
}

// FileResponseServiceEventParas 平台下发响应参数
type FileResponseServiceEventParas struct {
	Url            string      `json:"url"`
	BucketName     string      `json:"bucket_name"`
	ObjectName     string      `json:"object_name"`
	Expire         int         `json:"expire"`
	FileAttributes interface{} `json:"file_attributes"`
}

// FileResultServiceEventParas 上报文件上传下载结果参数
type FileResultServiceEventParas struct {
	ObjectName        string `json:"object_name"`
	ResultCode        int    `json:"result_code"`
	StatusCode        int    `json:"status_code"`
	StatusDescription string `json:"status_description"`
}

// FileRequest 设备获取文件上传下载请求体
type FileRequest struct {
	ObjectDeviceId string                    `json:"object_device_id"`
	Services       []FileRequestServiceEvent `json:"services"`
}

// FileResponse 平台下发文件上传和下载URL响应
type FileResponse struct {
	ObjectDeviceId string                     `json:"object_device_id"`
	Services       []FileResponseServiceEvent `json:"services"`
}

type FileResultResponse struct {
	ObjectDeviceId string                           `json:"object_device_id"`
	Services       []FileResultResponseServiceEvent `json:"services"`
}

type DeviceLog struct {
	Timestamp string `json:"timestamp"` // 日志产生时间
	Type      string `json:"type"`      // 日志类型：DEVICE_STATUS，DEVICE_PROPERTY ，DEVICE_MESSAGE ，DEVICE_COMMAND
	Content   string `json:"content"`   // 日志内容
}

// ReportDeviceLogRequest 上报设备日志请求
type ReportDeviceLogRequest struct {
	Services []ReportDeviceLogServiceEvent `json:"services,omitempty"`
}

type ReportDeviceLogServiceEvent struct {
	BaseServiceEvent
	Paras DeviceLog `json:"paras,omitempty"`
}

// DevicePropertyDownRequest 平台设置设备属性
type DevicePropertyDownRequest struct {
	DeviceId string                            `json:"device_id"`
	Services []DevicePropertyDownRequestEntity `json:"services"`
}
type DevicePropertyDownRequestEntity struct {
	ServiceId  string      `json:"service_id"`
	Properties interface{} `json:"properties"`
}

type ReportDeviceInfoServiceEvent struct {
	BaseServiceEvent
	Paras ReportDeviceInfoEventParas `json:"paras,omitempty"`
}

// ReportDeviceInfoEventParas 设备信息上报请求参数
type ReportDeviceInfoEventParas struct {
	DeviceSdkVersion string `json:"device_sdk_version,omitempty"`
	SwVersion        string `json:"sw_version,omitempty"`
	FwVersion        string `json:"fw_version,omitempty"`
}

// ReportDeviceInfoRequest 上报设备信息请求
type ReportDeviceInfoRequest struct {
	DeviceId string                         `json:"device_id,omitempty"`
	Services []ReportDeviceInfoServiceEvent `json:"services,omitempty"`
}

// SubDevicesStatus 网关更新子设备状态
type SubDevicesStatus struct {
	DeviceStatuses []DeviceStatus `json:"device_statuses"`
}

// DeviceStatus 子设备状态
type DeviceStatus struct {
	DeviceId string `json:"device_id"`
	Status   string `json:"status"` // 子设备状态。 OFFLINE：设备离线 ONLINE：设备上线
}
