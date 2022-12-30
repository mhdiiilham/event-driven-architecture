package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

var ()

func main() {
	rand.Seed(time.Now().Unix())
	log.Println("The NOP System Producer")
	run := true

	conn, err := amqp.Dial("amqp://noella:noella@localhost:5672/")
	if err != nil {
		log.Fatalf("failed establising RabbitMQ Connection: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed opening RabbitMQ channel: %v", err)
	}

	_, err = ch.QueueDeclare(
		"ch-testing", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("failed declaring RabbitMQ channel: %v", err)
	}

	doneCh := make(chan os.Signal, 1)
	signal.Notify(doneCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	for run {
		select {
		case <-doneCh:
			run = false
		default:
			go func() {
				correlationID := generateCorrelationID()
				db := rand.Intn(500)
				timeSleep := rand.Intn(5000) * int(time.Millisecond)
				time.Sleep(time.Duration(timeSleep))
				sendEvent(ch, correlationID, db)
			}()
		}
	}

	log.Println("system is closing...")
}

func generateCorrelationID() string {
	return uuid.NewString()
}

func sendEvent(channel *amqp.Channel, correlationID string, decibles int) {
	log.Printf("correlation-id=%s message=%d sending...", correlationID, decibles)
	msg := Event{CorrelationID: correlationID, Body: Decible{Value: decibles}, ContentType: "application/json"}
	b, _ := msg.ToBytes()

	if err := channel.PublishWithContext(context.Background(), "", "ch-testing", false, false, amqp.Publishing{
		ContentType:   msg.ContentType,
		Body:          b,
		CorrelationId: msg.CorrelationID,
	}); err != nil {
		log.Printf("failed pushing message %s to channel: %v", correlationID, err)
	} else {
		log.Printf("message %s published", correlationID)
	}
}
