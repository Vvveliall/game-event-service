package service

import (
	"context"

	"game-event-service/internal/model"
	"game-event-service/internal/queue"
	"game-event-service/internal/repository"
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

	if err := s.eventPublisher.PublishEvent(ctx, event); err != nil {
		return err
	}

	return nil
}
