package main

import (
	"fmt"

	"github.com/rabingaire/medium/rabbitmq"
)

func main() {
	q, err := rabbitmq.New(rabbitmq.Config{
		URL:        "amqp://localhost:5672",
		Exchange:   "example_exchange",
		QueueName:  "example_queuename",
		RoutingKey: "example",
		BindingKey: "example",
	})
	if err != nil {
		panic(err)
	}

	// Publish
	message := []byte("Hello")
	if err := q.Publish(message); err != nil {
		panic(err)
	}

	// Consumer
	messages, err := q.Consume()
	if err != nil {
		panic(err)
	}

	for m := range messages {
		fmt.Println(m)
	}
}
