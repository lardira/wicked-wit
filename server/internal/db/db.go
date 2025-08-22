package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	Conn *pgxpool.Pool
)

func Init(url string) error {
	if Conn != nil {
		return nil
	}

	connConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return fmt.Errorf("cannot parse config from %w\n", err)
	}

	pgPool, err := pgxpool.NewWithConfig(
		context.Background(),
		connConfig,
	)
	Conn = pgPool

	return err
}

func Close() {
	Conn.Close()
}
