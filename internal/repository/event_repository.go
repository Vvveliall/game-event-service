package repository

import (
	"context"
	"fmt"

	"game-event-service/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct {
	pool *pgxpool.Pool
}

func NewEventRepository(pool *pgxpool.Pool) *EventRepository {
	return &EventRepository{
		pool: pool,
	}
}

func (r *EventRepository) Create(ctx context.Context, event model.Event) error {
	query := `
		INSERT INTO events (player_id, type, payload)
		VALUES ($1, $2, $3)
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		event.PlayerID,
		event.Type,
		event.Payload,
	)
	if err != nil {
		return fmt.Errorf("create event: %w", err)
	}

	return nil
}
