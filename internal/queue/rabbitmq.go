package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"game-event-service/internal/model"

	amqp "github.com/rabbitmq/amqp091-go"
)

const eventQueue = "game.events"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("connect to rabbitmq: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("open rabbitmq channel: %w", err)
	}

	_, err = channel.QueueDeclare(
		eventQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("declare event queue: %w", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
	}, nil
}

func (r *RabbitMQ) PublishEvent(ctx context.Context, event model.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	err = r.channel.PublishWithContext(
		ctx,
		"",
		eventQueue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("publish event: %w", err)
	}

	return nil
}

func (r *RabbitMQ) Close() error {
	if err := r.channel.Close(); err != nil {
		r.conn.Close()
		return err
	}

	return r.conn.Close()
}
