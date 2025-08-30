package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lardira/wicked-wit/internal/domain/entity"
	"github.com/lardira/wicked-wit/internal/domain/interfaces"
	"github.com/lardira/wicked-wit/internal/helper/response"
)

type gameHandler struct {
	gameService interfaces.GameService
}

func NewGameHandler(gameService interfaces.GameService) *gameHandler {
	return &gameHandler{
		gameService: gameService,
	}
}

func GameRouter(gameService interfaces.GameService, cardService interfaces.CardService, roundService interfaces.RoundService) chi.Router {
	handler := NewGameHandler(
		gameService,
	)

	r := chi.NewRouter()

	r.Get("/", handler.GetGames)
	r.Get("/{id}", handler.GetGame)
	r.Post("/", handler.CreateGame)
	r.Delete("/{id}", handler.DeleteGame)

	r.Mount("/{gameId}/rounds", RoundRouter(roundService))
	r.Mount("/{gameId}/cards", CardRouter(cardService))

	return r
}

func (h *gameHandler) GetGames(w http.ResponseWriter, r *http.Request) {
	payload, err := h.gameService.GetGames()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *gameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	payload, err := h.gameService.GetGame(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *gameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var game entity.GameRequest
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(game.Users) == 0 {
		http.Error(w, "users must not be empty", http.StatusBadRequest)
		return
	}

	newId, err := h.gameService.CreateGame(&game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SimpleData(w, newId)
}

func (h *gameHandler) DeleteGame(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.gameService.DeleteGame(id)

	w.WriteHeader(http.StatusNoContent)
}
