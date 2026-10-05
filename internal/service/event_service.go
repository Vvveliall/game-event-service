package service

import (
	"context"

	"game-event-service/internal/model"
	"game-event-service/internal/repository"
)

type EventService struct {
	eventRepository repository.EventRepositoryInterface
}

func NewEventService(eventRepository repository.EventRepositoryInterface) *EventService {
	return &EventService{
		eventRepository: eventRepository,
	}
}

func (s *EventService) CreateEvent(ctx context.Context, event model.Event) error {
	return s.eventRepository.Create(ctx, event)
}
