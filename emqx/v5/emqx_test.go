package v5

import (
	"sailing.cn/emqx/v5/rule"
	"testing"
)

//var (
//	local = &config.Configuration{
//		Server:     "http://192.168.0.200:18083",
//		ApiSecret:  "i0Ql4JhKpUhevTPg9BNa2YIiZ9BBMtQRwDiSXgvw159AMH",
//		ApiKey:     "90977a6223e6a1be",
//		ApiVersion: "v5",
//	}
//	tencent = &config.Configuration{
//		Server:     "http://101.33.245.172:28083",
//		ApiSecret:  "oEuxKj8zR2JAXfdI9AHKqDsxVyGva029CggIMSm9BOlWZJ",
//		ApiKey:     "ac84a121a5677f0f",
//		ApiVersion: "v5",
//	}
//)

func TestEmqxClient(t *testing.T) {
	cfg := local
	emqx := NewEmqxClientWithCnf(cfg)
	result, _ := emqx.RuleService.ListRule(&rule.RuleQuery{})
	t.Logf("%v", result)
}

func TestBridges(t *testing.T) {
	cfg := local
	emqx := NewEmqxClientWithCnf(cfg)
	/*
		model := &bridge.HttpBridge{
			Body:             "body",
			ConnectTimeout:   "10s",
			Enable:           true,
			EnablePipelining: 10,
			LocalTopic:       "$sailing",
			MaxRetries:       11,
			Method:           "post",
			Name:             "second",
			PoolSize:         4,
			PoolType:         "random",
			RequestTimeout:   "20s",
			Ssl: struct {
				Enable bool `json:"enable"`
			}{},
			Type: "webhook",
			URL:  "http://localhost:8000",
		}
	*/
	result, err := emqx.BridgeService.ListBridge()
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Logf("%v", result)
}
