package iot

import "time"

// UpgradeResult 更新结果
type UpgradeResult int

// EventType 事件类型
type EventType string

// Service 服务
type Service struct {
	ServiceId   string     `json:"service_id" mapstructure:"ServiceId"`     //服务标识
	ServiceType string     `json:"service_type" mapstructure:"ServiceType"` //服务类型
	Description string     `json:"description" mapstructure:"Description"`  //服务描述
	Properties  []Property `json:"properties" mapstructure:"Properties"`    //属性列表
	Commands    []CommandX `json:"commands" mapstructure:"Commands"`        //命令列表
}

// Property 属性
type Property struct {
	Name     string `json:"name" mapstructure:"Name"`          //属性名称
	Required bool   `json:"required" mapstructure:"Required"`  //是否必须
	DataType string `json:"data_type" mapstructure:"DataType"` //数据类型
	//Method       string      `json:"method" mapstructure:"Method"`              //访问权限
	Method       Method      `json:"method" mapstructure:"Method"`              //访问权限
	Max          float32     `json:"max" mapstructure:"Max"`                    //最大值
	MaxLength    int32       `json:"max_length" mapstructure:"MaxLength"`       //最大长度
	Min          float32     `json:"min" mapstructure:"Min"`                    //最小长度
	Step         float32     `json:"step" mapstructure:"Step"`                  //步长
	Unit         string      `json:"unit" mapstructure:"Unit"`                  //单位
	Enums        []string    `json:"enums" mapstructure:"Enums"`                //枚举值
	Description  string      `json:"description" mapstructure:"Description"`    //描述信息
	DefaultValue interface{} `json:"default_value" mapstructure:"DefaultValue"` //默认值
}

// CommandX 命令
type CommandX struct {
	Name      string   `json:"name"`
	Params    []Param  `json:"params"`
	Responses Response `json:"responses"`
}

// Response 响应
type Response struct {
	Name   string  `json:"name"`
	Params []Param `json:"params"`
}

// Param 参数
type Param struct {
	Name        string   `json:"name"`                         //属性名称
	Required    bool     `json:"required"`                     //是否必须
	DataType    string   `json:"data_type"`                    //数据类型
	Method      Method   `json:"method" mapstructure:"Method"` //访问权限
	Max         float32  `json:"max"`                          //最大值
	MaxLength   int      `json:"max_length"`                   //最大长度
	Min         float32  `json:"min"`                          //最小长度
	Step        float32  `json:"step"`                         //步长
	Enums       []string `json:"enums"`                        //枚举值
	Unit        string   `json:"unit"`                         //单位
	Description string   `json:"description"`                  //描述信息
}

type Desired struct {
	EventTime  string                 `json:"event_time"`
	Properties map[string]interface{} `json:"properties"`
}

type Reported struct {
	EventTime  string                 `json:"event_time"`
	Properties map[string]interface{} `json:"properties"`
}

type SupportSourceVersion struct {
	SoftwareVersion string `json:"software_version"`
	FirmwareVersion string `json:"firmware_version"`
}
type SupportDevice struct {
	DeviceType       string `json:"device_type"`
	ManufacturerName string `json:"manufacturer_name"`
	ProtocolType     string `json:"protocol_type"`
}

type TaskFilter struct {
	GroupIdList  []string `json:"group_id_list"`
	DeviceIdList []string `json:"device_id_list"`
}

type TaskPolicy struct {
	//ScheduleTime  *time.Time `json:"schedule_time"`
	ScheduleTime  int64 `json:"schedule_time"`
	RetryCount    int32 `json:"retry_count"`
	RetryInterval int32 `json:"retry_interval"`
}

type TaskProgress struct {
	Total         int32 `json:"total"`
	Processing    int32 `json:"processing"`
	Success       int32 `json:"success"`
	Fail          int32 `json:"fail"`
	Waiting       int32 `json:"waiting"`
	FailWaitRetry int32 `json:"fail_wait_retry"`
	Stopped       int32 `json:"stopped"`
}
type Method struct {
	Write bool `json:"write"`
	Read  bool `json:"read"`
}

// DevicesService 网关批量上报子设备属性
type DevicesService struct {
	Devices []*DeviceService `json:"devices"` //设备集合
}

