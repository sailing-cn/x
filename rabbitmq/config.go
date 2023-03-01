package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

type ConnConfig struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	Vhost         string `yaml:"vhost"`
	PrefetchCount int    `yaml:"prefetch_count"`
}

type ExchangeConfig struct {
	Name       string
	Type       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type BindQueueConfig struct {
	QueueName  string
	Exchange   string
	RoutingKey string
	NoWait     bool
	Args       amqp.Table
}

type ConsumeConfig struct {
	QueueName         string
	Consumer          string
	Exchange          string
	AutoAck           bool
	Exclusive         bool
	NoLocal           bool
	NoWait            bool
	Args              amqp.Table
	ExecuteConcurrent bool
}

type PublishConfig struct {
	Exchange      string
	RoutingKey    string
	Mandatory     bool
	Immediate     bool
	Headers       amqp.Table
	ContentType   string
	Priority      uint8
	CorrelationID string
}
