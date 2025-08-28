package card

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lardira/wicked-wit/internal/db"
)

type CardStatus uint

const (
	CardStatusInUse   CardStatus = 0 //in use
	CardStatusDropped CardStatus = 1 //dropped from hand
	CardStatusPlayed  CardStatus = 2 //played in round
)

type TemplateCardModel struct {
	Id                int
	PlaceholdersCount int
	Text              string
}

type AnswerCardModel struct {
	Id   int
	Text string
}

func SelectUnusedAnswerCards(gameId string) ([]AnswerCardModel, error) {
	output := []AnswerCardModel{}

	query := `SELECT 
		ac.id,
		ac."text"
	FROM
		answer_card ac
	WHERE
		NOT EXISTS(
			SELECT guc.answer_card_id, guc.game_id  
				FROM game_used_card guc 
			WHERE guc.answer_card_id = ac.id AND guc.game_id = @game_id
		)`

	args := pgx.NamedArgs{
		"game_id": gameId,
	}

	rows, err := db.Conn.Query(context.Background(), query, args)
	if err != nil {
		return output, err
	}
	defer rows.Close()

	for rows.Next() {
		var c AnswerCardModel
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

func SelectUsedAnswerCards(gameId string, userId string, status CardStatus) ([]AnswerCardModel, error) {
	output := []AnswerCardModel{}

	query := `SELECT
		ac.id,
		ac."text"
	FROM
		answer_card ac
	JOIN game_used_card guc ON
		guc.answer_card_id = ac.id
	WHERE
		guc.status = @status
		AND guc.user_id = @user_id
		AND guc.game_id = @game_id`

	args := pgx.NamedArgs{
		"status":  status,
		"user_id": userId,
		"game_id": gameId,
	}

	rows, err := db.Conn.Query(context.Background(), query, args)
	if err != nil {
		return output, err
	}
	defer rows.Close()

	for rows.Next() {
		var c AnswerCardModel
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

func SelectCardTemplates() ([]TemplateCardModel, error) {
	output := []TemplateCardModel{}

	query := `SELECT id, text, placeholders_count FROM template_card`

	rows, err := db.Conn.Query(context.Background(), query)
	if err != nil {
		return output, err
	}
	defer rows.Close()

	for rows.Next() {
		var c TemplateCardModel
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

func SelectUnusedCardTemplates(gameId string) ([]TemplateCardModel, error) {
	output := []TemplateCardModel{}

	query := `SELECT
		tc.id, tc.text, tc.placeholders_count
	FROM
		template_card tc
	WHERE NOT EXISTS(
				SELECT r.template_card_id, r.game_id 
					FROM round r 
					WHERE r.template_card_id = tc.id AND r.game_id = @game_id
	)`

	args := pgx.NamedArgs{
		"game_id": gameId,
	}

	rows, err := db.Conn.Query(context.Background(), query, args)
	if err != nil {
		return output, err
	}
	defer rows.Close()

	for rows.Next() {
		var c TemplateCardModel
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
	query := `INSERT INTO template_card (text, placeholders_count)
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

func InsertCardBatchUsed(gameId string, userId string, cardIds ...int) error {
	batch := &pgx.Batch{}

	query := `INSERT INTO game_used_card
		(game_id, answer_card_id, status, user_id)
		VALUES (@game_id, @answer_card_id, @status, @user_id)`

	for _, cardId := range cardIds {
		batch.Queue(query, pgx.NamedArgs{
			"game_id":        gameId,
			"answer_card_id": cardId,
			"status":         CardStatusInUse,
			"user_id":        userId,
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

func InsertPlayedCard(answerId int, usedCardId int, placeholderIndex int) error {
	query := `INSERT INTO game_user_played_card
		(user_answer_id, game_used_card_id, placeholder_index)
		VALUES(@user_answer_id, @game_used_card_id, @placeholder_index)`

	args := pgx.NamedArgs{
		"user_answer_id":    answerId,
		"game_used_card_id": usedCardId,
		"placeholder_index": placeholderIndex,
	}

	_, err := db.Conn.Exec(
		context.Background(),
		query,
		args,
	)
	return err
}

func UpdateUsedCardStatus(usedCardId int, status CardStatus) error {
	query := `UPDATE
		game_used_card
	SET
		status = @status
	WHERE
		id = @id`

	args := pgx.NamedArgs{
		"id":     usedCardId,
		"status": status,
	}

	_, err := db.Conn.Exec(
		context.Background(),
		query,
		args,
	)
	return err
}
