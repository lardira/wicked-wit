package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lardira/wicked-wit/internal/db"
	"github.com/lardira/wicked-wit/internal/helper/response"
)

type GameStatus uint

const (
	GameStatusStarted    = 0
	GameStatusInProgress = 1
	GameStatusEnded      = 2
)

type GameModel struct {
	Id           string
	Title        string
	MaxPlayers   int
	MaxRound     int
	CurrentRound int
	Status       GameStatus
	UserHostId   string

	response.TimedModel
}

func SelectGame(id string) (GameModel, error) {
	var g GameModel
	query := `SELECT
			g.id,
			g.title,
			g.max_players,
			g.max_round,
			COUNT(r.id) AS current_round,
			g.status,
			g.user_host_id,
			g.created_at,
			g.updated_at
		FROM
			game g
		JOIN round r ON
			r.game_id = g.id
		WHERE g.id = $1
		GROUP BY
			g.id`

	err := db.Conn.QueryRow(context.Background(), query, id).Scan(
		&g.Id,
		&g.Title,
		&g.MaxPlayers,
		&g.MaxRound,
		&g.CurrentRound,
		&g.Status,
		&g.UserHostId,
		&g.CreatedAt,
		&g.UpdatedAt,
	)
	if err != nil {
		return g, err
	}

	return g, nil
}

func SelectGames() ([]GameModel, error) {
	output := []GameModel{}

	query := `SELECT
			g.id,
			g.title,
			g.max_players,
			g.max_round,
			COUNT(r.id) AS current_round,
			g.status,
			g.user_host_id,
			g.created_at,
			g.updated_at
		FROM
			game g
		JOIN round r ON
			r.game_id = g.id
		GROUP BY
			g.id`

	rows, err := db.Conn.Query(context.Background(), query)
	if err != nil {
		return output, err
	}
	defer rows.Close()

	for rows.Next() {
		var g GameModel
		err := rows.Scan(
			&g.Id,
			&g.Title,
			&g.MaxPlayers,
			&g.MaxRound,
			&g.CurrentRound,
			&g.Status,
			&g.UserHostId,
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

func InsertGame(title string, maxPlayers int, maxRound int, userHostId string) (string, error) {
	newGameId := uuid.NewString()

	query := `INSERT INTO game 
		(id, title, max_players, max_round, user_host_id)
		VALUES (@id, @title, @max_players, @max_round, @user_host_id)`

	args := pgx.NamedArgs{
		"id":           newGameId,
		"title":        title,
		"max_players":  maxPlayers,
		"max_round":    maxRound,
		"user_host_id": userHostId,
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

func BatchInsertGameUser(gameId string, userIds ...string) error {
	batch := &pgx.Batch{}

	query := `INSERT INTO
		games_users (game_id, user_id)
		VALUES (@game_id, @user_id)`

	for _, userId := range userIds {
		batch.Queue(query, pgx.NamedArgs{
			"game_id": gameId,
			"user_id": userId,
		})
	}

	batchResults := db.Conn.SendBatch(context.Background(), batch)
	defer batchResults.Close()

	for range batch.QueuedQueries {
		_, err := batchResults.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}
