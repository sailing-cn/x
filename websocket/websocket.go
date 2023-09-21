package websocket

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sailing.cn/v2/utils/id"
	"time"
)

var WebsocketManager = Manager{
	group:            make(map[string]map[string]*Client),
	groupCount:       0,
	clientCount:      0,
	Register:         make(chan *Client, 128),
	UnRegister:       make(chan *Client, 128),
	message:          make(chan *MessageData, 128),
	groupMessage:     make(chan *GroupMessageData, 128),
	broadcastMessage: make(chan *BroadcastMessageData, 128),
}

func (c *Client) Read() {
	defer func() {
		WebsocketManager.UnRegister <- c
		log.Infof("客户端 [%s] 断开连接", c.Id)
		if err := c.socket.Close(); err != nil {
			log.Errorf("客户端 [%s] 断开连接异常: %s", c.Id, err)
		}
	}()

	for {
		messageType, message, err := c.socket.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			break
		}
		//log.Printf("client [%s] 收到消息: %s", c.Id, string(message))
		if WebsocketManager.handle != nil {
			WebsocketManager.handle(c, message)
		}
		//c.message <- message
	}
}

func (c *Client) Write() {
	defer func() {
		log.Infof("客户端 [%s] 断开连接", c.Id)
		if err := c.socket.Close(); err != nil {
			log.Errorf("客户端 [%s] 断开连接异常: %s", c.Id, err)
		}
	}()

	for {
		select {
		case message, ok := <-c.message:
			if !ok {
				_ = c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Errorf("客户端 [%s] 写入消息异常:%s", c.Id, err)
				c.errCount++
			}
		}
	}
}

func (manager *Manager) Start() {
	log.Infof("websocket manager start")
	for {
		select {
		case client := <-manager.Register:
			log.Infof("客户端 [%s] 已连接", client.Id)
			log.Infof("客户端 [%s] 注册到群组 [%s]", client.Id, client.Group)
			manager.lock.Lock()
			if manager.group[client.Group] == nil {
				manager.group[client.Group] = make(map[string]*Client)
				manager.groupCount++
			}
			manager.group[client.Group][client.Id] = client
			manager.clientCount++
			manager.lock.Unlock()
		case client := <-manager.UnRegister:
			manager.lock.Lock()
			delete(manager.group[client.Group], client.Id)
			manager.clientCount--
			manager.lock.Unlock()
			log.Printf("客户端 [%s] 从群组 [%s] 中注销", client.Id, client.Group)
		}
	}
}

func (manager *Manager) Send() {
	for {
		select {
		case data := <-manager.message:
			if group, ok := manager.group[data.Group]; ok {
				if conn, ok := group[data.Id]; ok {
					conn.message <- data.Message
				}
			}
		}
	}
}

func (manager *Manager) SendGroup() {
	for {
		select {
		case data := <-manager.groupMessage:
			if group, ok := manager.group[data.Group]; ok {
				for _, client := range group {
					client.message <- data.Message
				}
			}
		}
	}
}

func (manager *Manager) Broadcast() {
	for {
		select {
		case data := <-manager.broadcastMessage:
			for _, group := range manager.group {
				for _, client := range group {
					client.message <- data.Message
				}
			}
		}
	}
}

// SendToClient Deprecated 下个版本中将放弃
func (manager *Manager) SendToClient(id string, group string, message []byte) {
	data := &MessageData{
		Id: id, Group: group, Message: message,
	}
	manager.message <- data
}

func (manager *Manager) SendToClient2(id string, group string, message interface{}) {

	data := &MessageData{Id: id, Group: group}
	if v, ok := message.([]byte); ok {
		data.Message = v
	} else if v, ok := message.(string); ok {
		data.Message = []byte(v)
	} else {
		bytes, _ := json.Marshal(message)
		data.Message = bytes
	}
	manager.message <- data
}

// SendToGroup Deprecated 下个版本中将放弃
func (manager *Manager) SendToGroup(group string, message []byte) {
	data := &GroupMessageData{
		Group: group, Message: message,
	}
	manager.groupMessage <- data
}

func (manager *Manager) SendToGroup2(group string, message interface{}) {

	data := &GroupMessageData{Group: group}
	if v, ok := message.([]byte); ok {
		data.Message = v
	} else if v, ok := message.(string); ok {
		data.Message = []byte(v)
	} else {
		bytes, _ := json.Marshal(message)
		data.Message = bytes
	}
	manager.groupMessage <- data
}

func (manager *Manager) BroadcastMessage(body interface{}) {
	m, ok := body.([]byte)
	var bytes []byte
	if ok {
		bytes = m
	} else {
		bytes, _ = json.Marshal(body)
	}
	data := &BroadcastMessageData{Message: bytes}
	manager.broadcastMessage <- data
}

func (manager *Manager) RegisterClient(client *Client) {
	manager.Register <- client
}

func (manager *Manager) UnRegisterClient(client *Client) {
	manager.UnRegister <- client
}

func (manager *Manager) GroupCount() uint {
	return manager.groupCount
}

func (manager *Manager) ClientCount() uint {
	return manager.clientCount
}

func (manager *Manager) Websocket(c *gin.Context) {
	grader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{c.GetHeader("Sec-WebSocket-Protocol")},
	}
	conn, err := grader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf("websocket 连接异常：%s", c.Param("channel"))
		return
	}
	clientId := c.Query("uid")
	if len(clientId) <= 0 {
		clientId = id.SnowflakeString()
	}
	client := &Client{
		Id:      clientId,
		Group:   c.Param("channel"),
		socket:  conn,
		message: make(chan []byte, 1024),
	}
	manager.RegisterClient(client)
	go client.Read()
	go client.Write()
	//time.Sleep(time.Second * 1)
	//manager.SendToClient(client.Id, client.Group, []byte("welcome"+utils.TimeString()))
}

func (manager *Manager) Cleanup(duration time.Duration) {
	ticker := time.NewTimer(duration)
	for range ticker.C {
		manager.lock.Lock()
		if manager.clientCount > 0 {
			for _, group := range manager.group {
				for _, client := range group {
					if client.errCount > 10 {
						manager.UnRegister <- client
					}
				}
			}
		}
		manager.lock.Unlock()
	}
}

func (manager *Manager) Shutdown() {
	if manager.clientCount > 0 {
		for _, group := range manager.group {
			for _, client := range group {
				if client.errCount > 10 {
					manager.UnRegister <- client
				}
			}
		}
	}
}

// SetHandle 设置消息处理
func (manager *Manager) SetHandle(handle func(c *Client, message []byte)) {
	manager.handle = handle
}

func (manager *Manager) FindClient(clientId string) *Client {
	for _, group := range manager.group {
		for _, client := range group {
			if client.Id == clientId {
				return client
			}
		}
	}
	return nil
}
