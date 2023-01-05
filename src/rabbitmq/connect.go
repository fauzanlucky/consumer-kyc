package rabbitmq

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func Start() *amqp.Connection {
	if os.Getenv("RABBITMQ_PORT") == "" {
		os.Setenv("RABBITMQ_PORT", "5672")
	}

	connString := fmt.Sprintf("amqp://%s:%s@%s:%s",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	conn, err := amqp.Dial(connString)

	if err != nil {
		log.Println(fmt.Sprintf("%s: %s", "Failed to connect to RabbitMQ", err.Error()))
	} else {
		fmt.Println("Successfully connected to RabbitMQ")
	}

	return conn
}
