package v5

import (
	log "github.com/sirupsen/logrus"
	"sailing.cn/v2/emqx/v5/bridge"
	"sailing.cn/v2/emqx/v5/client"
	"sailing.cn/v2/emqx/v5/config"
	"sailing.cn/v2/emqx/v5/http"
	"sailing.cn/v2/emqx/v5/rule"
)

type EmqxClient struct {
	RuleService   *rule.RuleService
	BridgeService *bridge.BridgeService
	ClientService *client.ClientService
}

//api key: 90977a6223e6a1be
//api secret: i0Ql4JhKpUhevTPg9BNa2YIiZ9BBMtQRwDiSXgvw159AMH

//tencent apikey: ac84a121a5677f0f
//tencent apisecret: oEuxKj8zR2JAXfdI9AHKqDsxVyGva029CggIMSm9BOlWZJ

func NewEmqxClient() *EmqxClient {
	cfg := config.NewConfiguration()
	err := cfg.Init()
	if err != nil {
		log.Errorf("初始化emax client 错误:%s", err)
		return nil
	}
	return NewEmqxClientWithCnf(cfg)
}

func NewEmqxClientWithCnf(cfg *config.Configuration) *EmqxClient {
	e := &EmqxClient{}
	_client := http.NewAPIClient(cfg)
	e.RuleService = &rule.RuleService{Client: _client}
	e.BridgeService = &bridge.BridgeService{Client: _client}
	e.ClientService = &client.ClientService{Client: _client}
	return e
}
