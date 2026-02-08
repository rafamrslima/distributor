package repositories

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rafamrslima/distributor/internal/db"
	"github.com/rafamrslima/distributor/internal/domain"
)

func SaveReceivedMessages(ctx context.Context, message domain.Message) error {
	pool, err := db.GetDB()
	if err != nil {
		return fmt.Errorf("failed to get database pool: %w", err)
	}

	// Create context with timeout for the database operation
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = pool.Exec(ctx,
		`INSERT INTO messages (client_name, report_name, client_email, message_received_at) VALUES ($1, $2, $3, $4)`,
		message.ClientName, message.ReportName, message.Email, message.MessageReceivedAt)

	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}

	log.Println("Row inserted successfully.")
	return nil
}
