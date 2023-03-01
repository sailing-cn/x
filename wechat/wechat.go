package wechat

import (
	"sailing.cn/wechat/work"
	"sailing.cn/wechat/work/config"
)

type Wechat struct {
}

func NewWechat() *Wechat {
	return &Wechat{}
}

func (w *Wechat) GetWork(cfg *config.Config) *work.Work {
	return work.NewWork(cfg)
}
