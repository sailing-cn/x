package restful

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
	"net/http"
	"sailing.cn/v2/conf"
)

// JSONResult 标准json结果
type JSONResult struct {
	Code int         `json:"code"`    //状态码
	Msg  string      `json:"message"` //描述信息
	Data interface{} `json:"data"`    //数据
}

func Success(c *gin.Context, data interface{}) {
	r := JSONResult{
		Code: 1,
		Msg:  "success",
		Data: data,
	}
	c.JSON(http.StatusOK, r)
}

func Fail(c *gin.Context, err error) {
	msg := "服务器错误"
	if conf.WebapiConf != nil && conf.WebapiConf.Webapi.Mode == "dev" {
		log.Error(err)
		msg = err.Error()
	} else if e, ok := status.FromError(err); ok {
		if e.Code() == 600 {
			msg = e.Message()
		}
		log.Error(e)
	} else {
		log.Error(err)
	}
	r := JSONResult{0, msg, nil}
	c.JSON(http.StatusOK, r)
}
