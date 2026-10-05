package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"game-event-service/internal/model"

	"github.com/redis/go-redis/v9"
)

const playerCacheTTL = 5 * time.Minute

type PlayerCache struct {
	client *redis.Client
}

func NewPlayerCache(client *redis.Client) *PlayerCache {
	return &PlayerCache{
		client: client,
	}
}

func (c *PlayerCache) Get(ctx context.Context, playerID int64) (model.Player, error) {
	key := fmt.Sprintf("player:%d", playerID)

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.Player{}, nil
		}

		return model.Player{}, fmt.Errorf("get player from redis: %w", err)
	}

	var player model.Player

	if err := json.Unmarshal(data, &player); err != nil {
		return model.Player{}, fmt.Errorf("unmarshal player from redis: %w", err)
	}

	return player, nil
}

func (c *PlayerCache) Set(ctx context.Context, player model.Player) error {
	data, err := json.Marshal(player)
	if err != nil {
		return fmt.Errorf("marshal player for redis: %w", err)
	}

	key := fmt.Sprintf("player:%d", player.ID)

	if err := c.client.Set(ctx, key, data, playerCacheTTL).Err(); err != nil {
		return fmt.Errorf("set player in redis: %w", err)
	}

	return nil
}
