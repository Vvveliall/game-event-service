package service

import (
	"context"
	"time"

	"game-event-service/internal/model"
	"game-event-service/internal/queue"
	"game-event-service/internal/repository"
	"game-event-service/internal/retry"
)

const (
	publishAttempts = 3
	publishDelay    = 200 * time.Millisecond
	publishTimeout  = 2 * time.Second
)

type EventService struct {
	eventRepository repository.EventRepositoryInterface
	eventPublisher  queue.EventPublisher
}

func NewEventService(
	eventRepository repository.EventRepositoryInterface,
	eventPublisher queue.EventPublisher,
) *EventService {
	return &EventService{
		eventRepository: eventRepository,
		eventPublisher:  eventPublisher,
	}
}

func (s *EventService) CreateEvent(ctx context.Context, event model.Event) error {
	if err := s.eventRepository.Create(ctx, event); err != nil {
		return err
	}

	publishCtx, cancel := context.WithTimeout(ctx, publishTimeout)
	defer cancel()

	if err := retry.Do(
		publishCtx,
		publishAttempts,
		publishDelay,
		func(ctx context.Context) error {
			return s.eventPublisher.PublishEvent(ctx, event)
		},
	); err != nil {
		return err
	}

	return nil
}
