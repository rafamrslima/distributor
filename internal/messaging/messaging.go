package messaging

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/joho/godotenv"
	"github.com/rafamrslima/distributor/internal/core"
)

const (
	batchSize     = 50              // messages fetched per poll
	maxWorkers    = 200             // max messages processed concurrently
	receiveWait   = 5 * time.Second // per-poll timeout
	settleTimeout = 3 * time.Second // per-complete timeout
)

func getClient() (*azservicebus.Client, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env.")
	}

	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")

	if connStr == "" {
		log.Println("No connection string found for service bus.")
		return nil, errors.New("connection string not found")
	}

	client, err := azservicebus.NewClientFromConnectionString(connStr, nil)
	if err != nil {
		return nil, errors.New("SERVICEBUS_QUEUE is empty")
	}

	return client, nil
}

func getQueueName() (string, error) {
	queue := os.Getenv("SERVICEBUS_QUEUE")
	if queue == "" {
		log.Println("SERVICEBUS_QUEUE config is not valid.")
		return "", errors.New("SERVICEBUS_QUEUE config is not valid")
	}
	return queue, nil
}

func StartMessageListener(ctx context.Context) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	queue, err := getQueueName()
	if err != nil {
		return err
	}

	defer client.Close(context.Background())

	receiver, err := client.NewReceiverForQueue(queue, nil)
	if err != nil {
		return fmt.Errorf("new receiver: %w", err)
	}

	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = receiver.Close(closeCtx)
	}()

	fmt.Println("Listening for messages... (Ctrl+C to quit)")

	// Worker pool
	jobs := make(chan *azservicebus.ReceivedMessage, 2*maxWorkers)
	var wg sync.WaitGroup
	for range maxWorkers {
		wg.Go(func() {
			for msg := range jobs {
				if err := core.Handle(ctx, msg); err != nil {
					// abandon or dead-letter
					abandonCtx, cancel := context.WithTimeout(ctx, settleTimeout)
					_ = receiver.DeadLetterMessage(abandonCtx, msg, nil)
					cancel()
					continue
				}
				ackCtx, cancel := context.WithTimeout(ctx, settleTimeout)
				_ = receiver.CompleteMessage(ackCtx, msg, nil)
				cancel()
			}
		})
	}

	for ctx.Err() == nil {
		callCtx, cancel := context.WithTimeout(ctx, receiveWait)
		messages, err := receiver.ReceiveMessages(callCtx, batchSize, nil)
		cancel()

		if ctx.Err() != nil {
			break
		}

		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				continue
			}
			log.Printf("receive error: %v", err)
			continue
		}

		for _, m := range messages {
			select {
			case jobs <- m:
			case <-ctx.Done():
				close(jobs)
				wg.Wait()
				return ctx.Err()
			}
		}
	}
	close(jobs)
	wg.Wait()
	return ctx.Err()
}
