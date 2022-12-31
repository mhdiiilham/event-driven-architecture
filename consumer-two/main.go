package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

var messages = []string{}

func main() {
	log.Println("consumer-two ready to receive messages")
	conn, err := amqp.Dial("amqp://noella:noella@localhost:5672/")
	if err != nil {
		log.Fatalf("failed establising RabbitMQ Connection: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed establising RabbitMQ channel: %v", err)
	}

	msgs, err := ch.Consume("telemetry", "consumer-two", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed consume RabbitMQ queue: %v", err)
	}

	doneCh := make(chan os.Signal, 1)
	signal.Notify(doneCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for msg := range msgs {
			log.Printf("message with correlation-id=%s received by consumer-two", msg.CorrelationId)
			messages = append(messages, string(msg.Body))
		}
	}()

	<-doneCh
	log.Println("consumer-two terminated")
	log.Printf("consumer-two received %d of messages", len(messages))

}
