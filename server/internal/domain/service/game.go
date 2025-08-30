package service

import (
	"github.com/lardira/wicked-wit/internal/domain/entity"
	"github.com/lardira/wicked-wit/internal/domain/interfaces"
	"github.com/lardira/wicked-wit/internal/domain/repository"
	"github.com/lardira/wicked-wit/internal/helper/response"
)

type gameService struct {
	roundService interfaces.RoundService
}

func NewGameService(roundService interfaces.RoundService) *gameService {
	return &gameService{
		roundService: roundService,
	}
}

func (s *gameService) GetGames() ([]entity.Game, error) {
	games := []entity.Game{}
	gameModels, err := repository.SelectGames()
	if err != nil {
		return nil, err
	}

	for _, model := range gameModels {
		game := entity.Game{
			Id:         model.Id,
			Title:      model.Title,
			MaxPlayers: model.MaxPlayers,
			MaxRound:   model.MaxRound,
			UserHostId: model.UserHostId,
			Timed:      response.TimedFromModel(&model.TimedModel),
		}

		games = append(games, game)
	}

	return games, nil
}

func (s *gameService) CreateGame(gameRequest *entity.GameRequest) (string, error) {
	// TODO: add transation

	newId, err := repository.InsertGame(
		gameRequest.Title,
		gameRequest.MaxPlayers,
		gameRequest.MaxRound,
		gameRequest.HostId,
	)
	if err != nil {
		return "", err
	}

	if err := repository.BatchInsertGameUser(newId, gameRequest.Users...); err != nil {
		return "", err
	}

	if _, err := s.roundService.AddRound(newId); err != nil {
		return "", err
	}

	return newId, nil
}

func (s *gameService) DeleteGame(id string) {
	repository.DeleteGame(id)
}
