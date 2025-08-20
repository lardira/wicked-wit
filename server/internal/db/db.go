package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func Init(url string) error {
	c, err := pgx.Connect(context.Background(), url)
	Conn = c
	return err
}

func Close() {
	Conn.Close(context.Background())
}
