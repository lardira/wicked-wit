package model

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lardira/wicked-wit/internal/db"
)

type Game struct {
	Id           string
	Title        string
	MaxPlayers   uint
	CurrentRound uint
	MaxRound     uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func SelectGames() ([]Game, error) {
	output := make([]Game, 0)

	query := "SELECT id, title, max_players, max_round, current_round, created_at, updated_at FROM game"
	rows, err := db.Conn.Query(context.Background(), query)
	if err != nil {
		return output, err
	}
	defer rows.Close()

	for rows.Next() {
		var g Game
		err := rows.Scan(
			&g.Id,
			&g.Title,
			&g.MaxPlayers,
			&g.MaxRound,
			&g.CurrentRound,
			&g.CreatedAt,
			&g.UpdatedAt,
		)
		if err != nil {
			log.Fatal(err)
		}
		output = append(output, g)
	}

	return output, nil
}

func InsertGame(title string, maxPlayers uint, maxRound uint, currentLeadingPlayerId string) (string, error) {
	newGame := Game{
		Id:         uuid.NewString(),
		Title:      title,
		MaxPlayers: maxPlayers,
		MaxRound:   maxRound,
	}

	return newGame.Id, nil
}

func DeleteGame(id string) {
	fmt.Println("TODO: DELETE GAME")
}
