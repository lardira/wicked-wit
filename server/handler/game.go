package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

	return r
}

func (h GameHandler) GetGames(w http.ResponseWriter, r *http.Request) {
	payload := []entity.Game{}
	games, err := model.SelectGames()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for _, model := range games {
		payload = append(payload, entity.Game{
			Id:           model.Id,
			Title:        model.Title,
			MaxPlayers:   model.MaxPlayers,
			CurrentRound: model.CurrentRound,
			MaxRound:     model.MaxRound,
		})
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var game entity.GameRequest
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	model.InsertGame(
		game.Title,
		game.MaxPlayers,
		game.MaxRound,
		uuid.NewString(),
	)
}

func (h GameHandler) DeleteGame(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if len(id) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	model.DeleteGame(id)
	w.WriteHeader(http.StatusNoContent)
}
