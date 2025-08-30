package service

import (
	"github.com/lardira/wicked-wit/internal/domain/entity"
	"github.com/lardira/wicked-wit/internal/domain/repository"
	"github.com/lardira/wicked-wit/internal/helper"
)

type cardService struct{}

func NewCardService() *cardService {
	return &cardService{}
}

func (s *cardService) UseCards(gameId string, userId string, cardIds ...int) error {
	return repository.BatchInsertCardUsed(gameId, userId, cardIds...)
}

func (s *cardService) GetCards(gameId string, userId string) ([]entity.Card, error) {
	cards := []entity.Card{}

	cardModels, err := repository.SelectUsedAnswerCards(gameId, userId, repository.CardStatusInUse)
	if err != nil {
		return nil, err
	}

	for _, model := range cardModels {
		cards = append(cards,
			entity.Card{
				Id:   model.Id,
				Text: model.Text,
			},
		)
	}

	return cards, nil
}

func (s *cardService) PlayCards(roundId int, userId string, cardIds ...int) (int, error) {
	//TODO: check if template is able to contain all of the cards
	//TODO: transaction

	answerId, err := repository.InsertUserAnswer(userId, roundId)
	if err != nil {
		return 0, err
	}

	for i, cardId := range cardIds {
		if err := repository.InsertPlayedCard(answerId, cardId, i); err != nil {
			return 0, err
		}

		if err := repository.UpdateUsedCardStatus(cardId, repository.CardStatusPlayed); err != nil {
			return 0, err
		}
	}

	return answerId, nil
}

func (s *cardService) GetUnusedTemplateCards(gameId string) ([]entity.TemplateCard, error) {
	cards := []entity.TemplateCard{}

	cardModels, err := repository.SelectUnusedTemplateCards(gameId)
	if err != nil {
		return nil, err
	}

	for _, model := range cardModels {
		cards = append(cards,
			entity.TemplateCard{
				Id:                model.Id,
				Text:              model.Text,
				PlaceholdersCount: model.PlaceholdersCount,
			},
		)
	}

	return cards, nil
}

func (s *cardService) GetRandomTemplateCard(gameId string) (*entity.TemplateCard, error) {
	cards, err := s.GetUnusedTemplateCards(gameId)
	if err != nil {
		return nil, err
	}

	randomCard := helper.RandomItem(cards)
	return &randomCard, nil
}
