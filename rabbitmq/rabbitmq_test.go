package rabbitmq

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Define the custom testify suite
type RabbitmqTestSuite struct {
	suite.Suite
	queue *Rabbitmq
}

func (s *RabbitmqTestSuite) SetupTest() {
	q, err := New(Config{
		URL:        "amqp://rabbitmq:5672",
		Exchange:   "example_exchange",
		QueueName:  "example_queuename",
		RoutingKey: "example",
		BindingKey: "example",
	})
	if err != nil {
		panic(err)
	}
	s.queue = q
}

func (s *RabbitmqTestSuite) TearDownTest() {
	s.queue.close()
}

func (s *RabbitmqTestSuite) TestPublishMessage() {
	s.T().Run("publish a message", func(t *testing.T) {
		// Producer
		message := []byte("Hello")
		err := s.queue.Publish(message)
		assert.NoError(t, err, "Publish() error:\nwant  nil\ngot  %v", err)
	})
}

func (s *RabbitmqTestSuite) TestConsumeMessage() {
	// Producer
	message := []byte("Hello")
	err := s.queue.Publish(message)
	if err != nil {
		panic(err)
	}

	// Consumer
	messages, err := s.queue.Consume()
	if err != nil {
		panic(err)
	}

	s.T().Run("consume a message", func(t *testing.T) {
		assert.NoError(t, err, "Consume() error:\nwant  nil\ngot  %v", err)
	})

	s.T().Run("expect a delivery", func(t *testing.T) {
		expected := []byte("Hello")
		select {
		case message := <-messages:
			{
				assert.Equal(t, expected, message, "Consume() error:\nwant  %v\ngot  %v", expected, message)
			}
		}
	})
}

func TestRabbitmqTestSuite(t *testing.T) {
	suite.Run(t, new(RabbitmqTestSuite))
}
