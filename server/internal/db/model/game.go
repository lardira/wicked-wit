package model

import (
	"slices"

	"github.com/google/uuid"
)

var (
	games []Game = []Game{
		{
			Id:                     "gipoujqaowiugjopaiwejf",
			Title:                  "Game 1",
			MaxPlayers:             8,
			CurrentRound:           0,
			MaxRounds:              5,
			CurrentLeadingPlayerId: "gjiwropqgjiwergjwei",
		},
		{
			Id:                     "gipoujqaowiugjopaiwgawergejf",
			Title:                  "Game 2",
			MaxPlayers:             3,
			CurrentRound:           0,
			MaxRounds:              5,
			CurrentLeadingPlayerId: "gjiwropqgjiwergjwei",
		},
		{
			Id:                     "gipoujahwhwehwehwerhqaowiugjopaiwejf",
			Title:                  "Game 3",
			MaxPlayers:             10,
			CurrentRound:           0,
			MaxRounds:              5,
			CurrentLeadingPlayerId: "gjiwropqgjiwergjwei",
		},
	}
)

type Game struct {
	Id                     string
	Title                  string
	MaxPlayers             uint
	CurrentRound           uint
	MaxRounds              uint
	CurrentLeadingPlayerId string
}

func SelectGames() []Game {
	return games
}

func InsertGame(title string, maxPlayers uint, maxRounds uint, currentLeadingPlayerId string) (string, error) {
	newGame := Game{
		Id:                     uuid.NewString(),
		Title:                  title,
		MaxPlayers:             maxPlayers,
		MaxRounds:              maxRounds,
		CurrentLeadingPlayerId: currentLeadingPlayerId,
	}

	games = append(games, newGame)

	return newGame.Id, nil
}

func DeleteGame(id string) {
	for i, g := range games {
		if g.Id == id {
			games = slices.Delete(games, i, i+1)
			return
		}
	}
}
