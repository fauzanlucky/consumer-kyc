package v1

import (
	"fmt"
	"log"

	"github.com/fauzanlucky/consumer-kyc/src/rabbitmq"
	"github.com/streadway/amqp"
)

type Route struct {
	ExchangeName string
	ExchangeType string
	QueueName    string
	RoutingKey   string
	ConsumerTag  string
	Handler      func(*Request, *amqp.Delivery)
}

type Request struct {
	Headers amqp.Table
	Body    []byte
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func (route *Route) Consume() {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     route.ConsumerTag,
		done:    make(chan error),
	}

	c.conn = rabbitmq.Start()
	defer c.conn.Close()
	defer c.Shutdown()
	var err error

	fmt.Println("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		log.Println(fmt.Sprintf("%s: %s", "Failed to open a channel", err))
		return
	}
	defer c.channel.Close()

	fmt.Printf("got Channel, declaring Exchange (%q)\n", route.ExchangeName)
	if err = c.channel.ExchangeDeclare(
		route.ExchangeName, // name
		"topic",            // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	); err != nil {
		log.Println(fmt.Sprintf("%s: %s", "Failed to declare an exchange", err))
		return
	}

	fmt.Printf("declared Exchange, declaring Queue %q\n", route.QueueName)
	args := amqp.Table{
		"x-queue-mode":   "lazy",
		"x-max-priority": 255,
	}
	q, err := c.channel.QueueDeclare(
		route.QueueName, // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		args,            // arguments
	)

	if err != nil {
		log.Println(fmt.Sprintf("%s: %s", "Failed to declare a queue", err))
		return
	}

	log.Printf("Binding queue %s to exchange %s with routing key %s\n", q.Name, route.ExchangeName, route.RoutingKey)
	if err = c.channel.QueueBind(
		route.QueueName,    // queue name
		route.RoutingKey,   // routing key
		route.ExchangeName, // exchange
		false,              // noWait
		nil,                // arguments
	); err != nil {
		log.Println(fmt.Sprintf("%s: %s", "Failed to bind a queue", err))
		return
	}

	if err = c.channel.Qos(
		1,
		0,
		false,
	); err != nil {
		log.Println("failed on setting up QoS:", err.Error())
		return
	}

	deliveries, err := c.channel.Consume(
		route.QueueName, // queue
		c.tag,           // consumerTag
		false,           // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		log.Println(fmt.Sprintf("%s: %s", "Failed to register a consumer", err))
	}

	handle(deliveries, route)
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Println("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, route *Route) {
	log.Printf("[%s][*] Waiting for messages. To exit press CTRL+C", route.QueueName)
	for d := range deliveries {
		route.Handler(&Request{
			Headers: d.Headers,
			Body:    d.Body,
		}, &d)
	}
}
