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

	_, err = pool.Exec(ctx,
		`INSERT INTO messages (client_name, report_name, client_email, message_received_at) VALUES ($1, $2, $3, $4)`,
		message.ClientName, message.ReportName, message.Email, message.MessageReceivedAt)

	if err != nil {
		return err
	}

	log.Println("Row inserted successfully.")
	return nil
}

func GetReportInfo(clientEmail string, reportName string) ([]domain.Results, error) {
	pool, err := connect()
	if err != nil {
		return nil, err
	}
	defer pool.Close()

	ctx := context.Background()

	rows, err := pool.Query(ctx,
		`SELECT client_email, report_name, gains, losses, info_date 
		FROM investment_results WHERE info_date::date = CURRENT_DATE AND client_email = $1 AND report_name = $2`,
		clientEmail, reportName)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Results

	for rows.Next() {
		var res domain.Results
		if err := rows.Scan(&res.ClientEmail, &res.ReportName, &res.Gains, &res.Losses, &res.InfoDate); err != nil {
			log.Fatal(err)
		}
		results = append(results, res)
	}

	return results, nil
}
