package round

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lardira/wicked-wit/pkg/response"
)

type Handler struct{}

func Router() chi.Router {
	var handler Handler
	r := chi.NewRouter()

	r.Get("/", handler.GetRounds)
	r.Post("/", handler.AddRound)
	r.Delete("/{id}", handler.DeleteRound)

	return r
}

func (h *Handler) GetRounds(w http.ResponseWriter, r *http.Request) {
	payload := []Round{}

	rounds, err := Select()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for _, model := range rounds {
		round := Round{
			Id:       model.Id,
			Position: model.Position,
			GameId:   model.GameId,
			Timed:    response.TimedFromModel(&model.TimedModel),
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

func (h *Handler) AddRound(w http.ResponseWriter, r *http.Request) {
	gameId := (chi.URLParam(r, "gameId"))
	if err := uuid.Validate(gameId); err != nil {
		response.SimpleError(w, err, http.StatusBadRequest)
		return
	}

	rounds, err := Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newId, err := Insert(
		len(rounds)+1,
		gameId,
	)
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

	Delete(id)
	w.WriteHeader(http.StatusNoContent)
}
