package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rafamrslima/distributor/internal/domain"
)

func connect() (*pgxpool.Pool, error) {
	connString := os.Getenv("DATABASE_CONNECTION_STRING")
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	config.MaxConns = 20
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Println("Error to connect to database.", err)
		return nil, err
	}
	return pool, nil
}

func SaveReceivedMessages(message domain.Message) error {
	pool, err := connect()
	if err != nil {
		return err
	}
	defer pool.Close()

	ctx := context.Background()

	_, err = pool.Exec(ctx, `INSERT INTO messages (client_name, report_name, email, content, messageReceivedAt) VALUES ($1, $2, $3, $4, $5)`,
		message.ClientName, message.ReportName, message.Email, message.Content, message.MessageReceivedAt)

	if err != nil {
		log.Println("Error when inserting row into database.", err.Error())
		return err
	}

	log.Println("Row inserted successfully.")
	return nil
}
