package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"sync"
)

var handles = sync.Map{} //make(sync.Map[string]func(delivery *amqp.Delivery))

func (r *rabbit) Subscribe(ctx context.Context, config ConsumeConfig, handle func(delivery *amqp.Delivery)) (err error) {
	err = r.chConsumer.ExchangeDeclare(config.Exchange, amqp.ExchangeTopic, true, false, false, false, nil)
	if err != nil {
		log.Errorf("定义exchange %s 失败", config.Exchange)
	}
	r.CreateQueue(QueueConfig{Name: config.QueueName, Durable: true})
	r.BindQueueExchange(BindQueueConfig{QueueName: config.QueueName, Exchange: config.Exchange, RoutingKey: config.QueueName})
	return r.subscribe(ctx, config, handle)
}

func (r *rabbit) SubscribeDelay(ctx context.Context, config ConsumeConfig, handle func(delivery *amqp.Delivery)) (err error) {
	err = r.chConsumer.ExchangeDeclare(config.Exchange, "x-delayed-message", true, false, false, false, amqp.Table{"x-delayed-type": "topic"})
	if err != nil {
		log.Errorf("定义exchange %s 失败", config.Exchange)
	}
	r.CreateQueue(QueueConfig{Name: config.QueueName, Durable: true})
	r.BindQueueExchange(BindQueueConfig{QueueName: config.QueueName, Exchange: config.Exchange, RoutingKey: config.QueueName})
	return r.subscribe(ctx, config, handle)
}

func (r *rabbit) subscribe(ctx context.Context, config ConsumeConfig, handle func(delivery *amqp.Delivery)) (err error) {
	r.wgChannel.Add(1)
	defer r.wgChannel.Done()
	var msgs <-chan amqp.Delivery
	msgs, err = r.chConsumer.Consume(
		config.QueueName,
		config.Consumer,
		config.AutoAck,
		config.Exclusive,
		config.NoLocal,
		config.NoWait,
		config.Args,
	)
	if err != nil {
		return
	}
	handles.Store(config.QueueName, handle)
	//handles[config.Consumer] = handle
	var allCanceled bool
	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				return
			}
			r.wgChannel.Add(1)
			if config.ExecuteConcurrent {
				go func() {
					handle(&msg)
					r.wgChannel.Done()
				}()
			} else {
				handles.Range(func(key, value interface{}) bool {
					if key == msg.RoutingKey {
						handle := value.(func(delivery *amqp.Delivery))
						handle(&msg)
					}
					return true
				})
				msg.Ack(true)
				r.wgChannel.Done()
			}
		case <-ctx.Done():
			if allCanceled {
				continue
			}
			err = r.chConsumer.Cancel(config.Consumer, false)
			allCanceled = true
			continue
		}
	}
}
