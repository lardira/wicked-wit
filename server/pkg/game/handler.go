package game

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lardira/wicked-wit/pkg/card"
	"github.com/lardira/wicked-wit/pkg/response"
	"github.com/lardira/wicked-wit/pkg/round"
)

type Handler struct {
	gameService *Service
}

func NewGameHandler(gameService *Service) *Handler {
	return &Handler{
		gameService: gameService,
	}
}

func Router() chi.Router {
	handler := NewGameHandler(
		&Service{},
	)

	r := chi.NewRouter()

	r.Get("/", handler.GetGames)
	r.Post("/", handler.CreateGame)
	r.Delete("/{id}", handler.DeleteGame)

	r.Mount("/{gameId}/rounds", round.Router())
	r.Mount("/{gameId}/cards", card.Router())

	return r
}

func (h *Handler) GetGames(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var game GameRequest
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

func (h *Handler) DeleteGame(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.gameService.DeleteGame(id)

	w.WriteHeader(http.StatusNoContent)
}
