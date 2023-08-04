package websocket

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Config struct {
	Websocket struct {
		Addr string `yaml:"addr" mapstructure:"addr"`
		Port int    `yaml:"port" mapstructure:"port"`
	} `json:"websocket" yaml:"websocket" mapstructure:"websocket"`
}

// Client 单个客户端信息
type Client struct {
	Id, Group string
	socket    *websocket.Conn
	message   chan []byte
	errCount  int
}

// Manager 所有客户端信息
type Manager struct {
	group                   map[string]map[string]*Client
	groupCount, clientCount uint
	lock                    sync.Mutex
	Register, UnRegister    chan *Client
	message                 chan *MessageData
	groupMessage            chan *GroupMessageData
	broadcastMessage        chan *BroadcastMessageData
	handle                  func(c *Client, message []byte)
}

// MessageData 单个发送数据信息
type MessageData struct {
	Id, Group string
	Message   []byte
}

// GroupMessageData 群组消息数据信息
type GroupMessageData struct {
	Group   string
	Message []byte
}

// BroadcastMessageData 广播消息数据信息
type BroadcastMessageData struct {
	Message []byte
}