type DeviceService struct {
	DeviceId string                 `json:"device_id"` //设备标识
	Services []DevicePropertyEntity `json:"services"`  //服务集合
}

// DeviceProperties 设备属性
type DeviceProperties struct {
	Services []DevicePropertyEntity `json:"services"` //服务集合
}

// DevicePropertyEntity 设备的一个属性
type DevicePropertyEntity struct {
	ServiceId  string                 `json:"service_id"` //服务标识
	Properties map[string]interface{} `json:"properties"` //属性集合
	EventTime  time.Time              `json:"event_time"` //时间
}

// ReportDeviceInfoPara 上报设备信息参数
type ReportDeviceInfoPara struct {
	SoftwareVersion string `json:"software_version"` //软件版本号
	FirmwareVersion string `json:"firmware_version"` //固件版本号
}

// Command 设备命令
type Command struct {
	DeviceId    string      `json:"device_id"`
	ServiceId   string      `json:"service_id"`
	CommandName string      `json:"command_name"`
	Paras       interface{} `json:"paras"`
	CommandId   string      `json:"command_id"`
}

// CommandResponse 命令响应
type CommandResponse struct {
	ResultCode   byte        `json:"result_code"`   //结果
	ResponseName string      `json:"response_name"` //响应名称
	Paras        interface{} `json:"paras"`         //参数
	RequestId    string      `json:"request_id"`    //请求标识
}

// UpgradeProgress 设备升级状态响应，用于设备向平台反馈进度，错误信息等
type UpgradeProgress struct {
	ResultCode  UpgradeResult `json:"result_code"` //结果
	Progress    int           `json:"progress"`    // 设备的升级进度，范围：0到100
	Version     string        `json:"version"`     // 设备当前版本号
	Description string        `json:"description"` // 升级状态描述信息，可以返回具体升级失败原因。
}

