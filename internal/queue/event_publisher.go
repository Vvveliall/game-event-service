package queue

import (
	"context"

	"game-event-service/internal/model"
)

type EventPublisher interface {
	PublishEvent(ctx context.Context, event model.Event) error
}
