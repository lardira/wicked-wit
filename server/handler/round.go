package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lardira/wicked-wit/entity"
	"github.com/lardira/wicked-wit/internal/db/model"
)

type RoundHandler struct{}

func RoundRouter() chi.Router {
	var handler RoundHandler
	r := chi.NewRouter()

	r.Get("/", handler.GetRounds)
	r.Post("/", handler.AddRound)
	r.Delete("/{id}", handler.DeleteRound)

	return r
}

func (h *RoundHandler) GetRounds(w http.ResponseWriter, r *http.Request) {
	payload := []entity.Round{}

	rounds, err := model.SelectRounds()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for _, model := range rounds {
		round := entity.Round{
			Id:       model.Id,
			Position: model.Position,
			GameId:   model.GameId,
			Timed:    entity.TimedFromModel(&model.Timed),
		}

		if model.WinnerId.Valid {
			round.WinnerId = &model.WinnerId.String
		}

		payload = append(payload, round)
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *RoundHandler) AddRound(w http.ResponseWriter, r *http.Request) {
	gameId := (chi.URLParam(r, "gameId"))
	if err := uuid.Validate(gameId); err != nil {
		entity.SimpleError(w, err, http.StatusBadRequest)
		return
	}

	rounds, err := model.SelectRounds()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	newId, err := model.InsertRound(
		len(rounds)+1,
		gameId,
	)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	entity.SimpleData(w, newId)
}

func (h *RoundHandler) DeleteRound(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	model.DeleteRound(id)
	w.WriteHeader(http.StatusNoContent)
}
