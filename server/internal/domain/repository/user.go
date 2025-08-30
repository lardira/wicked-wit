package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lardira/wicked-wit/internal/db"
	"github.com/lardira/wicked-wit/internal/helper/response"
)

type UserModel struct {
	Id         string
	Username   string
	ProfileImg sql.NullString

	response.TimedModel
}

func SelectUser(userId string) (UserModel, error) {
	var output UserModel

	query := `SELECT 
	id, username, profile_img, created_at, updated_at 
	FROM users WHERE id=$1`

	err := db.Conn.QueryRow(context.Background(), query, userId).Scan(
		&output.Id,
		&output.Username,
		&output.ProfileImg,
		&output.CreatedAt,
		&output.UpdatedAt,
	)
	if err != nil {
		return output, err
	}

	return output, nil
}

func UpdateUserImg(id string, imgUrl string) error {

	query := `UPDATE users
			SET profile_img=@profile_img
			WHERE id=@id`

	args := pgx.NamedArgs{
		"id":          id,
		"profile_img": imgUrl,
	}

	_, err := db.Conn.Exec(
		context.Background(),
		query,
		args,
	)
	if err != nil {
		return err
	}

	return nil
}

func InsertUser(username string, password string) (string, error) {
	newUserId := uuid.NewString()

	query := `INSERT INTO users 
		(id, username, password)
		VALUES (@id, @username, @password)`

	args := pgx.NamedArgs{
		"id":       newUserId,
		"username": username,
		"password": password,
	}

	_, err := db.Conn.Exec(
		context.Background(),
		query,
		args,
	)
	if err != nil {
		return "", err
	}

	return newUserId, nil
}

func DeleteUser(id string) {
	query := "DELETE FROM users WHERE id = $1"
	db.Conn.Exec(context.Background(), query, id)
}

func InsertUserAnswer(userId string, roundId int) (int, error) {
	query := `INSERT INTO user_answer 
		(user_id, round_id)
		VALUES (@user_id, @round_id)
	`

	args := pgx.NamedArgs{
		"user_id":  userId,
		"round_id": roundId,
	}

	_, err := db.Conn.Exec(
		context.Background(),
		query,
		args,
	)
	if err != nil {
		return 0, err
	}

	var answerId int
	query = `SELECT
			ua.id
		FROM
			user_answer ua
		JOIN round r ON
			r.id = ua.round_id
		WHERE
			r.id = $1`

	err = db.Conn.QueryRow(context.Background(), query, roundId).Scan(&answerId)
	if err != nil {
		return 0, err
	}

	return answerId, nil
}
