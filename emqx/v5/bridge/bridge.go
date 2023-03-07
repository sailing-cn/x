package bridge

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	sdk "sailing.cn/emqx/v5/http"
	"sailing.cn/wrong"
	"strings"
)

const bridgeURL = "/bridges"

type BridgeService sdk.Service

// Bridge  目前均是 WebHook形式
type Bridge struct {
	Connector   string `json:"connector"`
	Direction   string `json:"direction"`
	Payload     string `json:"payload"`
	RemoteQos   string `json:"remote_qos"`
	RemoteTopic string `json:"remote_topic"`
	Retain      bool   `json:"retain"`

	Body             string `json:"body"`              //请求体
	ConnectTimeout   string `json:"connect_timeout"`   //连接超时
	Enable           bool   `json:"enable"`            //启用
	EnablePipelining int    `json:"enable_pipelining"` //启用管道
	LocalTopic       string `json:"local_topic"`       //本地主题
	MaxRetries       int    `json:"max_retries"`       //错误重试次数
	Method           string `json:"method"`            //请求类型
	Name             string `json:"name"`              //名称
	PoolSize         int    `json:"pool_size"`         //http 管道
	PoolType         string `json:"pool_type"`         //http 管道大小
	RequestTimeout   string `json:"request_timeout"`   //请求超时 单位: s
	Ssl              struct {
		Enable bool `json:"enable"` //启用
	} `json:"ssl"` //tls配置
	Type string `json:"type"` //桥接类型
	URL  string `json:"url"`  //桥接地址

	Metrics     sdk.Metrics `json:"metrics"`
	NodeMetrics []struct {
		Metrics sdk.Metrics `json:"metrics"`
		Node    string      `json:"node"`
	} `json:"node_metrics"`
}

// HttpBridge webhook桥接
type HttpBridge struct {
	Body             string `json:"body"`              //请求体
	ConnectTimeout   string `json:"connect_timeout"`   //连接超时
	Enable           bool   `json:"enable"`            //启用
	EnablePipelining int    `json:"enable_pipelining"` //启用管道
	LocalTopic       string `json:"local_topic"`       //本地主题
	MaxRetries       int    `json:"max_retries"`       //错误重试次数
	Method           string `json:"method"`            //请求类型
	Name             string `json:"name"`              //名称
	PoolSize         int    `json:"pool_size"`         //http 管道
	PoolType         string `json:"pool_type"`         //http 管道大小
	RequestTimeout   string `json:"request_timeout"`   //请求超时 单位: s
	Ssl              struct {
		Enable bool `json:"enable"` //启用
	} `json:"ssl"` //tls配置
	Type string `json:"type"` //桥接类型
	URL  string `json:"url"`  //桥接地址
}
type UpdateHttpBridge struct {
	Body             string `json:"body"`              //请求体
	ConnectTimeout   string `json:"connect_timeout"`   //连接超时
	Enable           bool   `json:"enable"`            //启用
	EnablePipelining int    `json:"enable_pipelining"` //启用管道
	LocalTopic       string `json:"local_topic"`       //本地主题
	MaxRetries       int    `json:"max_retries"`       //错误重试次数
	Method           string `json:"method"`            //请求类型
	//Name             string `json:"name"`              //名称
	PoolSize       int    `json:"pool_size"`       //http 管道
	PoolType       string `json:"pool_type"`       //http 管道大小
	RequestTimeout string `json:"request_timeout"` //请求超时 单位: s
	Ssl            struct {
		Enable bool `json:"enable"` //启用
	} `json:"ssl"` //tls配置
	//Type string `json:"type"` //桥接类型
	URL string `json:"url"` //桥接地址
}

// CreateBridge 创建数据桥接
func (s *BridgeService) CreateBridge(model *HttpBridge) (interface{}, error) {
	result := new(HttpBridge)
	resp, err := s.Client.R().
		SetBody(*model).
		SetResult(result).
		Post(s.Client.RequestURL(bridgeURL))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		// 删除当前创建的规则
		err := s.DeleteRule(strings.ToLower(model.Type) + ":" + model.Name)
		if err != nil {
			log.Error("删除数据桥接失败")
		}
		return nil, errors.New(fmt.Sprintf("服务端返回状态码:%d", resp.StatusCode))
	}
	return model, nil
}

func (s *BridgeService) ListBridge() (list []Bridge, err error) {
	resp, err := s.Client.R().SetResult(&list).Get(s.Client.RequestURL(bridgeURL))
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.String())
	}
	return list, err
}

// GetBridge bridgeId 示例  type : id  == mqtt:mqtt_example || webhook :123456
func (s *BridgeService) GetBridge(bridgeId string) (bridge *Bridge, err error) {
	result := &Bridge{}
	r := s.Client.R().SetResult(result)
	resp, err := r.Get(s.Client.RequestURL(bridgeURL + "/" + bridgeId))
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return bridge, errors.New(fmt.Sprintf("服务端返回状态码:%d", resp.StatusCode))
	}
	bridge = result
	return
}

// Update 更新资源  bridgeId 示例  type : id  == mqtt:mqtt_example || webhook :123456
func (s *BridgeService) Update(bridgeId string, bridge *UpdateHttpBridge) (interface{}, error) {
	model, err := s.GetBridge(bridgeId)
	if err != nil {
		return nil, wrong.New("当前动作未找到，请重试")
	}
	bridge.Body = model.Body
	bridge.ConnectTimeout = model.ConnectTimeout
	bridge.EnablePipelining = model.EnablePipelining
	bridge.LocalTopic = model.LocalTopic
	bridge.MaxRetries = model.MaxRetries
	bridge.Method = model.Method
	bridge.PoolSize = model.PoolSize
	bridge.PoolType = model.PoolType
	bridge.RequestTimeout = model.RequestTimeout
	bridge.Ssl = model.Ssl
	xxx, err := json.Marshal(bridge)
	log.Print(xxx)
	var result interface{}
	resp, err := s.Client.R().SetBody(bridge).SetResult(&result).Put(s.Client.RequestURL(bridgeURL + "/" + bridgeId))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("请求错误:code%d msg::%s", resp.StatusCode, resp.String()))
	}
	return result, nil
}

// DeleteRule 删除资源  bridgeId 示例  type : id  == mqtt:mqtt_example || webhook :123456
func (s *BridgeService) DeleteRule(bridgeId string) error {
	resp, err := s.Client.R().Delete(s.Client.RequestURL(bridgeURL + "/" + bridgeId))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.New("EMQX删除规则出现错误:" + resp.String())
	}
	return nil
}
