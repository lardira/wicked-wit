package round

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/lardira/wicked-wit/internal/db"
	"github.com/lardira/wicked-wit/pkg/response"
)

type RoundModel struct {
	Id       int
	WinnerId sql.NullString
	Position int
	GameId   string

	response.TimedModel
}

func Select(gameId string) ([]RoundModel, error) {
	output := []RoundModel{}

	query := `SELECT
		id,
		winner_id,
		"position",
		game_id,
		created_at,
		updated_at
	FROM
		round
	WHERE
		game_id = $1`

	rows, err := db.Conn.Query(context.Background(), query, gameId)
	if err != nil {
		return output, err
	}
	defer rows.Close()

	for rows.Next() {
		var r RoundModel
		err := rows.Scan(
			&r.Id,
			&r.WinnerId,
			&r.Position,
			&r.GameId,
			&r.CreatedAt,
			&r.UpdatedAt,
		)
		if err != nil {
			return output, err
		}
		output = append(output, r)
	}

	return output, nil
}

func Insert(position int, gameId string, templateId int) (int, error) {
	query := `INSERT INTO round 
		(position, game_id, template_card_id) 
		VALUES (@position, @game_id, @template_card_id)`

	args := pgx.NamedArgs{
		"position":         position,
		"game_id":          gameId,
		"template_card_id": templateId,
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

func Delete(id int) {
	query := "DELETE FROM round WHERE id = $1"
	db.Conn.Exec(context.Background(), query, id)
}
