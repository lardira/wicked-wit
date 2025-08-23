package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lardira/wicked-wit/entity"
	"github.com/lardira/wicked-wit/internal/db/model"
)

type GameHandler struct{}

func GameRouter() chi.Router {
	var handler GameHandler
	r := chi.NewRouter()

	r.Get("/", handler.GetGames)
	r.Post("/", handler.CreateGame)
	r.Delete("/{id}", handler.DeleteGame)

	r.Mount("/{gameId}/rounds", RoundRouter())

	return r
}

func (h *GameHandler) GetGames(w http.ResponseWriter, r *http.Request) {
	payload := []entity.Game{}
	games, err := model.SelectGames()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for _, model := range games {
		game := entity.Game{
			Id:         model.Id,
			Title:      model.Title,
			MaxPlayers: model.MaxPlayers,
			MaxRound:   model.MaxRound,
			Timed:      entity.TimedFromModel(&model.Timed),
		}

		if model.CurrentRound.Valid {
			round := uint(model.CurrentRound.Int64)
			game.CurrentRound = &round
		}

		payload = append(payload, game)
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var game entity.GameRequest
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	newId, err := model.InsertGame(
		game.Title,
		game.MaxPlayers,
		game.MaxRound,
	)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	entity.SimpleData(w, newId)
}

func (h *GameHandler) DeleteGame(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	model.DeleteGame(id)
	w.WriteHeader(http.StatusNoContent)
}
