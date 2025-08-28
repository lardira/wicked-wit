package game

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lardira/wicked-wit/internal/db"
	"github.com/lardira/wicked-wit/pkg/response"
)

type GameModel struct {
	Id         string
	Title      string
	MaxPlayers uint
	MaxRound   uint
	UserHostId string

	response.TimedModel
}

func SelectGames() ([]GameModel, error) {
	output := []GameModel{}

	query := `SELECT 
			id, title, max_players, max_round, user_host_id, created_at, updated_at 
		FROM game`

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

func InsertGame(title string, maxPlayers uint, maxRound uint, userHostId string) (string, error) {
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
