package rabbitmq

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *rabbit) Publish(exchange string, routingKey string, body interface{}) (err error) {
	config := PublishConfig{
		Exchange:   exchange,
		RoutingKey: routingKey,
	}
	return publish(r, body, config)
}

func (r *rabbit) PublishDelay(exchange string, routingKey string, delay int, body interface{}) (err error) {
	config := PublishConfig{
		Exchange:   exchange,
		RoutingKey: routingKey,
		Headers:    amqp.Table{"x-delay": delay},
	}
	return publish(r, body, config)
}

func publish(r *rabbit, body interface{}, config PublishConfig) (err error) {
	r.wgChannel.Add(1)
	defer r.wgChannel.Done()
	_bytes, _ := json.Marshal(body)
	err = r.chProducer.Publish(
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
