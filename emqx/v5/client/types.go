package client

import (
	"sailing.cn/v2/emqx"
	"time"
)

type ClientResponse struct {

	// v4 api name [awaiting_rel] Number of awaiting PUBREC packet
	AwaitingRelCnt int32 `json:"awaiting_rel_cnt,omitempty"`

	// v4 api name [max_awaiting_rel]. Maximum allowed number of awaiting PUBREC packet
	AwaitingRelMax int32 `json:"awaiting_rel_max,omitempty"`

	// Indicate whether the client is using a brand new session
	CleanStart bool `json:"clean_start,omitempty"`

	// Client identifier
	Clientid string `json:"clientid,omitempty"`

	// Whether the client is connected
	Connected bool `json:"connected,omitempty"`

	// Client connection time, rfc3339 or timestamp(millisecond)
	ConnectedAt time.Time `json:"connected_at,omitempty"`

	// Session creation time, rfc3339 or timestamp(millisecond)
	CreatedAt time.Time `json:"created_at,omitempty"`

	// Client offline time. It's Only valid and returned when connected is false, rfc3339 or timestamp(millisecond)
	DisconnectedAt time.Time `json:"disconnected_at,omitempty"`

	// Session expiration interval, with the unit of second
	ExpiryInterval int32 `json:"expiry_interval,omitempty"`

	// Process heap size with the unit of byte
	HeapSize int32 `json:"heap_size,omitempty"`

	// Current length of inflight
	InflightCnt int32 `json:"inflight_cnt,omitempty"`

	// v4 api name [max_inflight]. Maximum length of inflight
	InflightMax int32 `json:"inflight_max,omitempty"`

	// Client's IP address
	IpAddress string `json:"ip_address,omitempty"`

	// Indicates whether the client is connectedvia bridge
	IsBridge bool `json:"is_bridge,omitempty"`

	// keepalive time, with the unit of second
	Keepalive int32 `json:"keepalive,omitempty"`

	// Process mailbox size
	MailboxLen int32 `json:"mailbox_len,omitempty"`

	// Number of messages dropped by the message queue due to exceeding the length
	MqueueDropped int32 `json:"mqueue_dropped,omitempty"`

	// Current length of message queue
	MqueueLen int32 `json:"mqueue_len,omitempty"`

	// v4 api name [max_mqueue]. Maximum length of message queue
	MqueueMax int32 `json:"mqueue_max,omitempty"`

	// Name of the node to which the client is connected
	Node string `json:"node,omitempty"`

	// Client's port
	Port int32 `json:"port,omitempty"`

	// Client protocol name
	ProtoName string `json:"proto_name,omitempty"`

	// Protocol version used by the client
	ProtoVer int32 `json:"proto_ver,omitempty"`

	// Number of TCP packets received
	RecvCnt int32 `json:"recv_cnt,omitempty"`

	// Number of PUBLISH packets received
	RecvMsg int32 `json:"recv_msg,omitempty"`

	// Number of dropped PUBLISH packets
	RecvMsgDropped int32 `json:"recv_msg.dropped,omitempty"`

	// Number of dropped PUBLISH packets due to expired
	RecvMsgDroppedAwaitPubrelTimeout int32 `json:"recv_msg.dropped.await_pubrel_timeout,omitempty"`

	// Number of PUBLISH QoS0 packets received
	RecvMsgQos0 int32 `json:"recv_msg.qos0,omitempty"`

	// Number of PUBLISH QoS1 packets received
	RecvMsgQos1 int32 `json:"recv_msg.qos1,omitempty"`

	// Number of PUBLISH QoS2 packets received
	RecvMsgQos2 int32 `json:"recv_msg.qos2,omitempty"`

	// Number of bytes received
	RecvOct int32 `json:"recv_oct,omitempty"`

	// Number of MQTT packets received
	RecvPkt int32 `json:"recv_pkt,omitempty"`

	// Erlang reduction
	Reductions int32 `json:"reductions,omitempty"`

	// Number of TCP packets sent
	SendCnt int32 `json:"send_cnt,omitempty"`

	// Number of PUBLISH packets sent
	SendMsg int32 `json:"send_msg,omitempty"`

	// Number of dropped PUBLISH packets
	SendMsgDropped int32 `json:"send_msg.dropped,omitempty"`

	// Number of dropped PUBLISH packets due to expired
	SendMsgDroppedExpired int32 `json:"send_msg.dropped.expired,omitempty"`

	// Number of dropped PUBLISH packets due to queue full
	SendMsgDroppedQueueFull int32 `json:"send_msg.dropped.queue_full,omitempty"`

	// Number of dropped PUBLISH packets due to packet length too large
	SendMsgDroppedTooLarge int32 `json:"send_msg.dropped.too_large,omitempty"`

	// Number of PUBLISH QoS0 packets sent
	SendMsgQos0 int32 `json:"send_msg.qos0,omitempty"`

	// Number of PUBLISH QoS1 packets sent
	SendMsgQos1 int32 `json:"send_msg.qos1,omitempty"`

	// Number of PUBLISH QoS2 packets sent
	SendMsgQos2 int32 `json:"send_msg.qos2,omitempty"`

	// Number of bytes sent
	SendOct int32 `json:"send_oct,omitempty"`

	// Number of MQTT packets sent
	SendPkt int32 `json:"send_pkt,omitempty"`

	// Number of subscriptions established by this client.
	SubscriptionsCnt int32 `json:"subscriptions_cnt,omitempty"`

	// v4 api name [max_subscriptions] Maximum number of subscriptions allowed by this client
	SubscriptionsMax string `json:"subscriptions_max,omitempty"`

	// User name of client when connecting
	Username string `json:"username,omitempty"`

	// Topic mountpoint
	Mountpoint string `json:"mountpoint,omitempty"`

	// Indicate the configuration group used by the client
	Zone string `json:"zone,omitempty"`
}

type ClientQuery struct {
	Node           string `json:"node,omitempty"`
	Username       string `json:"username,omitempty"`
	Zone           string `json:"zone,omitempty"`
	IpAddress      string `json:"ip_address,omitempty"`
	ConnState      string `json:"conn_state,omitempty"`
	CleanStart     string `json:"clean_start,omitempty"`
	ProtoVer       string `json:"proto_ver,omitempty"`
	LikeClientId   string `json:"like_clientid,omitempty"`
	LikeUsername   string `json:"like_username,omitempty"`
	GteCreatedAt   string `json:"gte_created_at,omitempty"`
	LteCreatedAt   string `json:"lte_created_at,omitempty"`
	GteConnectedAt string `json:"gte_connected_at,omitempty"`
	LteConnectedAt string `json:"lte_connected_at,omitempty"`
}

type ClientPageQuery struct {
	emqx.PageQuery
	ClientQuery
}
