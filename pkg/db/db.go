package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	//"github.com/jmoiron/sqlx"
	//_ "github.com/lib/pq"
	//"github.com/jmoiron/sqlx"
)

type DB struct {
	//Client *sqlx.DB
	Client *pgxpool.Pool
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
	return d.Close()
}

func get(connStr string) (*pgxpool.Pool, error) {
	//db, err := sql.Open("postgres", connStr)
	conn, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
