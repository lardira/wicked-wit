package model

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lardira/wicked-wit/internal/db"
)

type User struct {
	Id         string
	Username   string
	ProfileImg sql.NullString

	Timed
}

func SelectUser(userId string) (User, error) {
	var output User

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
