package emqx

type Code int

func (c *Code) Format() string {
	switch *c {
	case 0:
		return "成功"
	case 101:
		return "RPC 错误"
	case 102:
		return "未知错误"
	case 103:
		return "用户名或密码错误"
	case 104:
		return "空用户名或密码"
	case 105:
		return "用户不存在"
	case 106:
		return "管理员账户不可删除"
	case 107:
		return "关键请求参数缺失"
	case 108:
		return "请求参数错误"
	case 109:
		return "请求参数不是合法 JSON 格式"
	case 110:
		return "插件已开启"
	case 111:
		return "插件已关闭"
	case 112:
		return "客户端不在线"
	case 113:
		return "用户已存在"
	case 114:
		return "旧密码错误"
	case 115:
		return "不合法的主题"

	}
	return ""
}

type Result struct {
	Data interface{} `json:"data"`
	Code Code        `json:"code"`
}
type Resource struct {
	Id     string                 `json:"id"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
	//Status      []interface{}          `json:"status"`
	Description string `json:"description"`
}
type Rule struct {
	Id          string    `json:"id"`
	Actions     []*Action `json:"actions"`
	Rawsql      string    `json:"rawsql"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
}
type Action struct {
	Name   string                 `json:"name"`
	Params map[string]interface{} `json:"params"`
}
