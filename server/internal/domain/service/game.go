package service

import (
	"errors"

	"github.com/lardira/wicked-wit/internal/domain/entity"
	"github.com/lardira/wicked-wit/internal/domain/interfaces"
	"github.com/lardira/wicked-wit/internal/domain/repository"
	"github.com/lardira/wicked-wit/internal/helper"
	"github.com/lardira/wicked-wit/internal/helper/response"
)

const (
	maxPlayerHandSize = 5
)

type gameService struct {
	roundService interfaces.RoundService
	cardService  interfaces.CardService
}

func NewGameService(roundService interfaces.RoundService, cardService interfaces.CardService) *gameService {
	return &gameService{
		roundService: roundService,
		cardService:  cardService,
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
			Id:           model.Id,
			Title:        model.Title,
			MaxPlayers:   model.MaxPlayers,
			MaxRound:     model.MaxRound,
			CurrentRound: model.CurrentRound,
			Status:       model.Status,
			UserHostId:   model.UserHostId,
			Timed:        response.TimedFromModel(&model.TimedModel),
		}

		games = append(games, game)
	}

	return games, nil
}

func (s *gameService) GetGame(gameId string) (entity.Game, error) {

	model, err := repository.SelectGame(gameId)
	if err != nil {
		return entity.Game{}, err
	}

	game := entity.Game{
		Id:           model.Id,
		Title:        model.Title,
		MaxPlayers:   model.MaxPlayers,
		MaxRound:     model.MaxRound,
		CurrentRound: model.CurrentRound,
		Status:       model.Status,
		UserHostId:   model.UserHostId,
		Timed:        response.TimedFromModel(&model.TimedModel),
	}

	return game, nil

}

func (s *gameService) FillUserHand(gameId string, userId string) error {
	playerCards, err := s.cardService.GetCards(gameId, userId)
	if err != nil {
		return err
	}
	currentHandSize := len(playerCards)
	if currentHandSize == maxPlayerHandSize {
		return nil
	}

	unusedCards, err := s.cardService.GetUnusedAnswerCards(gameId)
	if err != nil {
		return err
	}

	needCards := helper.MinInt(
		maxPlayerHandSize-currentHandSize,
		len(unusedCards), //may not be enough in db
	)

	randomCards := helper.RandomSubset(unusedCards, needCards)
	cardIds := make([]int, 0, needCards)

	for _, c := range randomCards {
		cardIds = append(cardIds, c.Id)
	}

	if err := s.cardService.UseCards(gameId, userId, cardIds...); err != nil {
		return err
	}

	return nil
}

func (s *gameService) AppendRound(gameId string, templateCardid int) (roundId int, err error) {
	currentGame, err := s.GetGame(gameId)
	if err != nil {
		return 0, err
	} else if currentGame.CurrentRound >= currentGame.MaxRound {
		return 0, errors.New("max round reached")
	}

	templateCard, err := s.cardService.GetRandomTemplateCard(currentGame.Id)
	if err != nil {
		return 0, err
	}

	id, err := s.roundService.AddRound(currentGame.Id, templateCard.Id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *gameService) CreateGame(gameRequest *entity.GameRequest) (string, error) {
	// TODO: add transaction

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

	for _, uId := range gameRequest.Users {
		err := s.FillUserHand(newId, uId)
		if err != nil {
			return "", err
		}
	}

	return newId, nil
}

func (s *gameService) DeleteGame(id string) {
	repository.DeleteGame(id)
}
