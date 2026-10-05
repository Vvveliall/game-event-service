package repository

import (
	"context"
	"errors"
	"fmt"

	"game-event-service/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlayerRepository struct {
	pool *pgxpool.Pool
}

func NewPlayerRepository(pool *pgxpool.Pool) *PlayerRepository {
	return &PlayerRepository{
		pool: pool,
	}
}

func (r *PlayerRepository) GetByID(ctx context.Context, id int64) (model.Player, error) {
	query := `
		SELECT id, username, created_at
		FROM players
		WHERE id = $1
	`

	var player model.Player

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&player.ID,
		&player.Username,
		&player.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Player{}, fmt.Errorf("player not found")
		}

		return model.Player{}, fmt.Errorf("get player: %w", err)
	}

	return player, nil
}
