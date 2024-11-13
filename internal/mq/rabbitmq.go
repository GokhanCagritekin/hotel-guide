package mq

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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

// InitializeQueue initializes or ensures the existence of the specified queue.
func (r *RabbitMQ) InitializeQueue(queueName string) error {
	_, err := r.declareQueue(queueName)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	log.Printf("Queue %s initialized", queueName)
	return nil
}

// declareQueue is a helper function to declare a queue.
func (r *RabbitMQ) declareQueue(queueName string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Not deleted when unused
		false,     // Not exclusive
		false,     // No wait
		nil,       // Additional arguments
	)
}

// Publish sends a message to the specified queue.
func (r *RabbitMQ) Publish(queueName string, message []byte) error {
	_, err := r.declareQueue(queueName)
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

// NewRabbitMQURL loads RabbitMQ connection details from environment variables and returns the URL
func NewRabbitMQURL() (string, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("error loading .env file: %v", err)
	}

	// Retrieve RabbitMQ connection details from environment variables
	mqUser := os.Getenv("MQ_USER")
	mqPassword := os.Getenv("MQ_PASSWORD")
	mqHost := os.Getenv("MQ_HOST")
	mqPort := os.Getenv("MQ_PORT")

	if mqUser == "" || mqPassword == "" || mqHost == "" || mqPort == "" {
		return "", fmt.Errorf("RabbitMQ connection details (MQ_USER, MQ_PASSWORD, MQ_HOST, MQ_PORT) are not set in the environment")
	}

	// Build RabbitMQ URL
	rabbitMQURL := "amqp://" + mqUser + ":" + mqPassword + "@" + mqHost + ":" + mqPort + "/"
	return rabbitMQURL, nil
}
