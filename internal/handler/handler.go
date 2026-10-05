package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"game-event-service/internal/model"
	"game-event-service/internal/service"
)

type Handler struct {
	eventService  *service.EventService
	playerService *service.PlayerService
}

func New(
	eventService *service.EventService,
	playerService *service.PlayerService,
) *Handler {
	return &Handler{
		eventService:  eventService,
		playerService: playerService,
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("POST /api/v1/events", h.createEvent)
	mux.HandleFunc("GET /api/v1/players/{id}", h.getPlayer)

	return mux
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (h *Handler) createEvent(w http.ResponseWriter, r *http.Request) {
	var event model.Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.eventService.CreateEvent(r.Context(), event); err != nil {
		http.Error(w, "failed to create event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(event)
}

func (h *Handler) getPlayer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid player id", http.StatusBadRequest)
		return
	}

	player, err := h.playerService.GetPlayer(r.Context(), id)
	if err != nil {
		http.Error(w, "failed to get player", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(player)
}
