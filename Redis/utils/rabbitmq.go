package utils

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ() *amqp091.Connection {

	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		log.Fatalf("❌ RabbitMQ connection failed: %v", err)
	}
	log.Println("✅ Connected to RabbitMQ")
	return conn

}
