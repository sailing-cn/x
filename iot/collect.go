package iot

// CollectConfig 采集配置
type CollectConfig struct {
	Collects []Collect `json:"collects"` //采集配置集合
	Version  string    `json:"version"`  //采集配置版本
}

// Collect 设备采集
type Collect struct {
	Name    string   `json:"name"`    //采集名称
	Enabled bool     `json:"enabled"` //是否启用
	Version string   `json:"version"` //版本
	Devices []Device `json:"devices"` //设备集合
}

// Device 采集设备
type Device struct {
	DeviceId              string   `json:"device_id" mapstructure:"device_id"`                               //设备标识
	DeviceName            string   `json:"device_name" mapstructure:"device_name"`                           //设备名称
	ProtocolType          string   `json:"protocol_type" mapstructure:"protocol_type"`                       //采集协议
	SwitchInterval        int16    `json:"switch_interval" mapstructure:"switch_interval"`                   //设备切换时间
	FirstByteWaitInterval int16    `json:"first_byte_wait_interval" mapstructure:"first_byte_wait_interval"` //首字节等待时间
	SendWaitInterval      int16    `json:"send_wait_interval" mapstructure:"send_wait_interval"`             //等待时间
	CollectInterval       int16    `json:"collect_interval" mapstructure:"collect_interval"`                 //采集周期
	UART                  *UART    `json:"uart"`                                                             //串口配置
	Network               *Network `json:"network"`                                                          //网络配置
	Model                 Model    `json:"model"`                                                            //采集模型
}

// UART 端口配置
type UART struct {
	Addr     string `json:"addr" mapstructure:"addr"`           //串口地址
	DataBit  int    `json:"data_bit" mapstructure:"data_bit"`   //数据位
	Parity   string `json:"parity" mapstructure:"parity"`       //奇偶校验 N:无 O:偶 E:奇
	StopBit  int    `json:"stop_bit" mapstructure:"stop_bit"`   //停止位
	BaudRate int    `json:"baud_rate" mapstructure:"baud_rate"` //波特率
}

// Network 网络配置
type Network struct {
	IP   string `json:"ip"`   //IP地址
	Port int    `json:"port"` //端口号
}

// Model 采集模型
type Model struct {
	Groups []Group `json:"groups"` //分组集合
}

type Group struct {
	Name     string   `json:"name" mapstructure:"name"`           //分组名称
	Desc     string   `json:"desc" mapstructure:"desc"`           //分组简介
	FuncCode int      `json:"func_code" mapstructure:"func_code"` //功能码
	Addr     byte     `json:"addr" mapstructure:"addr"`           //寄存器地址
	Quantity uint16   `json:"quantity" mapstructure:"quantity"`   //数据长度
	Targets  []Target `json:"targets" mapstructure:"targets"`     //指标集合
}

// Target 采集指标
type Target struct {
	Name   string  `json:"name"`   //指标名称 必须唯一
	Desc   string  `json:"desc"`   //指标描述
	Type   string  `json:"type"`   //数据类型
	Unit   string  `json:"unit"`   //计量单位
	Scale  float32 `json:"scale"`  //缩放倍数
	Addr   uint16  `json:"addr"`   //数据地址
	Length uint16  `json:"length"` //数据长度
}
