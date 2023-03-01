package zentao

import (
	"sailing.cn/zentao/config"
	"sailing.cn/zentao/context"
	"sailing.cn/zentao/token"
)

type ZenTao struct {
	ctx   *context.Context
	token string
}

func NewZenTao(cfg *config.Config) *ZenTao {
	tokenHandle := token.NewToken(cfg.User, cfg.Password, cfg.Domain, cfg.Cache)
	_token, err := tokenHandle.GetToken()
	if err != nil {
		return nil
	}
	ctx := &context.Context{Token: _token, Config: cfg}
	ctx.SetBaseURL(cfg.Domain)
	ctx.SetClient()
	return &ZenTao{ctx: ctx}
}
