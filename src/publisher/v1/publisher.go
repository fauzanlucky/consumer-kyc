package v1

import (
	"github.com/forkyid/consumer-kyc-update/src/rabbitmq"
	"github.com/streadway/amqp"
)

type Route struct {
	ExchangeName string
	ExchangeType string
	RoutingKey   string
	QueueName    string
}

type Publish struct {
	Headers amqp.Table
	Body    string
}

func (route *Route) Publish(publish *Publish, priority uint8) (err error) {
	conn := rabbitmq.Start()
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		return
	}
	defer channel.Close()

	// Queue and exchange should already be declared by client
	err = channel.Publish(
		route.ExchangeName, // exchange name
		route.RoutingKey,   // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         []byte(publish.Body),
			Headers:      publish.Headers,
			Priority:     priority,
		},
	)

	return
}
