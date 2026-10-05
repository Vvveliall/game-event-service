package service

import (
	"context"
	"fmt"

	"game-event-service/internal/cache"
	"game-event-service/internal/model"
	"game-event-service/internal/repository"
)

type PlayerService struct {
	playerRepository repository.PlayerRepositoryInterface
	playerCache      *cache.PlayerCache
}

func NewPlayerService(
	playerRepository repository.PlayerRepositoryInterface,
	playerCache *cache.PlayerCache,
) *PlayerService {
	return &PlayerService{
		playerRepository: playerRepository,
		playerCache:      playerCache,
	}
}

func (s *PlayerService) GetPlayer(ctx context.Context, id int64) (model.Player, error) {
	player, err := s.playerCache.Get(ctx, id)
	if err != nil {
		return model.Player{}, err
	}

	if player.ID != 0 {
		return player, nil
	}

	player, err = s.playerRepository.GetByID(ctx, id)
	if err != nil {
		return model.Player{}, fmt.Errorf("get player from repository: %w", err)
	}

	if err := s.playerCache.Set(ctx, player); err != nil {
		return model.Player{}, err
	}

	return player, nil
}
