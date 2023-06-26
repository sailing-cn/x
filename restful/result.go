package restful

// JSONResult 标准json结果
type JSONResult struct {
	Code int         `json:"code"`    //状态码
	Msg  string      `json:"message"` //描述信息
	Data interface{} `json:"data"`    //数据
}
