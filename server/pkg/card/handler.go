package card

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lardira/wicked-wit/pkg/response"
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

	r.Get("/", handler.GetCards)
	r.Post("/played", handler.PlayCards)

	return r
}

func (h *Handler) GetCards(w http.ResponseWriter, r *http.Request) {
	gameId := chi.URLParam(r, "gameId")
	if gameId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userId := r.URL.Query().Get("userId")
	if userId == "" {
		http.Error(w, "userId must not be empty", http.StatusBadRequest)
		return
	}

	fmt.Println(userId)

	cards, err := h.cardService.GetCards(gameId, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(cards); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PlayCards(w http.ResponseWriter, r *http.Request) {
	var req PlayCardRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.CardIds) == 0 {
		http.Error(w, "card list must not be empty", http.StatusBadRequest)
		return
	}

	answerId, err := h.cardService.PlayCards(req.RoundId, req.UserId, req.CardIds...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.SimpleData(w, answerId)
}
