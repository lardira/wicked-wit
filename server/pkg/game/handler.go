package game

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lardira/wicked-wit/pkg/response"
	"github.com/lardira/wicked-wit/pkg/round"
)

type Handler struct{}

func Router() chi.Router {
	var handler Handler
	r := chi.NewRouter()

	r.Get("/", handler.GetGames)
	r.Post("/", handler.CreateGame)
	r.Delete("/{id}", handler.DeleteGame)

	r.Mount("/{gameId}/rounds", round.Router())

	return r
}

func (h *Handler) GetGames(w http.ResponseWriter, r *http.Request) {
	payload := []Game{}
	games, err := SelectGames()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for _, model := range games {
		game := Game{
			Id:         model.Id,
			Title:      model.Title,
			MaxPlayers: model.MaxPlayers,
			MaxRound:   model.MaxRound,
			UserHostId: model.UserHostId,
			Timed:      response.TimedFromModel(&model.TimedModel),
		}

		payload = append(payload, game)
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var game GameRequest
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	newId, err := InsertGame(
		game.Title,
		game.MaxPlayers,
		game.MaxRound,
	)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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

	DeleteGame(id)
	w.WriteHeader(http.StatusNoContent)
}
