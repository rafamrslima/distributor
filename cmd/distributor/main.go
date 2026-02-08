package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rafamrslima/distributor/internal/db"
	"github.com/rafamrslima/distributor/internal/messaging"
)

func main() {
	fmt.Println("hello project")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env.")
	}

	_, err := db.GetDB()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := messaging.StartMessageListener(ctx)

		if err != nil {
			return
		}
	}()

	<-ctx.Done()
	wg.Wait()
	db.CloseDB()
}
