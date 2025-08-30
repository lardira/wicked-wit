package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lardira/wicked-wit/internal/domain/interfaces"
)

type roundHandler struct {
	roundService interfaces.RoundService
}

func NewRoundHandler(roundService interfaces.RoundService) *roundHandler {
	return &roundHandler{
		roundService: roundService,
	}
}

func RoundRouter(roundService interfaces.RoundService) chi.Router {
	handler := NewRoundHandler(
		roundService,
	)

	r := chi.NewRouter()

	r.Get("/", handler.GetRounds)
	r.Delete("/{id}", handler.DeleteRound)

	return r
}

func (h *roundHandler) GetRounds(w http.ResponseWriter, r *http.Request) {
	gameId := chi.URLParam(r, "gameId")
	if gameId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	payload, err := h.roundService.GetRounds(gameId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *roundHandler) DeleteRound(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.roundService.DeleteRound(id)
	w.WriteHeader(http.StatusNoContent)
}
