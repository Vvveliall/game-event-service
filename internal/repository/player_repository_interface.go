package repository

import (
	"context"

	"game-event-service/internal/model"
)

type PlayerRepositoryInterface interface {
	GetByID(ctx context.Context, id int64) (model.Player, error)
}
