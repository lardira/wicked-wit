package round

import (
	"github.com/lardira/wicked-wit/pkg/card"
	"github.com/lardira/wicked-wit/pkg/response"
)

type Service struct {
	cardService *card.Service
}

func NewRoundService(cardService *card.Service) *Service {
	return &Service{
		cardService: cardService,
	}
}

func (s *Service) GetRounds(gameId string) ([]Round, error) {
	rounds := []Round{}

	roundModels, err := Select(gameId)
	if err != nil {
		return nil, err
	}

	for _, model := range roundModels {
		round := Round{
			Id:       model.Id,
			Position: model.Position,
			GameId:   model.GameId,
			Timed:    response.TimedFromModel(&model.TimedModel),
		}

		if model.WinnerId.Valid {
			round.WinnerId = &model.WinnerId.String
		}

		rounds = append(rounds, round)
	}

	return rounds, nil
}

func (s *Service) AddRound(gameId string) (int, error) {
	rounds, err := s.GetRounds(gameId)
	if err != nil {
		return 0, err
	}

	var position int
	if len(rounds) == 0 {
		position = 0
	} else {
		position = len(rounds) + 1
	}

	templateCard, err := s.cardService.GetRandomTemplateCard(gameId)
	if err != nil {
		return 0, err
	}

	// TODO: check last position for better validation
	newId, err := Insert(
		position,
		gameId,
		templateCard.Id,
	)
	if err != nil {
		return 0, err
	}

	// TODO: fill all of the players' card hands

	return newId, nil
}

func (s *Service) DeleteRound(id int) {
	Delete(id)
}