// UpgradePara 更新参数
type UpgradePara struct {
	Version  string `json:"version"`
	URL      string `json:"url"`
	FileSize int64  `json:"file_size"`
	Token    string `json:"token"`
	Expires  int64  `json:"expires"`
	Sign     string `json:"sign"`
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

// VersionInfo 版本信息
type VersionInfo struct {
	SoftwareVersion string `json:"software_version"` //软件版本
	FirmwareVersion string `json:"firmware_version"` //固件版本
	CollectVersion  string `json:"collect_version"`  //采集版本
	ServiceVersion  string `json:"service_version"`  //服务版本
}

// DeviceOnlineInfo 设备上线信息
type DeviceOnlineInfo struct {
	Username       string `json:"username"`        //用户名
	Ts             int64  `json:"ts"`              //时间
	SocketPort     int    `json:"socket_port"`     //端口号
	ProtoVer       int    `json:"proto_ver"`       //协议版本
	ProtoName      string `json:"proto_name"`      //协议名称
	Keepalive      int    `json:"keepalive"`       //保持连接
	IPAddress      string `json:"ip_address"`      //IP地址
	ExpiryInterval int    `json:"expiry_interval"` //过期时间
	ConnectedAt    int64  `json:"connected_at"`    //
	ConnACK        int    `json:"conn_ack"`        //连接应答
	ClientId       string `json:"client_id"`       //客户端标识
	CleanStart     bool   `json:"clean_start"`     //
}

// DeviceOfflineInfo 设备离线信息
type DeviceOfflineInfo struct {
	Username       string `json:"username"`        //用户名
	Ts             int64  `json:"ts"`              //时间
	SocketPort     int    `json:"socket_port"`     //端口号
	Reason         string `json:"reason"`          //离线原因
	ProtoVer       int    `json:"proto_ver"`       //协议版本
	ProtoName      string `json:"proto_name"`      //协议名称
	IPAddress      string `json:"ip_address"`      //IP地址
	DisconnectedAt int64  `json:"disconnected_at"` //
	ClientId       string `json:"client_id"`       //客户端标识
}

// 定义平台和设备之间的数据交换结构体

// EventData 事件
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

// SubDevicesStatus 网关更新子设备状态
type SubDevicesStatus struct {
	DeviceStatuses []SubDeviceStatus `json:"device_statuses"`
}

// SubDeviceStatus 子设备状态
type SubDeviceStatus struct {
	DeviceId string `json:"device_id"`
	Status   string `json:"status"` // 子设备状态。 OFFLINE：设备离线 ONLINE：设备上线
}

// DevicePropertyQueryRequest 设备属性查询请求
type DevicePropertyQueryRequest struct {
	DeviceId  string `json:"device_id"`
	ServiceId string `json:"service_id"`
}

// DevicePropertyQueryResponse 设备属性查询响应
type DevicePropertyQueryResponse struct {
	ObjectDeviceId string             `json:"object_device_id"`
	Shadow         []DeviceShadowData `json:"shadow"`
}

// DeviceShadowData 设备影子数据
type DeviceShadowData struct {
	ServiceId string                     `json:"service_id"`
	Desired   DeviceShadowPropertiesData `json:"desired"`
	Reported  DeviceShadowPropertiesData `json:"reported"`
	Version   int                        `json:"version"`
}

// DeviceShadowPropertiesData 设备影子属性数据
type DeviceShadowPropertiesData struct {
	Properties interface{} `json:"properties"`
	EventTime  string      `json:"event_time"`
}

// DevicePropertyDownRequest 平台设置设备属性
type DevicePropertyDownRequest struct {
	DeviceId string                            `json:"device_id"`
	Services []DevicePropertyDownRequestEntity `json:"services"`
}

// DevicePropertyDownRequestEntity 设备属性设置请求详情
type DevicePropertyDownRequestEntity struct {
	ServiceId  string      `json:"service_id"`
	Properties interface{} `json:"properties"`
}

// AddSubDeviceParas 子设备添加事件参数
type AddSubDeviceParas struct {
	Devices []SubDeviceInfo `json:"devices"`
	Version int64           `json:"version"`
}

// DeleteSubDeviceParas 删除子设备事件参数
type DeleteSubDeviceParas struct {
	Devices []SubDeviceInfo `json:"devices"`
	Version int64           `json:"version"`
}

// DeviceLog 设备日志
type DeviceLog struct {
	Timestamp string `json:"timestamp"` // 日志产生时间
	Type      string `json:"type"`      // 日志类型：DEVICE_STATUS，DEVICE_PROPERTY ，DEVICE_MESSAGE ，DEVICE_COMMAND
	Content   string `json:"content"`   // 日志内容
}

// BaseServiceEvent 服务事件--基础信息
type BaseServiceEvent struct {
	ServiceId string `json:"service_id"`
	EventType string `json:"event_type"`
	EventTime int64  `json:"event_time,omitempty"`
}

type FileRequestServiceEvent struct {
	BaseServiceEvent
	Paras FileRequestServiceEventParas `json:"paras"`
}

// FileResultResponse 文件操作响应结果
type FileResultResponse struct {
	ObjectDeviceId string                           `json:"object_device_id"`
	Services       []FileResultResponseServiceEvent `json:"services"`
}

// FileRequestServiceEventParas 设备获取文件上传下载URL参数
type FileRequestServiceEventParas struct {
	FileName       string      `json:"file_name"`
	FileAttributes interface{} `json:"file_attributes"`
}

// FileResultResponseServiceEvent 文件操作事件响应结果
type FileResultResponseServiceEvent struct {
	BaseServiceEvent
	Paras FileResultServiceEventParas `json:"paras"`
}

// FileResultServiceEventParas 上报文件上传下载结果参数
type FileResultServiceEventParas struct {
	ObjectName        string `json:"object_name"`
	ResultCode        int    `json:"result_code"`
	StatusCode        int    `json:"status_code"`
	StatusDescription string `json:"status_description"`
}

// FileResponseServiceEventParas 平台下发响应参数
type FileResponseServiceEventParas struct {
	Url            string      `json:"url"`
	BucketName     string      `json:"bucket_name"`
	ObjectName     string      `json:"object_name"`
	Expire         int         `json:"expire"`
	FileAttributes interface{} `json:"file_attributes"`
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
type FileResponseServiceEvent struct {
	BaseServiceEvent
	Paras FileResponseServiceEventParas `json:"paras"`
}

// ReportDeviceLogServiceEvent 设备日志上报事件
type ReportDeviceLogServiceEvent struct {
	BaseServiceEvent
	Paras DeviceLog `json:"paras,omitempty"`
}

// ReportDeviceLogRequest 上报设备日志请求
type ReportDeviceLogRequest struct {
	Services []ReportDeviceLogServiceEvent `json:"services,omitempty"`
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
