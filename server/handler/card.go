package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lardira/wicked-wit/entity"
	"github.com/lardira/wicked-wit/internal/db/model"
)

type CardHandler struct{}

func CardRouter() chi.Router {
	var handler CardHandler
	r := chi.NewRouter()

	r.Get("/answers", handler.GetCardAnswers)
	r.Post("/answers/used", handler.UseAnswerCards)
	r.Get("/answers/used", handler.GetUsedAnswerCards)

	r.Get("/templates", handler.GetCardTemplates)

	r.Post("/", handler.CreateCard)

	return r
}

func (h *CardHandler) GetCardAnswers(w http.ResponseWriter, r *http.Request) {
	payload := []entity.CardAnswer{}
	cards := []model.CardAnswer{}

	gameId := r.URL.Query().Get("gameId")
	if gameId == "" {
		answerCards, err := model.SelectCardAnswers()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		cards = answerCards
	} else {
		err := uuid.Validate(gameId)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		answerCards, err := model.SelectUnusedCardAnswers(gameId)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		cards = answerCards
	}

	for _, model := range cards {
		card := entity.CardAnswer{
			Id:   model.Id,
			Text: model.Text,
		}

		payload = append(payload, card)
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CardHandler) GetUsedAnswerCards(w http.ResponseWriter, r *http.Request) {
	payload := []entity.CardAnswer{}

	gameId := r.URL.Query().Get("gameId")
	userId := r.URL.Query().Get("userId")
	statusQuery := r.URL.Query().Get("status")

	if gameId == "" || userId == "" || statusQuery == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	status, err := strconv.Atoi(statusQuery)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	cards, err := model.SelectUsedAnswerCards(gameId, userId, model.CardStatus(status))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for _, model := range cards {
		card := entity.CardAnswer{
			Id:   model.Id,
			Text: model.Text,
		}

		payload = append(payload, card)
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CardHandler) GetCardTemplates(w http.ResponseWriter, r *http.Request) {
	payload := []entity.CardTemplate{}
	cards := []model.CardTemplate{}

	gameId := r.URL.Query().Get("gameId")
	if gameId == "" {
		templateCards, err := model.SelectCardTemplates()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		cards = templateCards
	} else {
		err := uuid.Validate(gameId)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		templateCards, err := model.SelectUnusedCardTemplates(gameId)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		cards = templateCards
	}

	for _, model := range cards {
		card := entity.CardTemplate{
			Id:                model.Id,
			Text:              model.Text,
			PlaceholdersCount: model.PlaceholdersCount,
		}

		payload = append(payload, card)
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	var card entity.CardRequest
	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	switch card.Type {
	case entity.CardTypeAnswer:
		err = model.InsertCardAnswer(card.Text)

	case entity.CardTypeTemplate:
		var placeholdersCount int
		for _, c := range card.Text {
			if c == entity.CardPlaceholder {
				placeholdersCount++
			}
		}

		if placeholdersCount < entity.CardPlaceholdersMinCount || placeholdersCount > entity.CardPlaceholdersMaxCount {
			http.Error(w, "incorrect placeholder count", http.StatusBadRequest)
			return
		}

		err = model.InsertCardTemplate(card.Text, placeholdersCount)

	default:
		err = errors.New("undefined card type")
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *CardHandler) UseAnswerCards(w http.ResponseWriter, r *http.Request) {
	var cardsUsed entity.CardsUsedRequest
	err := json.NewDecoder(r.Body).Decode(&cardsUsed)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = model.InsertCardBatchUsed(
		cardsUsed.GameId,
		cardsUsed.UserId,
		cardsUsed.Cards...,
	)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
