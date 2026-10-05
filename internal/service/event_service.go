package service

import (
	"context"

	"game-event-service/internal/model"
	"game-event-service/internal/repository"
	"game-event-service/internal/worker"
)

type EventService struct {
	eventRepository repository.EventRepositoryInterface
	eventWorker     *worker.EventWorker
}

func NewEventService(
	eventRepository repository.EventRepositoryInterface,
	eventWorker *worker.EventWorker,
) *EventService {
	return &EventService{
		eventRepository: eventRepository,
		eventWorker:     eventWorker,
	}
}

func (s *EventService) CreateEvent(ctx context.Context, event model.Event) error {
	if err := s.eventRepository.Create(ctx, event); err != nil {
		return err
	}

	s.eventWorker.Submit(event)

	return nil
}
