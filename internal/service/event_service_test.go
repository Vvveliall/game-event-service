package service

import (
	"context"
	"errors"
	"testing"

	"game-event-service/internal/model"
)

type mockEventRepository struct {
	called bool
	event  model.Event
	err    error
}

func (m *mockEventRepository) Create(ctx context.Context, event model.Event) error {
	m.called = true
	m.event = event
	return m.err
}

type mockEventPublisher struct {
	called      bool
	event       model.Event
	failures    int
	publishings int
}

func (m *mockEventPublisher) PublishEvent(ctx context.Context, event model.Event) error {
	m.called = true
	m.event = event
	m.publishings++

	if m.publishings <= m.failures {
		return errors.New("temporary rabbitmq error")
	}

	return nil
}

func TestCreateEvent(t *testing.T) {
	repository := &mockEventRepository{}
	publisher := &mockEventPublisher{}

	service := NewEventService(repository, publisher)

	event := model.Event{
		PlayerID: 1,
		Type:     "purchase.created",
		Payload:  "sword_01",
	}

	err := service.CreateEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repository.called {
		t.Fatal("repository Create was not called")
	}

	if repository.event != event {
		t.Fatalf("unexpected event: got %+v, want %+v", repository.event, event)
	}

	if !publisher.called {
		t.Fatal("event publisher was not called")
	}

	if publisher.event != event {
		t.Fatalf("unexpected published event: got %+v, want %+v", publisher.event, event)
	}
}

func TestCreateEventRepositoryError(t *testing.T) {
	repositoryError := errors.New("database error")

	repository := &mockEventRepository{
		err: repositoryError,
	}

	publisher := &mockEventPublisher{}

	service := NewEventService(repository, publisher)

	event := model.Event{
		PlayerID: 1,
		Type:     "purchase.created",
		Payload:  "repository_error_test",
	}

	err := service.CreateEvent(context.Background(), event)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, repositoryError) {
		t.Fatalf("unexpected error: got %v, want %v", err, repositoryError)
	}

	if publisher.called {
		t.Fatal("publisher should not be called when repository fails")
	}
}

func TestCreateEventRetriesPublishing(t *testing.T) {
	repository := &mockEventRepository{}
	publisher := &mockEventPublisher{
		failures: 2,
	}

	service := NewEventService(repository, publisher)

	event := model.Event{
		PlayerID: 1,
		Type:     "purchase.created",
		Payload:  "retry_test",
	}

	err := service.CreateEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if publisher.publishings != 3 {
		t.Fatalf("unexpected publish attempts: got %d, want 3", publisher.publishings)
	}
}

func TestCreateEventFailsAfterRetries(t *testing.T) {
	repository := &mockEventRepository{}
	publisher := &mockEventPublisher{
		failures: 10,
	}

	service := NewEventService(repository, publisher)

	event := model.Event{
		PlayerID: 1,
		Type:     "purchase.created",
		Payload:  "retry_fail_test",
	}

	err := service.CreateEvent(context.Background(), event)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if publisher.publishings != 3 {
		t.Fatalf("unexpected publish attempts: got %d, want 3", publisher.publishings)
	}
}
