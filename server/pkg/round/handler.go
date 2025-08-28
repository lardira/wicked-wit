package round

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lardira/wicked-wit/pkg/card"
	"github.com/lardira/wicked-wit/pkg/response"
)

type Handler struct {
	roundService *Service
}

func NewRoundHandler(roundService *Service) *Handler {
	return &Handler{
		roundService: roundService,
	}
}

func Router() chi.Router {
	handler := NewRoundHandler(
		NewRoundService(&card.Service{}),
	)

	r := chi.NewRouter()

	r.Get("/", handler.GetRounds)
	r.Post("/", handler.AddRound)
	r.Delete("/{id}", handler.DeleteRound)

	return r
}

func (h *Handler) GetRounds(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) AddRound(w http.ResponseWriter, r *http.Request) {
	gameId := (chi.URLParam(r, "gameId"))
	if err := uuid.Validate(gameId); err != nil {
		response.SimpleError(w, err, http.StatusBadRequest)
		return
	}

	newId, err := h.roundService.AddRound(gameId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SimpleData(w, newId)
}

func (h *Handler) DeleteRound(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.roundService.DeleteRound(id)
	w.WriteHeader(http.StatusNoContent)
}
