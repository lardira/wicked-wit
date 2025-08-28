package card

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	cardService *Service
}

func NewCardHandler(cardService *Service) *Handler {
	return &Handler{
		cardService: cardService,
	}
}

func Router() chi.Router {
	handler := NewCardHandler(
		&Service{},
	)

	r := chi.NewRouter()

	r.Post("/played", handler.PlayCards)

	return r
}

func (h *Handler) PlayCards(w http.ResponseWriter, r *http.Request) {
	var req PlayCardRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.cardService.PlayCards(req.RoundId, req.UserId, req.CardIds...); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
