package context

import (
	"sailing.cn/wechat/credential"
	"sailing.cn/wechat/work/config"
)

// Context struct
type Context struct {
	*config.Config
	credential.AccessTokenHandle
}
