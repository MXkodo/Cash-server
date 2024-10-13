package postgresql

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewDb(connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("unable to connect to database: %v\n", err)
		return nil, err
	}

	fmt.Println("Successfully connected to the database")
	return pool, nil
}

func CloseDb(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}
