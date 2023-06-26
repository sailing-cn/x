package restful

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sailing.cn/v2/conf"
)

// Engine 网关引擎
type Engine struct {
	*gin.Engine
	cfg *conf.WebapiConfig
}

// Option 网关配置项
type Option func(e *Engine) error

var instance = &Engine{}

// NewGinDefault 创建一个默认的GinEngine
func NewGinDefault(cfg *conf.WebapiConfig, opts ...Option) *Engine {
	instance = &Engine{cfg: cfg, Engine: gin.Default()}
	if cfg.Webapi.Mode != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}
	for i := range opts {
		err := opts[i](instance)
		if err != nil {
			return nil
		}
	}
	return instance
}

// WithCors 跨域配置
func WithCors() Option {
	return func(e *Engine) error {
		e.Use(func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH, HEAD")
			c.Header("Access-Control-Allow-Headers", "*")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
			c.Next()
		})
		return nil
	}
}

// Run 运行网关引擎
func (e *Engine) Run() {
	addr := e.cfg.Webapi.Addr
	if len(addr) == 0 {
		addr = "0.0.0.0"
	}
	addr = fmt.Sprintf("%s:%d", addr, e.cfg.Webapi.Port)
	s := endless.NewServer(addr, e)
	log.Infof("启动 %s 服务 %s", e.cfg.Webapi.Name, addr)
	_ = s.ListenAndServe()
}
