package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"game-event-service/internal/service"
)

func TestHealth(t *testing.T) {
	handler := New(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	handler.Routes().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d, want %d", recorder.Code, http.StatusOK)
	}

	if !strings.Contains(recorder.Body.String(), `"status":"ok"`) {
		t.Fatalf("unexpected response body: %s", recorder.Body.String())
	}
}

func TestCreateEventInvalidJSON(t *testing.T) {
	handler := New(nil, nil)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/events",
		strings.NewReader(`{"player_id":1`),
	)

	recorder := httptest.NewRecorder()

	handler.Routes().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: got %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestGetPlayerInvalidID(t *testing.T) {
	handler := New(nil, &service.PlayerService{})

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/players/not-a-number",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.Routes().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: got %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}
