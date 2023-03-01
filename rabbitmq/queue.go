package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

func (r *rabbit) CreateQueue(config QueueConfig) (queue amqp.Queue, err error) {
	queue, err = r.chConsumer.QueueDeclare(
		config.Name,
		config.Durable,
		config.AutoDelete,
		config.Exclusive,
		config.NoWait,
		config.Args,
	)
	return
}

func (r *rabbit) BindQueueExchange(config BindQueueConfig) (err error) {
	err = r.chConsumer.QueueBind(
		config.QueueName,
		config.RoutingKey,
		config.Exchange,
		config.NoWait,
		config.Args,
	)
	return
}
func (r *rabbit) UnbindQueueExchange(config BindQueueConfig) (err error) {
	err = r.chConsumer.QueueUnbind(
		config.QueueName,
		config.RoutingKey,
		config.Exchange,
		config.Args,
	)
	return
}
