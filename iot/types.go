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
