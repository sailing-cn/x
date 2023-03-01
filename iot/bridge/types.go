package bridge

// Forward 转发
type Forward struct {
	Http struct {
		Url     string                 `json:"url"`
		Headers map[string]interface{} `json:"headers"`
	} `json:"http"`
}
