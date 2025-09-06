package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rafamrslima/distributor/internal/domain"
)

func Connect() (*pgxpool.Pool, error) {
	dsn := "postgres://admin:mypassword@localhost:5432/mydb?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Println("error to connect to database", err)
		return nil, err
	}
	return pool, nil
}

func SaveReceivedMessages(message domain.Message) {
	//todo
}
