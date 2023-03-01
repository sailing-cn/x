package restful

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
	"net/http"
	"sailing.cn/global"
)

type JSONResult struct {
	Code int         `json:"code"`    //状态码
	Msg  string      `json:"message"` //描述信息
	Data interface{} `json:"data"`    //数据
}

type PagerResultStruct struct {
	Data  []interface{} `json:"data"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Count int           `json:"count"`
	Size  int           `json:"size"`
}

//Success 成功
func Success(response *gin.Context, data interface{}) {
	r := JSONResult{1, "success", data}
	response.JSON(http.StatusOK, r)
}

//Fail 失败
func Fail(response *gin.Context, error error) {
	msg := "服务器错误"
	if global.WebapiConf != nil && global.WebapiConf.Webapi.Mode == "dev" {
		log.Error(error)
		msg = error.Error()
	} else if e, ok := status.FromError(error); ok {
		if e.Code() == 600 {
			msg = e.Message()
		}
		log.Error(e)
	} else {
		log.Error(error)
	}
	r := JSONResult{0, msg, nil}
	response.JSON(http.StatusOK, r)
}
