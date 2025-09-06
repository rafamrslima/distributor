package messaging

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

func sendMessage() {
	queue := os.Getenv("SERVICEBUS_QUEUE")

	client, err := getClient()
	if err != nil {
		log.Fatal(err)
	}

	sender, err := client.NewSender(queue, nil)
	if err != nil {
		log.Fatal(err)
	}

	msg := &azservicebus.Message{Body: []byte("Hello from Go!")}
	err = sender.SendMessage(context.Background(), msg, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Message sent")
}

func getClient() (*azservicebus.Client, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env")
	}

	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")

	client, err := azservicebus.NewClientFromConnectionString(connStr, nil)
	if err != nil {
		return nil, errors.New("SERVICEBUS_QUEUE is empty")
	}

	return client, nil
}

func getQueueName() (string, error) {
	queue := os.Getenv("SERVICEBUS_QUEUE")
	if queue == "" {
		fmt.Println("SERVICEBUS_QUEUE is empty")
		return "", errors.New("SERVICEBUS_QUEUE is empty")
	}
	return queue, nil
}

func StartMessageListener() error {
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
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = receiver.Close(ctx)
	}()

	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Println("Listening for messages... (Ctrl+C to quit)")

	jobs := make(chan *azservicebus.ReceivedMessage, 2*maxWorkers)

	var wg sync.WaitGroup
	for range maxWorkers {
		wg.Go(func() {
			for msg := range jobs {
				// 1) process (idempotent!)
				if err := core.Handle(msg); err != nil {
					// choose policy: abandon or dead-letter
					_ = receiver.AbandonMessage(appCtx, msg, nil)
					continue
				}
				// 2) settle (bounded)
				ackCtx, cancel := context.WithTimeout(appCtx, settleTimeout)
				_ = receiver.CompleteMessage(ackCtx, msg, nil)
				cancel()
			}
		})
	}

	for appCtx.Err() == nil {

		callCtx, cancel := context.WithTimeout(appCtx, receiveWait)
		messages, err := receiver.ReceiveMessages(callCtx, batchSize, nil)
		cancel()

		if appCtx.Err() != nil {
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
			case jobs <- m: // backpressure if workers busy
			case <-appCtx.Done():
				break
			}
		}
	}
	close(jobs)
	wg.Wait()
	return appCtx.Err()
}
