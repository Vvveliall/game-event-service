package service

import (
	"context"
	"testing"

	"game-event-service/internal/model"
)

type mockEventRepository struct {
	called bool
	event  model.Event
}

func (m *mockEventRepository) Create(ctx context.Context, event model.Event) error {
	m.called = true
	m.event = event
	return nil
}

type mockEventPublisher struct {
	called bool
	event  model.Event
}

func (m *mockEventPublisher) PublishEvent(ctx context.Context, event model.Event) error {
	m.called = true
	m.event = event
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
