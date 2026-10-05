package repository

import (
	"context"

	"game-event-service/internal/model"
)

type EventRepositoryInterface interface {
	Create(ctx context.Context, event model.Event) error
}
