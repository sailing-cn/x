package restful

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http"
	"sailing.cn/global"
	"time"
)

type engine struct {
	*gin.Engine
	cfg *global.WebapiConfig
}

var instance = &engine{}

func NewGinEngin(cfg *global.WebapiConfig) *engine {
	instance = &engine{cfg: cfg, Engine: gin.Default()}
	return instance
}
func StartRestful(cfg *global.WebapiConfig) {
	if cfg == nil {
		panic("webapi配置不能为空")
	}
	NewGinEngin(cfg).
		UserCors().
		UseSwagger("/swagger/*any").
		Run()
}

// Deprecated: 请使用AddRouter()
func (e *engine) Register(group string, handler func(g *gin.RouterGroup)) *engine {
	var g *gin.RouterGroup
	if len(group) > 0 {
		g = e.Group(group)
	} else {
		g = e.Group("")
	}
	handler(g)
	return e
}

func (e *engine) UserCors() *engine {
	e.Use(cors())
	return e
}

func (e *engine) UseSwagger(path string) *engine {
	e.GET(path, ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.DocExpansion("none")))
	return e
}
func (e *engine) WithTrace() *engine {
	e.Use(otelgin.Middleware(e.cfg.Webapi.Name + " " + e.cfg.Webapi.Version))
	return e
}

func (e *engine) AddRouter(handler func(g *gin.RouterGroup)) *engine {
	var group *gin.RouterGroup
	group = e.Group(e.cfg.Webapi.Prefix)
	handler(group)
	return e
}

func (e *engine) Run() *engine {
	addr := e.cfg.Webapi.Addr
	if len(addr) == 0 {
		addr = "0.0.0.0"
	}
	addr = fmt.Sprintf("%s:%d", addr, e.cfg.Webapi.Port)
	s := endless.NewServer(addr, e)
	s.ReadHeaderTimeout = 20 * time.Second
	s.WriteTimeout = 20 * time.Second
	s.MaxHeaderBytes = 1 << 20
	log.Infof("启动 %s 服务 %s", e.cfg.Webapi.Name, addr)
	s.ListenAndServe()
	return e
}

func (e *engine) UseLog() *engine {
	e.Use(Logger())
	return e
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
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
	}
}
