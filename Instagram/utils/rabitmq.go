package utils

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeName = "likes.exchange"
	QueueName    = "likes.queue"
	RoutingKey   = "like.event"
)

type Client struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// NewClient connects to RabbitMQ and sets up the exchange/queue
func NewClient(url string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare Exchange
	if err := ch.ExchangeDeclare(
		ExchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-delete
		false,        // internal
		false,        // no-wait
		nil,          // args
	); err != nil {
		conn.Close()
		return nil, fmt.Errorf("exchange declare failed: %w", err)
	}

	// Declare Queue
	_, err = ch.QueueDeclare(
		QueueName, // name
		true,      // durable
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("queue declare failed: %w", err)
	}

	// Bind Queue to Exchange
	if err := ch.QueueBind(
		QueueName,
		RoutingKey,
		ExchangeName,
		false,
		nil,
	); err != nil {
		conn.Close()
		return nil, fmt.Errorf("queue bind failed: %w", err)
	}

	log.Println("✅ RabbitMQ connection established and queue bound successfully")

	return &Client{conn: conn, ch: ch}, nil
}

// PublishLikeEvent publishes a like action to the queue
func (c *Client) PublishLikeEvent(userID int64, postID int64) error {
	body := fmt.Sprintf(`{"user_id":%d, "post_id":%d}`, userID, postID)

	ctx := context.Background() // amqp091-go requires context
	return c.ch.PublishWithContext(
		ctx,
		ExchangeName, // exchange
		RoutingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         []byte(body),
			DeliveryMode: amqp.Persistent, // durable message
		},
	)
}

// Close cleans up the connection and channel
func (c *Client) Close() {
	if c.ch != nil {
		_ = c.ch.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
