package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lardira/wicked-wit/internal/domain/entity"
	"github.com/lardira/wicked-wit/internal/domain/interfaces"
	"github.com/lardira/wicked-wit/internal/domain/service"
	"github.com/lardira/wicked-wit/internal/helper/response"
)

type cardHandler struct {
	cardService interfaces.CardService
}

func NewCardHandler(cardService interfaces.CardService) *cardHandler {
	return &cardHandler{
		cardService: cardService,
	}
}

func CardRouter(cardService interfaces.CardService) chi.Router {
	handler := NewCardHandler(
		service.NewCardService(),
	)

	r := chi.NewRouter()

	r.Get("/", handler.GetCards)
	r.Post("/played", handler.PlayCards)

	return r
}

func (h *cardHandler) GetCards(w http.ResponseWriter, r *http.Request) {
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

func (h *cardHandler) PlayCards(w http.ResponseWriter, r *http.Request) {
	var req entity.PlayCardRequest
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
