package v1

import (
	"github.com/fauzanlucky/consumer-kyc/src/rabbitmq"
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

	err = channel.ExchangeDeclare(
		route.ExchangeName, // name
		route.ExchangeType, // type
		true,               // durable
		false,              // auto-delete
		false,              // internal
		false,              // no-wait
		nil,                // argument
	)
	if err != nil {
		return
	}

	args := amqp.Table{
		"x-queue-mode":   "lazy",
		"x-max-priority": 255,
	}
	_, err = channel.QueueDeclare(
		route.QueueName, // queue name
		true,            // durable
		false,           // delete when used
		false,           // exclusive
		false,           // no-wait
		args,            // arguments
	)
	if err != nil {
		return
	}

	err = channel.QueueBind(
		route.QueueName,    // queue name
		route.RoutingKey,   // routing key
		route.ExchangeName, // exchange
		false,              // noWait
		nil,                // arguments
	)
	if err != nil {
		return
	}

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
