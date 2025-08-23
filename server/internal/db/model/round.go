package model

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/lardira/wicked-wit/internal/db"
)

type Round struct {
	Id       int
	WinnerId sql.NullString
	Position int
	GameId   string

	Timed
}

func SelectRounds() ([]Round, error) {
	output := []Round{}

	query := `SELECT 
		id, winner_id, position, created_at, updated_at, game_id
		FROM round`
	rows, err := db.Conn.Query(context.Background(), query)
	if err != nil {
		return output, err
	}
	defer rows.Close()

	for rows.Next() {
		var r Round
		err := rows.Scan(
			&r.Id,
			&r.WinnerId,
			&r.Position,
			&r.CreatedAt,
			&r.UpdatedAt,
			&r.GameId,
		)
		if err != nil {
			return output, err
		}
		output = append(output, r)
	}

	return output, nil
}

func InsertRound(position int, gameId string) (int, error) {
	query := `INSERT INTO round (position, game_id) VALUES (@position, @game_id)`

	args := pgx.NamedArgs{
		"position": position,
		"game_id":  gameId,
	}

	_, err := db.Conn.Exec(
		context.Background(),
		query,
		args,
	)
	if err != nil {
		return 0, err
	}

	var newRoundId int
	query = `SELECT id FROM round WHERE game_id = $1 ORDER BY position DESC LIMIT 1`

	err = db.Conn.QueryRow(context.Background(), query, gameId).Scan(&newRoundId)
	if err != nil {
		return 0, err
	}

	return newRoundId, nil
}

func DeleteRound(id int) {
	query := "DELETE FROM round WHERE id = $1"
	db.Conn.Exec(context.Background(), query, id)
}
