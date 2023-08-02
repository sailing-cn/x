package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ interface {
	Connector
	Closer
	QueueCreator
	Producer
	Consumer
}

type Connector interface {
	Connect(config ConnConfig) (notify chan *amqp.Error, err error)
}

type Closer interface {
	Close(ctx context.Context) (done chan struct{})
}

type QueueCreator interface {
	CreateQueue(config QueueConfig) (queue amqp.Queue, err error)
	BindQueueExchange(config BindQueueConfig) (err error)
	UnbindQueueExchange(config BindQueueConfig) error
}

type Producer interface {
	Publish(exchange string, routingKey string, body interface{}) (err error)
	PublishDelay(exchange string, routingKey string, delay int, body interface{}) (err error)
}

type Consumer interface {
	// Subscribe 订阅
	Subscribe(ctx context.Context, config ConsumeConfig, f func(*amqp.Delivery)) (err error)
	//SubscribeDelay 延时订阅
	SubscribeDelay(ctx context.Context, config ConsumeConfig, f func(*amqp.Delivery)) (err error)
}

type RabbitSetup interface {
	Setup()
}

type Setup func()

func (s Setup) Setup() {
	s()
}
