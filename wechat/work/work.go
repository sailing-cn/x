package work

import (
	"sailing.cn/wechat/credential"
	"sailing.cn/wechat/work/addresslist"
	"sailing.cn/wechat/work/config"
	"sailing.cn/wechat/work/context"
	"sailing.cn/wechat/work/externalcontact"
	"sailing.cn/wechat/work/kf"
	"sailing.cn/wechat/work/material"
	"sailing.cn/wechat/work/msgaudit"
	"sailing.cn/wechat/work/oauth"
	"sailing.cn/wechat/work/robot"
)

type Work struct {
	ctx *context.Context
}

func NewWork(cfg *config.Config) *Work {
	defaultTokenHandle := credential.NewWorkAccessToken(cfg.CorpID, cfg.CorpSecret, credential.CacheKeyWorkPrefix, cfg.Cache)
	ctx := &context.Context{
		Config:            cfg,
		AccessTokenHandle: defaultTokenHandle,
	}
	return &Work{ctx: ctx}
}

// GetContext get Context
func (wk *Work) GetContext() *context.Context {
	return wk.ctx
}

// GetOauth get oauth
func (wk *Work) GetOauth() *oauth.Oauth {
	return oauth.NewOauth(wk.ctx)
}

// GetMsgAudit get msgAudit
func (wk *Work) GetMsgAudit() (*msgaudit.Client, error) {
	return msgaudit.NewClient(wk.ctx.Config)
}

// GetKF get kf
func (wk *Work) GetKF() (*kf.Client, error) {
	return kf.NewClient(wk.ctx.Config)
}

// GetExternalContact get external_contact
func (wk *Work) GetExternalContact() *externalcontact.Client {
	return externalcontact.NewClient(wk.ctx)
}

// GetAddressList get address_list
func (wk *Work) GetAddressList() *addresslist.Client {
	return addresslist.NewClient(wk.ctx)
}

// GetMaterial get material
func (wk *Work) GetMaterial() *material.Client {
	return material.NewClient(wk.ctx)
}

// GetRobot get robot
func (wk *Work) GetRobot() *robot.Client {
	return robot.NewClient(wk.ctx)
}
