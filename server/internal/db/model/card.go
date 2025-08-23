package model

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lardira/wicked-wit/internal/db"
)

type CardTemplate struct {
	Id                int
	PlaceholdersCount int
	Text              string
}

type CardAnswer struct {
	Id   int
	Text string
}

func SelectCardAnswers() ([]CardAnswer, error) {
	output := []CardAnswer{}

	query := `SELECT id, text FROM answer_card`

	rows, err := db.Conn.Query(context.Background(), query)
	if err != nil {
		return output, err
	}
	defer rows.Close()

	for rows.Next() {
		var c CardAnswer
		err := rows.Scan(
			&c.Id,
			&c.Text,
		)
		if err != nil {
			return output, err
		}
		output = append(output, c)
	}

	return output, nil
}

func SelectUnusedCardAnswers(gameId string) ([]CardAnswer, error) {
	output := []CardAnswer{}

	query := `SELECT 
		ac.id,
		ac."text"
	FROM
		answer_card ac
	WHERE
		NOT EXISTS(
			SELECT guc.answer_card_id, guc.game_id  
				FROM game_used_card guc 
			WHERE guc.answer_card_id = ac.id AND guc.game_id = @gameId
		)`

	args := pgx.NamedArgs{
		"gameId": gameId,
	}

	rows, err := db.Conn.Query(context.Background(), query, args)
	if err != nil {
		return output, err
	}
	defer rows.Close()

	for rows.Next() {
		var c CardAnswer
		err := rows.Scan(
			&c.Id,
			&c.Text,
		)
		if err != nil {
			return output, err
		}
		output = append(output, c)
	}

	return output, nil
}

func SelectCardTemplates() ([]CardTemplate, error) {
	output := []CardTemplate{}

	query := `SELECT id, text, placeholders_count FROM template_card`

	rows, err := db.Conn.Query(context.Background(), query)
	if err != nil {
		return output, err
	}
	defer rows.Close()

	for rows.Next() {
		var c CardTemplate
		err := rows.Scan(
			&c.Id,
			&c.Text,
			&c.PlaceholdersCount,
		)
		if err != nil {
			return output, err
		}
		output = append(output, c)
	}

	return output, nil
}

func SelectUnusedCardTemplates(gameId string) ([]CardTemplate, error) {
	output := []CardTemplate{}

	query := `SELECT
		tc.id, tc.text, tc.placeholders_count
	FROM
		template_card tc
	WHERE NOT EXISTS(
				SELECT r.template_card_id, r.game_id 
					FROM round r 
					WHERE r.template_card_id = tc.id AND r.game_id = @gameId
	)`

	args := pgx.NamedArgs{
		"gameId": gameId,
	}

	rows, err := db.Conn.Query(context.Background(), query, args)
	if err != nil {
		return output, err
	}
	defer rows.Close()

	for rows.Next() {
		var c CardTemplate
		err := rows.Scan(
			&c.Id,
			&c.Text,
			&c.PlaceholdersCount,
		)
		if err != nil {
			return output, err
		}
		output = append(output, c)
	}

	return output, nil
}

func InsertCardAnswer(text string) error {
	query := `INSERT INTO answer_card (text) VALUES ($1)`

	_, err := db.Conn.Exec(
		context.Background(),
		query,
		text,
	)
	return err
}

func InsertCardTemplate(text string, placeholders_count int) error {
	query := `INSERT INTO answer_card (text, placeholders_count)
	 	VALUES (@text, @placeholders_count)`

	args := pgx.NamedArgs{
		"text":               text,
		"placeholders_count": placeholders_count,
	}

	_, err := db.Conn.Exec(
		context.Background(),
		query,
		args,
	)
	return err
}
