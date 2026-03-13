package rabbitmq

import "github.com/rabbitmq/amqp091-go"

// ExchangeOptions contains exchange declaration options
type ExchangeOptions struct {
	Name       string
	Kind       string // direct, fanout, topic, headers
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp091.Table
}

// QueueOptions contains queue declaration options
type QueueOptions struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp091.Table
}

// DeclareExchange declares an exchange
func (c *Client) DeclareExchange(opts ExchangeOptions) error {
	return c.channel.ExchangeDeclare(
		opts.Name,
		opts.Kind,
		opts.Durable,
		opts.AutoDelete,
		opts.Internal,
		opts.NoWait,
		opts.Args,
	)
}

// DeclareQueue declares a queue
func (c *Client) DeclareQueue(opts QueueOptions) (amqp091.Queue, error) {
	return c.channel.QueueDeclare(
		opts.Name,
		opts.Durable,
		opts.AutoDelete,
		opts.Exclusive,
		opts.NoWait,
		opts.Args,
	)
}

// BindQueue binds a queue to an exchange
func (c *Client) BindQueue(queueName, routingKey, exchangeName string) error {
	return c.channel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
}
