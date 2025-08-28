package game

import (
	"github.com/lardira/wicked-wit/pkg/response"
	"github.com/lardira/wicked-wit/pkg/round"
)

type Service struct {
	roundService *round.Service
}

func NewGameService(roundService *round.Service) *Service {
	return &Service{
		roundService: roundService,
	}
}

func (s *Service) GetGames() ([]Game, error) {
	games := []Game{}
	gameModels, err := SelectGames()
	if err != nil {
		return nil, err
	}

	for _, model := range gameModels {
		game := Game{
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

func (s *Service) CreateGame(gameRequest *GameRequest) (string, error) {
	// TODO: add transation

	newId, err := InsertGame(
		gameRequest.Title,
		gameRequest.MaxPlayers,
		gameRequest.MaxRound,
		gameRequest.HostId,
	)
	if err != nil {
		return "", err
	}

	if err := BatchInsertGameUser(newId, gameRequest.Users...); err != nil {
		return "", err
	}

	if _, err := s.roundService.AddRound(newId); err != nil {
		return "", err
	}

	return newId, nil
}

func (s *Service) DeleteGame(id string) {
	DeleteGame(id)
}
