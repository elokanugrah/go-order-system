package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/elokanugrah/go-order-system/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Println("Starting Worker Service...")

	cfg := config.Load()

	// Connect to RabbitMQ
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the queue to make sure it exists
	q, err := ch.QueueDeclare(
		"orders.created",
		true, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Start consuming messages from the queue
	msgs, err := ch.Consume(
		q.Name,
		"order-worker", // consumer name
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Goroutine to process messages
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			processMessage(d.Body)
		}
	}()

	log.Printf("Worker is waiting for messages. To exit press CTRL+C")

	// Handles graceful shutdown on receiving SIGINT or SIGTERM signals.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")

	log.Println("Worker exited gracefully.")
}

// A helper function to process the message payload.
func processMessage(body []byte) {
	time.Sleep(2 * time.Second) // Simulate a 2-second task

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err == nil {
		if orderID, ok := payload["order_id"]; ok {
			log.Printf("[WORKER] Finished processing confirmation for Order ID: %.0f", orderID)
		}
	} else {
		log.Printf("[WORKER] ERROR: Failed to unmarshal message: %v", err)
	}
}
