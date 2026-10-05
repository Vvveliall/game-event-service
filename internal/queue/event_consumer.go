package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"game-event-service/internal/model"
	"game-event-service/internal/worker"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventConsumer struct {
	channel *amqp.Channel
	worker  *worker.EventWorker
}

func NewEventConsumer(rabbitMQ *RabbitMQ, eventWorker *worker.EventWorker) *EventConsumer {
	return &EventConsumer{
		channel: rabbitMQ.channel,
		worker:  eventWorker,
	}
}

func (c *EventConsumer) Run(ctx context.Context) error {
	messages, err := c.channel.Consume(
		eventQueue,
		"game-event-consumer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("start event consumer: %w", err)
	}

	log.Println("event consumer started")

	for {
		select {
		case message, ok := <-messages:
			if !ok {
				return fmt.Errorf("event consumer channel closed")
			}

			var event model.Event

			if err := json.Unmarshal(message.Body, &event); err != nil {
				log.Printf("failed to decode event: %v", err)
				_ = message.Nack(false, false)
				continue
			}

			c.worker.Submit(event)

			if err := message.Ack(false); err != nil {
				log.Printf("failed to ack event: %v", err)
			}

		case <-ctx.Done():
			log.Println("event consumer stopped")
			return nil
		}
	}
}
