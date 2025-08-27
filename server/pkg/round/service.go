package round

import (
	"github.com/lardira/wicked-wit/pkg/response"
)

type Service struct{}

func (s *Service) GetRounds() ([]Round, error) {
	rounds := []Round{}

	roundModels, err := Select()
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
	rounds, err := s.GetRounds()
	if err != nil {
		return 0, err
	}

	newId, err := Insert(
		len(rounds)+1,
		gameId,
	)
	if err != nil {
		return 0, err
	}

	return newId, nil
}

func (s *Service) DeleteRound(id int) {
	Delete(id)
}
