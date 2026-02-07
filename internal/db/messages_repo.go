package db

import (
	"context"
	"log"

	"github.com/rafamrslima/distributor/internal/domain"
)

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
