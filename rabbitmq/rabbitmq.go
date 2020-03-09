package rabbitmq

import "github.com/streadway/amqp"

// Config parameters for rabbitmq
type Config struct {
	URL        string
	Exchange   string
	QueueName  string
	RoutingKey string
	BindingKey string
}

// Rabbitmq struct holds conn, channel and config
type Rabbitmq struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  Config
}

// New instance of rabbitmq
func New(config Config) (*Rabbitmq, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	if err := newExchange(ch, config.Exchange); err != nil {
		return nil, err
	}

	if err := newQueue(ch, config.Exchange, config.QueueName, config.BindingKey); err != nil {
		return nil, err
	}

	return &Rabbitmq{
		conn:    conn,
		channel: ch,
		config:  config,
	}, nil
}

func newExchange(c *amqp.Channel, name string) error {
	return c.ExchangeDeclare(
		name,    // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
}

func newQueue(ch *amqp.Channel, exchange string, name string, bindingKey string) error {
	q, err := ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	return ch.QueueBind(
		q.Name,     // queue name
		bindingKey, // routing key
		exchange,   // exchange
		false,
		nil,
	)
}

// Publish message to rabbitmq
func (r *Rabbitmq) Publish(message []byte) error {
	return r.channel.Publish(
		r.config.Exchange,   // exchange
		r.config.RoutingKey, // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			Body: message,
		},
	)
}

// Consume message from rabbitmq
func (r *Rabbitmq) Consume() (<-chan []byte, error) {
	msgs, err := r.channel.Consume(
		r.config.QueueName, // queue
		"",                 // consumer
		true,               // auto ack
		false,              // exclusive
		false,              // no local
		false,              // no wait
		nil,                // args
	)
	if err != nil {
		return nil, err
	}

	deliveries := make(chan []byte)
	go func() {
		for msg := range msgs {
			deliveries <- msg.Body
		}
	}()
	return (<-chan []byte)(deliveries), nil
}

// close rabbitmq connection
func (r *Rabbitmq) close() error {
	if err := r.conn.Close(); err != nil {
		return err
	}
	if err := r.channel.Close(); err != nil {
		return err
	}
	return nil
}
