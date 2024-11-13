package mq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// MessageQueue interface abstracts RabbitMQ operations.
type MessageQueue interface {
	Publish(queueName string, message []byte) error
	Consume(queueName string) (<-chan amqp.Delivery, error)
	Close() error
	InitializeQueue(queueName string) error
}

// RabbitMQ struct represents the RabbitMQ configuration implementing MessageQueue.
type RabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

// NewRabbitMQ creates a new RabbitMQ configuration and initializes the connection.
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &RabbitMQ{
		connection: conn,
		channel:    ch,
	}, nil
}

// Publish sends a message to the specified queue.
func (r *RabbitMQ) Publish(queueName string, message []byte) error {
	// Declare queue or use existing one
	_, err := r.channel.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Not deleted when unused
		false,     // Not exclusive
		false,     // No wait
		nil,       // Additional arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Publish message to the queue
	err = r.channel.Publish(
		"",        // Exchange
		queueName, // Routing key
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Message published to queue %s", queueName)
	return nil
}

// Consume starts consuming messages from the specified queue.
func (r *RabbitMQ) Consume(queueName string) (<-chan amqp.Delivery, error) {
	msgs, err := r.channel.Consume(
		queueName, // Queue name
		"",        // Consumer name
		true,      // Auto-acknowledge
		false,     // Not exclusive
		false,     // Not local-only
		false,     // No wait
		nil,       // Additional arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %w", err)
	}

	return msgs, nil
}

// Close closes the RabbitMQ connection and channel.
func (r *RabbitMQ) Close() error {
	if err := r.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}
	if err := r.connection.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}

func (r *RabbitMQ) InitializeQueue(queueName string) error {
	// Declare queue or use existing one
	_, err := r.channel.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Not deleted when unused
		false,     // Not exclusive
		false,     // No wait
		nil,       // Additional arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	log.Printf("Queue %s initialized", queueName)
	return nil
}
