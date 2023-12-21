package rabbitmq

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

func (r *rabbit) Publish(ctx context.Context, exchange string, routingKey string, body interface{}) (err error) {
	config := PublishConfig{
		Exchange:   exchange,
		RoutingKey: routingKey,
	}
	return publish(ctx, r, body, config)
}

func (r *rabbit) PublishDelay(ctx context.Context, exchange string, routingKey string, delay int, body interface{}) (err error) {
	config := PublishConfig{
		Exchange:   exchange,
		RoutingKey: routingKey,
		Headers:    amqp.Table{"x-delay": delay},
	}
	return publish(ctx, r, body, config)
}

func publish(ctx context.Context, r *rabbit, body interface{}, config PublishConfig) (err error) {
	r.wgChannel.Add(1)
	defer r.wgChannel.Done()
	_bytes, _ := json.Marshal(body)
	if e := r.chConsumer.ExchangeDeclare(config.Exchange, amqp.ExchangeDirect, true, false, false, false, nil); e != nil {
		log.Errorf("定义exchange %s 失败", config.Exchange)
	}
	if _, e := r.CreateQueue(QueueConfig{Name: config.RoutingKey, Durable: true}); e != nil {
		log.Errorf("创建队列出错:%s", e.Error())
	}
	if e := r.BindQueueExchange(BindQueueConfig{QueueName: config.RoutingKey, Exchange: config.Exchange, RoutingKey: config.RoutingKey}); e != nil {
		log.Errorf("绑定队列出错:%s", e.Error())
	}

	err = r.chProducer.PublishWithContext(
		ctx,
		config.Exchange,
		config.RoutingKey,
		config.Mandatory,
		config.Immediate,
		amqp.Publishing{
			Headers:       config.Headers,
			ContentType:   config.ContentType,
			Priority:      config.Priority,
			CorrelationId: config.CorrelationID,
			Body:          _bytes,
		},
	)
	return
}
