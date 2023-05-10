package iot

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
