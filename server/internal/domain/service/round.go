package service

import (
	"github.com/lardira/wicked-wit/internal/domain/entity"
	"github.com/lardira/wicked-wit/internal/domain/repository"
	"github.com/lardira/wicked-wit/internal/helper/response"
)

type roundService struct {
}

func NewRoundService() *roundService {
	return &roundService{}
}

func (s *roundService) GetRounds(gameId string) ([]entity.Round, error) {
	rounds := []entity.Round{}

	roundModels, err := repository.SelectRounds(gameId)
	if err != nil {
		return nil, err
	}

	for _, model := range roundModels {
		round := entity.Round{
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

func (s *roundService) AddRound(gameId string, templateCardId int) (int, error) {

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

	// TODO: check last position for better validation
	newId, err := repository.InsertRound(
		position,
		gameId,
		templateCardId,
	)
	if err != nil {
		return 0, err
	}

	return newId, nil
}

func (s *roundService) DeleteRound(id int) {
	repository.DeleteRound(id)
}
