package model

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lardira/wicked-wit/internal/db"
)

type Game struct {
	Id           string
	Title        string
	MaxPlayers   uint
	CurrentRound sql.NullInt64
	MaxRound     uint

	Timed
}

func SelectGames() ([]Game, error) {
	output := []Game{}

	query := `SELECT 
			id, title, max_players, max_round, current_round, created_at, updated_at 
		FROM game`

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
			return output, err
		}
		output = append(output, g)
	}

	return output, nil
}

func InsertGame(title string, maxPlayers uint, maxRound uint) (string, error) {
	newGameId := uuid.NewString()

	query := `INSERT INTO game 
		(id, title, max_players, max_round)
		VALUES (@id, @title, @max_players, @max_round)`

	args := pgx.NamedArgs{
		"id":          newGameId,
		"title":       title,
		"max_players": maxPlayers,
		"max_round":   maxRound,
	}

	_, err := db.Conn.Exec(
		context.Background(),
		query,
		args,
	)
	if err != nil {
		return "", err
	}

	return newGameId, nil
}

func DeleteGame(id string) {
	query := "DELETE FROM game WHERE id = $1"
	db.Conn.Exec(context.Background(), query, id)
}
