package card

import "github.com/lardira/wicked-wit/pkg/user"

type Service struct{}

func (s *Service) UseCards(gameId string, userId string, cardIds ...int) error {
	return InsertCardBatchUsed(gameId, userId, cardIds...)
}

func (s *Service) PlayCards(roundId int, userId string, cardIds ...int) error {
	//TODO: check if template is able to contain all of the cards
	//TODO: transaction

	answerId, err := user.InsertUserAnswer(userId, roundId)
	if err != nil {
		return err
	}

	for i, cardId := range cardIds {
		if err := InsertPlayedCard(answerId, cardId, i); err != nil {
			return err
		}

		if err := UpdateUsedCardStatus(cardId, CardStatusPlayed); err != nil {
			return err
		}
	}

	return nil
}
