package iot

type CollectConfig struct {
	Collects []Collect `json:"collects"`
	Version  string    `json:"version"`
}

type Collect struct {
	Name    string   `json:"name"`
	Enabled bool     `json:"enabled"`
	Version string   `json:"version"`
	Devices []Device `json:"devices"`
}

type Device struct {
	DeviceId              string   `json:"device_id" mapstructure:"device_id"`
	DeviceName            string   `json:"device_name" mapstructure:"device_name"`
	ProtocolType          string   `json:"protocol_type" mapstructure:"protocol_type"`
	SwitchInterval        int16    `json:"switch_interval" mapstructure:"switch_interval"`
	FirstByteWaitInterval int16    `json:"first_byte_wait_interval" mapstructure:"first_byte_wait_interval"`
	SendWaitInterval      int16    `json:"send_wait_interval" mapstructure:"send_wait_interval"`
	CollectInterval       int16    `json:"collect_interval" mapstructure:"collect_interval"`
	UART                  *UART    `json:"uart"`
	Network               *Network `json:"network"`
	Model                 Model    `json:"model"`
}

// UART 端口配置
type UART struct {
	Addr     string `json:"addr" mapstructure:"addr"`
	DataBit  int    `json:"data_bit" mapstructure:"data_bit"`
	Parity   string `json:"parity" mapstructure:"parity"`
	StopBit  int    `json:"stop_bit" mapstructure:"stop_bit"`
	BaudRate int    `json:"baud_rate" mapstructure:"baud_rate"`
}

type Network struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type Model struct {
	Groups []Group `json:"groups"`
}

type Group struct {
	Name     string   `json:"name" mapstructure:"name"`
	Desc     string   `json:"desc" mapstructure:"desc"`
	FuncCode int      `json:"func_code" mapstructure:"func_code"`
	Addr     byte     `json:"addr" mapstructure:"addr"`
	Quantity uint16   `json:"quantity" mapstructure:"quantity"`
	Targets  []Target `json:"targets" mapstructure:"targets"`
}
type Target struct {
	Name   string  `json:"name"`
	Desc   string  `json:"desc"`
	Type   string  `json:"type"`
	Unit   string  `json:"unit"`
	Scale  float32 `json:"scale"`
	Addr   uint16  `json:"addr"`
	Length uint16  `json:"length"`
}
