package card

import (
	"github.com/lardira/wicked-wit/internal/helper"
	"github.com/lardira/wicked-wit/pkg/user"
)

type Service struct{}

func (s *Service) UseCards(gameId string, userId string, cardIds ...int) error {
	return BatchInsertCardUsed(gameId, userId, cardIds...)
}

func (s *Service) GetCards(gameId string, userId string) ([]Card, error) {
	cards := []Card{}

	cardModels, err := SelectUsedAnswerCards(gameId, userId, CardStatusInUse)
	if err != nil {
		return nil, err
	}

	for _, model := range cardModels {
		cards = append(cards,
			Card{
				Id:   model.Id,
				Text: model.Text,
			},
		)
	}

	return cards, nil
}

func (s *Service) PlayCards(roundId int, userId string, cardIds ...int) (int, error) {
	//TODO: check if template is able to contain all of the cards
	//TODO: transaction

	answerId, err := user.InsertUserAnswer(userId, roundId)
	if err != nil {
		return 0, err
	}

	for i, cardId := range cardIds {
		if err := InsertPlayedCard(answerId, cardId, i); err != nil {
			return 0, err
		}

		if err := UpdateUsedCardStatus(cardId, CardStatusPlayed); err != nil {
			return 0, err
		}
	}

	return answerId, nil
}

func (s *Service) GetUnusedTemplateCards(gameId string) ([]TemplateCard, error) {
	cards := []TemplateCard{}

	cardModels, err := SelectUnusedTemplateCards(gameId)
	if err != nil {
		return nil, err
	}

	for _, model := range cardModels {
		cards = append(cards,
			TemplateCard{
				Id:                model.Id,
				Text:              model.Text,
				PlaceholdersCount: model.PlaceholdersCount,
			},
		)
	}

	return cards, nil
}

func (s *Service) GetRandomTemplateCard(gameId string) (*TemplateCard, error) {
	cards, err := s.GetUnusedTemplateCards(gameId)
	if err != nil {
		return nil, err
	}

	randomCard := helper.RandomItem(cards)
	return &randomCard, nil
}
