package card

import (
	"github.com/lardira/wicked-wit/pkg/user"
)

type Service struct{}

func (s *Service) UseCards(gameId string, userId string, cardIds ...int) error {
	return InsertCardBatchUsed(gameId, userId, cardIds...)
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
