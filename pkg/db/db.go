package db

import (
	"context"

	//_ "github.com/lib/pq"

	_ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	Client *sqlx.DB
}

func Get(connStr string) (*DB, error) {
	db, err := get(connStr)
	if err != nil {
		return nil, err
	}

	return &DB{
		Client: db,
	}, nil
}

func (d *DB) Close() error {
	return d.Client.Close()
}

func get(connStr string) (*sqlx.DB, error) {
	//db, err := sql.Open("postgres", connStr)
	conn, err := sqlx.ConnectContext(context.Background(), "postgres", connStr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
